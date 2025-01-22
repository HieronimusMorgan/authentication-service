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
		LogoutTime:   nil,
		ExpiresAt:    time.Unix(tokenClaims.Exp, 0),
		CreatedBy:    user.FullName,
		UpdatedBy:    user.FullName,
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
	session.LogoutTime = nil
	session.UpdatedBy = user.FullName

	err = s.UserSessionRepository.UpdateSession(session)
	if err != nil {
		return err
	}

	_ = utils.SaveDataToRedis(utils.UserSession, user.ClientID, session)

	return nil
}

func (s UsersSessionService) GetUserSessionByUserID(userID uint) (*models.UserSession, error) {
	return s.UserSessionRepository.GetUserSessionByUserID(userID)
}

func (s UsersSessionService) LogoutSession(userID uint) error {
	currentTime := time.Now()
	user, err := s.UserRepository.GetUserByID(userID)
	if err != nil {
		return err
	}

	session, err := s.UserSessionRepository.GetUserSessionByUserID(userID)
	if err != nil {
		return nil
	}
	session.IsActive = false
	session.LogoutTime = &currentTime
	session.UpdatedBy = user.FullName

	_ = utils.SaveDataToRedis(utils.UserSession, user.ClientID, session)

	return s.UserSessionRepository.UpdateSession(session)
}
