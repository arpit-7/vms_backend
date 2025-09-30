package handlers

import (
	"encoding/json"
	"net/http"

	"go-auth/config"
	"go-auth/db"
	"go-auth/models"
	"go-auth/utils"
	"golang.org/x/crypto/bcrypt"
)

// Request/Response structures
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	GroupId  int    `json:"groupId"`
	AreaName string `json:"areaName"`
	Role     string `json:"role"`
}

// AuthHandler handles JSON login requests
func AuthHandler(w http.ResponseWriter, r *http.Request) {
	utils.SetCORSHeaders(w)

	if r.Method == http.MethodOptions {
		return
	}

	if r.Method != http.MethodPost {
		utils.SendError(w, "Only POST method allowed", http.StatusMethodNotAllowed)
		return
	}

	var req LoginRequest
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

	// Create session JWT
	sessionToken, err := utils.CreateSessionJWT(user)
	if err != nil {
		utils.SendError(w, "Failed to create session", http.StatusInternalServerError)
		return
	}

	// Set session cookie
	utils.SetSessionCookie(w, sessionToken)

	// Return user data
	response := LoginResponse{
		ID:       user.ID,
		Username: user.Username,
		GroupId:  user.GroupId,
		AreaName: user.AreaName,
		Role:     user.Role,
	}

	utils.SendJSON(w, response, http.StatusOK)
}

// LoginFormHandler handles form-based login requests
func LoginFormHandler(w http.ResponseWriter, r *http.Request) {
	utils.SetCORSHeaders(w)

	if r.Method == http.MethodOptions {
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse form data
	err := r.ParseForm()
	if err != nil {
		http.Redirect(w, r, config.FRONTEND_URL+"/login?error=invalid_form", http.StatusTemporaryRedirect)
		return
	}

	username := r.FormValue("username")
	password := r.FormValue("password")

	// Validate input
	if err := utils.ValidateLoginRequest(username, password); err != nil {
		http.Redirect(w, r, config.FRONTEND_URL+"/login?error=missing_credentials", http.StatusTemporaryRedirect)
		return
	}

	// Find user
	var user models.User
	result := db.DB.Where("username = ?", username).First(&user)
	if result.Error != nil {
		http.Redirect(w, r, config.FRONTEND_URL+"/login?error=invalid_credentials", http.StatusTemporaryRedirect)
		return
	}

	// Verify password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		http.Redirect(w, r, config.FRONTEND_URL+"/login?error=invalid_credentials", http.StatusTemporaryRedirect)
		return
	}

	// Create session JWT
	sessionToken, err := utils.CreateSessionJWT(user)
	if err != nil {
		http.Redirect(w, r, config.FRONTEND_URL+"/login?error=session_failed", http.StatusTemporaryRedirect)
		return
	}

	// Set session cookie
	utils.SetSessionCookie(w, sessionToken)

	// Redirect to home
	http.Redirect(w, r, config.FRONTEND_URL+"/", http.StatusTemporaryRedirect)
}

// LogoutHandler handles logout requests
func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	utils.SetCORSHeaders(w)

	if r.Method == http.MethodOptions {
		return
	}

	// Clear session cookie
	utils.ClearSessionCookie(w)

	http.Redirect(w, r, config.FRONTEND_URL+"/login", http.StatusTemporaryRedirect)
}

// SessionHandler returns current session information
func SessionHandler(w http.ResponseWriter, r *http.Request) {
	utils.SetCORSHeaders(w)

	if r.Method == http.MethodOptions {
		return
	}

	if r.Method != http.MethodGet {
		utils.SendError(w, "Only GET method allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get session cookie
	cookie, err := utils.GetSessionCookie(r)
	if err != nil {
		utils.SendJSON(w, map[string]interface{}{"user": nil}, http.StatusOK)
		return
	}

	// Verify session JWT
	payload, err := utils.VerifySessionJWT(cookie.Value)
	if err != nil {
		utils.SendJSON(w, map[string]interface{}{"user": nil}, http.StatusOK)
		return
	}

	// Return session data
	utils.SendJSON(w, map[string]interface{}{
		"user": map[string]interface{}{
			"id":       payload["id"],
			"name":     payload["username"],
			"username": payload["username"],
			"groupId":  payload["groupId"],
			"areaName": payload["areaName"],
			"role":     payload["role"],
		},
	}, http.StatusOK)
}