package services

import (
	"authentication/internal/dto/in"
	"authentication/internal/dto/out"
	"authentication/internal/models"
	"authentication/internal/repository"
	"authentication/internal/utils"
	"authentication/package/response"
	"net/http"
)

type UserService interface {
	GetProfile(clientID string) (*out.UserResponse, response.ErrorResponse)
	UpdateNameUserProfile(updateNameRequest *in.UpdateNameRequest, clientID string) (interface{}, error)
	UpdatePhotoUserProfile(req *in.UpdatePhotoRequest, clientID string) (interface{}, error)
	DeleteUserById(userID uint, clientID string) response.ErrorResponse
}

type userService struct {
	UserRepository         repository.UserRepository
	ResourceRepository     repository.ResourceRepository
	RoleRepository         repository.RoleRepository
	RoleResourceRepository repository.RoleResourceRepository
	UserRoleRepository     repository.UserRoleRepository
	UserSessionRepository  repository.UserSessionRepository
	RedisService           utils.RedisService
	JWTService             utils.JWTService
	Encryption             utils.Encryption
}

func NewUserService(
	userRepo repository.UserRepository,
	redis utils.RedisService,
	jwtService utils.JWTService,
	Encryption utils.Encryption) UserService {
	return userService{
		UserRepository: userRepo,
		RedisService:   redis,
		JWTService:     jwtService,
		Encryption:     Encryption,
	}
}

func (s userService) GetProfile(clientID string) (*out.UserResponse, response.ErrorResponse) {

	user, err := s.UserRepository.GetUserResponseByClientID(clientID)
	if err != nil {
		return nil, response.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "User not found",
			Error:   err.Error(),
		}
	}

	var phoneNumber string
	decrypt, err := s.Encryption.Decrypt(user.PhoneNumber)
	if err != nil {
		phoneNumber = user.PhoneNumber
	} else {
		phoneNumber = decrypt
	}

	user.PhoneNumber = phoneNumber

	return user, response.ErrorResponse{}
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

	err = s.UserRepository.UpdateProfile(user)
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

func (s userService) UpdatePhotoUserProfile(req *in.UpdatePhotoRequest, clientID string) (interface{}, error) {
	user, err := s.UserRepository.GetUserByClientID(clientID)
	if err != nil {
		return nil, err
	}

	user.ProfilePicture = req.ProfilePicture
	user.UpdatedBy = user.FullName

	err = s.UserRepository.UpdateProfile(user)
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
