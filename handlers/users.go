package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"go-auth/db"
	"go-auth/models"
	"go-auth/utils"
	"golang.org/x/crypto/bcrypt"
)

// Request structures
type CreateUserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	GroupId  int    `json:"groupId"`
	AreaName string `json:"areaName"`
	Role     string `json:"role"`
}

type UpdateUserRequest struct {
	Password string `json:"password,omitempty"`
	GroupId  int    `json:"groupId,omitempty"`
	AreaName string `json:"areaName,omitempty"`
	Role     string `json:"role,omitempty"`
}

// CreateUserHandler creates a new user
func CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	utils.SetCORSHeaders(w)

	if r.Method == http.MethodOptions {
		return
	}

	if r.Method != http.MethodPost {
		utils.SendError(w, "Only POST method allowed", http.StatusMethodNotAllowed)
		return
	}

	var req CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate input
	if err := utils.ValidateCreateUserRequest(req.Username, req.Password, req.GroupId, req.AreaName, req.Role); err != nil {
		utils.SendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Check for duplicate username
	var existingUser models.User
	result := db.DB.Unscoped().Where("username = ?", req.Username).First(&existingUser)
	if result.Error == nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(map[string]string{
			"error":   "Username already exists",
			"message": "User with this username is already registered",
		})
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		utils.SendError(w, "Failed to process password", http.StatusInternalServerError)
		return
	}

	// Create new user
	user := models.User{
		Username: req.Username,
		Password: string(hashedPassword),
		GroupId:  req.GroupId,
		AreaName: req.AreaName,
		Role:     req.Role,
	}

	// Insert into database
	createResult := db.DB.Create(&user)
	if createResult.Error != nil {
		utils.SendError(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	// Success response
	utils.SendJSON(w, map[string]interface{}{
		"message": "User created successfully!",
		"user": map[string]interface{}{
			"id":       user.ID,
			"username": user.Username,
			"groupId":  user.GroupId,
			"areaName": user.AreaName,
			"role":     user.Role,
		},
	}, http.StatusOK)
}

// GetUsersHandler retrieves all users
func GetUsersHandler(w http.ResponseWriter, r *http.Request) {
	utils.SetCORSHeaders(w)

	if r.Method == http.MethodOptions {
		return
	}

	if r.Method != http.MethodGet {
		utils.SendError(w, "Only GET method allowed", http.StatusMethodNotAllowed)
		return
	}

	var users []models.User
	// Don't return passwords in the response
	db.DB.Select("id, username, group_id, area_name, role, created_at, updated_at").Find(&users)

	utils.SendJSON(w, users, http.StatusOK)
}

// UpdateUserHandler updates an existing user
func UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	utils.SetCORSHeaders(w)

	if r.Method == http.MethodOptions {
		return
	}

	if r.Method != http.MethodPut {
		utils.SendError(w, "Only PUT method allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract user ID from URL path (/users/5)
	path := strings.TrimPrefix(r.URL.Path, "/users/")
	userID, err := strconv.ParseUint(path, 10, 32)
	if err != nil {
		utils.SendError(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	var req UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Find user in database
	var user models.User
	if db.DB.First(&user, userID).Error != nil {
		utils.SendError(w, "User not found", http.StatusNotFound)
		return
	}

	// Prepare update data
	updateData := make(map[string]interface{})

	// Only update fields that are provided
	if req.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			utils.SendError(w, "Failed to process password", http.StatusInternalServerError)
			return
		}
		updateData["password"] = string(hashedPassword)
	}
	if req.GroupId != 0 {
		updateData["group_id"] = req.GroupId
	}
	if req.AreaName != "" {
		updateData["area_name"] = req.AreaName
	}
	if req.Role != "" {
		updateData["role"] = req.Role
	}

	// Check if there's anything to update
	if len(updateData) == 0 {
		utils.SendError(w, "No fields provided for update", http.StatusBadRequest)
		return
	}

	// Perform update
	result := db.DB.Model(&user).Updates(updateData)
	if result.Error != nil {
		utils.SendError(w, "Failed to update user", http.StatusInternalServerError)
		return
	}

	// Success response
	utils.SendJSON(w, map[string]interface{}{
		"message":        "User updated successfully!",
		"updated_fields": len(updateData),
	}, http.StatusOK)
}

// DeleteUserHandler deletes a user
func DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	utils.SetCORSHeaders(w)

	if r.Method == http.MethodOptions {
		return
	}

	if r.Method != http.MethodDelete {
		utils.SendError(w, "Only DELETE method allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract user ID from URL path (/users/5)
	path := strings.TrimPrefix(r.URL.Path, "/users/")
	userID, err := strconv.ParseUint(path, 10, 32)
	if err != nil {
		utils.SendError(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	// Find user in database
	var user models.User
	if db.DB.Unscoped().First(&user, userID).Error != nil{
		utils.SendError(w, "User not found", http.StatusNotFound)
		return
	}

	// Delete user
	if err := db.DB.Exec("DELETE FROM users WHERE id = ?", user.ID).Error; err != nil {
		utils.SendError(w, "Failed to delete user", http.StatusInternalServerError)
		return
	}

	// Success response
	utils.SendJSON(w, map[string]interface{}{
		"message":      "User deleted successfully!",
		"deleted_user": user.Username,
	}, http.StatusOK)
}