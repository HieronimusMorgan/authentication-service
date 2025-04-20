package repository

import (
	"authentication/internal/models"
	"gorm.io/gorm"
)

type UserSessionRepository interface {
	GetUserSessionByUserID(userID uint) (*models.UserSession, error)
	GetUserSession() (*[]models.UserSession, error)
	GetUserSessionExpired() (*[]models.UserSession, error)
	AddUserSession(userSession *models.UserSession) error
	UpdateSession(session *models.UserSession) error
	GetUserSessionByRefreshTokenAndUserID(userID uint, refreshToken string) (*models.UserSession, error)
}

type userSessionRepository struct {
	db gorm.DB
}

func NewUserSessionRepository(db gorm.DB) UserSessionRepository {
	return &userSessionRepository{db: db}
}

func (r userSessionRepository) GetUserSessionByUserID(userID uint) (*models.UserSession, error) {
	var userSession *models.UserSession
	err := r.db.Where("user_id = ?", userID).First(&userSession).Error
	if err != nil {
		return nil, err
	}
	return userSession, nil
}

func (r userSessionRepository) GetUserSession() (*[]models.UserSession, error) {
	var userSessions *[]models.UserSession
	err := r.db.Find(&userSessions).Error
	if err != nil {
		return nil, err
	}
	return userSessions, nil
}

func (r userSessionRepository) GetUserSessionExpired() (*[]models.UserSession, error) {
	var userSessions *[]models.UserSession
	err := r.db.Where("expired_at < now()").Find(&userSessions).Error
	if err != nil {
		return nil, err
	}
	return userSessions, nil
}

func (r userSessionRepository) AddUserSession(userSession *models.UserSession) error {
	err := r.db.Create(userSession).Error
	if err != nil {
		return err
	}
	return nil
}

func (r userSessionRepository) UpdateSession(session *models.UserSession) error {
	err := r.db.Save(session).Error
	if err != nil {
		return err
	}
	return nil
}

func (r userSessionRepository) GetUserSessionByRefreshTokenAndUserID(userID uint, refreshToken string) (*models.UserSession, error) {
	var userSession *models.UserSession
	err := r.db.Where("user_id = ? AND refresh_token = ?", userID, refreshToken).First(&userSession).Error
	if err != nil {
		return nil, err
	}
	return userSession, nil
}
