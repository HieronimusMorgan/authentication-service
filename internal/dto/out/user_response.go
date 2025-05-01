package out

import "github.com/lib/pq"

type UserResponse struct {
	UserID         uint                `gorm:"primarykey" json:"user_id"`
	ClientID       string              `gorm:"unique" json:"client_id"`
	Username       string              `gorm:"unique" json:"username"`
	FirstName      string              `json:"first_name"`
	LastName       string              `json:"last_name"`
	PhoneNumber    string              `gorm:"unique" json:"phone_number"`
	ProfilePicture *string             `json:"profile_picture,omitempty"`
	UserSetting    UserSettingResponse `json:"user_setting"`
}

type UserRoleResourceSettingResponse struct {
	UserID         uint                `gorm:"primarykey" json:"user_id"`
	ClientID       string              `gorm:"unique" json:"client_id"`
	Username       string              `gorm:"unique" json:"username"`
	FirstName      string              `json:"first_name"`
	LastName       string              `json:"last_name"`
	PhoneNumber    string              `gorm:"unique" json:"phone_number"`
	ProfilePicture *string             `json:"profile_picture,omitempty"`
	Role           []RoleResponse      `json:"role"`
	Resource       []ResourceResponse  `json:"resource"`
	UserSetting    UserSettingResponse `json:"user_setting"`
}

type UserSettingResponse struct {
	SettingID             uint          `json:"setting_id" binding:"required"`
	GroupInviteType       int           `json:"group_invite_type"`
	GroupInviteDisallowed pq.Int32Array `json:"group_invite_disallowed"`
}

type VerifyPinCodeResponse struct {
	ClientID  string `json:"client_id"`
	RequestID string `json:"request_id"`
	Valid     bool   `json:"valid"`
}
