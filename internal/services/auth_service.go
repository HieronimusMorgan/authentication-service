package services

import (
	"authentication/internal/dto/in"
	"authentication/internal/dto/out"
	"authentication/internal/models"
	"authentication/internal/repository"
	"authentication/internal/utils"
	"authentication/package/response"
	"net/http"
	"strings"
)

type AuthService interface {
	Register(req *in.RegisterRequest) (out.RegisterResponse, response.ErrorResponse)
	Login(req *in.LoginRequest) (interface{}, response.ErrorResponse)
	GetProfile(clientID string) (*out.UserResponse, response.ErrorResponse)
	UpdateNameUserProfile(updateNameRequest *in.UpdateNameRequest, clientID string) (interface{}, error)
	UpdatePhotoUserProfile(req *in.UpdatePhotoRequest, clientID string) (interface{}, error)
	RegisterInternalToken(req *struct {
		ResourceName string `json:"resource_name" binding:"required"`
	}) (interface{}, response.ErrorResponse)
	DeleteUserById(userID uint, clientID string) response.ErrorResponse
	UpdateRole(userID uint, roleID uint, clientID string) response.ErrorResponse
	GetListUser(clientID string) (interface{}, response.ErrorResponse)
	ChangePassword(password *struct {
		OldPassword string `json:"old_password" binding:"required"`
		NewPassword string `json:"new_password" binding:"required"`
	}, clientID string) response.ErrorResponse
}

type authService struct {
	AuthRepository         repository.AuthRepository
	ResourceRepository     repository.ResourceRepository
	RoleRepository         repository.RoleRepository
	RoleResourceRepository repository.RoleResourceRepository
	UserRepository         repository.UserRepository
	UserRoleRepository     repository.UserRoleRepository
	UserSessionRepository  repository.UserSessionRepository
	RedisService           utils.RedisService
	JWTService             utils.JWTService
}

func NewAuthService(
	authRepo repository.AuthRepository,
	resourceRepo repository.ResourceRepository,
	roleRepo repository.RoleRepository,
	roleResourceRepo repository.RoleResourceRepository,
	userRepo repository.UserRepository,
	userRoleRepo repository.UserRoleRepository,
	userSessionRepo repository.UserSessionRepository,
	redis utils.RedisService,
	jwtService utils.JWTService) AuthService {
	return authService{
		AuthRepository:         authRepo,
		ResourceRepository:     resourceRepo,
		RoleRepository:         roleRepo,
		RoleResourceRepository: roleResourceRepo,
		UserRepository:         userRepo,
		UserRoleRepository:     userRoleRepo,
		UserSessionRepository:  userSessionRepo,
		RedisService:           redis,
		JWTService:             jwtService,
	}
}

func (s authService) checkUserIsAdmin(user *models.Users) response.ErrorResponse {
	role, err := s.RoleRepository.GetRoleByID(user.RoleID)
	if err != nil {
		return response.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "failed to get role",
			Error:   err.Error(),
		}
	}
	if strings.EqualFold(role.Name, "Admin") || strings.EqualFold(role.Name, "Super Admin") {
		return response.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "user is not an admin",
			Error:   err.Error(),
		}
	}
	return response.ErrorResponse{}
}

