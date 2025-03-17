package family

import (
	"gorm.io/gorm"
	"time"
)

type FamilyInvitationStatus struct {
	StatusID   uint           `json:"status_id" gorm:"primaryKey"`
	StatusName string         `json:"status_name" gorm:"unique;not null"`
	CreatedAt  time.Time      `gorm:"autoCreateTime" json:"created_at,omitempty"`
	CreatedBy  string         `json:"created_by,omitempty"`
	UpdatedAt  time.Time      `gorm:"autoUpdateTime" json:"updated_at,omitempty"`
	UpdatedBy  string         `json:"updated_by,omitempty"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	DeletedBy  string         `json:"deleted_by,omitempty"`
}
