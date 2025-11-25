package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"go-auth/db"
	"go-auth/models"
	"go-auth/utils"
)

type CreateCustomMapRequest struct {
	Name        string                  `json:"name"`
	Type        string                  `json:"type"`
	ImageData   string                  `json:"imageData,omitempty"`
	ImageWidth  int                     `json:"imageWidth"`
	ImageHeight int                     `json:"imageHeight"`
	TileURL     string                  `json:"tileUrlPattern,omitempty"`
	StyleURL    string                  `json:"styleUrl,omitempty"`
	Bounds      map[string]interface{}  `json:"bounds,omitempty"`
	GroupID     int                     `json:"groupId"`   
	AreaName    string                  `json:"areaName"`  
	Cameras     []CameraPositionRequest `json:"cameras"`
}

type CameraPositionRequest struct {
	CameraID   string `json:"cameraId"`
	CameraName string `json:"cameraName"`
	X          int    `json:"x"`
	Y          int    `json:"y"`
	Bearing    int    `json:"bearing"`
	FOV        int    `json:"fov"`
	Range      int    `json:"range"`
}


// GetCustomMapsHandler gets all maps filtered by user permissions
func GetCustomMapsHandler(w http.ResponseWriter, r *http.Request) {
	utils.SetCORSHeaders(w)
	if r.Method == http.MethodOptions {
		return
	}

	// Get user from session
	user, err := utils.GetUserFromSession(r)
	if err != nil {
		utils.SendError(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Get maps filtered by user permissions
	maps, err := utils.GetCustomMapsByUser(user)
	if err != nil {
		utils.SendError(w, "Failed to fetch maps", http.StatusInternalServerError)
		return
	}

	utils.SendJSON(w, maps, http.StatusOK)
}

// CreateCustomMapHandler creates a new custom map with permission checks
func CreateCustomMapHandler(w http.ResponseWriter, r *http.Request) {
	utils.SetCORSHeaders(w)
	if r.Method == http.MethodOptions {
		return
	}

	// Get user from session
	user, err := utils.GetUserFromSession(r)
	if err != nil {
		utils.SendError(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Parse request
	var req CreateCustomMapRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendError(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if req.Name == "" {
		utils.SendError(w, "Map name required", http.StatusBadRequest)
		return
	}

	// Check permissions
	if err := utils.CanCreateCustomMap(user, req.GroupID); err != nil {
		utils.SendError(w, err.Error(), http.StatusForbidden)
		return
	}

	// Convert bounds to JSON
	boundsJSON := ""
	if req.Bounds != nil {
		boundsBytes, _ := json.Marshal(req.Bounds)
		boundsJSON = string(boundsBytes)
	}

	// Create map
	customMap := models.CustomMap{
		Name:        req.Name,
		Type:        req.Type,
		ImageData:   req.ImageData,
		ImageWidth:  req.ImageWidth,
		ImageHeight: req.ImageHeight,
		Available:   true,
		TileURL:     req.TileURL,
		StyleURL:    req.StyleURL,
		BoundsJSON:  boundsJSON,
		GroupID:     req.GroupID,
		AreaName:    req.AreaName,
		CreatedBy:   user.Username, // Track who created it
		UpdatedBy:   user.Username,
	}

	if err := db.DB.Create(&customMap).Error; err != nil {
		utils.SendError(w, "Failed to create map", http.StatusInternalServerError)
		return
	}

	// Save cameras
	for _, cam := range req.Cameras {
		cameraPos := models.CameraPosition{
			CustomMapID: customMap.ID,
			CameraID:    cam.CameraID,
			CameraName:  cam.CameraName,
			X:           cam.X,
			Y:           cam.Y,
			Bearing:     cam.Bearing,
			FOV:         cam.FOV,
			Range:       cam.Range,
		}
		db.DB.Create(&cameraPos)
	}

	utils.SendJSON(w, map[string]interface{}{
		"message": "Map created",
		"map":     customMap,
	}, http.StatusCreated)
}

// GetCustomMapHandler gets one map with permission check
func GetCustomMapHandler(w http.ResponseWriter, r *http.Request) {
	utils.SetCORSHeaders(w)
	if r.Method == http.MethodOptions {
		return
	}

	// Get user from session
	user, err := utils.GetUserFromSession(r)
	if err != nil {
		utils.SendError(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	mapID, err := parseMapID(r)
	if err != nil {
		utils.SendError(w, "Invalid map ID", http.StatusBadRequest)
		return
	}

	var customMap models.CustomMap
	if err := db.DB.First(&customMap, mapID).Error; err != nil {
		utils.SendError(w, "Map not found", http.StatusNotFound)
		return
	}

	// Check if user has access to this map
	if user.Role != "admin" && customMap.GroupID != user.GroupId {
		utils.SendError(w, "Access denied", http.StatusForbidden)
		return
	}

	var cameras []models.CameraPosition
	db.DB.Where("custom_map_id = ?", mapID).Find(&cameras)

	utils.SendJSON(w, map[string]interface{}{
		"map":     customMap,
		"cameras": cameras,
	}, http.StatusOK)
}

// UpdateCustomMapHandler updates a map with permission checks
func UpdateCustomMapHandler(w http.ResponseWriter, r *http.Request) {
	utils.SetCORSHeaders(w)
	if r.Method == http.MethodOptions {
		return
	}

	// Get user from session
	user, err := utils.GetUserFromSession(r)
	if err != nil {
		utils.SendError(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	mapID, err := parseMapID(r)
	if err != nil {
		utils.SendError(w, "Invalid map ID", http.StatusBadRequest)
		return
	}

	var req CreateCustomMapRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendError(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Check if map exists
	var existing models.CustomMap
	if err := db.DB.First(&existing, mapID).Error; err != nil {
		utils.SendError(w, "Map not found", http.StatusNotFound)
		return
	}

	// Check permissions
	if err := utils.CanUpdateCustomMap(user, &existing); err != nil {
		utils.SendError(w, err.Error(), http.StatusForbidden)
		return
	}

	// Convert bounds
	boundsJSON := ""
	if req.Bounds != nil {
		boundsBytes, _ := json.Marshal(req.Bounds)
		boundsJSON = string(boundsBytes)
	}

	// Update map
	db.DB.Model(&existing).Updates(models.CustomMap{
		Name:        req.Name,
		Type:        req.Type,
		ImageData:   req.ImageData,
		ImageWidth:  req.ImageWidth,
		ImageHeight: req.ImageHeight,
		TileURL:     req.TileURL,
		StyleURL:    req.StyleURL,
		BoundsJSON:  boundsJSON,
		GroupID:     req.GroupID,
		AreaName:    req.AreaName,
		UpdatedBy:   user.Username, // Track who updated it
	})

	// Delete old cameras and add new ones
	db.DB.Where("custom_map_id = ?", mapID).Delete(&models.CameraPosition{})

	for _, cam := range req.Cameras {
		cameraPos := models.CameraPosition{
			CustomMapID: mapID,
			CameraID:    cam.CameraID,
			CameraName:  cam.CameraName,
			X:           cam.X,
			Y:           cam.Y,
			Bearing:     cam.Bearing,
			FOV:         cam.FOV,
			Range:       cam.Range,
		}
		db.DB.Create(&cameraPos)
	}

	utils.SendJSON(w, map[string]interface{}{
		"message": "Map updated",
	}, http.StatusOK)
}

// DeleteCustomMapHandler deletes a map with permission checks
func DeleteCustomMapHandler(w http.ResponseWriter, r *http.Request) {
	utils.SetCORSHeaders(w)
	if r.Method == http.MethodOptions {
		return
	}

	// Get user from session
	user, err := utils.GetUserFromSession(r)
	if err != nil {
		utils.SendError(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	mapID, err := parseMapID(r)
	if err != nil {
		utils.SendError(w, "Invalid map ID", http.StatusBadRequest)
		return
	}

	var customMap models.CustomMap
	if err := db.DB.First(&customMap, mapID).Error; err != nil {
		utils.SendError(w, "Map not found", http.StatusNotFound)
		return
	}

	// Check permissions
	if err := utils.CanDeleteCustomMap(user, &customMap); err != nil {
		utils.SendError(w, err.Error(), http.StatusForbidden)
		return
	}

	// Delete cameras first, then map
	db.DB.Where("custom_map_id = ?", mapID).Delete(&models.CameraPosition{})
	db.DB.Delete(&customMap)

	utils.SendJSON(w, map[string]interface{}{
		"message": "Map deleted",
	}, http.StatusOK)
}

// Helper to parse map ID from URL
func parseMapID(r *http.Request) (uint, error) {
	path := strings.TrimPrefix(r.URL.Path, "/custom-maps/")
	id, err := strconv.ParseUint(path, 10, 32)
	if err != nil {
		return 0, fmt.Errorf("invalid ID")
	}
	return uint(id), nil
}