package repository

import (
	"authentication/internal/models/users"
	"gorm.io/gorm"
)

type UserSessionRepository interface {
	GetUserSessionByUserID(userID uint) (*users.UserSession, error)
	GetUserSession() (*[]users.UserSession, error)
	GetUserSessionExpired() (*[]users.UserSession, error)
	AddUserSession(userSession *users.UserSession) error
	UpdateSession(session *users.UserSession) error
}

type userSessionRepository struct {
	db gorm.DB
}

func NewUserSessionRepository(db gorm.DB) UserSessionRepository {
	return &userSessionRepository{db: db}
}

func (r userSessionRepository) GetUserSessionByUserID(userID uint) (*users.UserSession, error) {
	var userSession *users.UserSession
	err := r.db.Where("user_id = ?", userID).First(&userSession).Error
	if err != nil {
		return nil, err
	}
	return userSession, nil
}

func (r userSessionRepository) GetUserSession() (*[]users.UserSession, error) {
	var userSessions *[]users.UserSession
	err := r.db.Find(&userSessions).Error
	if err != nil {
		return nil, err
	}
	return userSessions, nil
}

func (r userSessionRepository) GetUserSessionExpired() (*[]users.UserSession, error) {
	var userSessions *[]users.UserSession
	err := r.db.Where("expired_at < now()").Find(&userSessions).Error
	if err != nil {
		return nil, err
	}
	return userSessions, nil
}

func (r userSessionRepository) AddUserSession(userSession *users.UserSession) error {
	err := r.db.Create(userSession).Error
	if err != nil {
		return err
	}
	return nil
}

func (r userSessionRepository) UpdateSession(session *users.UserSession) error {
	err := r.db.Save(session).Error
	if err != nil {
		return err
	}
	return nil
}
