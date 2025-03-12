package out

type UserResponse struct {
	UserID         uint   `gorm:"primarykey" json:"user_id"`
	ClientID       string `gorm:"unique" json:"client_id"`
	Username       string `gorm:"unique" json:"username"`
	FirstName      string `json:"first_name"`
	LastName       string `json:"last_name"`
	PhoneNumber    string `gorm:"unique" json:"phone_number"`
	ProfilePicture string `json:"profile_picture"`
}

type VerifyPinCodeResponse struct {
	ClientID  string `json:"client_id"`
	RequestID string `json:"request_id"`
	Valid     bool   `json:"valid"`
}
