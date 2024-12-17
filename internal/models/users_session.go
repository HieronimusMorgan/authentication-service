package models

import (
	"gorm.io/gorm"
	"time"
)

type UserSession struct {
	UserSessionID uint           `gorm:"primaryKey;autoIncrement" json:"user_session_id"`
	UserID        uint           `gorm:"not null;index" json:"user_id"` // Foreign key reference to users table
	SessionToken  string         `gorm:"unique;not null" json:"session_token"`
	RefreshToken  string         `gorm:"unique" json:"refresh_token"`
	IPAddress     string         `json:"ip_address"`
	UserAgent     string         `gorm:"type:text" json:"user_agent"`
	LoginTime     time.Time      `gorm:"default:CURRENT_TIMESTAMP" json:"login_time"`
	ExpiresAt     time.Time      `gorm:"not null" json:"expires_at"`
	LogoutTime    *time.Time     `gorm:"null" json:"logout_time"` // Nullable field
	IsActive      bool           `gorm:"default:true" json:"is_active"`
	CreatedAt     time.Time      `gorm:"autoCreateTime" json:"created_at,omitempty"`
	CreatedBy     string         `json:"created_by,omitempty"`
	UpdatedAt     time.Time      `gorm:"autoUpdateTime" json:"updated_at,omitempty"`
	UpdatedBy     string         `json:"updated_by,omitempty"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	DeletedBy     string         `json:"deleted_by,omitempty"`
}
