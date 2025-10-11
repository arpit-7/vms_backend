package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"go-auth/db"
	"go-auth/models"
	"go-auth/utils"
)

// Request structures
type CreateViewGroupRequest struct {
	ID                   string   `json:"id"`
	Name                 string   `json:"name"`
	GroupID              int      `json:"groupId"`
	AreaName             string   `json:"areaName"`
	IsHQ                 bool     `json:"isHQ"`
	Cameras              []string `json:"cameras"`
	AutoRotationInterval *int     `json:"autoRotationInterval,omitempty"`
}

type UpdateViewGroupRequest struct {
	Name                 string   `json:"name,omitempty"`
	Cameras              []string `json:"cameras,omitempty"`
	AutoRotationInterval *int     `json:"autoRotationInterval,omitempty"`
}

// CreateViewGroupHandler creates a new view group
func CreateViewGroupHandler(w http.ResponseWriter, r *http.Request) {
	utils.SetCORSHeaders(w)

	if r.Method == http.MethodOptions {
		return
	}

	if r.Method != http.MethodPost {
		utils.SendError(w, "Only POST method allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get current user from session
	user, err := getUserFromSession(r)
	if err != nil {
		utils.SendError(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Check if user has permission to create view groups
	if user.Role == "Basic User" {
		utils.SendError(w, "Basic Users cannot create view groups", http.StatusForbidden)
		return
	}

	var req CreateViewGroupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate input
	if req.ID == "" || req.Name == "" || req.GroupID == 0 || req.AreaName == "" {
		utils.SendError(w, "ID, name, groupId, and areaName are required", http.StatusBadRequest)
		return
	}

	// Authorization check: Area Admin can only create for their area
	if user.Role == "Area Admin" && user.GroupId != req.GroupID {
		utils.SendError(w, "Area Admin can only create view groups for their own area", http.StatusForbidden)
		return
	}

	// Check if view group with this ID already exists
	var existingVG models.ViewGroup
	if db.DB.Where("id = ?", req.ID).First(&existingVG).Error == nil {
		utils.SendError(w, "View group with this ID already exists", http.StatusConflict)
		return
	}

	// Initialize cameras array if nil
	if req.Cameras == nil {
		req.Cameras = []string{}
	}

	// Create new view group
	viewGroup := models.ViewGroup{
		ID:                   req.ID,
		Name:                 req.Name,
		GroupID:              req.GroupID,
		AreaName:             req.AreaName,
		IsHQ:                 req.IsHQ,
		Cameras:              req.Cameras,
		AutoRotationInterval: req.AutoRotationInterval,
		CreatedBy:            user.Username,
		UpdatedBy:            user.Username,
		CreatedAt:            time.Now(),
		UpdatedAt:            time.Now(),
	}

	// Insert into database
	if err := db.DB.Create(&viewGroup).Error; err != nil {
		utils.SendError(w, "Failed to create view group", http.StatusInternalServerError)
		return
	}

	// Create audit record
	createAuditRecord(req.ID, "CREATE", user.Username, fmt.Sprintf(`{"name":"%s","groupId":%d}`, req.Name, req.GroupID))

	// Success response
	utils.SendJSON(w, map[string]interface{}{
		"message":   "View group created successfully",
		"viewGroup": viewGroup,
	}, http.StatusCreated)
}

// GetViewGroupsHandler retrieves view groups based on user role
func GetViewGroupsHandler(w http.ResponseWriter, r *http.Request) {
	utils.SetCORSHeaders(w)

	if r.Method == http.MethodOptions {
		return
	}

	if r.Method != http.MethodGet {
		utils.SendError(w, "Only GET method allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get current user from session
	user, err := getUserFromSession(r)
	if err != nil {
		utils.SendError(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var viewGroups []models.ViewGroup

	// Filter based on role
	if user.Role == "admin" {
		// Admin gets all view groups
		db.DB.Order("created_at DESC").Find(&viewGroups)
	} else {
		// Area Admin and Basic User get only their area's view groups
		db.DB.Where("group_id = ?", user.GroupId).Order("created_at DESC").Find(&viewGroups)
	}

	utils.SendJSON(w, viewGroups, http.StatusOK)
}

// UpdateViewGroupHandler updates an existing view group
func UpdateViewGroupHandler(w http.ResponseWriter, r *http.Request) {
	utils.SetCORSHeaders(w)

	if r.Method == http.MethodOptions {
		return
	}

	if r.Method != http.MethodPut {
		utils.SendError(w, "Only PUT method allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get current user from session
	user, err := getUserFromSession(r)
	if err != nil {
		utils.SendError(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Check if user has permission to update view groups
	if user.Role == "Basic User" {
		utils.SendError(w, "Basic Users cannot update view groups", http.StatusForbidden)
		return
	}

	// Extract view group ID from URL path (/view-groups/123)
	path := strings.TrimPrefix(r.URL.Path, "/view-groups/")
	viewGroupID := path
	if viewGroupID == "" {
		utils.SendError(w, "Invalid view group ID", http.StatusBadRequest)
		return
	}

	var req UpdateViewGroupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Find view group in database
	var viewGroup models.ViewGroup
	if db.DB.Where("id = ?", viewGroupID).First(&viewGroup).Error != nil {
		utils.SendError(w, "View group not found", http.StatusNotFound)
		return
	}

	// Debug logging for authorization
	fmt.Printf("Update attempt - User: %s, Role: %s, UserGroupId: %d, ViewGroupId: %d, IsHQ: %t\n", 
		user.Username, user.Role, user.GroupId, viewGroup.GroupID, viewGroup.IsHQ)

	// Authorization check: Area Admin can only update their area's view groups
	if user.Role == "Area Admin" && user.GroupId != viewGroup.GroupID {
		utils.SendError(w, "Area Admin can only update view groups in their own area", http.StatusForbidden)
		return
	}

	// Prepare update data
	updateData := make(map[string]interface{})
	var changes []string

	if req.Name != "" && req.Name != viewGroup.Name {
		updateData["name"] = req.Name
		changes = append(changes, fmt.Sprintf(`"name":"%s"`, req.Name))
	}

	if req.Cameras != nil {
		// Convert to JSON string to ensure proper database storage
		camerasJSON, err := json.Marshal(req.Cameras)
		if err != nil {
			utils.SendError(w, "Failed to serialize cameras", http.StatusInternalServerError)
			return
		}
		
		// Store as JSON string (database expects JSON type)
		updateData["cameras"] = string(camerasJSON)
		changes = append(changes, fmt.Sprintf(`"cameras":%s`, string(camerasJSON)))
	}

	// Always update auto_rotation_interval (even if null)
	updateData["auto_rotation_interval"] = req.AutoRotationInterval
	if req.AutoRotationInterval != nil {
		changes = append(changes, fmt.Sprintf(`"autoRotationInterval":%d`, *req.AutoRotationInterval))
	} else {
		changes = append(changes, `"autoRotationInterval":null`)
	}

	// Always update UpdatedBy and UpdatedAt
	updateData["updated_by"] = user.Username
	updateData["updated_at"] = time.Now()

	// Check if there's anything to update
	if len(changes) == 0 {
		utils.SendError(w, "No fields provided for update", http.StatusBadRequest)
		return
	}

	// Perform update (Last Write Wins - no version checking)
	if err := db.DB.Model(&viewGroup).Updates(updateData).Error; err != nil {
		fmt.Printf("Database update error: %v\n", err)
		utils.SendError(w, "Failed to update view group", http.StatusInternalServerError)
		return
	}

	// Create audit record
	changesJSON := fmt.Sprintf(`{%s}`, strings.Join(changes, ","))
	createAuditRecord(viewGroupID, "UPDATE", user.Username, changesJSON)

	// Fetch updated view group
	db.DB.Where("id = ?", viewGroupID).First(&viewGroup)

	// Success response
	utils.SendJSON(w, map[string]interface{}{
		"message":        "View group updated successfully",
		"viewGroup":      viewGroup,
		"updated_fields": len(changes),
	}, http.StatusOK)
}

// DeleteViewGroupHandler deletes a view group
func DeleteViewGroupHandler(w http.ResponseWriter, r *http.Request) {
	utils.SetCORSHeaders(w)

	if r.Method == http.MethodOptions {
		return
	}

	if r.Method != http.MethodDelete {
		utils.SendError(w, "Only DELETE method allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get current user from session
	user, err := getUserFromSession(r)
	if err != nil {
		utils.SendError(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Check if user has permission to delete view groups
	if user.Role == "Basic User" {
		utils.SendError(w, "Basic Users cannot delete view groups", http.StatusForbidden)
		return
	}

	// Extract view group ID from URL path
	path := strings.TrimPrefix(r.URL.Path, "/view-groups/")
	viewGroupID := path
	if viewGroupID == "" {
		utils.SendError(w, "Invalid view group ID", http.StatusBadRequest)
		return
	}

	// Find view group in database
	var viewGroup models.ViewGroup
	if db.DB.Where("id = ?", viewGroupID).First(&viewGroup).Error != nil {
		utils.SendError(w, "View group not found", http.StatusNotFound)
		return
	}

	// Authorization check: Area Admin can only delete their area's view groups
	if user.Role == "Area Admin" && user.GroupId != viewGroup.GroupID {
		utils.SendError(w, "Area Admin can only delete view groups in their own area", http.StatusForbidden)
		return
	}

	// Create audit record before deletion
	createAuditRecord(viewGroupID, "DELETE", user.Username, fmt.Sprintf(`{"name":"%s"}`, viewGroup.Name))

	// Delete view group
	if err := db.DB.Where("id = ?", viewGroupID).Delete(&models.ViewGroup{}).Error; err != nil {
		utils.SendError(w, "Failed to delete view group", http.StatusInternalServerError)
		return
	}

	// Success response
	utils.SendJSON(w, map[string]interface{}{
		"message":         "View group deleted successfully",
		"deleted_view":    viewGroup.Name,
		"deleted_view_id": viewGroupID,
	}, http.StatusOK)
}

// Helper function to create audit records
func createAuditRecord(viewGroupID, action, changedBy, changes string) {
	audit := models.ViewGroupAudit{
		ViewGroupID: viewGroupID,
		Action:      action,
		ChangedBy:   changedBy,
		Changes:     changes,
	}
	db.DB.Create(&audit)
}

// Helper function to get user from session (reuses your existing auth system)
func getUserFromSession(r *http.Request) (*models.User, error) {
	// Get session cookie using your existing utility
	cookie, err := utils.GetSessionCookie(r)
	if err != nil {
		return nil, fmt.Errorf("no session cookie found")
	}

	// Verify session JWT using your existing utility
	payload, err := utils.VerifySessionJWT(cookie.Value)
	if err != nil {
		return nil, fmt.Errorf("invalid session token")
	}

	// Extract username from payload
	username, ok := payload["username"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid session payload")
	}

	// Get user from database
	var user models.User
	if err := db.DB.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, fmt.Errorf("user not found")
	}

	return &user, nil
}