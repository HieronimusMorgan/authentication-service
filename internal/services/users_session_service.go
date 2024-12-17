package services

import (
	"authentication/internal/models"
	"authentication/internal/repository"
	"authentication/internal/utils"
	"gorm.io/gorm"
	"time"
)

type UsersSessionService struct {
	UserSessionRepository *repository.UserSessionRepository
	UserRepository        *repository.UserRepository
}

func NewUsersSessionService(db *gorm.DB) *UsersSessionService {
	userSessionRepo := repository.NewUserSessionRepository(db)
	userRepo := repository.NewUserRepository(db)
	return &UsersSessionService{UserSessionRepository: userSessionRepo, UserRepository: userRepo}
}

func (s UsersSessionService) AddUserSession(userID uint, token, refreshToken, ipAddress, userAgent string) error {
	user, err := s.UserRepository.GetUserByID(userID)
	if err != nil {
		return err
	}

	tokenClaims, err := utils.ExtractClaims(token)
	var userSession = &models.UserSession{
		UserID:       userID,
		SessionToken: token,
		RefreshToken: refreshToken,
		IPAddress:    ipAddress,
		UserAgent:    userAgent,
		LoginTime:    time.Now(),
		ExpiresAt:    time.Unix(tokenClaims.Exp, 0),
		CreatedBy:    user.FullName,
	}

	session, err := s.UserSessionRepository.GetUserSessionByUserID(userID)
	if err != nil || session == nil {
		s.UserSessionRepository.AddUserSession(userSession)
		return nil
	}

	session.SessionToken = token
	session.RefreshToken = refreshToken
	session.IPAddress = ipAddress
	session.UserAgent = userAgent
	session.LoginTime = time.Now()
	session.ExpiresAt = time.Unix(tokenClaims.Exp, 0)
	session.IsActive = true
	session.UpdatedBy = user.FullName

	err = s.UserSessionRepository.UpdateSession(session)
	if err != nil {
		return err
	}

	return nil
}
