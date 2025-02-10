package repository

import (
	"authentication/internal/models"
	"gorm.io/gorm"
)

type ResourceRepository interface {
	AddResource(resource *models.Resource) error
	GetResourceByID(resourceID uint) (*models.Resource, error)
	DeleteResourceById(resourceID uint) error
	DeleteResource(resource *models.Resource) error
	UpdateResource(resource *models.Resource) error
	GetAllResources() (*[]models.Resource, error)
	GetResourceByResourceID(resourceID uint) (*models.Resource, error)
	GetResourceByResourceName(resourceName string) (*models.Resource, error)
	GetResourceByName(resourceName string) (*models.Resource, error)
	CreateInternalToken(resourceID uint, internalToken string) error
}

type resourceRepository struct {
	db *gorm.DB
}

func NewResourceRepository(db *gorm.DB) ResourceRepository {
	return &resourceRepository{db: db}
}

func (r resourceRepository) AddResource(resource *models.Resource) error {
	err := r.db.Where("name LIKE ?", resource.Name).FirstOrCreate(resource).Error
	if err != nil {
		return err
	}
	return nil
}

func (r resourceRepository) GetResourceByID(resourceID uint) (*models.Resource, error) {
	var resource models.Resource
	err := r.db.Where("resource_id = ?", resourceID).First(&resource).Error
	if err != nil {
		return nil, err
	}
	return &resource, nil
}

func (r resourceRepository) DeleteResourceById(resourceID uint) error {
	err := r.db.Where("resource_id = ?", resourceID).Delete(&models.Resource{}).Error
	if err != nil {
		return err
	}
	return nil
}

func (r resourceRepository) DeleteResource(resource *models.Resource) error {
	err := r.db.Model(&resource).
		Update("deleted_by", resource.DeletedBy).
		Delete(&resource).Error
	if err != nil {
		return err
	}
	return nil
}

func (r resourceRepository) UpdateResource(resource *models.Resource) error {
	err := r.db.Save(resource).Error
	if err != nil {
		return err
	}
	return nil
}

func (r resourceRepository) GetAllResources() (*[]models.Resource, error) {
	var resources []models.Resource
	err := r.db.Find(&resources).Where("delete_at NOT NULL").Order("resource_id ASC").Error
	if err != nil {
		return nil, err
	}
	return &resources, nil
}

func (r resourceRepository) GetResourceByResourceID(resourceID uint) (*models.Resource, error) {
	var resource models.Resource
	err := r.db.Where("resource_id = ?", resourceID).First(&resource).Error
	if err != nil {
		return nil, err
	}
	return &resource, nil
}

func (r resourceRepository) GetResourceByResourceName(resourceName string) (*models.Resource, error) {
	var resource models.Resource
	err := r.db.Where("resource_name = ?", resourceName).First(&resource).Error
	if err != nil {
		return nil, err
	}
	return &resource, nil
}

func (r resourceRepository) GetResourceByName(resourceName string) (*models.Resource, error) {
	var resource models.Resource
	err := r.db.Where("name = ?", resourceName).First(&resource).Error
	if err != nil {
		return nil, err
	}
	return &resource, nil
}

func (r resourceRepository) CreateInternalToken(resourceID uint, internalToken string) error {
	internal := models.InternalToken{
		ResourceID: resourceID,
		Token:      internalToken,
	}
	err := r.db.Create(&internal).Error
	if err != nil {
		return err
	}
	return nil
}
