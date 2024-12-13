package models

import "time"

type RoleResource struct {
	RoleID     uint       `gorm:"primaryKey" json:"role_id,omitempty"`
	ResourceID uint       `gorm:"primaryKey" json:"resource_id,omitempty"`
	Role       Role       `gorm:"foreignKey:RoleID;constraint:OnDelete:CASCADE" json:"role,omitempty"`
	Resource   Resource   `gorm:"foreignKey:ResourceID;constraint:OnDelete:CASCADE" json:"resource,omitempty"`
	CreatedAt  time.Time  `gorm:"autoCreateTime" json:"created_at,omitempty"`
	CreatedBy  string     `json:"created_by,omitempty"`
	UpdatedAt  time.Time  `gorm:"autoUpdateTime" json:"updated_at,omitempty"`
	UpdatedBy  string     `json:"updated_by,omitempty"`
	DeletedAt  *time.Time `gorm:"index" json:"deleted_at,omitempty"` // Soft delete with timestamp
	DeletedBy  string     `json:"deleted_by,omitempty"`
}
