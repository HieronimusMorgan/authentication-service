package models

import (
	"github.com/lib/pq"
)

type UserRedis struct {
	UserID         uint             `gorm:"primaryKey" json:"user_id,omitempty"`
	ClientID       string           `gorm:"unique;not null" json:"client_id,omitempty"`
	Username       string           `gorm:"unique;not null" json:"username,omitempty"`
	Email          string           `gorm:"unique;not null" json:"email,omitempty"`
	Password       string           `gorm:"not null" json:"-"`
	PinCode        string           `gorm:"not null" json:"-"`
	PinAttempts    int              `gorm:"default:0" json:"-"`
	FirstName      string           `json:"first_name,omitempty"`
	LastName       string           `json:"last_name,omitempty"`
	FullName       string           `json:"full_name,omitempty"`
	PhoneNumber    string           `gorm:"unique" json:"phone_number,omitempty"`
	ProfilePicture string           `json:"profile_picture,omitempty"`
	Role           []RoleRedis      `json:"role,omitempty"`
	Resource       []ResourceRedis  `json:"resource,omitempty"`
	UserSetting    UserSettingRedis `json:"user_setting,omitempty"`
	DeviceID       *string          `json:"device_id,omitempty"`
}

type RoleRedis struct {
	RoleID      uint   `json:"role_id,omitempty"`
	Name        string `gorm:"unique;not null" json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}

type ResourceRedis struct {
	ResourceID  uint   `json:"resource_id,omitempty"`
	Name        string `gorm:"unique;not null" json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}

type UserSettingRedis struct {
	SettingID             uint          `gorm:"column:setting_id"`
	UserID                uint          `gorm:"uniqueIndex;not null;column:user_id"`
	GroupInviteType       int           `gorm:"column:group_invite_type;default:1"`
	GroupInviteDisallowed pq.Int32Array `gorm:"type:int[];column:group_invite_disallowed;default:{none}"`
}
