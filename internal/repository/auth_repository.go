package repository

import (
	"authentication/internal/models"
	"gorm.io/gorm"
)

type AuthRepository interface {
	UpdatePinCode(user *models.Users) error
	CreateInternalToken(resourceID uint, token string) error
}

type authRepository struct {
	db gorm.DB
}

func NewAuthRepository(db gorm.DB) AuthRepository {
	return &authRepository{db: db}
}

func (r authRepository) UpdatePinCode(user *models.Users) error {
	return r.db.Save(user).Error
}

func (r authRepository) CreateInternalToken(resourceID uint, token string) error {
	if err := r.db.Create(&models.InternalToken{
		ResourceID: resourceID,
		Token:      token,
	}).Error; err != nil {
		return err
	}
	return nil
}
