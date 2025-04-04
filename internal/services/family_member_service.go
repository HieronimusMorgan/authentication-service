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
	RemoveFamilyMember(req *in.FamilyMemberRequest, clientID string) response.ErrorResponse
	GetFamilyMemberByFamilyID(familyID uint, clientID string) ([]out.FamilyMembersResponse, response.ErrorResponse)
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

	user, err := s.UserRepository.GetUserByClientID(data.ClientID)
	if err != nil {
		return response.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "User not found",
			Error:   err.Error(),
		}
	}

	newMember, err := s.UserRepository.GetUserByID(req.UserID)
	if err != nil {
		return response.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "User not found",
			Error:   err.Error(),
		}
	}

	familyOwner, err := s.FamilyRepository.GetFamilyByOwnerID(user.UserID)
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
		CreatedBy: user.ClientID,
		UpdatedBy: user.ClientID,
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
		CreatedBy:    user.ClientID,
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
	member, err := s.UserRepository.GetUserByID(req.UserID)
	if err != nil {
		return response.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "User not found",
			Error:   err.Error(),
		}
	}

	_, err = s.FamilyPermissionRepository.GetFamilyPermissionByName("Manage")
	if err != nil {
		return response.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Permission not found",
			Error:   err.Error(),
		}
	}
	_, err = s.FamilyMemberPermissionRepository.GetFamilyMemberPermissionByUserID(userOwner.UserID)
	if err != nil {
		return response.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Permission not found",
			Error:   err.Error(),
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
