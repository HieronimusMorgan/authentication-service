package models

import (
	"gorm.io/gorm"
	"time"
)

type UserResource struct {
	UserID     uint           `gorm:"primaryKey" json:"user_id,omitempty"`
	ResourceID uint           `gorm:"primaryKey" json:"resource_id,omitempty"`
	User       Users          `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"user,omitempty"`
	Resource   Resource       `gorm:"foreignKey:ResourceID;constraint:OnDelete:CASCADE" json:"resource,omitempty"`
	CreatedAt  time.Time      `gorm:"autoCreateTime" json:"created_at,omitempty"`
	CreatedBy  string         `json:"created_by,omitempty"`
	UpdatedAt  time.Time      `gorm:"autoUpdateTime" json:"updated_at,omitempty"`
	UpdatedBy  string         `json:"updated_by,omitempty"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	DeletedBy  string         `json:"deleted_by,omitempty"`
}
