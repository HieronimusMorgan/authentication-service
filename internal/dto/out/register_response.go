package out

type RegisterResponse struct {
	UserID         uint     `json:"user_id"`
	Username       string   `json:"username"`
	FirstName      string   `json:"first_name"`
	LastName       string   `json:"last_name"`
	PhoneNumber    string   `json:"phone_number"`
	ProfilePicture string   `json:"profile_picture"`
	Role           string   `json:"role"`
	Resource       []string `json:"resource"`
	Token          string   `json:"token"`
	RefreshToken   string   `json:"refresh_token"`
}
