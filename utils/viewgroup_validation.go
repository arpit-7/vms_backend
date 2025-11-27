package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type CreateViewGroupRequest struct {
	ID                   string   `json:"id"`
	Name                 string   `json:"name"`
	GroupID              int      `json:"groupId"`
	AreaName             string   `json:"areaName"`
	IsHQ                 bool     `json:"isHQ"`
	Cameras              []string `json:"cameras"`
	CamerasMetadata      []CameraMetadataRequest `json:"camerasMetadata"`
	AutoRotationInterval *int     `json:"autoRotationInterval,omitempty"`
}


type UpdateViewGroupRequest struct {
	Name                 string   `json:"name,omitempty"`
	Cameras              []string `json:"cameras,omitempty"`
	CamerasMetadata      []CameraMetadataRequest `json:"camerasMetadata,omitempty"`
	AutoRotationInterval *int     `json:"autoRotationInterval,omitempty"`
}

type CameraMetadataRequest struct {
    ID      string `json:"id"`
    Name    string `json:"name"`
    GroupID int    `json:"groupId"`
}

// validates and parses create view group request
func ValidateCreateRequest(r *http.Request) (*CreateViewGroupRequest, error) {
	var req CreateViewGroupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, fmt.Errorf("invalid request body")
	}

	// Validate required fields
	if req.ID == "" || req.Name == "" || req.GroupID == 0 || req.AreaName == "" {
		return nil, fmt.Errorf("ID, name, groupId, and areaName are required")
	}

	// Initialize cameras array if nil
	if req.Cameras == nil {
		req.Cameras = []string{}
	}

	return &req, nil
}

//  validates and parses update view group request
func ValidateUpdateRequest(r *http.Request) (*UpdateViewGroupRequest, error) {
	var req UpdateViewGroupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, fmt.Errorf("invalid request body")
	}

	return &req, nil
}