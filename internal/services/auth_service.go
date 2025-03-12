package services

import (
	"authentication/internal/dto/in"
	"authentication/internal/dto/out"
	"authentication/internal/models"
	"authentication/internal/repository"
	"authentication/internal/utils"
	"authentication/package/response"
	"github.com/google/uuid"
	"log"
	"net/http"
	"strings"
	"time"
)

type AuthService interface {
	Register(req *in.RegisterRequest) (out.RegisterResponse, response.ErrorResponse)
	Login(req *in.LoginRequest) (interface{}, response.ErrorResponse)
	LoginPhoneNumber(req *in.LoginPhoneNumber) (interface{}, response.ErrorResponse)
	VerifyPinCode(req *struct {
		PinCode string `json:"pin_code" binding:"required"`
	}, clientID string) (interface{}, response.ErrorResponse)
	ChangePinCode(s *struct {
		OldPinCode string `json:"old_pin_code" binding:"required"`
		NewPinCode string `json:"new_pin_code" binding:"required"`
	}, clientID string) response.ErrorResponse
	RegisterInternalToken(req *struct {
		ResourceName string `json:"resource_name" binding:"required"`
	}) (interface{}, response.ErrorResponse)
	UpdateRole(userID uint, roleID uint, clientID string) response.ErrorResponse
	GetListUser(clientID string) (interface{}, response.ErrorResponse)
	ChangePassword(password *struct {
		OldPassword string `json:"old_password" binding:"required"`
		NewPassword string `json:"new_password" binding:"required"`
	}, clientID string) response.ErrorResponse
	ResetPinAttempts()
	ForgetPinCode(req *struct {
		Email   string `json:"email" binding:"required"`
		PinCode string `json:"pin_code" binding:"required"`
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
	Encryption             utils.Encryption
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
	jwtService utils.JWTService,
	Encryption utils.Encryption) AuthService {
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
		Encryption:             Encryption,
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

	hashedPassword, err := s.Encryption.HashPassword(req.Password)
	if err != nil {
		return out.RegisterResponse{}, response.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid Password",
			Error:   err.Error(),
		}
	}

	hashedPin, err := s.Encryption.HashPassword(req.PinCode)
	if err != nil {
		return out.RegisterResponse{}, response.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid Pin Code",
			Error:   err.Error(),
		}
	}

	hashPhoneNumber, err := s.Encryption.Encrypt(req.PhoneNumber)
	if err != nil {
		return out.RegisterResponse{}, response.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid Phone Number",
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

	if err := utils.ValidateEmail(req.Email); err != nil {
		return out.RegisterResponse{}, response.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Email is invalid",
			Error:   err.Error(),
		}
	}

	//check email is exist
	_, err = s.UserRepository.GetUserByEmail(req.Email)
	if err == nil {
		return out.RegisterResponse{}, response.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Email already exist",
			Error:   error(nil).Error(),
		}
	}

	user := &models.Users{
		ClientID:       utils.GenerateClientID(),
		Username:       req.Username,
		Password:       hashedPassword,
		PinCode:        hashedPin,
		PinAttempts:    0,
		PinLastUpdated: time.Now(),
		FirstName:      firstName,
		LastName:       lastName,
		FullName:       fullName,
		Email:          req.Email,
		PhoneNumber:    hashPhoneNumber,
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

	resource, err := s.ResourceRepository.GetResourceByUserID(user.UserID)
	if err != nil {
		return out.RegisterResponse{}, response.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Unable to get resource",
			Error:   err.Error(),
		}
	}

	var resourceName []string
	for _, res := range *resource {
		resourceName = append(resourceName, res.Name)
	}

	user.Role = *role
	token, err := s.JWTService.GenerateToken(*user, resourceName)
	if err != nil {
		return out.RegisterResponse{}, response.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Token is invalid",
			Error:   err.Error(),
		}
	}

	_ = s.RedisService.SaveData(utils.Token, user.ClientID, token)
	_ = s.RedisService.SaveData(utils.User, user.ClientID, user)

	phoneNumber, _ := s.Encryption.Decrypt(user.PhoneNumber)
	responses := out.RegisterResponse{
		UserID:         user.UserID,
		Username:       user.Username,
		FirstName:      user.FirstName,
		LastName:       user.LastName,
		PhoneNumber:    phoneNumber,
		ProfilePicture: user.ProfilePicture,
		Token:          token.AccessToken,
		RefreshToken:   token.RefreshToken,
	}
	return responses, response.ErrorResponse{}
}

