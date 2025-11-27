package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type SetDefaultViewRequest struct {
	DefaultViewID *string `json:"defaultViewId"` // nullable - can be null to clear default
}

// ValidateSetDefaultViewRequest validates and parses set default view request
func ValidateSetDefaultViewRequest(r *http.Request) (*SetDefaultViewRequest, error) {
	var req SetDefaultViewRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, fmt.Errorf("invalid request body")
	}

	return &req, nil
}