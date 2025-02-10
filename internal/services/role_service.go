package services

import (
	"authentication/internal/dto/out"
	"authentication/internal/models"
	"authentication/internal/repository"
	"encoding/json"
	"errors"
	"log"
)

type RoleService struct {
	RoleRepository repository.RoleRepository
	UserRepository repository.UserRepository
}

// NewRoleService initializes RoleService with repository dependencies
func NewRoleService(
	roleRepo repository.RoleRepository,
	userRepo repository.UserRepository,
) RoleService {
	return RoleService{
		RoleRepository: roleRepo,
		UserRepository: userRepo,
	}
}

func (s RoleService) RegisterRole(req *struct {
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
	err = s.RoleRepository.RegisterRole(&role)
	if err != nil {
		return nil, err
	}

	return out.RoleResponse{
		RoleID:      role.RoleID,
		Name:        role.Name,
		Description: role.Description,
	}, nil
}

func (s RoleService) UpdateRole(roleID uint, req *struct {
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
	err = s.RoleRepository.UpdateRole(&role)

	if err != nil {
		return nil, err
	}

	return out.RoleResponse{
		RoleID:      role.RoleID,
		Name:        role.Name,
		Description: role.Description,
	}, nil
}

func (s RoleService) GetListRole(clientID string) (interface{}, error) {
	_, err := s.UserRepository.GetUserByClientID(clientID)
	if err != nil {
		return nil, err
	}

	roles, err := s.RoleRepository.GetAllRoles()
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

func (s RoleService) GetRoleById(roleID uint, clientID string) (interface{}, error) {
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

func (s RoleService) DeleteRole(roleID uint, clientID string) error {
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
	err = s.RoleRepository.DeleteRole(&role)
	if err != nil {
		return err
	}

	return nil
}

func (s RoleService) GetListRoleUsers(clientID string) (interface{}, error) {
	_, err := s.UserRepository.GetUserByClientID(clientID)
	if err != nil {
		return nil, err
	}

	roles, err := s.RoleRepository.GetAllRoles()
	if err != nil {
		return nil, err
	}

	var roleResponses []struct {
		RoleID uint               `json:"role_id"`
		Name   string             `json:"name"`
		Users  []out.UserResponse `json:"users"`
	}
	for _, role := range *roles {
		users, err := s.UserRepository.GetUserByRole(role.RoleID)
		if err != nil {
			return nil, err
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

	return roleResponses, nil
}