func (s authService) Login(req *in.LoginRequest) (interface{}, response.ErrorResponse) {
	user, err := s.UserRepository.GetUserByUsername(req.Username)
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

	resource, err := s.ResourceRepository.GetResourceByUserID(user.UserID)
	if err != nil {
		return out.RegisterResponse{}, response.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Unable to get resource",
			Error:   err.Error(),
		}
	}

	var resourceName []string
	for _, res := range *resource {
		resourceName = append(resourceName, res.Name)
	}

	log.Printf("Resource Name: %v", resourceName)

	token, err := s.JWTService.GenerateToken(*user, resourceName)
	if err != nil {
		return nil, response.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "User or Password is incorrect",
			Error:   err.Error(),
		}
	}

	_ = s.RedisService.SaveData(utils.Token, user.ClientID, token)
	_ = s.RedisService.SaveData(utils.User, user.ClientID, user)

	var phoneNumber string
	decrypt, err := s.Encryption.Decrypt(user.PhoneNumber)
	if err != nil {
		phoneNumber = user.PhoneNumber
	} else {
		phoneNumber = decrypt
	}

	responses := out.LoginResponse{
		UserID:         user.UserID,
		Username:       user.Username,
		FirstName:      user.FirstName,
		LastName:       user.LastName,
		PhoneNumber:    phoneNumber,
		ProfilePicture: user.ProfilePicture,
		Token:          token.AccessToken,
		RefreshToken:   token.RefreshToken,
	}
	return responses, response.ErrorResponse{}
}

func (s authService) LoginPhoneNumber(req *in.LoginPhoneNumber) (interface{}, response.ErrorResponse) {
	hashPhoneNumber, err := s.Encryption.Encrypt(req.PhoneNumber)
	if err != nil {
		return nil, response.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Phone Number is invalid",
			Error:   err.Error(),
		}
	}

	user, err := s.UserRepository.GetUserByPhoneNumber(hashPhoneNumber)
	if err != nil {
		return nil, response.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Phone Number or Pin Code is incorrect",
			Error:   err.Error(),
		}
	}

	err = s.Encryption.CheckPassword(user.PinCode, req.PinCode)
	if err != nil {
		if updateErr := s.UserRepository.UpdatePinAttempts(user.ClientID); updateErr != nil {
			return out.RegisterResponse{}, response.ErrorResponse{
				Code:    http.StatusBadRequest,
				Message: "Invalid User",
				Error:   updateErr.Error(),
			}
		}
		return out.RegisterResponse{}, response.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid Pin Code",
			Error:   err.Error(),
		}
	}

	role, err := s.RoleRepository.GetRoleByID(user.RoleID)
	user.Role = *role

	resource, err := s.ResourceRepository.GetResourceByUserID(user.UserID)
	if err != nil {
		return out.RegisterResponse{}, response.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Unable to get resource",
			Error:   err.Error(),
		}
	}

	var resourceName []string
	for _, res := range *resource {
		resourceName = append(resourceName, res.Name)
	}

	token, err := s.JWTService.GenerateToken(*user, resourceName)
	if err != nil {
		return nil, response.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "User or Password is incorrect",
			Error:   err.Error(),
		}
	}

	_ = s.RedisService.SaveData(utils.Token, user.ClientID, token)
	_ = s.RedisService.SaveData(utils.User, user.ClientID, user)

	var phoneNumber string
	decrypt, err := s.Encryption.Decrypt(user.PhoneNumber)
	if err != nil {
		phoneNumber = user.PhoneNumber
	} else {
		phoneNumber = decrypt
	}

	responses := out.LoginResponse{
		UserID:         user.UserID,
		Username:       user.Username,
		FirstName:      user.FirstName,
		LastName:       user.LastName,
		PhoneNumber:    phoneNumber,
		ProfilePicture: user.ProfilePicture,
		Token:          token.AccessToken,
		RefreshToken:   token.RefreshToken,
	}
	return responses, response.ErrorResponse{}
}

