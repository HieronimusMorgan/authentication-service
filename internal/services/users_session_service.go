package services

import (
	"authentication/internal/models"
	"authentication/internal/repository"
	"authentication/internal/utils"
	"time"
)

// UsersSessionService manages user session operations
type UsersSessionService struct {
	UserSessionRepository repository.UserSessionRepository
	UserRepository        repository.UserRepository
	JWTService            utils.JWTService
	Redis                 utils.RedisService
}

// NewUsersSessionService initializes UsersSessionService with repository dependencies
func NewUsersSessionService(
	userSessionRepo repository.UserSessionRepository,
	userRepo repository.UserRepository,
	jwtService utils.JWTService,
	redis utils.RedisService,
) UsersSessionService {
	return UsersSessionService{
		UserSessionRepository: userSessionRepo,
		UserRepository:        userRepo,
		JWTService:            jwtService,
		Redis:                 redis,
	}
}
func (s UsersSessionService) AddUserSession(userID uint, token, refreshToken, ipAddress, userAgent string) error {
	user, err := s.UserRepository.GetUserByID(userID)
	if err != nil {
		return err
	}

	tokenClaims, err := s.JWTService.ExtractClaims(token)
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

	_ = s.Redis.SaveData(utils.UserSession, user.ClientID, session)

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

	_ = s.Redis.SaveData(utils.UserSession, user.ClientID, session)

	return s.UserSessionRepository.UpdateSession(session)
}
