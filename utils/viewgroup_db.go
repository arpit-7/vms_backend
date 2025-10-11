package utils

import (
	"go-auth/db"
	"go-auth/models"
)

// GetViewGroupByID retrieves a view group by ID
func GetViewGroupByID(id string) (*models.ViewGroup, error) {
	var viewGroup models.ViewGroup
	if err := db.DB.Where("id = ?", id).First(&viewGroup).Error; err != nil {
		return nil, err
	}
	return &viewGroup, nil
}

// CheckViewGroupExists checks if a view group with given ID already exists
func CheckViewGroupExists(id string) bool {
	var viewGroup models.ViewGroup
	return db.DB.Where("id = ?", id).First(&viewGroup).Error == nil
}

// CreateViewGroupInDB creates a new view group in database
func CreateViewGroupInDB(viewGroup *models.ViewGroup) error {
	return db.DB.Create(viewGroup).Error
}

// UpdateViewGroupInDB updates a view group in database
func UpdateViewGroupInDB(id string, updateData map[string]interface{}) error {
	return db.DB.Model(&models.ViewGroup{}).Where("id = ?", id).Updates(updateData).Error
}

// DeleteViewGroupFromDB deletes a view group from database
func DeleteViewGroupFromDB(id string) error {
	return db.DB.Where("id = ?", id).Delete(&models.ViewGroup{}).Error
}

// GetViewGroupsByUser retrieves view groups based on user role
func GetViewGroupsByUser(user *models.User) ([]models.ViewGroup, error) {
	var viewGroups []models.ViewGroup
	query := db.DB.Order("created_at DESC")

	if user.Role != "admin" {
		// Area Admin and Basic User get only their area's view groups
		query = query.Where("group_id = ?", user.GroupId)
	}

	if err := query.Find(&viewGroups).Error; err != nil {
		return nil, err
	}

	return viewGroups, nil
}