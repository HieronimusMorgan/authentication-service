package family

import (
	"gorm.io/gorm"
	"time"
)

type FamilyPermission struct {
	PermissionID   uint           `gorm:"primaryKey" json:"permission_id,omitempty"`
	PermissionName string         `gorm:"unique;not null" json:"permission_name,omitempty"`
	Description    string         `json:"description,omitempty"`
	CreatedAt      time.Time      `gorm:"autoCreateTime" json:"created_at,omitempty"`
	CreatedBy      string         `json:"created_by,omitempty"`
	UpdatedAt      time.Time      `gorm:"autoUpdateTime" json:"updated_at,omitempty"`
	UpdatedBy      string         `json:"updated_by,omitempty"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	DeletedBy      string         `json:"deleted_by,omitempty"`
}
