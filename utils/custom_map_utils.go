package utils

import (
	"fmt"
	"go-auth/db"
	"go-auth/models"
)

// GetCustomMapsByUser retrieves custom maps based on user role and permissions
func GetCustomMapsByUser(user *models.User) ([]models.CustomMap, error) {
	var customMaps []models.CustomMap
	query := db.DB.Order("created_at DESC")

	if user.Role != "admin" {
		query = query.Where("group_id = ?", user.GroupId)
	}

	if err := query.Find(&customMaps).Error; err != nil {
		return nil, err
	}

	return customMaps, nil
}

// CanCreateCustomMap checks if user can create a custom map for target area
func CanCreateCustomMap(user *models.User, targetGroupId int) error {
	if user.Role == "Basic User" {
		return fmt.Errorf("basic users cannot create custom maps")
	}

	if user.Role == "Area Admin" && user.GroupId != targetGroupId {
		return fmt.Errorf("area admin can only create maps for their own area")
	}

	return nil
}

// CanUpdateCustomMap checks if user can update a custom map
func CanUpdateCustomMap(user *models.User, customMap *models.CustomMap) error {
	if user.Role == "Basic User" {
		return fmt.Errorf("basic users cannot update custom maps")
	}

	if user.Role == "Area Admin" && user.GroupId != customMap.GroupID {
		return fmt.Errorf("area admin can only update maps in their own area")
	}

	return nil
}

// CanDeleteCustomMap checks if user can delete a custom map
func CanDeleteCustomMap(user *models.User, customMap *models.CustomMap) error {
	if user.Role == "Basic User" {
		return fmt.Errorf("basic users cannot delete custom maps")
	}

	if user.Role == "Area Admin" && user.GroupId != customMap.GroupID {
		return fmt.Errorf("area admin can only delete maps in their own area")
	}

	return nil
}