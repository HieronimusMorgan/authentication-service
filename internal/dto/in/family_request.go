package in

type FamilyRequest struct {
	FamilyName string `json:"family_name" binding:"required"`
}

type UpdateFamilyRequest struct {
	FamilyID   uint   `json:"family_id" binding:"required"`
	FamilyName string `json:"family_name" binding:"required"`
}

type FamilyMemberRequest struct {
	FamilyID    uint   `json:"family_id,omitempty" binding:"required"`
	PhoneNumber string `json:"phone_number,omitempty" binding:"required"`
}

type UpdateFamilyMemberPermissionsRequest struct {
	FamilyID     uint   `json:"family_id,omitempty" binding:"required"`
	PermissionID uint   `json:"permission_id,omitempty" binding:"required"`
	PhoneNumber  string `json:"phone_number,omitempty" binding:"required"`
}

type FamilyPermissionRequest struct {
	PermissionName string `json:"permission_name,omitempty" binding:"required"`
	Description    string `json:"description,omitempty" binding:"required"`
}

type FamilyMemberPermissionRequest struct {
	FamilyID     uint `json:"family_id,omitempty" binding:"required"`
	UserID       uint `json:"user_id,omitempty" binding:"required"`
	PermissionID uint `json:"permission_id,omitempty" binding:"required"`
}

type ChangeFamilyMemberPermissionRequest struct {
	FamilyID     uint   `json:"family_id,omitempty" binding:"required"`
	PhoneNumber  string `json:"phone_number,omitempty" binding:"required"`
	PermissionID uint   `json:"permission_id,omitempty" binding:"required"`
}
