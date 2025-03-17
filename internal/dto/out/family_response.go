package out

type FamilyResponse struct {
	FamilyID   uint   `json:"family_id,omitempty"`
	FamilyName string `json:"family_name,omitempty"`
	OwnerID    uint   `json:"owner_id,omitempty"`
}

type FamilyPermissionResponse struct {
	PermissionID   uint   `json:"permission_id,omitempty"`
	PermissionName string `json:"permission_name,omitempty"`
	Description    string `json:"description,omitempty"`
}

type FamilyMemberResponse struct {
	FamilyID       uint   `json:"family_id,omitempty"`
	UserID         uint   `json:"user_id,omitempty"`
	ClientID       string `json:"client_id"`
	Username       string `json:"username"`
	FirstName      string `json:"first_name"`
	LastName       string `json:"last_name"`
	PhoneNumber    string `json:"phone_number"`
	ProfilePicture string `json:"profile_picture"`
}

type FamilyMemberPermissionResponse struct {
	FamilyID       uint   `gorm:"primaryKey" json:"family_id,omitempty"`
	UserID         uint   `gorm:"primaryKey" json:"user_id,omitempty"`
	PermissionID   uint   `gorm:"primaryKey" json:"permission_id,omitempty"`
	PermissionName string `gorm:"unique;not null" json:"permission_name,omitempty"`
	Description    string `json:"description,omitempty"`
}
