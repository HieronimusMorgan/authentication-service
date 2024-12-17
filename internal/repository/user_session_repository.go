package repository

import (
	"authentication/internal/models"
	"gorm.io/gorm"
)

type UserSessionRepository struct {
	DB *gorm.DB
}

func NewUserSessionRepository(db *gorm.DB) *UserSessionRepository {
	return &UserSessionRepository{DB: db}
}

func (r UserSessionRepository) GetUserSessionByUserID(userID uint) (*models.UserSession, error) {
	var userSession *models.UserSession
	err := r.DB.Where("user_id = ?", userID).First(&userSession).Error
	if err != nil {
		return nil, err
	}
	return userSession, nil
}

func (r UserSessionRepository) AddUserSession(userSession *models.UserSession) error {
	err := r.DB.Create(userSession).Error
	if err != nil {
		return err
	}
	return nil
}

func (r UserSessionRepository) UpdateSession(session *models.UserSession) error {
	err := r.DB.Save(session).Error
	if err != nil {
		return err
	}
	return nil
}
