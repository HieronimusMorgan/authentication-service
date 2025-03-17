package services

import (
	"authentication/internal/dto/out"
	"authentication/internal/models/resource"
	"authentication/internal/models/role"
	"authentication/internal/models/users"
	"authentication/internal/repository"
	"errors"
	"strings"
)

type ResourceService interface {
	AddResource(name *string, description *string, clientID string) (interface{}, error)
	UpdateResource(resourceID uint, name *string, description *string, clientID string) (interface{}, error)
	GetResources(clientID string) (interface{}, error)
	AssignResourceToRole(roleID uint, resourceID uint, clientID string) (interface{}, error)
	GetResourceById(resourceID uint, clientID string) (interface{}, error)
	DeleteResourceById(resourceID uint, clientID string) error
	GetResourceUserById(resourceID uint, clientID string) (interface{}, error)
	GetResourceRoles(clientID string) (interface{}, error)
}

type resourceService struct {
	ResourceRepository     repository.ResourceRepository
	RoleResourceRepository repository.RoleResourceRepository
	RoleRepository         repository.RoleRepository
	UserRepository         repository.UserRepository
}

func NewResourceService(
	resourceRepo repository.ResourceRepository,
	roleResourceRepo repository.RoleResourceRepository,
	roleRepo repository.RoleRepository,
	userRepo repository.UserRepository,
) ResourceService {
	return resourceService{
		ResourceRepository:     resourceRepo,
		RoleResourceRepository: roleResourceRepo,
		RoleRepository:         roleRepo,
		UserRepository:         userRepo,
	}
}

func (s resourceService) checkUserIsAdmin(user *users.Users) error {
	role, err := s.RoleRepository.GetRoleByID(user.RoleID)
	if err != nil {
		return errors.New("role not found")
	}
	if strings.EqualFold(role.Name, "Admin") || strings.EqualFold(role.Name, "Super Admin") {
		return nil
	}
	return errors.New("user is not an admin")
}

func (s resourceService) AddResource(name *string, description *string, clientID string) (interface{}, error) {
	user, err := s.UserRepository.GetUserByClientID(clientID)
	if err != nil {
		return nil, err
	}

	err = s.checkUserIsAdmin(user)
	if err != nil {
		return nil, err
	}

	var resource = resource.Resource{
		Name:        *name,
		Description: *description,
		CreatedBy:   user.FullName,
		UpdatedBy:   user.FullName,
	}

	err = s.ResourceRepository.AddResource(&resource)
	if err != nil {
		return nil, err
	}

	return out.ResourceResponse{
		ResourceID:  resource.ResourceID,
		Name:        resource.Name,
		Description: resource.Description,
	}, nil
}

func (s resourceService) UpdateResource(resourceID uint, name *string, description *string, clientID string) (interface{}, error) {
	user, err := s.UserRepository.GetUserByClientID(clientID)
	if err != nil {
		return nil, err
	}

	err = s.checkUserIsAdmin(user)
	if err != nil {
		return nil, err
	}

	resource, err := s.ResourceRepository.GetResourceByID(resourceID)
	if err != nil {
		return nil, err
	}

	resource.Name = *name
	resource.Description = *description
	resource.UpdatedBy = user.FullName
	err = s.ResourceRepository.UpdateResource(resource)
	if err != nil {
		return nil, err
	}

	return out.ResourceResponse{
		ResourceID:  resource.ResourceID,
		Name:        resource.Name,
		Description: resource.Description,
	}, nil
}

func (s resourceService) GetResources(clientID string) (interface{}, error) {
	user, err := s.UserRepository.GetUserByClientID(clientID)
	if err != nil {
		return nil, err
	}

	err = s.checkUserIsAdmin(user)
	if err != nil {
		return nil, err
	}

	resources, err := s.ResourceRepository.GetAllResources()
	if err != nil {
		return nil, err
	}

	var resourceResponses []out.ResourceResponse
	for _, resource := range *resources {
		resourceResponses = append(resourceResponses, out.ResourceResponse{
			ResourceID:  resource.ResourceID,
			Name:        resource.Name,
			Description: resource.Description,
		})
	}

	return resources, nil
}

