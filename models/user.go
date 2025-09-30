package models

import "gorm.io/gorm"

type User struct{
	gorm.Model
	Username string `gorm:"unique;not null" json:"username"`
	Password string `gorm:"not null" json:"password,omitempty"`
	GroupId int `json:"groupId"` 
	AreaName string `json:"areaName"`
	Role string `json:"role"`
}