func (s authService) Register(req *in.RegisterRequest) (out.RegisterResponse, response.ErrorResponse) {
	if err := utils.ValidateUsername(req.Username); err != nil {
		return out.RegisterResponse{}, response.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Validation Username",
			Error:   err.Error(),
		}
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return out.RegisterResponse{}, response.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid Password",
			Error:   err.Error(),
		}
	}

	role, err := s.RoleRepository.GetRoleByName("User")
	if err != nil {
		return out.RegisterResponse{}, response.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Unable to get role",
			Error:   err.Error(),
		}
	}

	firstName := utils.ValidationTrimSpace(req.FirstName)
	lastName := utils.ValidationTrimSpace(req.LastName)
	fullName := firstName + " " + lastName

	user := &models.Users{
		ClientID:       utils.GenerateClientID(),
		Username:       req.Username,
		Password:       hashedPassword,
		FirstName:      firstName,
		LastName:       lastName,
		FullName:       fullName,
		PhoneNumber:    req.PhoneNumber,
		ProfilePicture: req.ProfilePicture,
		RoleID:         role.RoleID,
		CreatedBy:      "system",
		UpdatedBy:      "system",
	}

	if err := s.UserRepository.RegisterUser(&user); err != nil {
		return out.RegisterResponse{}, response.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Unable to register user",
			Error:   err.Error(),
		}
	}

	userRole := &models.UserRole{
		UserID:    user.UserID,
		RoleID:    role.RoleID,
		CreatedBy: "system",
		UpdatedBy: "system",
	}

	if err := s.UserRoleRepository.RegisterUserRole(&userRole); err != nil {
		return out.RegisterResponse{}, response.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Unable to register user role",
			Error:   err.Error(),
		}
	}

	user.Role = *role
	token, err := s.JWTService.GenerateToken(*user)
	if err != nil {
		return out.RegisterResponse{}, response.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Token is invalid",
			Error:   err.Error(),
		}
	}

	_ = s.RedisService.SaveData(utils.Token, user.ClientID, token)
	_ = s.RedisService.SaveData(utils.User, user.ClientID, user)

	responses := out.RegisterResponse{
		UserID:         user.UserID,
		Username:       user.Username,
		FirstName:      user.FirstName,
		LastName:       user.LastName,
		PhoneNumber:    user.PhoneNumber,
		ProfilePicture: user.ProfilePicture,
		Token:          token.AccessToken,
		RefreshToken:   token.RefreshToken,
	}
	return responses, response.ErrorResponse{}
}

func (s authService) Login(req *in.LoginRequest) (interface{}, response.ErrorResponse) {
	user, err := s.AuthRepository.GetUserByUsername(req.Username)
	if err != nil {
		return nil, response.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Username or Password is incorrect",
			Error:   err.Error(),
		}
	}
	if err := utils.CheckPassword(user.Password, req.Password); err != nil {
		return nil, response.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Username or Password is incorrect",
			Error:   err.Error(),
		}
	}

	role, err := s.RoleRepository.GetRoleByID(user.RoleID)
	user.Role = *role

	token, err := s.JWTService.GenerateToken(*user)
	if err != nil {
		return nil, response.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "User or Password is incorrect",
			Error:   err.Error(),
		}
	}

	_ = s.RedisService.SaveData("token", user.ClientID, token)
	_ = s.RedisService.SaveData("user", user.ClientID, user)

	responses := out.LoginResponse{
		UserID:         user.UserID,
		Username:       user.Username,
		FirstName:      user.FirstName,
		LastName:       user.LastName,
		PhoneNumber:    user.PhoneNumber,
		ProfilePicture: user.ProfilePicture,
		Token:          token.AccessToken,
		RefreshToken:   token.RefreshToken,
	}
	return responses, response.ErrorResponse{}
}

func (s authService) GetProfile(clientID string) (*out.UserResponse, response.ErrorResponse) {

	user, err := s.AuthRepository.GetUserResponseByClientID(clientID)
	if err != nil {
		return nil, response.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "User not found",
			Error:   err.Error(),
		}
	}

	return user, response.ErrorResponse{}
}

func (s authService) UpdateNameUserProfile(updateNameRequest *in.UpdateNameRequest, clientID string) (interface{}, error) {
	user, err := s.AuthRepository.GetUserByClientID(clientID)
	if err != nil {
		return nil, err
	}

	user.FirstName = utils.ValidationTrimSpace(updateNameRequest.FirstName)
	user.LastName = utils.ValidationTrimSpace(updateNameRequest.LastName)
	user.FullName = user.FirstName + " " + user.LastName
	user.UpdatedBy = user.FullName

	err = s.AuthRepository.UpdateProfile(user)
	if err != nil {
		return nil, err
	}

	return out.UserResponse{
		UserID:         user.UserID,
		ClientID:       user.ClientID,
		Username:       user.Username,
		FirstName:      user.FirstName,
		LastName:       user.LastName,
		PhoneNumber:    user.PhoneNumber,
		ProfilePicture: user.ProfilePicture,
	}, nil
}

