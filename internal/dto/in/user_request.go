package in

import "github.com/lib/pq"

type UserSettingsRequest struct {
	SettingID             uint          `json:"setting_id" binding:"required"`
	ArchivedEnabled       bool          `json:"archived_enabled"`
	ArchivedExceptions    pq.Int32Array `json:"archived_exceptions"`
	GroupInviteType       int           `json:"group_invite_type"`
	GroupInviteDisallowed pq.Int32Array `json:"group_invite_disallowed"`
}
