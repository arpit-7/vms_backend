package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"go-auth/config"
	"go-auth/db"
	"go-auth/models"
	"go-auth/utils"
	"golang.org/x/crypto/bcrypt"
)

// Request/Response structures
type GenerateTokenRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type GenerateTokenResponse struct {
	Success   bool   `json:"success"`
	Token     string `json:"token"`
	LoginLink string `json:"loginLink"`
}

type VerifyTokenRequest struct {
	Token string `json:"token"`
}

type VerifyTokenResponse struct {
	Success  bool   `json:"success"`
	ID       uint   `json:"id"`
	Username string `json:"username"`
	GroupId  int    `json:"groupId"`
	AreaName string `json:"areaName"`
	Role     string `json:"role"`
}

// GenerateTokenHandler generates a login token
func GenerateTokenHandler(w http.ResponseWriter, r *http.Request) {
	utils.SetCORSHeaders(w)

	if r.Method == http.MethodOptions {
		return
	}

	if r.Method != http.MethodPost {
		utils.SendError(w, "Only POST method allowed", http.StatusMethodNotAllowed)
		return
	}

	var req GenerateTokenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate input
	if err := utils.ValidateLoginRequest(req.Username, req.Password); err != nil {
		utils.SendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Find user in database
	var user models.User
	result := db.DB.Where("username = ?", req.Username).First(&user)
	if result.Error != nil {
		utils.SendError(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Verify password
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		utils.SendError(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Generate JWT token
	expiresAt := time.Now().Add(time.Hour * 24 * 365 * config.TOKEN_EXPIRY_YEARS)

	jwtToken, err := utils.CreateTokenJWT(user, expiresAt)
	if err != nil {
		utils.SendError(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	// Store token in database
	tokenRecord := models.Token{
		UserID:    user.ID,
		Username:  user.Username,
		GroupId:   user.GroupId,
		AreaName:  user.AreaName,
		Role:      user.Role,
		Token:     jwtToken,
		IsUsed:    false,
		ExpiresAt: expiresAt,
	}

	if db.DB.Create(&tokenRecord).Error != nil {
		utils.SendError(w, "Failed to store token", http.StatusInternalServerError)
		return
	}

	// Return response
	response := GenerateTokenResponse{
		Success:   true,
		Token:     jwtToken,
		LoginLink: fmt.Sprintf("%s/verify?token=%s", config.BACKEND_URL, jwtToken),
	}

	utils.SendJSON(w, response, http.StatusOK)
}

// VerifyTokenHandler verifies a token (API endpoint)
func VerifyTokenHandler(w http.ResponseWriter, r *http.Request) {
	utils.SetCORSHeaders(w)

	if r.Method == http.MethodOptions {
		return
	}

	if r.Method != http.MethodPost {
		utils.SendError(w, "Only POST method allowed", http.StatusMethodNotAllowed)
		return
	}

	var req VerifyTokenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate token
	if err := utils.ValidateToken(req.Token); err != nil {
		utils.SendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Find token in database
	var tokenRecord models.Token
	result := db.DB.Where("token = ?", req.Token).First(&tokenRecord)
	if result.Error != nil {
		utils.SendError(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	// Check if token expired
	if time.Now().After(tokenRecord.ExpiresAt) {
		utils.SendError(w, "Token expired", http.StatusUnauthorized)
		return
	}

	// Verify JWT signature
	valid, err := utils.VerifyTokenJWT(req.Token)
	if err != nil || !valid {
		utils.SendError(w, "Invalid token signature", http.StatusUnauthorized)
		return
	}

	// Track first use
	if !tokenRecord.IsUsed {
		now := time.Now()
		tokenRecord.IsUsed = true
		tokenRecord.UsedAt = &now
		db.DB.Save(&tokenRecord)
	}

	// Return user data
	response := VerifyTokenResponse{
		Success:  true,
		ID:       tokenRecord.UserID,
		Username: tokenRecord.Username,
		GroupId:  tokenRecord.GroupId,
		AreaName: tokenRecord.AreaName,
		Role:     tokenRecord.Role,
	}

	utils.SendJSON(w, response, http.StatusOK)
}

// VerifyAndLoginHandler verifies token and creates session (Browser endpoint)
func VerifyAndLoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET method allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get token from URL query parameter
	token := r.URL.Query().Get("token")
	if err := utils.ValidateToken(token); err != nil {
		http.Redirect(w, r, config.FRONTEND_URL+"/login?error=missing_token", http.StatusTemporaryRedirect)
		return
	}

	// Find token in database
	var tokenRecord models.Token
	result := db.DB.Where("token = ?", token).First(&tokenRecord)
	if result.Error != nil {
		http.Redirect(w, r, config.FRONTEND_URL+"/login?error=invalid_token", http.StatusTemporaryRedirect)
		return
	}

	// Check if token expired
	if time.Now().After(tokenRecord.ExpiresAt) {
		http.Redirect(w, r, config.FRONTEND_URL+"/login?error=token_expired", http.StatusTemporaryRedirect)
		return
	}

	// Verify JWT signature
	valid, err := utils.VerifyTokenJWT(token)
	if err != nil || !valid {
		http.Redirect(w, r, config.FRONTEND_URL+"/login?error=invalid_signature", http.StatusTemporaryRedirect)
		return
	}

	// Track first use
	if !tokenRecord.IsUsed {
		now := time.Now()
		tokenRecord.IsUsed = true
		tokenRecord.UsedAt = &now
		db.DB.Save(&tokenRecord)
	}

	// Create NextAuth JWT token
	nextAuthToken, err := utils.CreateNextAuthJWT(tokenRecord)
	if err != nil {
		http.Error(w, "Failed to create session", http.StatusInternalServerError)
		return
	}

	// Set NextAuth session cookie
	utils.SetSessionCookie(w, nextAuthToken)

	// Redirect to home page
	http.Redirect(w, r, config.FRONTEND_URL+"/", http.StatusTemporaryRedirect)
}