func (s resourceService) AssignResourceToRole(roleID uint, resourceID uint, clientID string) (interface{}, error) {
	user, err := s.UserRepository.GetUserByClientID(clientID)
	if err != nil {
		return nil, err
	}

	err = s.checkUserIsAdmin(user)
	if err != nil {
		return nil, err
	}

	_, err = s.ResourceRepository.GetResourceByID(resourceID)
	if err != nil {
		return nil, err
	}

	var roleResource = &role.RoleResource{
		RoleID:     roleID,
		ResourceID: resourceID,
		CreatedBy:  user.FullName,
		UpdatedBy:  user.FullName,
	}

	err = s.RoleResourceRepository.RegisterRoleResource(&roleResource)
	if err != nil {
		return nil, err
	}
	return struct {
		RoleID     uint
		ResourceID uint
	}{
		RoleID:     roleResource.RoleID,
		ResourceID: roleResource.ResourceID,
	}, nil
}

func (s resourceService) GetResourceById(resourceID uint, clientID string) (interface{}, error) {
	user, err := s.UserRepository.GetUserByClientID(clientID)
	if err != nil {
		return nil, err
	}

	err = s.checkUserIsAdmin(user)
	if err != nil {
		return nil, err
	}

	resource, err := s.ResourceRepository.GetResourceByID(resourceID)
	if err != nil {
		return nil, err
	}

	return out.ResourceResponse{
		ResourceID:  resource.ResourceID,
		Name:        resource.Name,
		Description: resource.Description,
	}, nil
}

func (s resourceService) DeleteResourceById(resourceID uint, clientID string) error {
	user, err := s.UserRepository.GetUserByClientID(clientID)
	if err != nil {
		return err
	}
	resource, err := s.ResourceRepository.GetResourceByID(resourceID)
	if err != nil {
		return err
	}
	resource.DeletedBy = user.FullName
	err = s.ResourceRepository.DeleteResource(resource)
	if err != nil {
		return err
	}

	return nil
}

func (s resourceService) GetResourceUserById(resourceID uint, clientID string) (interface{}, error) {
	user, err := s.UserRepository.GetUserByClientID(clientID)
	if err != nil {
		return nil, err
	}

	err = s.checkUserIsAdmin(user)
	if err != nil {
		return nil, err
	}

	resource, err := s.ResourceRepository.GetResourceByID(resourceID)
	if err != nil {
		return nil, err
	}

	users, err := s.UserRepository.GetUserByResourceID(resourceID)
	if err != nil {
		return nil, err
	}

	var userResponses []struct {
		UserID      uint
		ClientID    string
		FullName    string
		PhoneNumber string
		RoleID      uint
		CreatedAt   string
		UpdatedAt   string
	}

	for _, user := range *users {
		userResponses = append(userResponses, struct {
			UserID      uint
			ClientID    string
			FullName    string
			PhoneNumber string
			RoleID      uint
			CreatedAt   string
			UpdatedAt   string
		}{
			UserID:      user.UserID,
			ClientID:    user.ClientID,
			FullName:    user.FullName,
			PhoneNumber: user.PhoneNumber,
			RoleID:      user.RoleID,
			CreatedAt:   user.CreatedAt.String(),
			UpdatedAt:   user.UpdatedAt.String(),
		})
	}
	data := struct {
		ResourceID   uint   `json:"resource_id"`
		ResourceName string `json:"resource_name"`
		Users        []struct {
			UserID      uint
			ClientID    string
			FullName    string
			PhoneNumber string
			RoleID      uint
			CreatedAt   string
			UpdatedAt   string
		} `json:"users"`
	}{
		ResourceID:   resourceID,
		ResourceName: resource.Name,
		Users:        userResponses,
	}

	return data, nil
}

func (s resourceService) GetResourceRoles(clientID string) (interface{}, error) {
	user, err := s.UserRepository.GetUserByClientID(clientID)
	if err != nil {
		return nil, err
	}

	err = s.checkUserIsAdmin(user)
	if err != nil {
		return nil, err
	}

	resources, err := s.ResourceRepository.GetAllResources()
	if err != nil {
		return nil, err
	}

	var resourceResponses []struct {
		ResourceID   uint
		ResourceName string
		Roles        []struct {
			RoleID   uint
			RoleName string
		}
	}

	for _, resource := range *resources {
		roles, err := s.RoleRepository.GetAllRolesByResourceId(&resource)
		if err != nil {
			return nil, err
		}
		var roleResponses []struct {
			RoleID   uint
			RoleName string
		}
		for _, role := range *roles {
			roleResponses = append(roleResponses, struct {
				RoleID   uint
				RoleName string
			}{
				RoleID:   role.RoleID,
				RoleName: role.Name,
			})
		}
		resourceResponses = append(resourceResponses, struct {
			ResourceID   uint
			ResourceName string
			Roles        []struct {
				RoleID   uint
				RoleName string
			}
		}{
			ResourceID:   resource.ResourceID,
			ResourceName: resource.Name,
			Roles:        roleResponses,
		})
	}

	return resourceResponses, nil
}
