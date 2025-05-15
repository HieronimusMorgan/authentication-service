package repository

import (
	"authentication/internal/models"
	"gorm.io/gorm"
)

type RoleRepository interface {
	RegisterRole(role *models.Role) error
	GetRoleByID(id uint) (*models.Role, error)
	GetRoleByName(name string) (*models.Role, error)
	GetAllRoles(index, size int) (*[]models.Role, error)
	UpdateRole(role *models.Role) error
	DeleteRole(role models.Role) error
	GetCountRole() (int64, error)
}

type roleRepository struct {
	db gorm.DB
}

func NewRoleRepository(db gorm.DB) RoleRepository {
	return &roleRepository{db: db}
}

func (r roleRepository) RegisterRole(role *models.Role) error {
	err := r.db.Where("name LIKE ?", role.Name).FirstOrCreate(&role).Error
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

func (r roleRepository) GetAllRoles(index, size int) (*[]models.Role, error) {
	var roles []models.Role
	err := r.db.Find(&roles).
		Where("delete_at NOT NULL").
		Order("role_id ASC").
		Limit(size).Offset((index - 1) * size).
		Error
	if err != nil {
		return nil, err
	}
	return &roles, nil
}

func (r roleRepository) UpdateRole(role *models.Role) error {
	err := r.db.Save(&role).Error
	if err != nil {
		return err
	}
	return nil
}

func (r roleRepository) DeleteRole(role models.Role) error {
	err := r.db.Model(role).
		Update("deleted_by", role.DeletedBy).
		Delete(role).Error
	if err != nil {
		return err
	}
	return nil
}

func (r roleRepository) GetCountRole() (int64, error) {
	var count int64
	err := r.db.Model(&models.UserRole{}).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}
