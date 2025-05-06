package models

import "time"

type Notification struct {
	ID          uint      `json:"id"`
	Title       string    `json:"title"`
	Body        string    `json:"body"`
	TargetToken string    `json:"target_token"`
	Platform    string    `json:"platform"` // android or web
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
}
