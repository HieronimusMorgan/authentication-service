package services

import (
	"authentication/internal/dto/out"
	"authentication/internal/models"
	"authentication/internal/repository"
	"encoding/json"
	"errors"
	"log"
)

type RoleService interface {
	RegisterRole(req *struct {
		Name        string `json:"name" binding:"required"`
		Description string `json:"description" binding:"required"`
	}, clientID string) (interface{}, error)
	UpdateRole(roleID uint, req *struct {
		Name        string `json:"name" binding:"required"`
		Description string `json:"description" binding:"optional"`
	}, clientID string) (interface{}, error)
	GetListRole(clientID string) (interface{}, error)
	GetRoleById(roleID uint, clientID string) (interface{}, error)
	DeleteRole(roleID uint, clientID string) error
	GetListRoleUsers(clientID string, index int, size int) (interface{}, int64, error)
	GetListUserRole(clientID string, roleID uint, index int, size int) (interface{}, int64, error)
}

type roleService struct {
	RoleRepository repository.RoleRepository
	UserRepository repository.UserRepository
}

func NewRoleService(
	roleRepo repository.RoleRepository,
	userRepo repository.UserRepository,
) RoleService {
	return roleService{
		RoleRepository: roleRepo,
		UserRepository: userRepo,
	}
}

func (s roleService) RegisterRole(req *struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description" binding:"required"`
}, clientID string) (interface{}, error) {
	user, err := s.UserRepository.GetUserByClientID(clientID)
	if err != nil {
		return nil, err
	}

	var role = &models.Role{
		Name:        req.Name,
		Description: req.Description,
		CreatedBy:   user.FullName,
		UpdatedBy:   user.FullName,
	}
	err = s.RoleRepository.RegisterRole(role)
	if err != nil {
		return nil, err
	}

	return out.RoleResponse{
		RoleID:      role.RoleID,
		Name:        role.Name,
		Description: role.Description,
	}, nil
}

func (s roleService) UpdateRole(roleID uint, req *struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description" binding:"optional"`
}, clientID string) (interface{}, error) {
	admin, err := s.UserRepository.GetUserByClientID(clientID)
	if err != nil {
		return nil, err
	}

	role, err := s.RoleRepository.GetRoleByID(roleID)
	if err != nil {
		return nil, err
	}

	role.Name = req.Name
	if req.Description != "" {
		role.Description = req.Description
	}
	role.UpdatedBy = admin.FullName
	err = s.RoleRepository.UpdateRole(role)

	if err != nil {
		return nil, err
	}

	return out.RoleResponse{
		RoleID:      role.RoleID,
		Name:        role.Name,
		Description: role.Description,
	}, nil
}

func (s roleService) GetListRole(clientID string) (interface{}, error) {
	_, err := s.UserRepository.GetUserByClientID(clientID)
	if err != nil {
		return nil, err
	}

	roles, err := s.RoleRepository.GetAllRoles(0, 0)
	if err != nil {
		return nil, err
	}

	var roleResponses []out.RoleResponse
	for _, role := range *roles {
		roleResponses = append(roleResponses, out.RoleResponse{
			RoleID:      role.RoleID,
			Name:        role.Name,
			Description: role.Description,
		})
	}

	return roleResponses, nil
}

func (s roleService) GetRoleById(roleID uint, clientID string) (interface{}, error) {
	_, err := s.UserRepository.GetUserByClientID(clientID)
	if err != nil {
		return nil, err
	}

	role, err := s.RoleRepository.GetRoleByID(roleID)
	if err != nil {
		return nil, err
	}

	return out.RoleResponse{
		RoleID:      role.RoleID,
		Name:        role.Name,
		Description: role.Description,
	}, nil
}

func (s roleService) DeleteRole(roleID uint, clientID string) error {
	user, err := s.UserRepository.GetUserByClientID(clientID)
	if err != nil {
		return err
	}

	role, err := s.RoleRepository.GetRoleByID(roleID)
	if err != nil {
		return err
	}

	users, err := s.UserRepository.GetUserByRole(roleID)
	if err != nil {
		return err
	}

	if len(*users) > 0 {
		usersJSON, err := json.Marshal(users)
		if err != nil {
			return err
		}
		log.Println("Role is still being used by users+ " + string(usersJSON))
		return errors.New("role is still being used by users")
	}

	role.DeletedBy = user.FullName
	err = s.RoleRepository.DeleteRole(*role)
	if err != nil {
		return err
	}

	return nil
}

func (s roleService) GetListRoleUsers(clientID string, index, size int) (interface{}, int64, error) {
	_, err := s.UserRepository.GetUserByClientID(clientID)
	if err != nil {
		return nil, 0, err
	}

	roles, err := s.RoleRepository.GetAllRoles(index, size)
	if err != nil {
		return nil, 0, err
	}

	roleCount, err := s.RoleRepository.GetCountRole()
	if err != nil {
		return nil, 0, err
	}

	var roleResponses []struct {
		RoleID uint               `json:"role_id"`
		Name   string             `json:"name"`
		Users  []out.UserResponse `json:"users"`
	}
	for _, role := range *roles {
		users, err := s.UserRepository.GetUserByRole(role.RoleID)
		if err != nil {
			return nil, 0, err
		}

		var userResponses []out.UserResponse
		for _, user := range *users {
			userResponses = append(userResponses, out.UserResponse{
				UserID:         user.UserID,
				ClientID:       user.ClientID,
				Username:       user.Username,
				FirstName:      user.FirstName,
				LastName:       user.LastName,
				PhoneNumber:    user.PhoneNumber,
				ProfilePicture: user.ProfilePicture,
			})
		}

		roleResponses = append(roleResponses, struct {
			RoleID uint               `json:"role_id"`
			Name   string             `json:"name"`
			Users  []out.UserResponse `json:"users"`
		}{
			RoleID: role.RoleID,
			Name:   role.Name,
			Users:  userResponses,
		})
	}

	return roleResponses, roleCount, nil
}

func (s roleService) GetListUserRole(clientID string, roleID uint, index, size int) (interface{}, int64, error) {
	_, err := s.UserRepository.GetUserByClientID(clientID)
	if err != nil {
		return nil, 0, err
	}

	users, err := s.UserRepository.GetUserByRolePagination(roleID, index, size)
	if err != nil {
		return nil, 0, err
	}

	userCount, err := s.UserRepository.GetCountUserByRole(roleID)
	if err != nil {
		return nil, 0, err
	}

	var userResponses []out.UserResponse
	for _, user := range *users {
		userResponses = append(userResponses, out.UserResponse{
			UserID:         user.UserID,
			ClientID:       user.ClientID,
			Username:       user.Username,
			FirstName:      user.FirstName,
			LastName:       user.LastName,
			PhoneNumber:    user.PhoneNumber,
			ProfilePicture: &user.ProfilePicture,
		})
	}

	return userResponses, userCount, nil
}
