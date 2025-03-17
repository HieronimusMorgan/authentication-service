package family

import (
	"gorm.io/gorm"
	"time"
)

type Family struct {
	FamilyID   uint           `gorm:"primaryKey" json:"family_id,omitempty"`
	FamilyName string         `gorm:"not null" json:"family_name,omitempty"`
	OwnerID    uint           `gorm:"not null" json:"owner_id,omitempty"`
	CreatedAt  time.Time      `gorm:"autoCreateTime" json:"created_at,omitempty"`
	CreatedBy  string         `json:"created_by,omitempty"`
	UpdatedAt  time.Time      `gorm:"autoUpdateTime" json:"updated_at,omitempty"`
	UpdatedBy  string         `json:"updated_by,omitempty"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	DeletedBy  string         `json:"deleted_by,omitempty"`
}
