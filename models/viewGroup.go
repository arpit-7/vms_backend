package models

import (
	"time"
	"gorm.io/gorm"
	"database/sql/driver"
	"encoding/json"
	"errors"
)

type StringArray []string

func (s *StringArray) Scan(value interface{}) error {
	if value == nil {
		*s  = []string{}
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("failed to unmarshal StringArray value ")
	}
	return json.Unmarshal(bytes, s)
}

func (s StringArray)  Value() (driver.Value, error) {
	if len(s) == 0 {
		return "[]",nil
	}
	return json.Marshal(s)
}

type ViewGroup struct {
	ID                   string      `gorm:"primaryKey" json:"id"`
	Name                 string      `gorm:"not null" json:"name"`
	GroupID              int         `gorm:"not null;index" json:"groupId"`
	AreaName             string      `json:"areaName"`
    IsHQ                 bool        `gorm:"default:false" json:"isHQ"`
	Cameras              StringArray `gorm:"type:json" json:"cameras"`
	AutoRotationInterval *int        `json:"autoRotationInterval,omitempty"`
	CreatedBy            string      `json:"createdBy"`
	UpdatedBy            string      `json:"updatedBy,omitempty"`
	CreatedAt            time.Time   `json:"createdAt"`
	UpdatedAt            time.Time   `json:"updatedAt"`
}

//store trail of changes

type ViewGroupAudit struct {
	gorm.Model
	ViewGroupID string `gorm:"index" json:"viewGroupId"`
	Action      string `json:"action"` // create , update , delete etc.
	ChangedBy   string `json:"changedBy"`
	Changes     string `gorm:"type:json" json:"changes"` //json styring of what changed 
}