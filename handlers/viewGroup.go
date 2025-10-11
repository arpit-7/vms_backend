package handlers

import (
	"fmt"
	"net/http"

	"go-auth/utils"
)

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
	user, err := utils.GetUserFromSession(r)
	if err != nil {
		utils.SendError(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Validate request
	req, err := utils.ValidateCreateRequest(r)
	if err != nil {
		utils.SendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Authorization check
	if err := utils.CanCreateViewGroup(user, req.GroupID); err != nil {
		utils.SendError(w, err.Error(), http.StatusForbidden)
		return
	}

	// Check if view group already exists
	if utils.CheckViewGroupExists(req.ID) {
		utils.SendError(w, "View group with this ID already exists", http.StatusConflict)
		return
	}

	// Build view group from request
	viewGroup := utils.BuildViewGroupFromRequest(req, user.Username)

	// Insert into database
	if err := utils.CreateViewGroupInDB(viewGroup); err != nil {
		utils.SendError(w, "Failed to create view group", http.StatusInternalServerError)
		return
	}

	// Create audit record
	utils.CreateAuditRecord(req.ID, "CREATE", user.Username, 
		fmt.Sprintf(`{"name":"%s","groupId":%d}`, req.Name, req.GroupID))

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
	user, err := utils.GetUserFromSession(r)
	if err != nil {
		utils.SendError(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Get view groups based on user role
	viewGroups, err := utils.GetViewGroupsByUser(user)
	if err != nil {
		utils.SendError(w, "Failed to fetch view groups", http.StatusInternalServerError)
		return
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
	user, err := utils.GetUserFromSession(r)
	if err != nil {
		utils.SendError(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Extract view group ID from URL
	viewGroupID, err := utils.ParseViewGroupID(r)
	if err != nil {
		utils.SendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Validate request
	req, err := utils.ValidateUpdateRequest(r)
	if err != nil {
		utils.SendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Find view group in database
	viewGroup, err := utils.GetViewGroupByID(viewGroupID)
	if err != nil {
		utils.SendError(w, "View group not found", http.StatusNotFound)
		return
	}

	// Authorization check
	if err := utils.CanUpdateViewGroup(user, viewGroup); err != nil {
		utils.SendError(w, err.Error(), http.StatusForbidden)
		return
	}

	// Build update data
	updateData, changes, err := utils.BuildUpdateData(req, viewGroup, user.Username)
	if err != nil {
		utils.SendError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Check if there's anything to update
	if len(changes) == 0 {
		utils.SendError(w, "No fields provided for update", http.StatusBadRequest)
		return
	}

	// Perform update
	if err := utils.UpdateViewGroupInDB(viewGroupID, updateData); err != nil {
		utils.SendError(w, "Failed to update view group", http.StatusInternalServerError)
		return
	}

	// Create audit record
	changesJSON := utils.FormatChangesJSON(changes)
	utils.CreateAuditRecord(viewGroupID, "UPDATE", user.Username, changesJSON)

	// Fetch updated view group
	updatedViewGroup, _ := utils.GetViewGroupByID(viewGroupID)

	// Success response
	utils.SendJSON(w, map[string]interface{}{
		"message":        "View group updated successfully",
		"viewGroup":      updatedViewGroup,
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
	user, err := utils.GetUserFromSession(r)
	if err != nil {
		utils.SendError(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Extract view group ID from URL
	viewGroupID, err := utils.ParseViewGroupID(r)
	if err != nil {
		utils.SendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Find view group in database
	viewGroup, err := utils.GetViewGroupByID(viewGroupID)
	if err != nil {
		utils.SendError(w, "View group not found", http.StatusNotFound)
		return
	}

	// Authorization check
	if err := utils.CanDeleteViewGroup(user, viewGroup); err != nil {
		utils.SendError(w, err.Error(), http.StatusForbidden)
		return
	}

	// Create audit record before deletion
	utils.CreateAuditRecord(viewGroupID, "DELETE", user.Username, 
		fmt.Sprintf(`{"name":"%s"}`, viewGroup.Name))

	// Delete view group
	if err := utils.DeleteViewGroupFromDB(viewGroupID); err != nil {
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