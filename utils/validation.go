package utils

import "fmt"

// ValidateLoginRequest validates login credentials
func ValidateLoginRequest(username, password string) error {
	if username == "" {
		return fmt.Errorf("username is required")
	}
	if password == "" {
		return fmt.Errorf("password is required")
	}
	return nil
}

// validate roles
func ValidateRole(role string) error {
	allowedRole := []string{"Basic User", "Area Admin", "admin"}
	for _, r := range allowedRole {
		if role == r {
			return nil
		}
	}
	return fmt.Errorf("invalid role: must be one of 'Basic User', 'Area Admin', or 'admin'")
}

// ValidateCreateUserRequest validates user creation data
func ValidateCreateUserRequest(username, password string, groupId int, areaName, role string) error {
	if username == "" {
		return fmt.Errorf("username is required")
	}
	if password == "" {
		return fmt.Errorf("password is required")
	}
	if groupId == 0 {
		return fmt.Errorf("groupId is required")
	}
	if areaName == "" {
		return fmt.Errorf("areaName is required")
	}
	if role == "" {
		return fmt.Errorf("role is required")
	}
	return nil
}

// ValidateToken validates token string
func ValidateToken(token string) error {
	if token == "" {
		return fmt.Errorf("token is required")
	}
	return nil
}