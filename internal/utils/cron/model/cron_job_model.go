package model

import (
	"time"
)

type CronJob struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	Name           string    `gorm:"type:varchar(100);not null" json:"name"`
	Schedule       string    `gorm:"type:varchar(100);not null" json:"schedule"`
	IsActive       bool      `gorm:"not null" json:"is_active"`
	Description    string    `gorm:"type:text" json:"description"`
	LastExecutedAt time.Time `gorm:"type:datetime" json:"last_executed_at"`
	CreatedAt      time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt      time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
