package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"go-auth/db"
	"go-auth/models"
)

// GetUserFromSession retrieves user from session cookie
func GetUserFromSession(r *http.Request) (*models.User, error) {
	// Get session cookie
	cookie, err := GetSessionCookie(r)
	if err != nil {
		return nil, fmt.Errorf("no session cookie found")
	}

	// Verify session JWT
	payload, err := VerifySessionJWT(cookie.Value)
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

// CreateAuditRecord creates an audit record for view group changes
func CreateAuditRecord(viewGroupID, action, changedBy, changes string) {
	audit := models.ViewGroupAudit{
		ViewGroupID: viewGroupID,
		Action:      action,
		ChangedBy:   changedBy,
		Changes:     changes,
	}
	db.DB.Create(&audit)
}

// ParseViewGroupID extracts view group ID from URL path
func ParseViewGroupID(r *http.Request) (string, error) {
	path := strings.TrimPrefix(r.URL.Path, "/view-groups/")
	viewGroupID := path
	if viewGroupID == "" || viewGroupID == "/view-groups" {
		return "", fmt.Errorf("invalid view group ID")
	}
	return viewGroupID, nil
}

// BuildViewGroupFromRequest creates a ViewGroup model from request
func BuildViewGroupFromRequest(req *CreateViewGroupRequest, username string) *models.ViewGroup {
	return &models.ViewGroup{
		ID:                   req.ID,
		Name:                 req.Name,
		GroupID:              req.GroupID,
		AreaName:             req.AreaName,
		IsHQ:                 req.IsHQ,
		Cameras:              req.Cameras,
		AutoRotationInterval: req.AutoRotationInterval,
		CreatedBy:            username,
		UpdatedBy:            username,
		CreatedAt:            time.Now(),
		UpdatedAt:            time.Now(),
	}
}

// BuildUpdateData prepares update data map and change log from request
func BuildUpdateData(req *UpdateViewGroupRequest, viewGroup *models.ViewGroup, username string) (map[string]interface{}, []string, error) {
	updateData := make(map[string]interface{})
	var changes []string

	// Update name if provided and different
	if req.Name != "" && req.Name != viewGroup.Name {
		updateData["name"] = req.Name
		changes = append(changes, fmt.Sprintf(`"name":"%s"`, req.Name))
	}

	// Update cameras if provided
	if req.Cameras != nil {
		// Convert to JSON string for database storage
		camerasJSON, err := json.Marshal(req.Cameras)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to serialize cameras")
		}

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

	// Always update metadata
	updateData["updated_by"] = username
	updateData["updated_at"] = time.Now()

	return updateData, changes, nil
}

// FormatChangesJSON formats changes array into JSON string
func FormatChangesJSON(changes []string) string {
	return fmt.Sprintf(`{%s}`, strings.Join(changes, ","))
}