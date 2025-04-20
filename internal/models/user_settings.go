package models

import (
	"github.com/lib/pq"
	"time"
)

type UserSetting struct {
	SettingID             uint          `gorm:"primaryKey;column:setting_id"`
	UserID                uint          `gorm:"uniqueIndex;not null;column:user_id"`
	GroupInviteType       int           `gorm:"column:group_invite_type;default:1"`
	GroupInviteDisallowed pq.Int32Array `gorm:"type:int[];column:group_invite_disallowed;default:{none}"`
	CreatedAt             time.Time     `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt             time.Time     `gorm:"column:updated_at;autoUpdateTime"`
}
