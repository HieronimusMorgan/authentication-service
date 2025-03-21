package services

import (
	"authentication/internal/dto/in"
	"authentication/internal/dto/out"
	"authentication/internal/models"
	"authentication/internal/repository"
	"authentication/internal/utils"
	"authentication/package/response"
	"net/http"
	"time"
)

type FamilyService interface {
	CreateFamily(req *in.FamilyRequest, clientID string) response.ErrorResponse
	AddFamilyMember(req *in.FamilyMemberRequest, clientID string) response.ErrorResponse
	AddFamilyMemberPermission(req *in.ChangeFamilyMemberPermissionRequest, clientID string) response.ErrorResponse
	RemoveFamilyMemberPermission(req *in.ChangeFamilyMemberPermissionRequest, clientID string) response.ErrorResponse
	GetFamilyMemberByFamilyID(familyID uint, clientID string) ([]out.FamilyMembersResponse, response.ErrorResponse)
	RemoveFamilyMember(req *in.FamilyMemberRequest, clientID string) response.ErrorResponse
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

	checkUser, _ := s.FamilyMemberRepository.GetFamilyMembersByUserID(user.UserID)
	if checkUser.UserID != 0 {
		return response.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "User already has a family",
			Error:   "User already has a family",
		}
	}

	f := &models.Family{
		FamilyName: req.FamilyName,
		OwnerID:    user.UserID,
		CreatedBy:  user.ClientID,
		UpdatedBy:  user.ClientID,
	}

	allFamilyPermission, err := s.FamilyPermissionRepository.GetAllFamilyPermissions()
	if err != nil {
		return response.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Permission not found",
			Error:   err.Error(),
		}
	}

	familyMember := &models.FamilyMember{
		FamilyID:  f.FamilyID,
		UserID:    user.UserID,
		JoinedAt:  time.Now(),
		CreatedBy: user.ClientID,
		UpdatedBy: user.ClientID,
	}

	if err := s.FamilyRepository.CreateFamily(f, allFamilyPermission, familyMember); err != nil {
		return response.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Failed to create family",
			Error:   err.Error(),
		}
	}

	return response.ErrorResponse{}
}

