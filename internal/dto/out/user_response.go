package out

type UserResponse struct {
	ID             uint   `gorm:"primarykey" json:"id"`
	ClientID       string `gorm:"unique" json:"client_id"`
	Username       string `gorm:"unique" json:"username"`
	FirstName      string `json:"first_name"`
	LastName       string `json:"last_name"`
	PhoneNumber    string `gorm:"unique" json:"phone_number"`
	ProfilePicture string `json:"profile_picture"`
}
