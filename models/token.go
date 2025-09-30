package models

import (
	"time"
	"gorm.io/gorm"
)

type Token struct {
	gorm.Model
	UserID    uint       `gorm:"column:user_id" json:"userId"`
	Username  string     `gorm:"column:username;type:varchar(255)" json:"username"`
	GroupId   int        `gorm:"column:group_id" json:"groupId"`
	AreaName  string     `gorm:"column:area_name;type:varchar(255)" json:"areaName"`
	Role      string     `gorm:"column:role;type:varchar(100)" json:"role"`
	Token     string     `gorm:"column:token;type:text;not null" json:"token"`
	IsUsed    bool       `gorm:"column:is_used;default:false" json:"isUsed"`
	UsedAt    *time.Time `gorm:"column:used_at" json:"usedAt,omitempty"`
	ExpiresAt time.Time  `gorm:"column:expires_at" json:"expiresAt"`
}