package services

import (
	"authentication/internal/dto/out"
	"authentication/internal/models"
	"authentication/internal/repository"
	"errors"
	"gorm.io/gorm"
	"strings"
)

type ResourceService struct {
	ResourceRepository     *repository.ResourceRepository
	RoleResourceRepository *repository.RoleResourceRepository
	RoleRepository         *repository.RoleRepository
	UserRepository         *repository.UserRepository
}

func NewResourceService(db *gorm.DB) *ResourceService {
	resourceRepo := repository.NewResourceRepository(db)
	roleResourceRepo := repository.NewRoleResourceRepository(db)
	roleRepo := repository.NewRoleRepository(db)
	userRepo := repository.NewUserRepository(db)
	return &ResourceService{ResourceRepository: resourceRepo, UserRepository: userRepo,
		RoleRepository: roleRepo, RoleResourceRepository: roleResourceRepo}
}

func (s ResourceService) checkUserIsAdmin(user *models.Users) error {
	role, err := s.RoleRepository.GetRoleByID(user.RoleID)
	if err != nil {
		return errors.New("role not found")
	}
	if strings.EqualFold(role.Name, "Admin") {
		return nil
	}
	return errors.New("user is not an admin")
}

func (s ResourceService) AddResource(name *string, description *string, clientID string) (interface{}, error) {
	user, err := s.UserRepository.GetUserByClientID(clientID)
	if err != nil {
		return nil, err
	}

	err = s.checkUserIsAdmin(user)
	if err != nil {
		return nil, err
	}

	var resource = models.Resource{
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

func (s ResourceService) UpdateResource(resourceID uint, name *string, description *string, clientID string) (interface{}, error) {
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

func (s ResourceService) GetResources(clientID string) (interface{}, error) {
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

func (s ResourceService) AssignResourceToRole(roleID uint, resourceID uint, clientID string) (interface{}, error) {
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

	var roleResource = &models.RoleResource{
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

func (s ResourceService) GetResourceById(resourceID uint, clientID string) (interface{}, error) {
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

func (s ResourceService) DeleteResourceById(resourceID uint, clientID string) error {
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
