package services

import (
	"authentication/internal/dto/in"
	"authentication/internal/dto/out"
	"authentication/internal/models/family"
	"authentication/internal/repository"
	"authentication/internal/utils"
	"authentication/package/response"
	"net/http"
	"time"
)

type FamilyService interface {
	CreateFamily(req *in.FamilyRequest, clientID string) response.ErrorResponse
	GetFamilyByID(familyID uint, clientID string) (*family.Family, response.ErrorResponse)
	UpdateFamily(req *in.FamilyRequest, clientID string) response.ErrorResponse
	DeleteFamily(id uint, clientID string) response.ErrorResponse
	GetAllFamilies() ([]family.Family, error)
	GetFamilyResponseByClientID(clientID string) (*out.FamilyResponse, error)
	ChangeFamilyOwner(familyID uint, clientID string, newOwnerID uint) response.ErrorResponse

	//Family Permissions
	CreateFamilyPermission(req *in.FamilyPermissionRequest, clientID string) response.ErrorResponse
	GetFamilyPermissionByID(id uint) (*family.FamilyPermission, error)
	UpdateFamilyPermission(permission *family.FamilyPermission) error
	DeleteFamilyPermission(permission *family.FamilyPermission) error
	GetAllFamilyPermissions() ([]family.FamilyPermission, error)
	GetPermissionsByUserID(userID uint) ([]family.FamilyPermission, error)

	// Family Members Permissions
	GetFamilyMemberPermissionByID(familyID uint, clientID string) (*family.FamilyMemberPermission, response.ErrorResponse)
	UpdateFamilyMemberPermission(clientID string, req *in.FamilyMemberPermissionRequest) response.ErrorResponse
	DeleteFamilyMemberPermission(req *in.FamilyMemberPermissionRequest) response.ErrorResponse
	GetAllFamilyMemberPermissions(clientID string) ([]family.FamilyMemberPermission, response.ErrorResponse)

	// Family Members
	CreateFamilyMember(req *in.FamilyMemberRequest, clientID string) response.ErrorResponse
	GetFamilyMemberByID(id uint) (*family.FamilyMember, error)
	UpdateFamilyMember(family *family.FamilyMember) error
	DeleteFamilyMember(family *family.FamilyMember) error
	GetAllFamilyMembers() ([]family.FamilyMember, error)
	GetFamilyMembersByFamilyID(familyID uint) ([]family.FamilyMember, error)
	GetFamilyMembersByMemberID(memberID uint) ([]family.FamilyMember, error)
	GetFamilyMembersByFamilyIDAndMemberID(familyID uint, memberID uint) (*family.FamilyMember, error)
}

type familyService struct {
	UserRepository                   repository.UserRepository
	FamilyPermissionRepository       repository.FamilyPermissionRepository
	FamilyRepository                 repository.FamilyRepository
	FamilyMemberPermissionRepository repository.FamilyMemberPermissionRepository
	FamilyMemberRepository           repository.FamilyMemberRepository
	RedisService                     utils.RedisService
	JWTService                       utils.JWTService
	Encryption                       utils.Encryption
}

func NewFamilyService(
	UserRepository repository.UserRepository,
	FamilyPermissionRepository repository.FamilyPermissionRepository,
	FamilyRepository repository.FamilyRepository,
	FamilyMemberPermissionRepository repository.FamilyMemberPermissionRepository,
	FamilyMemberRepository repository.FamilyMemberRepository,
	RedisService utils.RedisService,
	JWTService utils.JWTService,
	Encryption utils.Encryption) FamilyService {
	return &familyService{
		UserRepository:                   UserRepository,
		FamilyPermissionRepository:       FamilyPermissionRepository,
		FamilyRepository:                 FamilyRepository,
		FamilyMemberPermissionRepository: FamilyMemberPermissionRepository,
		FamilyMemberRepository:           FamilyMemberRepository,
		RedisService:                     RedisService,
		JWTService:                       JWTService,
		Encryption:                       Encryption,
	}
}

