package family

import (
	"gorm.io/gorm"
	"time"
)

type FamilyMemberPermission struct {
	FamilyID     uint           `gorm:"primaryKey" json:"family_id,omitempty"`
	UserID       uint           `gorm:"primaryKey" json:"user_id,omitempty"`
	PermissionID uint           `gorm:"primaryKey" json:"permission_id,omitempty"`
	CreatedAt    time.Time      `gorm:"autoCreateTime" json:"created_at,omitempty"`
	CreatedBy    string         `json:"created_by,omitempty"`
	UpdatedAt    time.Time      `gorm:"autoUpdateTime" json:"updated_at,omitempty"`
	UpdatedBy    string         `json:"updated_by,omitempty"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	DeletedBy    string         `json:"deleted_by,omitempty"`
}
