package family

import "time"

type FamilyInvitation struct {
	InvitationID   uint      `gorm:"primaryKey" json:"invitation_id,omitempty"`
	FamilyID       uint      `gorm:"not null" json:"family_id,omitempty"`
	SenderUserID   uint      `gorm:"not null" json:"sender_user_id,omitempty"`
	ReceiverUserID uint      `gorm:"not null" json:"receiver_user_id,omitempty"`
	StatusID       uint      `gorm:"not null" json:"status_id,omitempty"`
	InvitedAt      time.Time `gorm:"autoCreateTime" json:"invited_at,omitempty"`
	RespondedAt    time.Time `json:"responded_at,omitempty"`
}
