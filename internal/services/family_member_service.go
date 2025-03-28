package services

import (
	"authentication/internal/dto/in"
	"authentication/internal/dto/out"
	"authentication/internal/models"
	"authentication/internal/repository"
	"authentication/internal/utils"
	"authentication/package/response"
	"net/http"
)

type FamilyMemberService interface {
	AddFamilyMember(req *in.FamilyMemberRequest, clientID string) response.ErrorResponse
	UpdateFamilyMemberPermissions(req *in.UpdateFamilyMemberPermissionsRequest, clientID string) response.ErrorResponse
	GetFamilyMemberByFamilyID(familyID uint, clientID string) ([]out.FamilyMembersResponse, response.ErrorResponse)
	RemoveFamilyMember(req *in.FamilyMemberRequest, clientID string) response.ErrorResponse
}

type familyMemberService struct {
	UserRepository                   repository.UserRepository
	FamilyPermissionRepository       repository.FamilyPermissionRepository
	FamilyRepository                 repository.FamilyRepository
	FamilyMemberPermissionRepository repository.FamilyMemberPermissionRepository
	FamilyMemberRepository           repository.FamilyMemberRepository
	RedisService                     utils.RedisService
	JWTService                       utils.JWTService
	Encryption                       utils.Encryption
}

func NewFamilyMemberService(
	UserRepository repository.UserRepository,
	FamilyPermissionRepository repository.FamilyPermissionRepository,
	FamilyRepository repository.FamilyRepository,
	FamilyMemberPermissionRepository repository.FamilyMemberPermissionRepository,
	FamilyMemberRepository repository.FamilyMemberRepository,
	RedisService utils.RedisService,
	JWTService utils.JWTService,
	Encryption utils.Encryption) FamilyMemberService {
	return &familyMemberService{
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

func (s *familyMemberService) AddFamilyMember(req *in.FamilyMemberRequest, clientID string) response.ErrorResponse {
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

func (s *familyMemberService) UpdateFamilyMemberPermissions(req *in.UpdateFamilyMemberPermissionsRequest, clientID string) response.ErrorResponse {
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
		return response.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Permission denied",
			Error:   "You don't have permission to change the permission",
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

	permission, err := s.FamilyPermissionRepository.GetFamilyPermissionByID(req.PermissionID)
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
			Error:   "You cannot change the owner's permission",
		}
	}

	_, err = s.FamilyMemberRepository.GetFamilyMembersByFamilyIDAndMemberID(req.FamilyID, member.UserID)
	if err != nil {
		return response.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Member not found",
			Error:   err.Error(),
		}
	}

	permissionUser.PermissionID = permission.PermissionID
	err = s.FamilyMemberPermissionRepository.UpdateFamilyMemberPermission(permissionUser)
	if err != nil {
		return response.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Failed to update permission",
			Error:   err.Error(),
		}
	}

	return response.ErrorResponse{}
}

func (s *familyMemberService) GetFamilyMemberByFamilyID(familyID uint, clientID string) ([]out.FamilyMembersResponse, response.ErrorResponse) {
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

func (s *familyMemberService) RemoveFamilyMember(req *in.FamilyMemberRequest, clientID string) response.ErrorResponse {
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