func (s *familyService) CreateFamily(req *in.FamilyRequest, clientID string) response.ErrorResponse {
	data, err := utils.GetUserRedis(s.RedisService, utils.User, clientID)
	if err != nil {
		return response.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "User not found",
			Error:   err.Error(),
		}
	}

	user, err := s.UserRepository.GetUserByClientID(data.ClientID)
	if err != nil {
		return response.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "User not found",
			Error:   err.Error(),
		}
	}

	// check if user already has a family in family member by user id
	checkUser, _ := s.FamilyMemberRepository.GetFamilyMembersByUserID(user.UserID)
	if checkUser.UserID != 0 {
		return response.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "User already has a family",
			Error:   "User already has a family",
		}
	}

	f := &family.Family{
		FamilyName: req.FamilyName,
		OwnerID:    user.UserID,
		CreatedBy:  user.ClientID,
		UpdatedBy:  user.ClientID,
	}

	fullAccessPermission, err := s.FamilyPermissionRepository.GetFamilyPermissionByName("Manage")
	if err != nil {
		return response.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Permission not found",
			Error:   err.Error(),
		}
	}

	//create family member permission for owner
	permission := &family.FamilyMemberPermission{
		FamilyID:     f.FamilyID,
		UserID:       user.UserID,
		PermissionID: fullAccessPermission.PermissionID,
		CreatedBy:    user.ClientID,
		UpdatedBy:    user.ClientID,
	}

	// Create family member for owner
	familyMember := &family.FamilyMember{
		FamilyID:  f.FamilyID,
		UserID:    user.UserID,
		JoinedAt:  time.Now(),
		CreatedBy: user.ClientID,
		UpdatedBy: user.ClientID,
	}

	if err := s.FamilyRepository.CreateFamily(f, permission, familyMember); err != nil {
		return response.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Failed to create family",
			Error:   err.Error(),
		}
	}

	return response.ErrorResponse{}
}

func (s *familyService) GetFamilyByID(familyID uint, clientID string) (*family.Family, response.ErrorResponse) {
	data, err := utils.GetUserRedis(s.RedisService, utils.User, clientID)
	if err != nil {
		return nil, response.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "User not found",
			Error:   err.Error(),
		}
	}

	user, err := s.UserRepository.GetUserByClientID(data.ClientID)
	if err != nil {
		return nil, response.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "User not found",
			Error:   err.Error(),
		}
	}

	f, err := s.FamilyRepository.GetFamilyByFamilyIdAndOwnerID(familyID, user.UserID)
	if err != nil {
		return nil, response.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Family not found",
			Error:   err.Error(),
		}
	}

	return f, response.ErrorResponse{}
}

func (s *familyService) UpdateFamily(req *in.FamilyRequest, clientID string) response.ErrorResponse {
	data, err := utils.GetUserRedis(s.RedisService, utils.User, clientID)
	if err != nil {
		return response.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "User not found",
			Error:   err.Error(),
		}
	}

	f, err := s.FamilyRepository.GetFamilyByOwnerID(data.UserID)
	if err != nil {
		return response.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Family not found",
			Error:   err.Error(),
		}
	}

	f.FamilyName = req.FamilyName
	f.UpdatedBy = data.ClientID

	err = s.FamilyRepository.UpdateFamily(f)
	if err != nil {
		return response.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Failed to update family",
			Error:   err.Error(),
		}
	}
	return response.ErrorResponse{}
}

func (s *familyService) DeleteFamily(id uint, clientID string) response.ErrorResponse {
	data, err := utils.GetUserRedis(s.RedisService, utils.User, clientID)
	if err != nil {
		return response.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "User not found",
			Error:   err.Error(),
		}
	}

	f, err := s.FamilyRepository.GetFamilyByFamilyIdAndOwnerID(id, data.UserID)
	if err != nil {
		return response.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Family not found",
			Error:   err.Error(),
		}
	}

	f.DeletedBy = data.ClientID
	err = s.FamilyRepository.DeleteFamily(f)
	if err != nil {
		return response.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Failed to delete family",
			Error:   err.Error(),
		}
	}
	return response.ErrorResponse{}
}

func (s *familyService) GetAllFamilies() ([]family.Family, error) {
	return s.FamilyRepository.GetAllFamilies()
}

func (s *familyService) GetFamilyResponseByClientID(clientID string) (*out.FamilyResponse, error) {
	return s.FamilyRepository.GetFamilyResponseByClientID(clientID)
}

func (s *familyService) ChangeFamilyOwner(familyID uint, clientID string, newOwnerID uint) response.ErrorResponse {
	data, err := utils.GetUserRedis(s.RedisService, utils.User, clientID)
	if err != nil {
		return response.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "User not found",
			Error:   err.Error(),
		}
	}

	newOwner, err := s.UserRepository.GetUserByID(newOwnerID)
	if err != nil {
		return response.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "User not found",
			Error:   err.Error(),
		}
	}

	f, err := s.FamilyRepository.GetFamilyByFamilyIdAndOwnerID(familyID, data.UserID)
	if err != nil {
		return response.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Family not found",
			Error:   err.Error(),
		}
	}

	err = s.FamilyRepository.ChangeFamilyOwner(f.FamilyID, newOwner.UserID)
	if err != nil {
		return response.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Failed to change family owner",
			Error:   err.Error(),
		}
	}

	return response.ErrorResponse{}
}

