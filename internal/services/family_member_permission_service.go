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

type FamilyMemberPermissionService interface {
	AddFamilyMemberPermissions(req *in.UpdateFamilyMemberPermissionsRequest, clientID string) response.ErrorResponse
	RemoveFamilyMemberPermissions(req *in.UpdateFamilyMemberPermissionsRequest, clientID string) response.ErrorResponse
	GetListFamilyMemberPermissions(req in.FamilyMemberRequest, clientID string) ([]out.FamilyPermissionResponse, response.ErrorResponse)
}

type familyMemberPermissionService struct {
	UserRepository                   repository.UserRepository
	FamilyRepository                 repository.FamilyRepository
	FamilyMemberRepository           repository.FamilyMemberRepository
	FamilyPermissionRepository       repository.FamilyPermissionRepository
	FamilyMemberPermissionRepository repository.FamilyMemberPermissionRepository
	RedisService                     utils.RedisService
	JWTService                       utils.JWTService
}

func NewFamilyMemberPermissionService(
	UserRepository repository.UserRepository,
	FamilyRepository repository.FamilyRepository,
	FamilyMemberRepository repository.FamilyMemberRepository,
	FamilyPermissionRepository repository.FamilyPermissionRepository,
	FamilyMemberPermissionRepository repository.FamilyMemberPermissionRepository,
	RedisService utils.RedisService,
	JWTService utils.JWTService,
) FamilyMemberPermissionService {
	return &familyMemberPermissionService{
		UserRepository:                   UserRepository,
		FamilyRepository:                 FamilyRepository,
		FamilyMemberRepository:           FamilyMemberRepository,
		FamilyPermissionRepository:       FamilyPermissionRepository,
		FamilyMemberPermissionRepository: FamilyMemberPermissionRepository,
		RedisService:                     RedisService,
		JWTService:                       JWTService,
	}
}

func (s *familyMemberPermissionService) AddFamilyMemberPermissions(req *in.UpdateFamilyMemberPermissionsRequest, clientID string) response.ErrorResponse {
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

	family, err := s.FamilyRepository.GetFamilyByFamilyIdAndOwnerID(req.FamilyID, user.UserID)
	if err != nil {
		return response.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Family not found",
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

	member, err := s.UserRepository.GetUserByID(req.UserID)
	if err != nil {
		return response.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "User not found",
			Error:   err.Error(),
		}
	}

	familyMember, err := s.FamilyMemberRepository.GetFamilyMembersByFamilyIDAndMemberID(family.FamilyID, member.UserID)
	if err != nil {
		return response.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Family member not found",
			Error:   err.Error(),
		}
	}

	if familyMember.UserID == 0 {
		return response.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Family member not found",
			Error:   "This user is not a member of the family",
		}
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
		return response.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Permission denied",
			Error:   "You don't have permission to change the permission",
		}
	}

	permissionMember, err := s.FamilyMemberPermissionRepository.GetFamilyMemberPermissionByUserID(member.UserID)
	if err != nil {
		return response.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Permission not found",
			Error:   err.Error(),
		}
	}

	hasPermission = false
	for _, v := range permissionMember {
		if v.PermissionID == permission.PermissionID {
			hasPermission = true
			break
		}
	}

	if hasPermission {
		return response.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Permission already used",
			Error:   "This permission level is already assigned to the user",
		}
	}

	newPermission := &models.FamilyMemberPermission{
		PermissionID: permission.PermissionID,
		FamilyID:     family.FamilyID,
		UserID:       member.UserID,
		CreatedBy:    user.ClientID,
	}

	err = s.FamilyMemberPermissionRepository.AddFamilyMemberPermission(newPermission)
	if err != nil {
		return response.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Failed to update permission",
			Error:   err.Error(),
		}
	}

	return response.ErrorResponse{}
}

func (s *familyMemberPermissionService) RemoveFamilyMemberPermissions(req *in.UpdateFamilyMemberPermissionsRequest, clientID string) response.ErrorResponse {
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

	family, err := s.FamilyRepository.GetFamilyByFamilyIdAndOwnerID(req.FamilyID, user.UserID)
	if err != nil {
		return response.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Family not found",
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

	member, err := s.UserRepository.GetUserByID(req.UserID)
	if err != nil {
		return response.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "User not found",
			Error:   err.Error(),
		}
	}

	familyMember, err := s.FamilyMemberRepository.GetFamilyMembersByFamilyIDAndMemberID(family.FamilyID, member.UserID)
	if err != nil {
		return response.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Family member not found",
			Error:   err.Error(),
		}
	}

	if familyMember.UserID == 0 {
		return response.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Family member not found",
			Error:   "This user is not a member of the family",
		}
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
		return response.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Permission denied",
			Error:   "You don't have permission to change the permission",
		}
	}

	permissionMember, err := s.FamilyMemberPermissionRepository.GetFamilyMemberPermissionByUserID(member.UserID)
	if err != nil {
		return response.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Permission not found",
			Error:   err.Error(),
		}
	}

	hasPermission = false
	for _, v := range permissionMember {
		if v.PermissionID == permission.PermissionID {
			hasPermission = true
			break
		}
	}

	if hasPermission {
		newPermission := &models.FamilyMemberPermission{
			PermissionID: permission.PermissionID,
			FamilyID:     family.FamilyID,
			UserID:       member.UserID,
			CreatedBy:    user.ClientID,
		}

		err = s.FamilyMemberPermissionRepository.RemoveFamilyMemberPermission(newPermission)
		if err != nil {
			return response.ErrorResponse{
				Code:    http.StatusInternalServerError,
				Message: "Failed to update permission",
				Error:   err.Error(),
			}
		}

	} else {
		return response.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Permission not found",
			Error:   "This permission level is not assigned to the user",
		}
	}

	return response.ErrorResponse{}
}

func (s *familyMemberPermissionService) GetListFamilyMemberPermissions(req in.FamilyMemberRequest, clientID string) ([]out.FamilyPermissionResponse, response.ErrorResponse) {
	data, err := utils.GetUserRedis(s.RedisService, utils.User, clientID)
	if err != nil {
		return nil, response.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "User not found",
			Error:   err.Error(),
		}
	}

	_, err = s.UserRepository.GetUserByClientID(data.ClientID)
	if err != nil {
		return nil, response.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "User not found",
			Error:   err.Error(),
		}
	}

	permissions, err := s.FamilyMemberPermissionRepository.GetListFamilyMemberPermissionByFamilyIDAndUserID(req.FamilyID, req.UserID)
	if err != nil {
		return nil, response.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Permission not found",
			Error:   err.Error(),
		}
	}

	return permissions, response.ErrorResponse{}
}
