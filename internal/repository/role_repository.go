package repository

import (
	"authentication/internal/models/resource"
	"authentication/internal/models/role"
	"gorm.io/gorm"
)

type RoleRepository interface {
	RegisterRole(role **role.Role) error
	GetRoleByID(id uint) (*role.Role, error)
	GetRoleByName(name string) (*role.Role, error)
	GetAllRoles() (*[]role.Role, error)
	UpdateRole(role **role.Role) error
	DeleteRole(role **role.Role) error
	GetAllRolesByResourceId(resource *resource.Resource) (*[]role.Role, error)
}

type roleRepository struct {
	db gorm.DB
}

func NewRoleRepository(db gorm.DB) RoleRepository {
	return &roleRepository{db: db}
}

func (r roleRepository) RegisterRole(role **role.Role) error {
	err := r.db.Where("name LIKE ?", (*role).Name).FirstOrCreate(role).Error
	if err != nil {
		return err
	}
	return nil
}

func (r roleRepository) GetRoleByID(id uint) (*role.Role, error) {
	var role role.Role
	err := r.db.First(&role, id).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}

func (r roleRepository) GetRoleByName(name string) (*role.Role, error) {
	var role role.Role
	err := r.db.Where("name = ?", name).First(&role).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}

func (r roleRepository) GetAllRoles() (*[]role.Role, error) {
	var roles []role.Role
	err := r.db.Find(&roles).Where("delete_at NOT NULL").Order("role_id ASC").Error
	if err != nil {
		return nil, err
	}
	return &roles, nil
}

func (r roleRepository) UpdateRole(role **role.Role) error {
	err := r.db.Save(role).Error
	if err != nil {
		return err
	}
	return nil
}

func (r roleRepository) DeleteRole(role **role.Role) error {
	err := r.db.Model(role).
		Update("deleted_by", (*role).DeletedBy).
		Delete(role).Error
	if err != nil {
		return err
	}
	return nil
}

func (r roleRepository) GetAllRolesByResourceId(resource *resource.Resource) (*[]role.Role, error) {
	var roles []role.Role
	err := r.db.Table("authentication.roles").
		Select("roles.role_id, roles.name").
		Joins("JOIN authentication.role_resources ON roles.role_id = role_resources.role_id").
		Joins("JOIN authentication.resources ON role_resources.resource_id = resources.resource_id").
		Where("resources.resource_id = ?", resource.ResourceID).
		Find(&roles).Error
	if err != nil {
		return nil, err
	}
	return &roles, nil
}
