package models

import (
	"time"
)

type FamilyMemberPermission struct {
	FamilyID     uint      `gorm:"primaryKey" json:"family_id,omitempty"`
	UserID       uint      `gorm:"primaryKey" json:"user_id,omitempty"`
	PermissionID uint      `gorm:"primaryKey" json:"permission_id,omitempty"`
	CreatedAt    time.Time `gorm:"autoCreateTime" json:"created_at,omitempty"`
	CreatedBy    string    `json:"created_by,omitempty"`
}
