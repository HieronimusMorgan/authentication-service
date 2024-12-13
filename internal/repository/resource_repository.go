package repository

import (
	"Authentication/internal/models"
	"gorm.io/gorm"
)

type ResourceRepository struct {
	DB *gorm.DB
}

func NewResourceRepository(db *gorm.DB) *ResourceRepository {
	return &ResourceRepository{DB: db}
}

func (r ResourceRepository) GetResourceByName(resourceName string) (*models.Resource, error) {
	var resource models.Resource
	err := r.DB.Where("name = ?", resourceName).First(&resource).Error
	if err != nil {
		return nil, err
	}
	return &resource, nil
}

func (r ResourceRepository) CreateInternalToken(resourceID uint, internalToken string) error {
	internal := models.InternalToken{
		ResourceID: resourceID,
		Token:      internalToken,
	}
	err := r.DB.Create(&internal).Error
	if err != nil {
		return err
	}
	return nil
}
