package models

import (
	"gorm.io/gorm"
	"time"
)

type Users struct {
	UserID         uint           `gorm:"primaryKey" json:"user_id,omitempty"`
	ClientID       string         `gorm:"unique;not null" json:"client_id,omitempty"`
	Username       string         `gorm:"unique;not null" json:"username,omitempty"`
	Email          string         `gorm:"unique;not null" json:"email,omitempty"`
	Password       string         `gorm:"not null" json:"-"`
	PinCode        string         `gorm:"not null" json:"-"`
	PinAttempts    int            `gorm:"default:0" json:"-"`
	PinLastUpdated time.Time      `json:"-"`
	FirstName      string         `json:"first_name,omitempty"`
	LastName       string         `json:"last_name,omitempty"`
	FullName       string         `json:"full_name,omitempty"`
	PhoneNumber    string         `gorm:"unique" json:"phone_number,omitempty"`
	ProfilePicture string         `json:"profile_picture,omitempty"`
	RoleID         uint           `gorm:"not null" json:"role_id,omitempty"`
	Role           Role           `gorm:"foreignKey:RoleID;constraint:OnDelete:CASCADE" json:"role,omitempty"`
	CreatedAt      time.Time      `gorm:"autoCreateTime" json:"created_at,omitempty"`
	CreatedBy      string         `json:"created_by,omitempty"`
	UpdatedAt      time.Time      `gorm:"autoUpdateTime" json:"updated_at,omitempty"`
	UpdatedBy      string         `json:"updated_by,omitempty"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	DeletedBy      string         `json:"deleted_by,omitempty"`
}

type TokenDetails struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	AccessUUID   string `json:"access_uuid"`
	RefreshUUID  string `json:"refresh_uuid"`
	AtExpires    int64  `json:"at_expires"`
	RtExpires    int64  `json:"rt_expires"`
}
