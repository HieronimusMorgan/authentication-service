package in

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginPhoneNumber struct {
	PhoneNumber string `json:"phone_number" binding:"required"`
	PinCode     string `json:"pin_code" binding:"required"`
}
