package in

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	DeviceID string `json:"device_id"`
}

type LoginPhoneNumber struct {
	PhoneNumber string `json:"phone_number" binding:"required"`
	PinCode     string `json:"pin_code" binding:"required"`
	DeviceID    string `json:"device_id"`
}
