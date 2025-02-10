package repository

import (
	"authentication/internal/models"
	"gorm.io/gorm"
)

type RoleResourceRepository interface {
	RegisterRoleResource(roleResource **models.RoleResource) error
	GetRoleResourceByRoleID(roleID uint) (*models.RoleResource, error)
	UpdateRoleResource(roleResource **models.RoleResource) error
	DeleteRoleResource(roleResource **models.RoleResource) error
	GetRoleResourceByResourceID(resourceID uint) (*models.RoleResource, error)
	GetRoleResourceByRoleIDAndResourceID(roleID, resourceID uint) (*models.RoleResource, error)
}

type roleResourceRepository struct {
	db *gorm.DB
}

func NewRoleResourceRepository(db *gorm.DB) RoleResourceRepository {
	return &roleResourceRepository{db: db}
}

func (r roleResourceRepository) RegisterRoleResource(roleResource **models.RoleResource) error {
	err := r.db.Create(roleResource).Error
	if err != nil {
		return err
	}
	return nil
}

func (r roleResourceRepository) GetRoleResourceByRoleID(roleID uint) (*models.RoleResource, error) {
	var roleResource models.RoleResource
	err := r.db.Where("role_id = ?", roleID).First(&roleResource).Error
	if err != nil {
		return nil, err
	}
	return &roleResource, nil
}

func (r roleResourceRepository) UpdateRoleResource(roleResource **models.RoleResource) error {
	err := r.db.Save(roleResource).Error
	if err != nil {
		return err
	}
	return nil
}

func (r roleResourceRepository) DeleteRoleResource(roleResource **models.RoleResource) error {
	err := r.db.Model(roleResource).
		Update("deleted_by", (*roleResource).DeletedBy).
		Delete(roleResource).Error
	if err != nil {
		return err
	}
	return nil
}

func (r roleResourceRepository) GetRoleResourceByResourceID(resourceID uint) (*models.RoleResource, error) {
	var roleResource models.RoleResource
	err := r.db.Where("resource_id = ?", resourceID).First(&roleResource).Error
	if err != nil {
		return nil, err
	}
	return &roleResource, nil
}

func (r roleResourceRepository) GetRoleResourceByRoleIDAndResourceID(roleID, resourceID uint) (*models.RoleResource, error) {
	var roleResource models.RoleResource
	err := r.db.Where("role_id = ? AND resource_id = ?", roleID, resourceID).First(&roleResource).Error
	if err != nil {
		return nil, err
	}
	return &roleResource, nil
}
