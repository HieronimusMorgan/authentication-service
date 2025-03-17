package repository

import (
	"authentication/internal/models/role"
	"gorm.io/gorm"
)

type RoleResourceRepository interface {
	RegisterRoleResource(roleResource **role.RoleResource) error
	GetRoleResourceByRoleID(roleID uint) (*role.RoleResource, error)
	UpdateRoleResource(roleResource **role.RoleResource) error
	DeleteRoleResource(roleResource **role.RoleResource) error
	GetRoleResourceByResourceID(resourceID uint) (*role.RoleResource, error)
	GetRoleResourceByRoleIDAndResourceID(roleID, resourceID uint) (*role.RoleResource, error)
}

type roleResourceRepository struct {
	db gorm.DB
}

func NewRoleResourceRepository(db gorm.DB) RoleResourceRepository {
	return &roleResourceRepository{db: db}
}

func (r roleResourceRepository) RegisterRoleResource(roleResource **role.RoleResource) error {
	err := r.db.Create(roleResource).Error
	if err != nil {
		return err
	}
	return nil
}

func (r roleResourceRepository) GetRoleResourceByRoleID(roleID uint) (*role.RoleResource, error) {
	var roleResource role.RoleResource
	err := r.db.Where("role_id = ?", roleID).First(&roleResource).Error
	if err != nil {
		return nil, err
	}
	return &roleResource, nil
}

func (r roleResourceRepository) UpdateRoleResource(roleResource **role.RoleResource) error {
	err := r.db.Save(roleResource).Error
	if err != nil {
		return err
	}
	return nil
}

func (r roleResourceRepository) DeleteRoleResource(roleResource **role.RoleResource) error {
	err := r.db.Model(roleResource).
		Update("deleted_by", (*roleResource).DeletedBy).
		Delete(roleResource).Error
	if err != nil {
		return err
	}
	return nil
}

func (r roleResourceRepository) GetRoleResourceByResourceID(resourceID uint) (*role.RoleResource, error) {
	var roleResource role.RoleResource
	err := r.db.Where("resource_id = ?", resourceID).First(&roleResource).Error
	if err != nil {
		return nil, err
	}
	return &roleResource, nil
}

func (r roleResourceRepository) GetRoleResourceByRoleIDAndResourceID(roleID, resourceID uint) (*role.RoleResource, error) {
	var roleResource role.RoleResource
	err := r.db.Where("role_id = ? AND resource_id = ?", roleID, resourceID).First(&roleResource).Error
	if err != nil {
		return nil, err
	}
	return &roleResource, nil
}
