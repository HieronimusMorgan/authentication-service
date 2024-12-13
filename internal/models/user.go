package models

import "time"

type User struct {
	UserID         uint       `gorm:"primaryKey" json:"user_id,omitempty"`
	ClientID       string     `gorm:"unique;not null" json:"client_id,omitempty"`
	Username       string     `gorm:"unique;not null" json:"username,omitempty"`
	Password       string     `gorm:"not null" json:"-"` // Hashed password, omit from JSON
	FirstName      string     `json:"first_name,omitempty"`
	LastName       string     `json:"last_name,omitempty"`
	FullName       string     `json:"full_name,omitempty"`
	PhoneNumber    string     `gorm:"unique" json:"phone_number,omitempty"`
	ProfilePicture string     `json:"profile_picture,omitempty"`
	RoleID         uint       `gorm:"not null" json:"role_id,omitempty"`
	Role           Role       `gorm:"foreignKey:RoleID;constraint:OnDelete:CASCADE" json:"role,omitempty"`
	CreatedAt      time.Time  `gorm:"autoCreateTime" json:"created_at,omitempty"`
	CreatedBy      string     `json:"created_by,omitempty"`
	UpdatedAt      time.Time  `gorm:"autoUpdateTime" json:"updated_at,omitempty"`
	UpdatedBy      string     `json:"updated_by,omitempty"`
	DeletedAt      *time.Time `gorm:"index" json:"deleted_at,omitempty"` // Soft delete with timestamp
	DeletedBy      string     `json:"deleted_by,omitempty"`
}

type TokenDetails struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	AccessUUID   string `json:"access_uuid"`
	RefreshUUID  string `json:"refresh_uuid"`
	AtExpires    int64  `json:"at_expires"`
	RtExpires    int64  `json:"rt_expires"`
}
