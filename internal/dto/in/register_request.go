package in

type RegisterRequest struct {
	Username       string `json:"username" binding:"required"`
	Password       string `json:"password" binding:"required"`
	FirstName      string `json:"first_name" binding:"required"`
	LastName       string `json:"last_name" binding:"required"`
	PhoneNumber    string `json:"phone_number" binding:"required"`
	ProfilePicture string `json:"profile_picture,omitempty"`
}