// Permissions
func (s *familyService) CreateFamilyPermission(req *in.FamilyPermissionRequest, clientID string) response.ErrorResponse {
	data, err := utils.GetUserRedis(s.RedisService, utils.User, clientID)
	if err != nil {
		return response.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "User not found",
			Error:   err.Error(),
		}
	}

	user, err := s.UserRepository.GetUserByClientID(data.ClientID)
	if err != nil {
		return response.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "User not found",
			Error:   err.Error(),
		}
	}

	permission := &family.FamilyPermission{
		PermissionName: req.PermissionName,
		Description:    req.Description,
		CreatedBy:      user.ClientID,
		UpdatedBy:      user.ClientID,
	}

	err = s.FamilyPermissionRepository.CreateFamilyPermission(permission)
	if err != nil {
		return response.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Failed to create family permission",
			Error:   err.Error(),
		}
	}
	return response.ErrorResponse{}
}

func (s *familyService) GetFamilyPermissionByID(id uint) (*family.FamilyPermission, error) {
	return s.FamilyPermissionRepository.GetFamilyPermissionByID(id)
}

func (s *familyService) UpdateFamilyPermission(permission *family.FamilyPermission) error {
	return s.FamilyPermissionRepository.UpdateFamilyPermission(permission)
}

func (s *familyService) DeleteFamilyPermission(permission *family.FamilyPermission) error {
	return s.FamilyPermissionRepository.DeleteFamilyPermission(permission)
}

func (s *familyService) GetAllFamilyPermissions() ([]family.FamilyPermission, error) {
	return s.FamilyPermissionRepository.GetAllFamilyPermissions()
}

func (s *familyService) GetPermissionsByUserID(userID uint) ([]family.FamilyPermission, error) {
	return s.FamilyPermissionRepository.GetPermissionsByUserID(userID)
}

// Family Members Permissions
func (s *familyService) GetFamilyMemberPermissionByID(familyID uint, clientID string) (*family.FamilyMemberPermission, response.ErrorResponse) {
	_, err := utils.GetUserRedis(s.RedisService, utils.User, clientID)
	if err != nil {
		return nil, response.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "User not found",
			Error:   err.Error(),
		}
	}

	permission, err := s.FamilyMemberPermissionRepository.GetFamilyMemberPermissionByID(familyID)
	if err != nil {
		return nil, response.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Permission not found",
			Error:   err.Error(),
		}
	}

	return permission, response.ErrorResponse{}
}

func (s *familyService) UpdateFamilyMemberPermission(clientID string, req *in.FamilyMemberPermissionRequest) response.ErrorResponse {
	data, err := utils.GetUserRedis(s.RedisService, utils.User, clientID)
	if err != nil {
		return response.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "User not found",
			Error:   err.Error(),
		}
	}

	owner, err := s.FamilyRepository.GetFamilyByOwnerID(data.UserID)
	if err != nil {
		return response.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Family not found",
			Error:   err.Error(),
		}
	}

	if owner.FamilyID != req.FamilyID {
		return response.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Permission denied",
			Error:   "You are not the owner of this family",
		}
	}

	isMember, err := s.FamilyMemberRepository.GetFamilyMembersByFamilyIDAndMemberID(req.FamilyID, req.UserID)
	if err != nil {
		return response.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "User not found",
			Error:   err.Error(),
		}
	}

	if isMember == nil {
		return response.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "User not found",
			Error:   "User is not a member of this family",
		}
	}

	permission, err := s.FamilyMemberPermissionRepository.GetFamilyMemberPermissionByUserID(isMember.UserID)
	if err != nil {
		return response.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Permission not found",
			Error:   err.Error(),
		}
	}

	permission.PermissionID = req.PermissionID
	permission.UpdatedBy = data.ClientID

	err = s.FamilyMemberPermissionRepository.UpdateFamilyMemberPermission(permission)
	if err != nil {
		return response.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Failed to update permission",
			Error:   err.Error(),
		}
	}
	return response.ErrorResponse{}
}

func (s *familyService) DeleteFamilyMemberPermission(req *in.FamilyMemberPermissionRequest) response.ErrorResponse {
	permission, err := s.FamilyMemberPermissionRepository.GetFamilyMemberPermissionByID(req.FamilyID)
	if err != nil {
		return response.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Permission not found",
			Error:   err.Error(),
		}
	}

	err = s.FamilyMemberPermissionRepository.DeleteFamilyMemberPermission(permission)
	if err != nil {
		return response.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Failed to delete permission",
			Error:   err.Error(),
		}
	}
	return response.ErrorResponse{}
}

