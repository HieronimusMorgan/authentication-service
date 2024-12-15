package repository

import (
	"authentication/internal/models"
	"gorm.io/gorm"
)

type UserRoleRepository struct {
	DB *gorm.DB
}

func NewUserRoleRepository(db *gorm.DB) *UserRoleRepository {
	return &UserRoleRepository{DB: db}
}

func (r UserRoleRepository) RegisterUserRole(userRole **models.UserRole) error {
	err := r.DB.Create(userRole).Error
	if err != nil {
		return err
	}
	return nil
}

func (r UserRoleRepository) GetUserRoleByUserID(userID uint) (*models.UserRole, error) {
	var userRole models.UserRole
	err := r.DB.Where("user_id = ?", userID).First(&userRole).Error
	if err != nil {
		return nil, err
	}
	return &userRole, nil
}

func (r UserRoleRepository) UpdateUserRole(userRole **models.UserRole) error {
	err := r.DB.Save(userRole).Error
	if err != nil {
		return err
	}
	return nil
}

func (r UserRoleRepository) DeleteUserRole(userRole **models.UserRole) error {
	err := r.DB.Delete(userRole).Error
	if err != nil {
		return err
	}
	return nil
}

func (r UserRoleRepository) GetAllUserRole() (*[]models.UserRole, error) {
	var userRoles []models.UserRole
	err := r.DB.Find(&userRoles).Error
	if err != nil {
		return nil, err
	}
	return &userRoles, nil
}

func (r UserRoleRepository) GetUserRoleByID(id uint) (*models.UserRole, error) {
	var userRole models.UserRole
	err := r.DB.Where("id = ?", id).First(&userRole).Error
	if err != nil {
		return nil, err
	}
	return &userRole, nil
}

func (r UserRoleRepository) GetUserRoleByRoleID(roleID uint) (*models.UserRole, error) {
	var userRole models.UserRole
	err := r.DB.Where("role_id = ?", roleID).First(&userRole).Error
	if err != nil {
		return nil, err
	}
	return &userRole, nil
}
