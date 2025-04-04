package services

import (
	"authentication/internal/dto/in"
	"authentication/internal/models"
	"authentication/internal/repository"
	"authentication/internal/utils"
	"authentication/package/response"
	"net/http"
	"time"
)

type FamilyService interface {
	CreateFamily(req *in.FamilyRequest, clientID string) response.ErrorResponse
	UpdateFamily(req *in.FamilyRequest, clientID string) response.ErrorResponse
	DeleteFamily(req *in.FamilyMemberRequest, clientID string) response.ErrorResponse
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

func (s *familyService) UpdateFamily(req *in.FamilyRequest, clientID string) response.ErrorResponse {
	data, err := utils.GetUserRedis(s.RedisService, utils.User, clientID)
	if err != nil {
		return response.ErrorResponse{Code: http.StatusBadRequest, Message: "User not found", Error: err.Error()}
	}

	user, err := s.UserRepository.GetUserByClientID(data.ClientID)
	if err != nil {
		return response.ErrorResponse{Code: http.StatusBadRequest, Message: "User not found", Error: err.Error()}
	}

	checkUser, err := s.FamilyMemberRepository.GetFamilyMembersByUserID(user.UserID)
	if err != nil {
		return response.ErrorResponse{Code: http.StatusBadRequest, Message: "User not found", Error: err.Error()}
	}

	permissionOwner, err := s.FamilyPermissionRepository.GetListFamilyPermissionAccess(checkUser.UserID)
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
		return response.ErrorResponse{Code: http.StatusBadRequest, Message: "Permission denied", Error: "You do not have permission to delete the family"}
	}

	family, err := s.FamilyRepository.GetFamilyByID(checkUser.FamilyID)
	if err != nil {
		return response.ErrorResponse{Code: http.StatusBadRequest, Message: "Family not found", Error: err.Error()}
	}

	if family.OwnerID != user.UserID {
		return response.ErrorResponse{Code: http.StatusBadRequest, Message: "Permission denied", Error: "You do not have permission to update the family"}
	}

	family.FamilyName = req.FamilyName
	family.UpdatedBy = user.ClientID

	err = s.FamilyRepository.UpdateFamily(family)
	if err != nil {
		return response.ErrorResponse{Code: http.StatusInternalServerError, Message: "Failed to update family", Error: err.Error()}
	}

	return response.ErrorResponse{}
}

func (s *familyService) DeleteFamily(req *in.FamilyMemberRequest, clientID string) response.ErrorResponse {
	data, err := utils.GetUserRedis(s.RedisService, utils.User, clientID)
	if err != nil {
		return response.ErrorResponse{Code: http.StatusBadRequest, Message: "User not found", Error: err.Error()}
	}

	user, err := s.UserRepository.GetUserByClientID(data.ClientID)
	if err != nil {
		return response.ErrorResponse{Code: http.StatusBadRequest, Message: "User not found", Error: err.Error()}
	}

	newMember, err := s.UserRepository.GetUserByID(req.UserID)
	if err != nil {
		return response.ErrorResponse{Code: http.StatusBadRequest, Message: "User not found", Error: err.Error()}
	}

	if user.UserID == newMember.UserID {
		return response.ErrorResponse{Code: http.StatusBadRequest, Message: "Permission denied", Error: "You cannot delete your own family"}
	}

	family, err := s.FamilyMemberRepository.GetFamilyMembersByFamilyIDAndMemberID(req.FamilyID, newMember.UserID)
	if err != nil {
		return response.ErrorResponse{Code: http.StatusBadRequest, Message: "Member not found", Error: err.Error()}
	}

	permissionOwner, err := s.FamilyPermissionRepository.GetListFamilyPermissionAccess(user.UserID)
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
		return response.ErrorResponse{Code: http.StatusBadRequest, Message: "Permission denied", Error: "You do not have permission to delete the family"}
	}

	err = s.FamilyRepository.DeleteFamilyByID(family.FamilyID)
	if err != nil {
		return response.ErrorResponse{Code: http.StatusInternalServerError, Message: "Failed to delete family", Error: err.Error()}
	}

	return response.ErrorResponse{}
}
