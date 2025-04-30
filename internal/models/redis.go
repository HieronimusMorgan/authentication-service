package models

import (
	"github.com/lib/pq"
)

type UserRedis struct {
	UserID         uint             `json:"user_id,omitempty"`
	ClientID       string           `json:"client_id,omitempty"`
	Username       string           `json:"username,omitempty"`
	Email          string           `json:"email,omitempty"`
	Password       string           `json:"-"`
	PinCode        string           `json:"-"`
	PinAttempts    int              `json:"-"`
	FirstName      string           `json:"first_name,omitempty"`
	LastName       string           `json:"last_name,omitempty"`
	FullName       string           `json:"full_name,omitempty"`
	PhoneNumber    string           `json:"phone_number,omitempty"`
	ProfilePicture string           `json:"profile_picture,omitempty"`
	Role           []RoleRedis      `json:"role,omitempty"`
	Resource       []ResourceRedis  `json:"resource,omitempty"`
	UserSetting    UserSettingRedis `json:"user_setting,omitempty"`
	DeviceID       *string          `json:"device_id,omitempty"`
}

type RoleRedis struct {
	RoleID      uint   `json:"role_id,omitempty"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}

type ResourceRedis struct {
	ResourceID  uint   `json:"resource_id,omitempty"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}

type UserSettingRedis struct {
	SettingID             uint          `json:"setting_id,omitempty"`
	UserID                uint          `json:"user_id,omitempty"`
	GroupInviteType       int           `json:"group_invite_type,omitempty"`
	GroupInviteDisallowed pq.Int32Array `json:"group_invite_disallowed,omitempty"`
}
