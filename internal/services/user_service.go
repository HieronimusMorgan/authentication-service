package services

import (
	"authentication/internal/dto/in"
	"authentication/internal/dto/out"
	"authentication/internal/models"
	"authentication/internal/repository"
	"authentication/internal/utils"
	"authentication/package/response"
	"net/http"
	"time"
)

type UserService interface {
	GetProfile(clientID string) (*out.UserResponse, response.ErrorResponse)
	AddUserKey(clientID string) error
	UpdateNameUserProfile(updateNameRequest *in.UpdateNameRequest, clientID string) (interface{}, error)
	UpdatePhotoUserProfile(req string, clientID string) (interface{}, error)
	UpdateUserSetting(userSetting *in.UserSettingsRequest, clientID string) response.ErrorResponse
	DeleteUserById(userID uint, clientID string) response.ErrorResponse
}

type userService struct {
	UserRepository         repository.UserRepository
	UserKeyRepository      repository.UserKeyRepository
	UserSettingRepository  repository.UserSettingRepository
	ResourceRepository     repository.ResourceRepository
	RoleRepository         repository.RoleRepository
	UserResourceRepository repository.UserResourceRepository
	UserRoleRepository     repository.UserRoleRepository
	UserSessionRepository  repository.UserSessionRepository
	RedisService           utils.RedisService
	JWTService             utils.JWTService
	Encryption             utils.Encryption
}

func NewUserService(
	userRepo repository.UserRepository,
	userKeyRepo repository.UserKeyRepository,
	userSettingRepository repository.UserSettingRepository,
	redis utils.RedisService,
	jwtService utils.JWTService,
	Encryption utils.Encryption) UserService {
	return userService{
		UserRepository:        userRepo,
		UserKeyRepository:     userKeyRepo,
		UserSettingRepository: userSettingRepository,
		RedisService:          redis,
		JWTService:            jwtService,
		Encryption:            Encryption,
	}
}

func (s userService) GetProfile(clientID string) (*out.UserResponse, response.ErrorResponse) {
	var userResponse out.UserResponse
	user, err := s.UserRepository.GetUserByClientID(clientID)
	if err != nil {
		return nil, response.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "User not found",
			Error:   err.Error(),
		}
	}
	userResponse.UserID = user.UserID
	userResponse.ClientID = user.ClientID
	userResponse.Username = user.Username
	userResponse.FirstName = user.FirstName
	userResponse.LastName = user.LastName
	userResponse.ProfilePicture = user.ProfilePicture

	var phoneNumber string
	decrypt, err := s.Encryption.Decrypt(user.PhoneNumber)
	if err != nil {
		phoneNumber = user.PhoneNumber
	} else {
		phoneNumber = decrypt
	}

	userResponse.PhoneNumber = phoneNumber
	userSetting, err := s.UserSettingRepository.GetUserSettingByUserID(user.UserID)
	if err != nil {
		return nil, response.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "User setting not found",
			Error:   err.Error(),
		}
	}

	userSettingModel := out.UserSettingResponse{
		SettingID:             userSetting.SettingID,
		GroupInviteType:       userSetting.GroupInviteType,
		GroupInviteDisallowed: userSetting.GroupInviteDisallowed,
	}

	userResponse.UserSetting = userSettingModel

	return &userResponse, response.ErrorResponse{}
}

func (s userService) AddUserKey(clientID string) error {
	user, err := s.UserRepository.GetUserByClientID(clientID)
	if err != nil {
		return err
	}

	userKey, err := s.UserKeyRepository.GetUserKeyByUserID(user.UserID)
	if userKey == nil && err != nil {
		userKeys, err := utils.GenerateUserKey(user)
		if err != nil {
			return err
		}

		if err := s.UserRepository.SaveUserKey(userKeys); err != nil {
			return err
		}
	}
	return nil
}

func (s userService) UpdateNameUserProfile(updateNameRequest *in.UpdateNameRequest, clientID string) (interface{}, error) {
	user, err := s.UserRepository.GetUserByClientID(clientID)
	if err != nil {
		return nil, err
	}

	user.FirstName = utils.ValidationTrimSpace(updateNameRequest.FirstName)
	user.LastName = utils.ValidationTrimSpace(updateNameRequest.LastName)
	user.FullName = user.FirstName + " " + user.LastName
	user.UpdatedBy = user.FullName

	err = s.UserRepository.UpdateUser(user)
	if err != nil {
		return nil, err
	}

	var phoneNumber string
	decrypt, err := s.Encryption.Decrypt(user.PhoneNumber)
	if err != nil {
		phoneNumber = user.PhoneNumber
	} else {
		phoneNumber = decrypt
	}
	return out.UserResponse{
		UserID:         user.UserID,
		ClientID:       user.ClientID,
		Username:       user.Username,
		FirstName:      user.FirstName,
		LastName:       user.LastName,
		PhoneNumber:    phoneNumber,
		ProfilePicture: user.ProfilePicture,
	}, nil
}

func (s userService) UpdatePhotoUserProfile(req string, clientID string) (interface{}, error) {
	user, err := s.UserRepository.GetUserByClientID(clientID)
	if err != nil {
		return nil, err
	}

	user.ProfilePicture = &req
	user.UpdatedBy = user.FullName

	err = s.UserRepository.UpdateUser(user)
	if err != nil {
		return nil, err
	}

	var phoneNumber string
	decrypt, err := s.Encryption.Decrypt(user.PhoneNumber)
	if err != nil {
		phoneNumber = user.PhoneNumber
	} else {
		phoneNumber = decrypt
	}
	return out.UserResponse{
		UserID:         user.UserID,
		ClientID:       user.ClientID,
		Username:       user.Username,
		FirstName:      user.FirstName,
		LastName:       user.LastName,
		PhoneNumber:    phoneNumber,
		ProfilePicture: user.ProfilePicture,
	}, nil
}

func (s userService) UpdateUserSetting(userSetting *in.UserSettingsRequest, clientID string) response.ErrorResponse {
	user, err := s.UserRepository.GetUserByClientID(clientID)
	if err != nil {
		return response.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "user not found",
			Error:   err.Error(),
		}
	}

	userSettingModel, err := s.UserSettingRepository.GetUserSettingByUserIDAndSettingID(user.UserID, userSetting.SettingID)
	if err != nil {
		return response.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "user setting not found",
			Error:   err.Error(),
		}
	}

	userSettingModel.GroupInviteType = userSetting.GroupInviteType

	if len(userSetting.GroupInviteDisallowed) > 0 {
		userSettingModel.GroupInviteDisallowed = userSetting.GroupInviteDisallowed
	} else {
		userSettingModel.GroupInviteDisallowed = nil
	}
	userSettingModel.UpdatedAt = time.Now()

	err = s.UserSettingRepository.UpdateUserSetting(userSettingModel)
	if err != nil {
		return response.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "failed to update user settings",
			Error:   err.Error(),
		}
	}

	return response.ErrorResponse{}
}

func (s userService) DeleteUserById(userID uint, clientID string) response.ErrorResponse {
	admin, err := s.UserRepository.GetUserByClientID(clientID)
	if err != nil {
		return response.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "user is not an admin",
			Error:   err.Error(),
		}
	}

	var user *models.Users
	user, err = s.UserRepository.GetUserByID(userID)
	if err != nil {
		return response.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "User not found",
			Error:   err.Error(),
		}
	}

	user.DeletedBy = admin.FullName
	err = s.UserRepository.DeleteUser(user)
	if err != nil {
		return response.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "User ",
			Error:   err.Error(),
		}
	}
	return response.ErrorResponse{}
}