func (s authService) VerifyPinCode(req *struct {
	PinCode string `json:"pin_code" binding:"required"`
}, clientID string) (interface{}, response.ErrorResponse) {
	data, err := utils.GetUserRedis(s.RedisService, utils.User, clientID)
	if err != nil {
		return nil, response.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "User not found",
			Error:   err.Error(),
		}
	}

	user, err := s.UserRepository.GetUserByClientID(data.ClientID)
	if err != nil {
		return nil, response.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "User not found",
			Error:   err.Error(),
		}
	}

	err = s.Encryption.CheckPassword(user.PinCode, req.PinCode)
	if err != nil {
		if updateErr := s.UserRepository.UpdatePinAttempts(data.ClientID); updateErr != nil {
			return out.RegisterResponse{}, response.ErrorResponse{
				Code:    http.StatusBadRequest,
				Message: "Invalid User",
				Error:   updateErr.Error(),
			}
		}
		return out.RegisterResponse{}, response.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid Pin Code",
			Error:   err.Error(),
		}
	}

	var requestID = uuid.New().String()
	responseModel := out.VerifyPinCodeResponse{
		ClientID:  user.ClientID,
		RequestID: requestID,
		Valid:     true,
	}

	err = s.RedisService.SaveDataExpired(utils.PinVerify, requestID, 5, responseModel)
	if err != nil {
		return out.RegisterResponse{}, response.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid Pin Code",
			Error:   err.Error(),
		}
	}

	return responseModel, response.ErrorResponse{}
}

func (s authService) ChangePinCode(req *struct {
	OldPinCode string `json:"old_pin_code" binding:"required"`
	NewPinCode string `json:"new_pin_code" binding:"required"`
}, clientID string) response.ErrorResponse {
	data, err := utils.GetUserRedis(s.RedisService, utils.User, clientID)
	if err != nil {
		return response.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "User not found",
			Error:   err.Error(),
		}
	}

	user, err := s.UserRepository.GetUserByClientID(data.ClientID)
	if err != nil {
		return response.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "User not found",
			Error:   err.Error(),
		}
	}

	err = s.Encryption.CheckPassword(user.PinCode, req.OldPinCode)
	if err != nil {
		return response.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid Pin Code",
			Error:   err.Error(),
		}
	}

	hashedNewPin, err := s.Encryption.HashPassword(req.NewPinCode)
	if err != nil {
		return response.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid Pin Code",
			Error:   err.Error(),
		}
	}

	if hashedNewPin == user.PinCode {
		return response.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid Pin Code",
			Error:   "Old Pin and New Pin is same",
		}
	}

	user.PinCode = hashedNewPin
	user.PinLastUpdated = time.Now()
	user.PinAttempts = 0
	user.UpdatedBy = user.ClientID

	err = s.AuthRepository.UpdatePinCode(user)
	if err != nil {
		return response.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid Pin Code",
			Error:   err.Error(),
		}
	}
	return response.ErrorResponse{}
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
		var phoneNumber string
		decrypt, err := s.Encryption.Decrypt(user.PhoneNumber)
		if err != nil {
			phoneNumber = user.PhoneNumber
		} else {
			phoneNumber = decrypt
		}
		userResponse = append(userResponse, out.UserResponse{
			UserID:         user.UserID,
			ClientID:       user.ClientID,
			Username:       user.Username,
			FirstName:      user.FirstName,
			LastName:       user.LastName,
			PhoneNumber:    phoneNumber,
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

func (s authService) ResetPinAttempts() {
	users, _ := s.UserRepository.GetListUser()
	for _, user := range *users {
		go func(u models.Users) {
			if u.PinAttempts > 0 {
				u.PinAttempts = 0
				_ = s.UserRepository.ResetPinAttempts(&u)
			}
		}(user)
	}
}

func (s authService) ForgetPinCode(req *struct {
	Email   string `json:"email" binding:"required"`
	PinCode string `json:"pin_code" binding:"required"`
}, clientID string) response.ErrorResponse {
	user, err := s.UserRepository.GetUserByClientID(clientID)
	if err != nil {
		return response.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "User not found",
			Error:   err.Error(),
		}
	}

	if err := utils.ValidateEmail(req.Email); err != nil {
		return response.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Email is invalid",
			Error:   err.Error(),
		}
	}

	user, err = s.UserRepository.GetUserByEmail(req.Email)
	if err != nil {
		return response.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Email not found",
			Error:   err.Error(),
		}
	}

	//generate send email

	hashedPin, err := s.Encryption.HashPassword(req.PinCode)

	log.Printf("Pin Code: %s", req.PinCode)
	log.Printf("Pin Code: %s", hashedPin)
	if err != nil {
		return response.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid Pin Code",
			Error:   err.Error(),
		}
	}

	user.PinCode = hashedPin
	user.PinLastUpdated = time.Now()
	user.PinAttempts = 0
	user.UpdatedBy = user.ClientID
	err = s.UserRepository.UpdateProfile(user)
	if err != nil {
		return response.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Failed to update pin code",
			Error:   err.Error(),
		}
	}

	_ = s.RedisService.SaveData(utils.User, user.ClientID, user)

	return response.ErrorResponse{}
}
