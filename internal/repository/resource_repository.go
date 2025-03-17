package repository

import (
	"authentication/internal/models/resource"
	"gorm.io/gorm"
)

type ResourceRepository interface {
	AddResource(resource *resource.Resource) error
	GetResourceByID(resourceID uint) (*resource.Resource, error)
	GetResourceByUserID(userID uint) (*[]resource.Resource, error)
	DeleteResourceById(resourceID uint) error
	DeleteResource(resource *resource.Resource) error
	UpdateResource(resource *resource.Resource) error
	GetAllResources() (*[]resource.Resource, error)
	GetResourceByResourceID(resourceID uint) (*resource.Resource, error)
	GetResourceByResourceName(resourceName string) (*resource.Resource, error)
	GetResourceByName(resourceName string) (*resource.Resource, error)
}

type resourceRepository struct {
	db gorm.DB
}

func NewResourceRepository(db gorm.DB) ResourceRepository {
	return &resourceRepository{db: db}
}

func (r resourceRepository) AddResource(resource *resource.Resource) error {
	err := r.db.Where("name LIKE ?", resource.Name).FirstOrCreate(resource).Error
	if err != nil {
		return err
	}
	return nil
}

func (r resourceRepository) GetResourceByID(resourceID uint) (*resource.Resource, error) {
	var resource resource.Resource
	err := r.db.Where("resource_id = ?", resourceID).First(&resource).Error
	if err != nil {
		return nil, err
	}
	return &resource, nil
}

func (r resourceRepository) GetResourceByUserID(userID uint) (*[]resource.Resource, error) {
	var resources []resource.Resource
	query := `
		SELECT res.*
		FROM "authentication"."resources" res
		WHERE EXISTS (
			SELECT 1
			FROM "authentication"."role_resources" rr
			JOIN "authentication"."user_roles" ur 
				ON rr.role_id = ur.role_id
			WHERE ur.user_id = ? 
				AND rr.resource_id = res.resource_id
		);
	`

	err := r.db.Raw(query, userID).Scan(&resources).Error
	if err != nil {
		return nil, err
	}

	return &resources, nil
}

func (r resourceRepository) DeleteResourceById(resourceID uint) error {
	err := r.db.Where("resource_id = ?", resourceID).Delete(&resource.Resource{}).Error
	if err != nil {
		return err
	}
	return nil
}

func (r resourceRepository) DeleteResource(resource *resource.Resource) error {
	err := r.db.Model(&resource).
		Update("deleted_by", resource.DeletedBy).
		Delete(&resource).Error
	if err != nil {
		return err
	}
	return nil
}

func (r resourceRepository) UpdateResource(resource *resource.Resource) error {
	err := r.db.Save(resource).Error
	if err != nil {
		return err
	}
	return nil
}

func (r resourceRepository) GetAllResources() (*[]resource.Resource, error) {
	var resources []resource.Resource
	err := r.db.Find(&resources).Where("delete_at NOT NULL").Order("resource_id ASC").Error
	if err != nil {
		return nil, err
	}
	return &resources, nil
}

func (r resourceRepository) GetResourceByResourceID(resourceID uint) (*resource.Resource, error) {
	var resource resource.Resource
	err := r.db.Where("resource_id = ?", resourceID).First(&resource).Error
	if err != nil {
		return nil, err
	}
	return &resource, nil
}

func (r resourceRepository) GetResourceByResourceName(resourceName string) (*resource.Resource, error) {
	var resource resource.Resource
	err := r.db.Where("resource_name = ?", resourceName).First(&resource).Error
	if err != nil {
		return nil, err
	}
	return &resource, nil
}

func (r resourceRepository) GetResourceByName(resourceName string) (*resource.Resource, error) {
	var resource resource.Resource
	err := r.db.Where("name = ?", resourceName).First(&resource).Error
	if err != nil {
		return nil, err
	}
	return &resource, nil
}
