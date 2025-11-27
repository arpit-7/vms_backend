package handlers

import (
	"net/http"

	"go-auth/utils"
)

// HandleDefaultView handles GET and PUT for default view
func HandleDefaultView(w http.ResponseWriter, r *http.Request) {
	utils.SetCORSHeaders(w)

	if r.Method == http.MethodOptions {
		return
	}

	switch r.Method {
	case http.MethodGet:
		GetDefaultViewHandler(w, r)
	case http.MethodPut:
		SetDefaultViewHandler(w, r)
	default:
		utils.SendError(w, "Only GET and PUT methods allowed", http.StatusMethodNotAllowed)
	}
}

// GetDefaultViewHandler retrieves the current user's default view
func GetDefaultViewHandler(w http.ResponseWriter, r *http.Request) {
	// Get current user from session
	user, err := utils.GetUserFromSession(r)
	if err != nil {
		utils.SendError(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Get user preference
	preference, err := utils.GetUserPreference(user.ID)
	if err != nil {
		// No preference set yet - return null
		utils.SendJSON(w, map[string]interface{}{
			"defaultViewId": nil,
		}, http.StatusOK)
		return
	}

	utils.SendJSON(w, map[string]interface{}{
		"defaultViewId": preference.DefaultViewID,
	}, http.StatusOK)
}

// SetDefaultViewHandler sets the default view for current user
func SetDefaultViewHandler(w http.ResponseWriter, r *http.Request) {
	// Get current user from session
	user, err := utils.GetUserFromSession(r)
	if err != nil {
		utils.SendError(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Validate request
	req, err := utils.ValidateSetDefaultViewRequest(r)
	if err != nil {
		utils.SendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	// If defaultViewID is provided (not null), verify the view exists and user has access
	if req.DefaultViewID != nil {
		viewGroup, err := utils.GetViewGroupByID(*req.DefaultViewID)
		if err != nil {
			utils.SendError(w, "View not found", http.StatusNotFound)
			return
		}

		// Check if user has access to this view
		if user.Role != "admin" && user.GroupId != viewGroup.GroupID {
			utils.SendError(w, "You don't have access to this view", http.StatusForbidden)
			return
		}
	}

	// Upsert user preference
	if err := utils.UpsertUserPreference(user.ID, user.Username, req.DefaultViewID); err != nil {
		utils.SendError(w, "Failed to set default view", http.StatusInternalServerError)
		return
	}

	utils.SendJSON(w, map[string]interface{}{
		"message":       "Default view updated successfully",
		"defaultViewId": req.DefaultViewID,
	}, http.StatusOK)
}