func (s *familyService) GetAllFamilyMemberPermissions(clientID string) ([]family.FamilyMemberPermission, response.ErrorResponse) {
	_, err := utils.GetUserRedis(s.RedisService, utils.User, clientID)
	if err != nil {
		return nil, response.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "User not found",
			Error:   err.Error(),
		}
	}

	permissions, err := s.FamilyMemberPermissionRepository.GetAllFamilyMemberPermissions()
	if err != nil {
		return nil, response.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Failed to get permissions",
			Error:   err.Error(),
		}
	}

	return permissions, response.ErrorResponse{}
}

// Family Members
func (s *familyService) CreateFamilyMember(req *in.FamilyMemberRequest, clientID string) response.ErrorResponse {
	data, err := utils.GetUserRedis(s.RedisService, utils.User, clientID)
	if err != nil {
		return response.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "User not found",
			Error:   err.Error(),
		}
	}

	userOwner, err := s.UserRepository.GetUserByClientID(data.ClientID)
	if err != nil {
		return response.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "User not found",
			Error:   err.Error(),
		}
	}

	if err := utils.ValidatePhoneNumber(req.PhoneNumber); err != nil {
		return response.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Validation Phone Number",
			Error:   err.Error(),
		}
	}

	hashPhoneNumber, err := s.Encryption.Encrypt(req.PhoneNumber)
	if err != nil {
		return response.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Failed to encrypt phone number",
			Error:   err.Error(),
		}
	}

	newMember, err := s.UserRepository.GetUserByPhoneNumber(hashPhoneNumber)
	if err != nil {
		return response.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "User not found",
			Error:   err.Error(),
		}
	}

	// Check if the user is the owner of the family
	familyOwner, err := s.FamilyRepository.GetFamilyByOwnerID(userOwner.UserID)
	if err != nil {
		return response.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Family not found",
			Error:   err.Error(),
		}
	}

	if familyOwner.FamilyID != req.FamilyID {
		return response.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Permission denied",
			Error:   "You are not the owner of this family",
		}
	}

	// Check if the user is already a member of the family
	isMember, _ := s.FamilyMemberRepository.GetFamilyMembersByFamilyIDAndMemberID(req.FamilyID, newMember.UserID)

	if isMember != nil {
		return response.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "User already a member",
			Error:   "User is already a member of this family",
		}
	}

	newMemberFamily := &family.FamilyMember{
		FamilyID:  familyOwner.FamilyID,
		UserID:    newMember.UserID,
		CreatedBy: userOwner.ClientID,
		UpdatedBy: userOwner.ClientID,
	}

	// Create family member permission
	readOnly, err := s.FamilyPermissionRepository.GetFamilyPermissionByName("Read")
	if err != nil {
		return response.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Permission not found",
			Error:   err.Error(),
		}
	}

	permission := &family.FamilyMemberPermission{
		PermissionID: readOnly.PermissionID,
		FamilyID:     familyOwner.FamilyID,
		UserID:       newMember.UserID,
		CreatedBy:    userOwner.ClientID,
		UpdatedBy:    userOwner.ClientID,
	}

	err = s.FamilyMemberRepository.CreateFamilyMember(newMemberFamily, permission)
	if err != nil {
		return response.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Failed to create family member",
			Error:   err.Error(),
		}
	}

	//return s.FamilyMemberRepository.CreateFamilyMember(family)
	return response.ErrorResponse{}
}

func (s *familyService) GetFamilyMemberByID(id uint) (*family.FamilyMember, error) {
	return s.FamilyMemberRepository.GetFamilyMemberByID(id)
}

func (s *familyService) UpdateFamilyMember(family *family.FamilyMember) error {
	return s.FamilyMemberRepository.UpdateFamilyMember(family)
}

func (s *familyService) DeleteFamilyMember(family *family.FamilyMember) error {
	return s.FamilyMemberRepository.DeleteFamilyMember(family)
}

func (s *familyService) GetAllFamilyMembers() ([]family.FamilyMember, error) {
	return s.FamilyMemberRepository.GetAllFamilyMembers()
}

func (s *familyService) GetFamilyMembersByFamilyID(familyID uint) ([]family.FamilyMember, error) {
	return s.FamilyMemberRepository.GetFamilyMembersByFamilyID(familyID)
}

func (s *familyService) GetFamilyMembersByMemberID(memberID uint) ([]family.FamilyMember, error) {
	return s.FamilyMemberRepository.GetFamilyMembersByMemberID(memberID)
}

func (s *familyService) GetFamilyMembersByFamilyIDAndMemberID(familyID uint, memberID uint) (*family.FamilyMember, error) {
	return s.FamilyMemberRepository.GetFamilyMembersByFamilyIDAndMemberID(familyID, memberID)
}
