package services

import (
	"authentication/internal/dto/out"
	"authentication/internal/models"
	"authentication/internal/repository"
	nt "authentication/internal/utils/nats"
	"errors"
	"log"
	"strings"
)

type ResourceService interface {
	AddResource(name *string, description *string, clientID string) (interface{}, error)
	UpdateResource(resourceID uint, name *string, description *string, clientID string) (interface{}, error)
	GetResources(clientID string) (interface{}, error)
	AssignUserResource(userID uint, resourceID uint, clientID string) (interface{}, error)
	RemoveAssignUserResource(userID uint, resourceID uint, clientID string) error
	GetResourceById(resourceID uint, clientID string) (interface{}, error)
	DeleteResourceById(resourceID uint, clientID string) error
	GetResourceUserById(resourceID uint, clientID string) (interface{}, error)
	GetUserResources(clientID string) (interface{}, error)
}

type resourceService struct {
	ResourceRepository     repository.ResourceRepository
	UserResourceRepository repository.UserResourceRepository
	RoleRepository         repository.RoleRepository
	UserRepository         repository.UserRepository
	NatsService            nt.Service
	AuthService            AuthService
}

func NewResourceService(resourceRepo repository.ResourceRepository, roleResourceRepo repository.UserResourceRepository, roleRepo repository.RoleRepository, userRepo repository.UserRepository, service nt.Service, authService AuthService) ResourceService {
	return resourceService{
		ResourceRepository:     resourceRepo,
		UserResourceRepository: roleResourceRepo,
		RoleRepository:         roleRepo,
		UserRepository:         userRepo,
		NatsService:            service,
		AuthService:            authService,
	}
}

func (s resourceService) checkUserIsAdmin(user *models.Users) error {
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

func (s resourceService) AssignUserResource(userID uint, resourceID uint, clientID string) (interface{}, error) {
	admin, err := s.UserRepository.GetUserByClientID(clientID)
	if err != nil {
		return nil, err
	}

	err = s.checkUserIsAdmin(admin)
	if err != nil {
		return nil, err
	}

	_, err = s.ResourceRepository.GetResourceByID(resourceID)
	if err != nil {
		return nil, err
	}
	user, err := s.UserRepository.GetUserByID(userID)
	if err != nil {
		return nil, err
	}

	var userResource = models.UserResource{
		UserID:     userID,
		ResourceID: resourceID,
		CreatedBy:  admin.FullName,
		UpdatedBy:  admin.FullName,
	}

	err = s.UserResourceRepository.RegisterUserResource(userResource)
	if err != nil {
		return nil, err
	}

	token, err := s.AuthService.UpdateToken(userID, admin.ClientID)
	if err != nil {
		return nil, err
	}

	if user.DeviceToken != nil {
		notification := models.Notification{
			TargetToken:   *user.DeviceToken,
			Title:         "Assign User Resource",
			Body:          "You have been assigned a new resource",
			Priority:      "high",
			Color:         "#1E88E5",
			Platform:      "android",
			ServiceSource: "authentication",
			EventType:     "assign_user_resource",
			ClickAction:   "OPEN_ACTIVITY",
			Payload: map[string]string{
				"token":         token.AccessToken,
				"refresh_token": token.RefreshToken,
			},
		}

		err = s.NatsService.RequestNotification("authentication", notification)
	} else {
		log.Println("Device token is nil, skipping notification")
	}

	return struct {
		UserID     uint
		ResourceID uint
	}{
		UserID:     userResource.UserID,
		ResourceID: userResource.ResourceID,
	}, nil
}

func (s resourceService) RemoveAssignUserResource(userID uint, resourceID uint, clientID string) error {
	admin, err := s.UserRepository.GetUserByClientID(clientID)
	if err != nil {
		return err
	}

	err = s.checkUserIsAdmin(admin)
	if err != nil {
		return err
	}

	user, err := s.UserRepository.GetUserByID(userID)
	if err != nil {
		return err
	}

	_, err = s.ResourceRepository.GetResourceByID(resourceID)
	if err != nil {
		return err
	}

	userResource, err := s.UserResourceRepository.GetUserResourceByUserIDAndResourceID(userID, resourceID)
	if err != nil {
		return err
	}

	err = s.UserResourceRepository.DeleteUserResource(userResource)
	if err != nil {
		return err
	}

	token, err := s.AuthService.UpdateToken(userID, admin.ClientID)
	if err != nil {
		return err
	}

	if user.DeviceToken != nil {
		notification := models.Notification{
			TargetToken:   *user.DeviceToken,
			Title:         "Remove User Resource",
			Body:          "You have been removed from a resource",
			Priority:      "high",
			Color:         "#1E88E5",
			Platform:      "android",
			ServiceSource: "authentication",
			EventType:     "remove_user_resource",
			ClickAction:   "OPEN_ACTIVITY",
			Payload: map[string]string{
				"token":         token.AccessToken,
				"refresh_token": token.RefreshToken,
			},
		}
		err = s.NatsService.RequestNotification("authentication", notification)
	} else {
		log.Println("Device token is nil, skipping notification")
	}

	return nil
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

func (s resourceService) GetUserResources(clientID string) (interface{}, error) {
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
		User         []struct {
			UserID    uint
			Username  string
			FirstName string
			LastName  string
		}
	}

	for _, resource := range *resources {
		users, err := s.UserRepository.GetAllUsersByResourceId(&resource)
		if err != nil {
			return nil, err
		}
		var userResponse []struct {
			UserID    uint
			Username  string
			FirstName string
			LastName  string
		}
		for _, u := range *users {
			userResponse = append(userResponse, struct {
				UserID    uint
				Username  string
				FirstName string
				LastName  string
			}{
				UserID:    u.UserID,
				Username:  u.Username,
				FirstName: u.FirstName,
				LastName:  u.LastName,
			})
		}
		resourceResponses = append(resourceResponses, struct {
			ResourceID   uint
			ResourceName string
			User         []struct {
				UserID    uint
				Username  string
				FirstName string
				LastName  string
			}
		}{
			ResourceID:   resource.ResourceID,
			ResourceName: resource.Name,
			User:         userResponse,
		})
	}

	return resourceResponses, nil
}
