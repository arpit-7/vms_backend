package utils

import (
	"go-auth/db"
	"go-auth/models"
)

// GetUserPreference retrieves user preference by user ID
func GetUserPreference(userID uint) (*models.UserPreference, error) {
	var preference models.UserPreference
	err := db.DB.Where("user_id = ?", userID).First(&preference).Error
	if err != nil {
		return nil, err
	}
	return &preference, nil
}

// CreateUserPreference creates a new user preference
func CreateUserPreference(preference *models.UserPreference) error {
	return db.DB.Create(preference).Error
}

// UpdateUserPreference updates user preference
func UpdateUserPreference(userID uint, defaultViewID *string) error {
	return db.DB.Model(&models.UserPreference{}).
		Where("user_id = ?", userID).
		Updates(map[string]interface{}{
			"default_view_id": defaultViewID,
		}).Error
}

// UpsertUserPreference creates or updates user preference
func UpsertUserPreference(userID uint, username string, defaultViewID *string) error {
	var preference models.UserPreference
	err := db.DB.Where("user_id = ?", userID).First(&preference).Error
	
	if err != nil {
		// Record doesn't exist, create new
		newPreference := models.UserPreference{
			UserID:        userID,
			Username:      username,
			DefaultViewID: defaultViewID,
		}
		return CreateUserPreference(&newPreference)
	}
	
	// Record exists, update
	return UpdateUserPreference(userID, defaultViewID)
}