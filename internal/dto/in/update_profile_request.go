package in

type UpdateNameRequest struct {
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
}

type UpdatePhotoRequest struct {
	ProfilePicture string `json:"profile_picture"  binding:"required"`
}