func (s *familyService) AddFamilyMember(req *in.FamilyMemberRequest, clientID string) response.ErrorResponse {
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

	newMemberFamily := &models.FamilyMember{
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

	permission := &models.FamilyMemberPermission{
		PermissionID: readOnly.PermissionID,
		FamilyID:     familyOwner.FamilyID,
		UserID:       newMember.UserID,
		CreatedBy:    userOwner.ClientID,
	}

	err = s.FamilyMemberRepository.CreateFamilyMember(newMemberFamily, permission)
	if err != nil {
		return response.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Failed to create family member",
			Error:   err.Error(),
		}
	}

	return response.ErrorResponse{}
}
func (s *familyService) AddFamilyMemberPermission(req *in.ChangeFamilyMemberPermissionRequest, clientID string) response.ErrorResponse {
	data, err := utils.GetUserRedis(s.RedisService, utils.User, clientID)
	if err != nil {
		return response.ErrorResponse{Code: http.StatusBadRequest, Message: "User not found", Error: err.Error()}
	}

	user, err := s.UserRepository.GetUserByClientID(data.ClientID)
	if err != nil {
		return response.ErrorResponse{Code: http.StatusBadRequest, Message: "User not found", Error: err.Error()}
	}

	if err := utils.ValidatePhoneNumber(req.PhoneNumber); err != nil {
		return response.ErrorResponse{Code: http.StatusBadRequest, Message: "Invalid phone number", Error: err.Error()}
	}

	hashPhoneNumber, err := s.Encryption.Encrypt(req.PhoneNumber)
	if err != nil {
		return response.ErrorResponse{Code: http.StatusInternalServerError, Message: "Encryption failed", Error: err.Error()}
	}

	newMember, err := s.UserRepository.GetUserByPhoneNumber(hashPhoneNumber)
	if err != nil {
		return response.ErrorResponse{Code: http.StatusBadRequest, Message: "User not found", Error: err.Error()}
	}

	if user.UserID == newMember.UserID {
		return response.ErrorResponse{Code: http.StatusBadRequest, Message: "Permission denied", Error: "You cannot change your own permission"}
	}

	permission, err := s.FamilyPermissionRepository.GetFamilyPermissionByID(req.PermissionID)
	if err != nil {
		return response.ErrorResponse{Code: http.StatusBadRequest, Message: "Permission not found", Error: err.Error()}
	}

	family, err := s.FamilyMemberRepository.GetFamilyMembersByFamilyIDAndMemberID(req.FamilyID, newMember.UserID)
	if err != nil {
		return response.ErrorResponse{Code: http.StatusBadRequest, Message: "Member not found", Error: err.Error()}
	}

	permissionOwner, err := s.FamilyPermissionRepository.GetListFamilyPermissionAccess(user.UserID)
	if err != nil {
		return response.ErrorResponse{Code: http.StatusBadRequest, Message: "Permission not found", Error: err.Error()}
	}

	for _, v := range permissionOwner {
		if v.PermissionName == "Admin" || v.PermissionName == "Manage" {
			break
		}
		return response.ErrorResponse{Code: http.StatusBadRequest, Message: "Permission denied", Error: "You do not have permission to change the permission"}
	}

	permissionUser, err := s.FamilyMemberPermissionRepository.GetFamilyMemberPermissionByUserID(user.UserID)
	if err != nil {
		return response.ErrorResponse{Code: http.StatusBadRequest, Message: "Permission not found", Error: err.Error()}
	}

	if permissionUser.PermissionID == permission.PermissionID {
		return response.ErrorResponse{Code: http.StatusBadRequest, Message: "Permission denied", Error: "You cannot change the permission to the same permission"}
	}

	memberPermission := &models.FamilyMemberPermission{
		FamilyID:     family.FamilyID,
		UserID:       newMember.UserID,
		PermissionID: permission.PermissionID,
		CreatedBy:    user.ClientID,
	}

	if err := s.FamilyMemberPermissionRepository.CreateFamilyMemberPermission(memberPermission); err != nil {
		return response.ErrorResponse{Code: http.StatusInternalServerError, Message: "Failed to update permission", Error: err.Error()}
	}

	return response.ErrorResponse{}
}

