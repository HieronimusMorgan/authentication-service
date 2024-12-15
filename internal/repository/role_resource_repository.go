package repository

import (
	"authentication/internal/models"
	"gorm.io/gorm"
)

type RoleResourceRepository struct {
	DB *gorm.DB
}

func NewRoleResourceRepository(db *gorm.DB) *RoleResourceRepository {
	return &RoleResourceRepository{DB: db}
}

func (r RoleResourceRepository) RegisterRoleResource(roleResource **models.RoleResource) error {
	err := r.DB.Create(roleResource).Error
	if err != nil {
		return err
	}
	return nil
}

func (r RoleResourceRepository) GetRoleResourceByRoleID(roleID uint) (*models.RoleResource, error) {
	var roleResource models.RoleResource
	err := r.DB.Where("role_id = ?", roleID).First(&roleResource).Error
	if err != nil {
		return nil, err
	}
	return &roleResource, nil
}

func (r RoleResourceRepository) UpdateRoleResource(roleResource **models.RoleResource) error {
	err := r.DB.Save(roleResource).Error
	if err != nil {
		return err
	}
	return nil
}

func (r RoleResourceRepository) DeleteRoleResource(roleResource **models.RoleResource) error {
	err := r.DB.Model(roleResource).
		Update("deleted_by", (*roleResource).DeletedBy).
		Delete(roleResource).Error
	if err != nil {
		return err
	}
	return nil
}

func (r RoleResourceRepository) GetRoleResourceByResourceID(resourceID uint) (*models.RoleResource, error) {
	var roleResource models.RoleResource
	err := r.DB.Where("resource_id = ?", resourceID).First(&roleResource).Error
	if err != nil {
		return nil, err
	}
	return &roleResource, nil
}

func (r RoleResourceRepository) GetRoleResourceByRoleIDAndResourceID(roleID, resourceID uint) (*models.RoleResource, error) {
	var roleResource models.RoleResource
	err := r.DB.Where("role_id = ? AND resource_id = ?", roleID, resourceID).First(&roleResource).Error
	if err != nil {
		return nil, err
	}
	return &roleResource, nil
}
