package models

import "gorm.io/gorm"

type UserPreference struct {
	gorm.Model
	UserID        uint    `gorm:"uniqueIndex" json:"userId"`
	DefaultViewID *string `json:"defaultViewId"` 
	Username      string  `gorm:"index" json:"username"`
}