func (s authService) UpdatePhotoUserProfile(req *in.UpdatePhotoRequest, clientID string) (interface{}, error) {
	user, err := s.AuthRepository.GetUserByClientID(clientID)
	if err != nil {
		return nil, err
	}

	user.ProfilePicture = req.ProfilePicture
	user.UpdatedBy = user.FullName

	err = s.AuthRepository.UpdateProfile(user)
	if err != nil {
		return nil, err
	}

	return out.UserResponse{
		UserID:         user.UserID,
		ClientID:       user.ClientID,
		Username:       user.Username,
		FirstName:      user.FirstName,
		LastName:       user.LastName,
		PhoneNumber:    user.PhoneNumber,
		ProfilePicture: user.ProfilePicture,
	}, nil
}

func (s authService) RegisterInternalToken(req *struct {
	ResourceName string `json:"resource_name" binding:"required"`
}) (interface{}, response.ErrorResponse) {
	resource, err := s.ResourceRepository.GetResourceByName(req.ResourceName)
	if err != nil {
		return nil, response.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "user is not an admin",
			Error:   err.Error(),
		}
	}

	token, err := s.JWTService.GenerateInternalToken(resource.Name)
	if err != nil {
		return nil, response.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "user is not an admin",
			Error:   err.Error(),
		}
	}

	if err := s.ResourceRepository.CreateInternalToken(resource.ResourceID, token); err != nil {
		return nil, response.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "user is not an admin",
			Error:   err.Error(),
		}
	}

	return token, response.ErrorResponse{}
}

func (s authService) DeleteUserById(userID uint, clientID string) response.ErrorResponse {
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

func (s authService) UpdateRole(userID uint, roleID uint, clientID string) response.ErrorResponse {
	admin, err := s.UserRepository.GetUserByClientID(clientID)
	if err != nil {
		return response.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "user is not an admin",
			Error:   err.Error(),
		}
	}

	user, err := s.UserRepository.GetUserByID(userID)
	if err != nil {
		return response.ErrorResponse{}
	}

	user.RoleID = roleID
	user.UpdatedBy = admin.FullName
	err = s.UserRepository.UpdateRole(user)
	if err != nil {
		return response.ErrorResponse{}
	}
	return response.ErrorResponse{}
}

func (s authService) GetListUser(clientID string) (interface{}, response.ErrorResponse) {
	_, err := s.UserRepository.GetUserByClientID(clientID)
	if err != nil {
		return nil, response.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "user is not an admin",
			Error:   err.Error(),
		}
	}

	users, err := s.UserRepository.GetListUser()
	if err != nil {
		return nil, response.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "user is not an admin",
			Error:   err.Error(),
		}
	}

	var userResponse []out.UserResponse

	for _, user := range *users {
		userResponse = append(userResponse, out.UserResponse{
			UserID:         user.UserID,
			ClientID:       user.ClientID,
			Username:       user.Username,
			FirstName:      user.FirstName,
			LastName:       user.LastName,
			PhoneNumber:    user.PhoneNumber,
			ProfilePicture: user.ProfilePicture,
		})
	}

	return userResponse, response.ErrorResponse{}
}

func (s authService) ChangePassword(password *struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required"`
}, clientID string) response.ErrorResponse {
	user, err := s.UserRepository.GetUserByClientID(clientID)
	if err != nil {
		return response.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "User not found",
			Error:   err.Error(),
		}
	}

	if err := utils.CheckPassword(user.Password, password.OldPassword); err != nil {
		return response.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Old password is incorrect",
			Error:   err.Error(),
		}
	}

	hashedPassword, err := utils.HashPassword(password.NewPassword)
	if err != nil {
		return response.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Password is invalid",
			Error:   err.Error(),
		}
	}

	user.Password = hashedPassword
	user.UpdatedBy = user.FullName
	err = s.UserRepository.ChangePassword(user)
	if err != nil {
		return response.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Password is invalid",
			Error:   err.Error(),
		}
	}
	return response.ErrorResponse{}
}