func (s *familyService) RemoveFamilyMemberPermission(req *in.ChangeFamilyMemberPermissionRequest, clientID string) response.ErrorResponse {
	data, err := utils.GetUserRedis(s.RedisService, utils.User, clientID)
	if err != nil {
		return response.ErrorResponse{Code: http.StatusBadRequest, Message: "User not found", Error: err.Error()}
	}

	userOwner, err := s.UserRepository.GetUserByClientID(data.ClientID)
	if err != nil {
		return response.ErrorResponse{Code: http.StatusBadRequest, Message: "User not found", Error: err.Error()}
	}

	if err := utils.ValidatePhoneNumber(req.PhoneNumber); err != nil {
		return response.ErrorResponse{Code: http.StatusBadRequest, Message: "Invalid phone number", Error: err.Error()}
	}

	hashPhoneNumber, err := s.Encryption.Encrypt(req.PhoneNumber)
	if err != nil {
		return response.ErrorResponse{Code: http.StatusInternalServerError, Message: "Encryption failed", Error: err.Error()}
	}

	newMember, err := s.UserRepository.GetUserByPhoneNumber(hashPhoneNumber)
	if err != nil {
		return response.ErrorResponse{Code: http.StatusBadRequest, Message: "User not found", Error: err.Error()}
	}

	if userOwner.UserID == newMember.UserID {
		return response.ErrorResponse{Code: http.StatusBadRequest, Message: "Permission denied", Error: "You cannot remove your own permission"}
	}

	permission, err := s.FamilyPermissionRepository.GetFamilyPermissionByID(req.PermissionID)
	if err != nil {
		return response.ErrorResponse{Code: http.StatusBadRequest, Message: "Permission not found", Error: err.Error()}
	}

	_, err = s.FamilyMemberRepository.GetFamilyMembersByFamilyIDAndMemberID(req.FamilyID, newMember.UserID)
	if err != nil {
		return response.ErrorResponse{Code: http.StatusBadRequest, Message: "Member not found", Error: err.Error()}
	}

	permissionOwner, err := s.FamilyPermissionRepository.GetListFamilyPermissionAccess(userOwner.UserID)
	if err != nil {
		return response.ErrorResponse{Code: http.StatusBadRequest, Message: "Permission not found", Error: err.Error()}
	}

	hasPermission := false
	for _, v := range permissionOwner {
		if v.PermissionName == "Admin" || v.PermissionName == "Manage" {
			hasPermission = true
			break
		}
	}

	if !hasPermission {
		return response.ErrorResponse{Code: http.StatusBadRequest, Message: "Permission denied", Error: "You do not have permission to change the permission"}
	}

	err = s.FamilyMemberPermissionRepository.DeleteFamilyMemberPermissionByFamilyAndUserAndPermission(req.FamilyID, newMember.UserID, permission.PermissionID)
	if err != nil {
		return response.ErrorResponse{Code: http.StatusInternalServerError, Message: "Failed to delete permission", Error: err.Error()}
	}
	return response.ErrorResponse{}
}

func (s *familyService) GetFamilyMemberByFamilyID(familyID uint, clientID string) ([]out.FamilyMembersResponse, response.ErrorResponse) {
	_, err := utils.GetUserRedis(s.RedisService, utils.User, clientID)
	if err != nil {
		return nil, response.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "User not found",
			Error:   err.Error(),
		}
	}

	fm, err := s.FamilyMemberRepository.GetFamilyMembersByFamilyID(familyID)
	if err != nil {
		return nil, response.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Family not found",
			Error:   err.Error(),
		}
	}

	for i, v := range fm {
		decryptedPhoneNumber, err := s.Encryption.Decrypt(v.PhoneNumber)
		if err != nil {
			return nil, response.ErrorResponse{
				Code:    http.StatusInternalServerError,
				Message: "Failed to decrypt phone number",
				Error:   err.Error(),
			}
		}
		fm[i].PhoneNumber = decryptedPhoneNumber
	}

	return fm, response.ErrorResponse{}
}

func (s *familyService) RemoveFamilyMember(req *in.FamilyMemberRequest, clientID string) response.ErrorResponse {
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

	member, err := s.UserRepository.GetUserByPhoneNumber(hashPhoneNumber)
	if err != nil {
		return response.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "User not found",
			Error:   err.Error(),
		}
	}

	permission, err := s.FamilyPermissionRepository.GetFamilyPermissionByName("Manage")
	if err != nil {
		return response.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Permission not found",
			Error:   err.Error(),
		}
	}
	permissionUser, err := s.FamilyMemberPermissionRepository.GetFamilyMemberPermissionByUserID(userOwner.UserID)
	if err != nil {
		return response.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Permission not found",
			Error:   err.Error(),
		}
	}

	if permission.PermissionID != permissionUser.PermissionID {
		return response.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Permission denied",
			Error:   "You cannot remove the owner of the family",
		}
	}

	memberFamily, err := s.FamilyMemberRepository.GetFamilyMembersByFamilyIDAndMemberID(req.FamilyID, member.UserID)
	if err != nil {
		return response.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Member not found",
			Error:   err.Error(),
		}
	}

	memberFamily.DeletedBy = data.ClientID
	err = s.FamilyMemberRepository.DeleteFamilyMember(memberFamily)
	if err != nil {
		return response.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Failed to delete member",
			Error:   err.Error(),
		}
	}
	return response.ErrorResponse{}
}
