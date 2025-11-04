package utils

import (
	"fmt"
	"go-auth/models"
)


func CanCreateViewGroup(user *models.User, targetGroupId int) error {
	// Basic Users cannot create view groups
	if user.Role == "Basic User" {
		return fmt.Errorf("basic users cannot create view groups")
	}

	// Area Admin can only create for their area
	if user.Role == "Area Admin" && user.GroupId != targetGroupId {
		return fmt.Errorf("area admin can only create view groups for their own area")
	}

	return nil
}

//  update view groups
func CanUpdateViewGroup(user *models.User, viewGroup *models.ViewGroup) error {
	// Basic Users cannot update view groups
	if user.Role == "Basic User" {
		return fmt.Errorf("basic users cannot update view groups")
	}

	// Area Admin can only update their area's view groups
	if user.Role == "Area Admin" && user.GroupId != viewGroup.GroupID {
		return fmt.Errorf("area admin can only update view groups in their own area")
	}

	return nil
}

//  delete view groups
func CanDeleteViewGroup(user *models.User, viewGroup *models.ViewGroup) error {
	// Basic Users cannot delete view groups
	if user.Role == "Basic User" {
		return fmt.Errorf("basic Users cannot delete view groups")
	}

	// Area Admin can only delete their area's view groups
	if user.Role == "Area Admin" && user.GroupId != viewGroup.GroupID {
		return fmt.Errorf("area admin can only delete view groups in their own area")
	}

	return nil
}