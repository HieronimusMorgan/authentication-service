package repository

import (
	"authentication/internal/models"
	"gorm.io/gorm"
)

type ResourceRepository struct {
	DB *gorm.DB
}

func NewResourceRepository(db *gorm.DB) *ResourceRepository {
	return &ResourceRepository{DB: db}
}

func (r ResourceRepository) AddResource(resource *models.Resource) error {
	err := r.DB.Where("name LIKE ?", resource.Name).FirstOrCreate(resource).Error
	if err != nil {
		return err
	}
	return nil
}

func (r ResourceRepository) GetResourceByID(resourceID uint) (*models.Resource, error) {
	var resource models.Resource
	err := r.DB.Where("resource_id = ?", resourceID).First(&resource).Error
	if err != nil {
		return nil, err
	}
	return &resource, nil
}

func (r ResourceRepository) DeleteResourceById(resourceID uint) error {
	err := r.DB.Where("resource_id = ?", resourceID).Delete(&models.Resource{}).Error
	if err != nil {
		return err
	}
	return nil
}

func (r ResourceRepository) DeleteResource(resource *models.Resource) error {
	err := r.DB.Model(&resource).
		Update("deleted_by", resource.DeletedBy).
		Delete(&resource).Error
	if err != nil {
		return err
	}
	return nil
}

func (r ResourceRepository) UpdateResource(resource *models.Resource) error {
	err := r.DB.Save(resource).Error
	if err != nil {
		return err
	}
	return nil
}

func (r ResourceRepository) GetAllResources() (*[]models.Resource, error) {
	var resources []models.Resource
	err := r.DB.Find(&resources).Where("delete_at NOT NULL").Order("resource_id ASC").Error
	if err != nil {
		return nil, err
	}
	return &resources, nil
}

func (r ResourceRepository) GetResourceByResourceID(resourceID uint) (*models.Resource, error) {
	var resource models.Resource
	err := r.DB.Where("resource_id = ?", resourceID).First(&resource).Error
	if err != nil {
		return nil, err
	}
	return &resource, nil
}

func (r ResourceRepository) GetResourceByResourceName(resourceName string) (*models.Resource, error) {
	var resource models.Resource
	err := r.DB.Where("resource_name = ?", resourceName).First(&resource).Error
	if err != nil {
		return nil, err
	}
	return &resource, nil
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
