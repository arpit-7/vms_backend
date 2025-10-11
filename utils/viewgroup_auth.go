package utils

import (
	"fmt"
	"go-auth/models"
)

// CanCreateViewGroup checks if user has permission to create view groups
func CanCreateViewGroup(user *models.User, targetGroupId int) error {
	// Basic Users cannot create view groups
	if user.Role == "Basic User" {
		return fmt.Errorf("Basic Users cannot create view groups")
	}

	// Area Admin can only create for their area
	if user.Role == "Area Admin" && user.GroupId != targetGroupId {
		return fmt.Errorf("Area Admin can only create view groups for their own area")
	}

	return nil
}

// CanUpdateViewGroup checks if user has permission to update view groups
func CanUpdateViewGroup(user *models.User, viewGroup *models.ViewGroup) error {
	// Basic Users cannot update view groups
	if user.Role == "Basic User" {
		return fmt.Errorf("Basic Users cannot update view groups")
	}

	// Area Admin can only update their area's view groups
	if user.Role == "Area Admin" && user.GroupId != viewGroup.GroupID {
		return fmt.Errorf("Area Admin can only update view groups in their own area")
	}

	return nil
}

// CanDeleteViewGroup checks if user has permission to delete view groups
func CanDeleteViewGroup(user *models.User, viewGroup *models.ViewGroup) error {
	// Basic Users cannot delete view groups
	if user.Role == "Basic User" {
		return fmt.Errorf("Basic Users cannot delete view groups")
	}

	// Area Admin can only delete their area's view groups
	if user.Role == "Area Admin" && user.GroupId != viewGroup.GroupID {
		return fmt.Errorf("Area Admin can only delete view groups in their own area")
	}

	return nil
}