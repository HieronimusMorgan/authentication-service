package repository

import (
	"authentication/internal/models"
	"gorm.io/gorm"
)

type UserRoleRepository interface {
	RegisterUserRole(userRole **models.UserRole) error
	GetUserRoleByUserID(userID uint) (*models.UserRole, error)
	UpdateUserRole(userRole **models.UserRole) error
	DeleteUserRole(userRole **models.UserRole) error
	GetAllUserRole() (*[]models.UserRole, error)
	GetUserRoleByID(id uint) (*models.UserRole, error)
	GetUserRoleByRoleID(roleID uint) (*models.UserRole, error)
}

type userRoleRepository struct {
	db gorm.DB
}

func NewUserRoleRepository(db gorm.DB) UserRoleRepository {
	return &userRoleRepository{db: db}
}

func (r userRoleRepository) RegisterUserRole(userRole **models.UserRole) error {
	err := r.db.Create(userRole).Error
	if err != nil {
		return err
	}
	return nil
}

func (r userRoleRepository) GetUserRoleByUserID(userID uint) (*models.UserRole, error) {
	var userRole models.UserRole
	err := r.db.Where("user_id = ?", userID).First(&userRole).Error
	if err != nil {
		return nil, err
	}
	return &userRole, nil
}

func (r userRoleRepository) UpdateUserRole(userRole **models.UserRole) error {
	err := r.db.Save(userRole).Error
	if err != nil {
		return err
	}
	return nil
}

func (r userRoleRepository) DeleteUserRole(userRole **models.UserRole) error {
	err := r.db.Model(userRole).
		Update("deleted_by", (*userRole).DeletedBy).
		Delete(userRole).Error
	if err != nil {
		return err
	}
	return nil
}

func (r userRoleRepository) GetAllUserRole() (*[]models.UserRole, error) {
	var userRoles []models.UserRole
	err := r.db.Find(&userRoles).Error
	if err != nil {
		return nil, err
	}
	return &userRoles, nil
}

func (r userRoleRepository) GetUserRoleByID(id uint) (*models.UserRole, error) {
	var userRole models.UserRole
	err := r.db.Where("id = ?", id).First(&userRole).Error
	if err != nil {
		return nil, err
	}
	return &userRole, nil
}

func (r userRoleRepository) GetUserRoleByRoleID(roleID uint) (*models.UserRole, error) {
	var userRole models.UserRole
	err := r.db.Where("role_id = ?", roleID).First(&userRole).Error
	if err != nil {
		return nil, err
	}
	return &userRole, nil
}
