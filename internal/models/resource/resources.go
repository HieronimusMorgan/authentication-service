package resource

import (
	"gorm.io/gorm"
	"time"
)

type Resource struct {
	ResourceID  uint           `gorm:"primaryKey" json:"resource_id,omitempty"`
	Name        string         `gorm:"unique;not null" json:"name,omitempty"`
	Description string         `json:"description,omitempty"`
	CreatedAt   time.Time      `gorm:"autoCreateTime" json:"created_at,omitempty"`
	CreatedBy   string         `json:"created_by,omitempty"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime" json:"updated_at,omitempty"`
	UpdatedBy   string         `json:"updated_by,omitempty"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"` // Soft delete with timestamp
	DeletedBy   string         `json:"deleted_by,omitempty"`
}
