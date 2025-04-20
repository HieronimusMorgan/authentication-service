package out

type LoginResponse struct {
	UserID         uint                `json:"user_id"`
	ClientID       string              `json:"client_id"`
	Username       string              `json:"username"`
	FirstName      string              `json:"first_name"`
	LastName       string              `json:"last_name"`
	PhoneNumber    string              `json:"phone_number"`
	ProfilePicture string              `json:"profile_picture"`
	UserSetting    UserSettingResponse `json:"user_setting"`
	RefreshToken   string              `json:"refresh_token"`
	Token          string              `json:"token"`
}
