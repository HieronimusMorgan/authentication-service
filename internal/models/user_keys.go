package models

import (
	"gorm.io/gorm"
	"time"
)

type UserKey struct {
	UserID              uint           `gorm:"primaryKey;column:user_id"`
	PublicKey           string         `gorm:"column:public_key"`
	EncryptedPrivateKey string         `gorm:"column:encrypted_private_key"`
	EncryptionAlgorithm string         `gorm:"column:encryption_algorithm"`
	Salt                string         `gorm:"column:salt"`
	CreatedAt           time.Time      `gorm:"column:created_at"`
	CreatedBy           *string        `gorm:"column:created_by"`
	UpdatedAt           time.Time      `gorm:"column:updated_at"`
	UpdatedBy           *string        `gorm:"column:updated_by"`
	DeletedAt           gorm.DeletedAt `gorm:"column:deleted_at"`
	DeletedBy           *string        `gorm:"column:deleted_by"`
}
