package models

type Email struct {
	To       string `json:"to" binding:"required,email"`
	FullName string `json:"full_name" binding:"required"`
	URL      string `json:"url" binding:"required"`
	Subject  string `json:"subject" binding:"required"`
}
