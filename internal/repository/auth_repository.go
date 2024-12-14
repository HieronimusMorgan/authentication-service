package repository

import (
	"authentication/internal/models"
	"encoding/json"
	"errors"
	"gorm.io/gorm"
	"log"
	"time"
)

type AuthRepository struct {
	DB *gorm.DB
}

func NewAuthRepository(db *gorm.DB) *AuthRepository {
	return &AuthRepository{DB: db}
}

func (r AuthRepository) CreateUser(user *models.User) error {
	return r.DB.Create(user).Error
}

func (r AuthRepository) GetUserByUsername(username string) (interface{}, error) {
	var user models.User
	err := r.DB.Where("username = ?", username).First(&user).Error
	return user, err
}

func (r AuthRepository) GetUserByClientID(clientID string) (*models.User, error) {
	var user models.User
	err := r.DB.Where("client_id = ?", clientID).First(&user).Error
	return &user, err
}

func (r AuthRepository) AssignUserResource(userID uint, resourceID uint) (*AssignResource, error) {
	// Validate user
	var user models.User
	if err := r.DB.Preload("Role").First(&user, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	// Validate resource
	var resource models.Resource
	if err := r.DB.First(&resource, resourceID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("resource not found")
		}
		return nil, err
	}

	var role models.Role
	if err := r.DB.First(&role, user.RoleID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("role not found")
		}
		return nil, err
	}

	// Check if the user already has the resource assigned
	var existingAssignment models.RoleResource
	err := r.DB.Where("role_id = ? AND resource_id = ?", user.RoleID, resource.ResourceID).First(&existingAssignment).Error
	if err == nil {
		jsonData, _ := json.Marshal(existingAssignment)
		log.Printf(string(jsonData))
		return &AssignResource{
			UserID:     userID,
			ResourceID: existingAssignment.ResourceID,
			RoleID:     existingAssignment.RoleID,
			Role:       role.Name,
			Resource:   resource.Name,
		}, nil
	}

	// Create role-resource assignment
	roleResource := models.RoleResource{
		RoleID:     user.RoleID,
		ResourceID: resource.ResourceID,
		CreatedAt:  time.Now(),
		CreatedBy:  "system",
	}

	if err := r.DB.Create(&roleResource).Error; err != nil {
		return nil, err
	}

	// Log success and return response
	log.Printf("Resource '%s' assigned to user '%s' successfully!", resource.Name, user.Username)
	return &AssignResource{
		UserID:     userID,
		ResourceID: resourceID,
		RoleID:     user.RoleID,
		Role:       role.Name,
		Resource:   resource.Name,
	}, nil
}

func (r AuthRepository) AssignUserResourceByName(userID uint, resourceName string) (*AssignResource, error) {
	// Validate user
	var user models.User
	if err := r.DB.Preload("Role").First(&user, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	// Validate resource by name
	var resource models.Resource
	if err := r.DB.Where("name = ?", resourceName).First(&resource).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("resource not found")
		}
		return nil, err
	}

	// Check if the user already has the resource assigned
	var existingAssignment models.RoleResource
	err := r.DB.Where("role_id = ? AND resource_id = ?", user.RoleID, resource.ResourceID).First(&existingAssignment).Error
	if err == nil {
		return &AssignResource{
			UserID:     userID,
			ResourceID: existingAssignment.ResourceID,
			RoleID:     existingAssignment.RoleID,
			Role:       existingAssignment.Role.Name,
			Resource:   existingAssignment.Resource.Name,
		}, nil
	}

	// Create role-resource assignment
	roleResource := models.RoleResource{
		RoleID:     user.RoleID,
		ResourceID: resource.ResourceID,
		CreatedAt:  time.Now(),
		CreatedBy:  "system",
	}

	if err := r.DB.Create(&roleResource).Error; err != nil {
		return nil, err
	}

	log.Printf("Resource '%s' assigned to user '%s' successfully!", resource.Name, user.Username)
	return &AssignResource{
		UserID:     userID,
		ResourceID: resource.ResourceID,
		RoleID:     user.RoleID,
		Role:       user.Role.Name,
		Resource:   resource.Name,
	}, nil
}

type AssignResource struct {
	UserID     uint   `json:"user_id"`
	ResourceID uint   `json:"resource_id"`
	RoleID     uint   `json:"role_id"`
	Role       string `json:"role"`
	Resource   string `json:"resource"`
}
