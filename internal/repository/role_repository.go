package repository

import (
	"authentication/internal/models"
	"gorm.io/gorm"
)

type RoleRepository struct {
	DB *gorm.DB
}

func NewRoleRepository(db *gorm.DB) *RoleRepository {
	return &RoleRepository{DB: db}
}

func (r RoleRepository) RegisterRole(role **models.Role) error {
	err := r.DB.Where("name LIKE ?", (*role).Name).FirstOrCreate(role).Error
	if err != nil {
		return err
	}
	return nil
}

func (r RoleRepository) GetRoleByID(id uint) (*models.Role, error) {
	var role models.Role
	err := r.DB.First(&role, id).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}

func (r RoleRepository) GetRoleByName(name string) (*models.Role, error) {
	var role models.Role
	err := r.DB.Where("name = ?", name).First(&role).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}

func (r RoleRepository) GetAllRoles() (*[]models.Role, error) {
	var roles []models.Role
	err := r.DB.Find(&roles).Where("delete_at NOT NULL").Order("role_id ASC").Error
	if err != nil {
		return nil, err
	}
	return &roles, nil
}

func (r RoleRepository) UpdateRole(role **models.Role) error {
	err := r.DB.Save(role).Error
	if err != nil {
		return err
	}
	return nil
}

func (r RoleRepository) DeleteRole(role **models.Role) error {
	err := r.DB.Model(role).
		Update("deleted_by", (*role).DeletedBy).
		Delete(role).Error
	if err != nil {
		return err
	}
	return nil
}

func (r RoleRepository) GetAllRolesByResourceId(resource *models.Resource) (*[]models.Role, error) {
	var roles []models.Role
	err := r.DB.Table("authentication.roles").
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
