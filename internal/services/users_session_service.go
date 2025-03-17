package services

import (
	"authentication/internal/models/users"
	"authentication/internal/repository"
	"authentication/internal/utils"
	"github.com/rs/zerolog/log"
	"time"
)

type UsersSessionService interface {
	AddUserSession(userID uint, token, refreshToken, ipAddress, userAgent string) error
	GetUserSessionByUserID(userID uint) (*users.UserSession, error)
	LogoutSession(userID uint) error
	CheckUser()
}

type usersSessionService struct {
	UserSessionRepository repository.UserSessionRepository
	UserRepository        repository.UserRepository
	JWTService            utils.JWTService
	Redis                 utils.RedisService
}

func NewUsersSessionService(
	userSessionRepo repository.UserSessionRepository,
	userRepo repository.UserRepository,
	jwtService utils.JWTService,
	redis utils.RedisService,
) UsersSessionService {
	return usersSessionService{
		UserSessionRepository: userSessionRepo,
		UserRepository:        userRepo,
		JWTService:            jwtService,
		Redis:                 redis,
	}
}
func (s usersSessionService) AddUserSession(userID uint, token, refreshToken, ipAddress, userAgent string) error {
	user, err := s.UserRepository.GetUserByID(userID)
	if err != nil {
		return err
	}

	tokenClaims, err := s.JWTService.ExtractClaims(token)
	var userSession = &users.UserSession{
		UserID:       userID,
		SessionToken: token,
		RefreshToken: refreshToken,
		IPAddress:    ipAddress,
		UserAgent:    userAgent,
		LoginTime:    time.Now(),
		LogoutTime:   nil,
		ExpiresAt:    time.Unix(tokenClaims.Exp, 0),
		CreatedBy:    user.ClientID,
		UpdatedBy:    user.ClientID,
	}

	session, err := s.UserSessionRepository.GetUserSessionByUserID(userID)
	if err != nil || session == nil {
		_ = s.UserSessionRepository.AddUserSession(userSession)
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
	session.UpdatedBy = user.ClientID

	err = s.UserSessionRepository.UpdateSession(session)
	if err != nil {
		return err
	}

	_ = s.Redis.SaveData(utils.UserSession, user.ClientID, session)

	return nil
}

func (s usersSessionService) GetUserSessionByUserID(userID uint) (*users.UserSession, error) {
	return s.UserSessionRepository.GetUserSessionByUserID(userID)
}

func (s usersSessionService) LogoutSession(userID uint) error {
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
	session.UpdatedBy = user.ClientID

	_ = s.Redis.SaveData(utils.UserSession, user.ClientID, session)

	return s.UserSessionRepository.UpdateSession(session)
}

func (s usersSessionService) CheckUser() {
	userSession, err := s.UserSessionRepository.GetUserSessionExpired()
	if err != nil {
		log.Error().Err(err).Msg("Failed to get user session")
		return
	}

	for _, session := range *userSession {
		go func(session users.UserSession) {
			if time.Now().After(session.ExpiresAt) && session.IsActive {
				session.IsActive = false
				session.UpdatedBy = "system"
				_ = s.UserSessionRepository.UpdateSession(&session)
				user, _ := s.UserRepository.GetUserByID(session.UserID)
				if user != nil {
					log.Info().Str("client_id", user.ClientID).Msg("User session expired")
					err = s.Redis.DeleteData(utils.UserSession, user.ClientID)
					if err != nil {
						log.Error().Err(err).Msg("Failed to delete data from Redis")
					}
					err = s.Redis.DeleteToken(user.ClientID)
					if err != nil {
						log.Error().Err(err).Msg("Failed to delete token from Redis")
					}
					err = s.Redis.DeleteData(utils.User, user.ClientID)
					if err != nil {
						log.Error().Err(err).Msg("Failed to delete user from Redis")
					}
				} else {
					log.Error().Err(err).Msg("Failed to get user by ID")
				}
			}
		}(session)
	}
}
