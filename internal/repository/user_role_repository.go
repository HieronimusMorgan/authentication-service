package repository

import (
	"authentication/internal/models/users"
	"gorm.io/gorm"
)

type UserRoleRepository interface {
	RegisterUserRole(userRole **users.UserRole) error
	GetUserRoleByUserID(userID uint) (*users.UserRole, error)
	UpdateUserRole(userRole **users.UserRole) error
	DeleteUserRole(userRole **users.UserRole) error
	GetAllUserRole() (*[]users.UserRole, error)
	GetUserRoleByID(id uint) (*users.UserRole, error)
	GetUserRoleByRoleID(roleID uint) (*users.UserRole, error)
}

type userRoleRepository struct {
	db gorm.DB
}

func NewUserRoleRepository(db gorm.DB) UserRoleRepository {
	return &userRoleRepository{db: db}
}

func (r userRoleRepository) RegisterUserRole(userRole **users.UserRole) error {
	err := r.db.Create(userRole).Error
	if err != nil {
		return err
	}
	return nil
}

func (r userRoleRepository) GetUserRoleByUserID(userID uint) (*users.UserRole, error) {
	var userRole users.UserRole
	err := r.db.Where("user_id = ?", userID).First(&userRole).Error
	if err != nil {
		return nil, err
	}
	return &userRole, nil
}

func (r userRoleRepository) UpdateUserRole(userRole **users.UserRole) error {
	err := r.db.Save(userRole).Error
	if err != nil {
		return err
	}
	return nil
}

func (r userRoleRepository) DeleteUserRole(userRole **users.UserRole) error {
	err := r.db.Model(userRole).
		Update("deleted_by", (*userRole).DeletedBy).
		Delete(userRole).Error
	if err != nil {
		return err
	}
	return nil
}

func (r userRoleRepository) GetAllUserRole() (*[]users.UserRole, error) {
	var userRoles []users.UserRole
	err := r.db.Find(&userRoles).Error
	if err != nil {
		return nil, err
	}
	return &userRoles, nil
}

func (r userRoleRepository) GetUserRoleByID(id uint) (*users.UserRole, error) {
	var userRole users.UserRole
	err := r.db.Where("id = ?", id).First(&userRole).Error
	if err != nil {
		return nil, err
	}
	return &userRole, nil
}

func (r userRoleRepository) GetUserRoleByRoleID(roleID uint) (*users.UserRole, error) {
	var userRole users.UserRole
	err := r.db.Where("role_id = ?", roleID).First(&userRole).Error
	if err != nil {
		return nil, err
	}
	return &userRole, nil
}
