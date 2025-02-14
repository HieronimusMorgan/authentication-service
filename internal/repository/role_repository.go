package repository

import (
	"authentication/internal/models"
	"gorm.io/gorm"
)

type RoleRepository interface {
	RegisterRole(role **models.Role) error
	GetRoleByID(id uint) (*models.Role, error)
	GetRoleByName(name string) (*models.Role, error)
	GetAllRoles() (*[]models.Role, error)
	UpdateRole(role **models.Role) error
	DeleteRole(role **models.Role) error
	GetAllRolesByResourceId(resource *models.Resource) (*[]models.Role, error)
}

type roleRepository struct {
	db gorm.DB
}

func NewRoleRepository(db gorm.DB) RoleRepository {
	return &roleRepository{db: db}
}

func (r roleRepository) RegisterRole(role **models.Role) error {
	err := r.db.Where("name LIKE ?", (*role).Name).FirstOrCreate(role).Error
	if err != nil {
		return err
	}
	return nil
}

func (r roleRepository) GetRoleByID(id uint) (*models.Role, error) {
	var role models.Role
	err := r.db.First(&role, id).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}

func (r roleRepository) GetRoleByName(name string) (*models.Role, error) {
	var role models.Role
	err := r.db.Where("name = ?", name).First(&role).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}

func (r roleRepository) GetAllRoles() (*[]models.Role, error) {
	var roles []models.Role
	err := r.db.Find(&roles).Where("delete_at NOT NULL").Order("role_id ASC").Error
	if err != nil {
		return nil, err
	}
	return &roles, nil
}

func (r roleRepository) UpdateRole(role **models.Role) error {
	err := r.db.Save(role).Error
	if err != nil {
		return err
	}
	return nil
}

func (r roleRepository) DeleteRole(role **models.Role) error {
	err := r.db.Model(role).
		Update("deleted_by", (*role).DeletedBy).
		Delete(role).Error
	if err != nil {
		return err
	}
	return nil
}

func (r roleRepository) GetAllRolesByResourceId(resource *models.Resource) (*[]models.Role, error) {
	var roles []models.Role
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
