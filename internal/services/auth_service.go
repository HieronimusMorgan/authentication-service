package services

import (
	"authentication/internal/dto/in"
	"authentication/internal/dto/out"
	"authentication/internal/models"
	"authentication/internal/repository"
	"authentication/internal/utils"
	nt "authentication/internal/utils/nats"
	"errors"
	"github.com/google/uuid"
	"log"
	"regexp"
	"strings"
	"time"
)

type AuthService interface {
	Register(req *in.RegisterRequest, deviceID string) (out.RegisterResponse, error)
	RegisterDeviceToken(req *struct {
		DeviceToken string `json:"device_token" binding:"required"`
	}, clientID string) error
	Login(req *in.LoginRequest, deviceID string) (interface{}, error)
	ReLogin(req struct {
		UserID       uint   `json:"user_id" binding:"required"`
		RefreshToken string `json:"refresh_token" binding:"required"`
	}) (interface{}, error)
	LoginPhoneNumber(req *in.LoginPhoneNumber, deviceID string) (interface{}, error)
	ChangeDeviceID(s *struct {
		PhoneNumber string `json:"phone_number" binding:"required"`
		DeviceID    string `json:"device_id" binding:"required"`
	}) (interface{}, error)
	VerifyDeviceID(req *struct {
		RequestID string `json:"request_id" binding:"required"`
		PinCode   string `json:"pin_code" binding:"required"`
	}) (interface{}, error)
	VerifyPinCode(req *struct {
		PinCode string `json:"pin_code" binding:"required"`
	}, clientID string) (interface{}, error)
	ChangePinCode(s *struct {
		OldPinCode string `json:"old_pin_code" binding:"required"`
		NewPinCode string `json:"new_pin_code" binding:"required"`
	}, clientID string) error
	UpdateToken(userID uint, clientID string) (*models.TokenDetails, error)
	RefreshToken(req *struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}, id string) (interface{}, error)
	RegisterInternalToken(req *struct {
		ResourceName string `json:"resource_name" binding:"required"`
	}) (interface{}, error)
	UpdateRole(userID uint, roleID uint, clientID string) error
	GetListUser(clientID string) (interface{}, error)
	ChangePassword(password *struct {
		OldPassword string `json:"old_password" binding:"required"`
		NewPassword string `json:"new_password" binding:"required"`
	}, clientID string) error
	ResetPinAttempts()
	ForgetPinCode(req *struct {
		Email   string `json:"email" binding:"required"`
		PinCode string `json:"pin_code" binding:"required"`
	}, clientID string) error
	GetUserByID(userID uint, clientID string) (interface{}, error)
	GenerateCredentialKey(clientID string) (interface{}, error)
}

type authService struct {
	AuthRepository            repository.AuthRepository
	ResourceRepository        repository.ResourceRepository
	RoleRepository            repository.RoleRepository
	UserResourceRepository    repository.UserResourceRepository
	UserRepository            repository.UserRepository
	UserKeyRepository         repository.UserKeyRepository
	UserRoleRepository        repository.UserRoleRepository
	UserSessionRepository     repository.UserSessionRepository
	UserTransactionRepository repository.UserTransactionalRepository
	UserSettingRepository     repository.UserSettingRepository
	RedisService              utils.RedisService
	JWTService                utils.JWTService
	Encryption                utils.Encryption
	NatsService               nt.Service
}

func NewAuthService(authRepo repository.AuthRepository, resourceRepo repository.ResourceRepository, roleRepo repository.RoleRepository, roleResourceRepo repository.UserResourceRepository, userRepo repository.UserRepository, userKeyRepo repository.UserKeyRepository, userRoleRepo repository.UserRoleRepository, userSessionRepo repository.UserSessionRepository, userTransactionRepo repository.UserTransactionalRepository, userSetting repository.UserSettingRepository, redis utils.RedisService, jwtService utils.JWTService, Encryption utils.Encryption, service nt.Service) AuthService {
	return authService{
		AuthRepository:            authRepo,
		ResourceRepository:        resourceRepo,
		RoleRepository:            roleRepo,
		UserResourceRepository:    roleResourceRepo,
		UserRepository:            userRepo,
		UserKeyRepository:         userKeyRepo,
		UserRoleRepository:        userRoleRepo,
		UserSessionRepository:     userSessionRepo,
		UserTransactionRepository: userTransactionRepo,
		UserSettingRepository:     userSetting,
		RedisService:              redis,
		JWTService:                jwtService,
		Encryption:                Encryption,
		NatsService:               service,
	}
}

func (s authService) checkUserIsAdmin(user *models.Users) error {
	role, err := s.RoleRepository.GetRoleByID(user.RoleID)
	if err != nil {
		return errors.New("failed to get role")
	}
	if strings.EqualFold(role.Name, "Admin") || strings.EqualFold(role.Name, "Super Admin") {
		return errors.New("user is not an admin")
	}
	return nil
}

func (s authService) Register(req *in.RegisterRequest, deviceID string) (out.RegisterResponse, error) {
	//if err := utils.ValidateUsername(req.Username); err != nil {
	//	return interface{}, error{
	//		Code:    http.StatusBadRequest,
	//		Message: "Validation Username",
	//		Error:   err.Error(),
	//	}
	//}

	if err := utils.ValidatePhoneNumber(req.PhoneNumber); err != nil {
		return out.RegisterResponse{}, errors.New("phone Number is invalid")
	}

	hashedPassword, err := s.Encryption.HashPassword(req.Password)
	if err != nil {
		return out.RegisterResponse{}, errors.New("invalid Password")
	}

	var hashedPin *string
	if req.PinCode != nil && *req.PinCode != "" {
		if len(*req.PinCode) < 6 {
			return out.RegisterResponse{}, errors.New("pin Code must be 6 digits")
		}
		if len(*req.PinCode) > 6 {
			return out.RegisterResponse{}, errors.New("pin Code must be 6 digits")
		}
		if !regexp.MustCompile(`^\d+$`).MatchString(*req.PinCode) {
			return out.RegisterResponse{}, errors.New("pin Code must be numeric")
		}
		hashedPin, err = s.Encryption.HashPassword(*req.PinCode)
	}

	var hashDeviceID string
	if deviceID == "MOBILE" {
		if req.DeviceID == nil {
			return out.RegisterResponse{}, errors.New("device ID is required")
		}
		hashDeviceID, err = s.Encryption.Encrypt(*req.DeviceID)
		if err != nil {
			return out.RegisterResponse{}, errors.New("invalid Device ID")
		}
	}

	hashPhoneNumber, err := s.Encryption.Encrypt(req.PhoneNumber)
	if err != nil {
		return out.RegisterResponse{}, errors.New("phone Number is invalid")
	}

	role, err := s.RoleRepository.GetRoleByName("User")
	if err != nil {
		return out.RegisterResponse{}, errors.New("unable to get role")
	}

	firstName := utils.ValidationTrimSpace(req.FirstName)
	lastName := utils.ValidationTrimSpace(req.LastName)
	fullName := firstName + " " + lastName

	if err := utils.ValidateEmail(req.Email); err != nil {
		return out.RegisterResponse{}, errors.New("email is invalid")
	}

	//check email exists
	_, err = s.UserRepository.GetUserByEmail(req.Email)
	if err == nil {
		return out.RegisterResponse{}, errors.New("email already exist")
	}

	phone, err := s.UserRepository.GetUserByPhoneNumber(hashPhoneNumber)
	if phone != nil && phone.PhoneNumber != "" {
		return out.RegisterResponse{}, errors.New("phone Number already exist")
	}

	user := &models.Users{
		ClientID:       utils.GenerateClientID(),
		Username:       req.Username,
		Password:       *hashedPassword,
		PinCode:        hashedPin,
		PinAttempts:    0,
		PinLastUpdated: time.Now(),
		FirstName:      firstName,
		LastName:       lastName,
		FullName:       fullName,
		Email:          req.Email,
		PhoneNumber:    hashPhoneNumber,
		RoleID:         role.RoleID,
		DeviceID: func() *string {
			if hashDeviceID == "" {
				return nil
			} else {
				return &hashDeviceID
			}
		}(),
		CreatedBy: "system",
		UpdatedBy: "system",
	}

	if err := s.UserTransactionRepository.RegistrationUser(user); err != nil {
		return out.RegisterResponse{}, errors.New("Unable to register user")
	}

	resource, err := s.ResourceRepository.GetResourceByUserID(user.UserID)
	if err != nil {
		return out.RegisterResponse{}, errors.New("Unable to get resource")
	}

	role, err = s.RoleRepository.GetRoleByID(user.RoleID)
	if err != nil {
		return out.RegisterResponse{}, errors.New("Unable to get role")
	}

	var resourceName []string
	for _, res := range *resource {
		resourceName = append(resourceName, res.Name)
	}

	userSetting, err := s.UserSettingRepository.GetUserSettingByUserID(user.UserID)
	if err != nil {
		return out.RegisterResponse{}, errors.New("User setting not found")
	}

	userSettingModel := out.UserSettingResponse{
		SettingID:             userSetting.SettingID,
		GroupInviteType:       userSetting.GroupInviteType,
		GroupInviteDisallowed: userSetting.GroupInviteDisallowed,
	}

	userKeys, err := utils.GenerateUserKey(user)
	if err != nil {
		return out.RegisterResponse{}, errors.New("Unable to generate user key")
	}
	// Save user key to database
	if err := s.UserRepository.SaveUserKey(userKeys); err != nil {
		return out.RegisterResponse{}, errors.New("Unable to save user key")
	}

	token, err := s.JWTService.GenerateToken(*user, resourceName, role.Name)
	if err != nil {
		return out.RegisterResponse{}, errors.New("User or Password is incorrect")
	}

	userRedis, err := s.UserRepository.GetUserRedisByClientID(user.ClientID)
	if err != nil {
		return out.RegisterResponse{}, errors.New("User not found")
	}

	if userRedis == nil {
		return out.RegisterResponse{}, errors.New("User not found")
	}
	_ = s.RedisService.SaveData(utils.Token, user.ClientID, token)
	_ = s.RedisService.SaveData(utils.User, user.ClientID, userRedis)
	_ = s.RedisService.SaveData(utils.UserKey, user.ClientID, userKeys)

	phoneNumber, _ := s.Encryption.Decrypt(user.PhoneNumber)
	responses := out.RegisterResponse{
		UserID:         user.UserID,
		Username:       user.Username,
		FirstName:      user.FirstName,
		LastName:       user.LastName,
		PhoneNumber:    phoneNumber,
		Role:           role.Name,
		Resource:       resourceName,
		UserSetting:    userSettingModel,
		ProfilePicture: user.ProfilePicture,
		Token:          token.AccessToken,
		RefreshToken:   token.RefreshToken,
	}
	return responses, nil
}

func (s authService) RegisterDeviceToken(req *struct {
	DeviceToken string `json:"device_token" binding:"required"`
}, clientID string) error {
	user, err := s.UserRepository.GetUserByClientID(clientID)
	if err != nil {
		return errors.New("user not found")
	}

	if user.DeviceToken == nil {
		user.DeviceToken = &req.DeviceToken
	} else {
		user.DeviceToken = &req.DeviceToken
	}

	if err := s.UserRepository.UpdateUser(user); err != nil {
		return errors.New("unable to update user")
	}
	return nil
}

func (s authService) Login(req *in.LoginRequest, deviceID string) (interface{}, error) {
	user, err := s.UserRepository.GetUserByUsername(req.Username)
	if err != nil {
		return nil, errors.New("username or Password is incorrect")
	}
	if err := utils.CheckPassword(user.Password, req.Password); err != nil {
		return nil, errors.New("username or Password is incorrect")
	}

	if deviceID == "MOBILE" && req.DeviceID != "" && (user.DeviceID == nil || *user.DeviceID != req.DeviceID) {
		hashDeviceID, err := s.Encryption.Encrypt(req.DeviceID)
		if err != nil {
			return nil, errors.New("device ID is invalid")
		}
		if err := s.UserRepository.UpdateDeviceID(user.UserID, hashDeviceID); err != nil {
			return nil, errors.New("unable to update device ID")
		}
	}

	role, err := s.RoleRepository.GetRoleByID(user.RoleID)
	if err != nil {
		return nil, errors.New("unable to get role")
	}

	resource, err := s.ResourceRepository.GetResourceByUserID(user.UserID)
	if err != nil {
		return nil, errors.New("unable to get resource")
	}

	var resourceName []string
	for _, res := range *resource {
		resourceName = append(resourceName, res.Name)
	}

	userSetting, err := s.UserSettingRepository.GetUserSettingByUserID(user.UserID)
	if err != nil {
		newUserSetting := models.UserSetting{
			UserID:          user.UserID,
			GroupInviteType: 1,
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
		}
		if addErr := s.UserSettingRepository.AddUserSetting(&newUserSetting); addErr != nil {
			return nil, errors.New("unable to create user setting")
		}
		userSetting = &newUserSetting
	}

	userSettingModel := out.UserSettingResponse{
		SettingID:             userSetting.SettingID,
		GroupInviteType:       userSetting.GroupInviteType,
		GroupInviteDisallowed: userSetting.GroupInviteDisallowed,
	}

	token, err := s.JWTService.GenerateToken(*user, resourceName, role.Name)
	if err != nil {
		return nil, errors.New("user or Password is incorrect")
	}

	userRedis, err := s.UserRepository.GetUserRedisByClientID(user.ClientID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	if userRedis == nil {
		return nil, errors.New("user not found")
	}

	_ = s.RedisService.SaveData(utils.Token, user.ClientID, token)
	_ = s.RedisService.SaveData(utils.User, user.ClientID, userRedis)

	var phoneNumber string
	decrypt, err := s.Encryption.Decrypt(user.PhoneNumber)
	if err != nil {
		phoneNumber = user.PhoneNumber
	} else {
		phoneNumber = decrypt
	}
	var deviceIdResponse string
	decrypt, err = s.Encryption.Decrypt(*user.DeviceID)
	if err != nil {
		deviceIdResponse = user.PhoneNumber
	} else {
		deviceIdResponse = decrypt
	}

	userKey, err := s.UserKeyRepository.GetUserKeyByUserID(user.UserID)
	if userKey == nil && err != nil {
		userKeys, err := utils.GenerateUserKey(user)
		if err != nil {
			return out.RegisterResponse{}, errors.New("unable to generate user key")
		}

		if err := s.UserRepository.SaveUserKey(userKeys); err != nil {
			return out.RegisterResponse{}, errors.New("unable to save user key")
		}
	}

	responses := out.LoginResponse{
		UserID:         user.UserID,
		ClientID:       user.ClientID,
		Username:       user.Username,
		FirstName:      user.FirstName,
		LastName:       user.LastName,
		PhoneNumber:    phoneNumber,
		Email:          user.Email,
		DeviceID:       &deviceIdResponse,
		DeviceToken:    user.DeviceToken,
		ProfilePicture: user.ProfilePicture,
		UserSetting:    userSettingModel,
		Token:          token.AccessToken,
		RefreshToken:   token.RefreshToken,
	}
	return responses, nil
}

func (s authService) ReLogin(req struct {
	UserID       uint   `json:"user_id" binding:"required"`
	RefreshToken string `json:"refresh_token" binding:"required"`
}) (interface{}, error) {
	userSession, err := s.UserSessionRepository.GetUserSessionByRefreshTokenAndUserID(req.UserID, req.RefreshToken)
	if err != nil {
		return nil, errors.New("invalid Refresh Token")
	}

	user, err := s.UserRepository.GetUserByID(userSession.UserID)
	if err != nil {
		return nil, errors.New("invalid User")
	}

	return s.Login(&in.LoginRequest{
		Username: user.Username,
	}, *user.DeviceID)
}

func (s authService) LoginPhoneNumber(req *in.LoginPhoneNumber, deviceID string) (interface{}, error) {
	hashPhoneNumber, err := s.Encryption.Encrypt(req.PhoneNumber)
	if err != nil {
		return nil, errors.New("phone Number is invalid")
	}

	user, err := s.UserRepository.GetUserByPhoneNumber(hashPhoneNumber)
	if err != nil {
		return nil, errors.New("phone Number is invalid")
	}
	if user.PinCode == nil {
		return nil, errors.New("pin Code is not set")
	}

	err = s.Encryption.CheckPassword(*user.PinCode, req.PinCode)
	if err != nil {
		if updateErr := s.UserRepository.UpdatePinAttempts(user.ClientID); updateErr != nil {
			return nil, errors.New("invalid User")
		}
		return nil, errors.New("invalid Pin Code")
	}

	if deviceID == "MOBILE" && req.DeviceID != "" && (user.DeviceID == nil || *user.DeviceID != req.DeviceID) {
		hashDeviceID, err := s.Encryption.Encrypt(req.DeviceID)
		if err != nil {
			return nil, errors.New("device ID is invalid")
		}
		if err := s.UserRepository.UpdateDeviceID(user.UserID, hashDeviceID); err != nil {
			return nil, errors.New("unable to update device ID")
		}
	}

	resource, err := s.ResourceRepository.GetResourceByUserID(user.UserID)
	if err != nil {
		return nil, errors.New("Unable to get resource")
	}

	var resourceName []string
	for _, res := range *resource {
		resourceName = append(resourceName, res.Name)
	}

	role, err := s.RoleRepository.GetRoleByID(user.RoleID)
	if err != nil {
		return nil, errors.New("Unable to get role")
	}

	token, err := s.JWTService.GenerateToken(*user, resourceName, role.Name)
	if err != nil {
		return nil, errors.New("User or Password is incorrect")
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

	decryptDeviceID := utils.DecryptOptionalString(user.DeviceID, s.Encryption)
	decryptDeviceToken := utils.DecryptOptionalString(user.DeviceToken, s.Encryption)

	responses := out.LoginResponse{
		UserID:         user.UserID,
		Username:       user.Username,
		FirstName:      user.FirstName,
		LastName:       user.LastName,
		PhoneNumber:    phoneNumber,
		Email:          user.Email,
		DeviceID:       decryptDeviceID,
		DeviceToken:    decryptDeviceToken,
		ProfilePicture: user.ProfilePicture,
		Token:          token.AccessToken,
		RefreshToken:   token.RefreshToken,
	}
	return responses, nil
}

func (s authService) ChangeDeviceID(req *struct {
	PhoneNumber string `json:"phone_number" binding:"required"`
	DeviceID    string `json:"device_id" binding:"required"`
}) (interface{}, error) {
	hashPhoneNumber, err := s.Encryption.Encrypt(req.PhoneNumber)
	if err != nil {
		return nil, errors.New("phone Number is invalid")
	}

	user, err := s.UserRepository.GetUserByPhoneNumber(hashPhoneNumber)
	if err != nil {
		return nil, errors.New("phone Number is invalid")
	}
	hashDeviceID, err := s.Encryption.Encrypt(req.DeviceID)
	if err != nil {
		return nil, errors.New("Device ID is invalid")
	}
	user.DeviceID = &hashDeviceID

	requestID := uuid.New().String()

	_ = s.RedisService.SaveDataExpired(utils.DeviceVerify, requestID, 1, user)

	return map[string]interface{}{
		"request_id": requestID,
		"timestamp":  time.Now().Unix(),
	}, nil
}

func (s authService) VerifyDeviceID(req *struct {
	RequestID string `json:"request_id" binding:"required"`
	PinCode   string `json:"pin_code" binding:"required"`
}) (interface{}, error) {
	data, err := utils.GetUserRedis(s.RedisService, utils.DeviceVerify, req.RequestID)
	if err != nil {
		return nil, errors.New("User not found")
	}

	_ = s.RedisService.DeleteData(utils.DeviceVerify, req.RequestID)

	user, err := s.UserRepository.GetUserByClientID(data.ClientID)
	if err != nil {
		return nil, errors.New("User not found")
	}

	if user.PinCode == nil {
		return nil, errors.New("Pin Code is not set")
	}

	err = s.Encryption.CheckPassword(*user.PinCode, req.PinCode)
	if err != nil {
		if updateErr := s.UserRepository.UpdatePinAttempts(data.ClientID); updateErr != nil {
			return nil, errors.New("Invalid User")
		}
		return nil, errors.New("Invalid Pin Code")
	}

	user.DeviceID = data.DeviceID
	user.UpdatedBy = user.ClientID

	if err = s.UserRepository.UpdateUser(user); err != nil {
		return nil, errors.New("Unable to update user")
	}

	deviceID, err := s.Encryption.Decrypt(*user.DeviceID)
	if err != nil {
		deviceID = *user.DeviceID
	}

	phoneNumber, err := s.Encryption.Decrypt(user.PhoneNumber)
	if err != nil {
		phoneNumber = user.PhoneNumber
	}

	user.DeviceID = &deviceID
	user.PhoneNumber = phoneNumber

	_ = s.RedisService.SaveData(utils.User, user.ClientID, user)

	return user, nil
}

func (s authService) VerifyPinCode(req *struct {
	PinCode string `json:"pin_code" binding:"required"`
}, clientID string) (interface{}, error) {
	data, err := utils.GetUserRedis(s.RedisService, utils.User, clientID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	user, err := s.UserRepository.GetUserByClientID(data.ClientID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	if user.PinCode == nil {
		return nil, errors.New("pin Code is not set")
	}

	err = s.Encryption.CheckPassword(*user.PinCode, req.PinCode)
	if err != nil {
		if updateErr := s.UserRepository.UpdatePinAttempts(data.ClientID); updateErr != nil {
			return nil, errors.New("invalid User")
		}
		return nil, errors.New("invalid Pin Code")
	}

	var requestID = uuid.New().String()
	responseModel := out.VerifyPinCodeResponse{
		ClientID:  user.ClientID,
		RequestID: requestID,
		Valid:     true,
	}

	err = s.RedisService.SaveDataExpired(utils.PinVerify, user.ClientID, 5, responseModel)
	if err != nil {
		return nil, errors.New("unable to save data")
	}

	return responseModel, nil
}

func (s authService) ChangePinCode(req *struct {
	OldPinCode string `json:"old_pin_code" binding:"required"`
	NewPinCode string `json:"new_pin_code" binding:"required"`
}, clientID string) error {
	data, err := utils.GetUserRedis(s.RedisService, utils.User, clientID)
	if err != nil {
		return errors.New("user not found")
	}

	user, err := s.UserRepository.GetUserByClientID(data.ClientID)
	if err != nil {
		return errors.New("user not found")
	}

	if user.PinCode == nil {
		return errors.New("pin Code is not set")
	}

	err = s.Encryption.CheckPassword(*user.PinCode, req.OldPinCode)
	if err != nil {
		return errors.New("old Pin Code is incorrect")
	}

	hashedNewPin, err := s.Encryption.HashPassword(req.NewPinCode)
	if err != nil {
		return errors.New("invalid Pin Code")
	}

	if hashedNewPin == user.PinCode {
		return errors.New("old Pin and New Pin is same")
	}

	user.PinCode = hashedNewPin
	user.PinLastUpdated = time.Now()
	user.PinAttempts = 0
	user.UpdatedBy = user.ClientID

	err = s.AuthRepository.UpdatePinCode(user)
	if err != nil {
		return errors.New("Unable to update pin code")
	}
	return nil
}

func (s authService) UpdateToken(userID uint, clientID string) (*models.TokenDetails, error) {
	data, err := utils.GetUserRedis(s.RedisService, utils.User, clientID)
	if err != nil {
		return nil, errors.New("user not found")
	}
	admin, err := s.UserRepository.GetUserByClientID(data.ClientID)

	user, err := s.UserRepository.GetUserByID(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	resource, err := s.ResourceRepository.GetResourceByUserID(user.UserID)
	if err != nil {
		return nil, errors.New("unable to get resource")
	}

	var resourceName []string
	for _, res := range *resource {
		resourceName = append(resourceName, res.Name)
	}

	role, err := s.RoleRepository.GetRoleByID(user.RoleID)
	if err != nil {
		return nil, errors.New("unable to get role")
	}

	token, err := s.JWTService.GenerateToken(*user, resourceName, role.Name)
	if err != nil {
		return nil, errors.New("user or Password is incorrect")
	}

	userRedis, err := s.UserRepository.GetUserRedisByClientID(user.ClientID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	if userRedis == nil {
		return nil, errors.New("user not found")
	}
	_ = s.RedisService.SaveData(utils.Token, user.ClientID, token)
	_ = s.RedisService.SaveData(utils.User, user.ClientID, userRedis)

	userSession, err := s.UserSessionRepository.GetUserSessionByUserID(user.UserID)
	if err != nil {
		return nil, errors.New("user session not found")
	}

	if userSession == nil {
		userSession = &models.UserSession{
			UserID:       user.UserID,
			SessionToken: token.AccessToken,
			RefreshToken: token.RefreshToken,
			ExpiresAt:    time.Unix(token.AtExpires, 0),
			LoginTime:    time.Now(),
			CreatedBy:    admin.ClientID,
			UpdatedBy:    admin.ClientID,
		}
		err = s.UserSessionRepository.AddUserSession(userSession)
		if err != nil {
			return nil, errors.New("unable to add session")
		}
	} else {
		userSession.SessionToken = token.AccessToken
		userSession.RefreshToken = token.RefreshToken
		userSession.ExpiresAt = time.Unix(token.AtExpires, 0)
		userSession.LoginTime = time.Now()
		userSession.UpdatedBy = admin.ClientID

		err = s.UserSessionRepository.UpdateSession(userSession)
		if err != nil {
			return nil, errors.New("unable to update session")
		}
	}

	return &token, nil
}

func (s authService) RefreshToken(req *struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}, id string) (interface{}, error) {
	user, err := s.UserRepository.GetUserByClientID(id)
	if err != nil {
		return nil, errors.New("user not found")
	}

	userSession, err := s.UserSessionRepository.GetUserSessionByRefreshTokenAndUserID(user.UserID, req.RefreshToken)
	if err != nil {
		return nil, errors.New("invalid Refresh Token")
	}
	if userSession == nil {
		return nil, errors.New("invalid Refresh Token")
	}

	if userSession.ExpiresAt.Before(time.Now()) {
		return nil, errors.New("refresh Token is expired")
	}

	_ = s.RedisService.DeleteData(utils.Token, user.ClientID)
	_ = s.RedisService.DeleteData(utils.User, user.ClientID)

	resource, err := s.ResourceRepository.GetResourceByUserID(user.UserID)
	if err != nil {
		return nil, errors.New("unable to get resource")
	}

	var resourceName []string
	for _, res := range *resource {
		resourceName = append(resourceName, res.Name)
	}

	role, err := s.RoleRepository.GetRoleByID(user.RoleID)
	if err != nil {
		return nil, errors.New("unable to get role")
	}

	token, err := s.JWTService.GenerateToken(*user, resourceName, role.Name)
	if err != nil {
		return nil, errors.New("user or Password is incorrect")
	}

	_ = s.RedisService.SaveData(utils.Token, user.ClientID, token)
	_ = s.RedisService.SaveData(utils.User, user.ClientID, user)

	userSession.SessionToken = token.AccessToken
	userSession.RefreshToken = token.RefreshToken
	userSession.ExpiresAt = time.Unix(token.AtExpires, 0)
	userSession.LoginTime = time.Unix(token.AtExpires, 0)
	userSession.UpdatedAt = time.Now()
	userSession.UpdatedBy = user.ClientID

	err = s.UserSessionRepository.UpdateSession(userSession)
	if err != nil {
		return nil, errors.New("unable to update session")
	}

	return token, nil
}

func (s authService) RegisterInternalToken(req *struct {
	ResourceName string `json:"resource_name" binding:"required"`
}) (interface{}, error) {
	resource, err := s.ResourceRepository.GetResourceByName(req.ResourceName)
	if err != nil {
		return nil, errors.New("Resource not found")
	}

	token, err := s.JWTService.GenerateInternalToken(resource.Name)
	if err != nil {
		return nil, errors.New("Unable to generate token")
	}

	if err := s.AuthRepository.CreateInternalToken(resource.ResourceID, token); err != nil {
		return nil, errors.New("Unable to create token")
	}

	return token, nil
}

func (s authService) UpdateRole(userID uint, roleID uint, clientID string) error {
	admin, err := s.UserRepository.GetUserByClientID(clientID)
	if err != nil {
		return errors.New("user is not an admin")
	}

	user, err := s.UserRepository.GetUserByID(userID)
	if err != nil {
		return nil
	}

	user.RoleID = roleID
	user.UpdatedBy = admin.FullName
	err = s.UserRepository.UpdateRole(user)
	if err != nil {
		return nil
	}
	return nil
}

func (s authService) GetListUser(clientID string) (interface{}, error) {
	_, err := s.UserRepository.GetUserByClientID(clientID)
	if err != nil {
		return nil, errors.New("user is not an admin")
	}

	users, err := s.UserRepository.GetListUserResponse()
	if err != nil {
		return nil, errors.New("user is not an admin")
	}

	var userResponse []out.UserRoleResourceSettingResponse

	for _, user := range *users {
		var phoneNumber string
		decrypt, err := s.Encryption.Decrypt(user.PhoneNumber)
		if err != nil {
			phoneNumber = user.PhoneNumber
		} else {
			phoneNumber = decrypt
		}
		userResponse = append(userResponse, out.UserRoleResourceSettingResponse{
			UserID:         user.UserID,
			ClientID:       user.ClientID,
			Username:       user.Username,
			FirstName:      user.FirstName,
			LastName:       user.LastName,
			PhoneNumber:    phoneNumber,
			ProfilePicture: user.ProfilePicture,
			Role:           user.Role,
			Resource:       user.Resource,
			UserSetting:    user.UserSetting,
		})
	}

	return userResponse, nil
}

func (s authService) ChangePassword(password *struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required"`
}, clientID string) error {
	user, err := s.UserRepository.GetUserByClientID(clientID)
	if err != nil {
		return errors.New("user not found")
	}

	if err := utils.CheckPassword(user.Password, password.OldPassword); err != nil {
		return errors.New("Old Password is incorrect")
	}

	hashedPassword, err := utils.HashPassword(password.NewPassword)
	if err != nil {
		return errors.New("Invalid Password")
	}

	user.Password = hashedPassword
	user.UpdatedBy = user.FullName
	err = s.UserRepository.ChangePassword(user)
	if err != nil {
		return errors.New("Unable to change password")
	}
	return nil
}

func (s authService) ResetPinAttempts() {
	listUsers, _ := s.UserRepository.GetListUser()
	for _, user := range *listUsers {
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
}, clientID string) error {
	user, err := s.UserRepository.GetUserByClientID(clientID)
	if err != nil {
		return errors.New("user not found")
	}

	if err := utils.ValidateEmail(req.Email); err != nil {
		return errors.New("Email is invalid")
	}

	user, err = s.UserRepository.GetUserByEmail(req.Email)
	if err != nil {
		return errors.New("Email not found")
	}

	hashedPin, err := s.Encryption.HashPassword(req.PinCode)

	log.Printf("Pin Code: %s", req.PinCode)
	log.Printf("Pin Code: %s", hashedPin)
	if err != nil {
		return errors.New("Invalid Pin Code")
	}

	user.PinCode = hashedPin
	user.PinLastUpdated = time.Now()
	user.PinAttempts = 0
	user.UpdatedBy = user.ClientID
	err = s.UserRepository.UpdateUser(user)
	if err != nil {
		return errors.New("Unable to update pin code")
	}

	_ = s.RedisService.SaveData(utils.User, user.ClientID, user)

	return nil
}

func (s authService) GetUserByID(userID uint, clientID string) (interface{}, error) {
	_, err := s.UserRepository.GetUserByClientID(clientID)
	if err != nil {
		return nil, errors.New("user is not an admin")
	}

	users, err := s.UserRepository.GetListUserByUserIDResponse(userID)
	if err != nil {
		return nil, errors.New("user is not an admin")
	}

	var userResponse []out.UserRoleResourceSettingResponse

	for _, user := range *users {
		var phoneNumber string
		decrypt, err := s.Encryption.Decrypt(user.PhoneNumber)
		if err != nil {
			phoneNumber = user.PhoneNumber
		} else {
			phoneNumber = decrypt
		}
		userResponse = append(userResponse, out.UserRoleResourceSettingResponse{
			UserID:         user.UserID,
			ClientID:       user.ClientID,
			Username:       user.Username,
			FirstName:      user.FirstName,
			LastName:       user.LastName,
			PhoneNumber:    phoneNumber,
			ProfilePicture: user.ProfilePicture,
			Role:           user.Role,
			Resource:       user.Resource,
			UserSetting:    user.UserSetting,
		})
	}

	return userResponse, nil
}

func (s authService) GenerateCredentialKey(clientID string) (interface{}, error) {
	data, err := utils.GetUserRedis(s.RedisService, utils.User, clientID)
	if err != nil {
		return nil, errors.New("user not found")
	}
	credentialKey := uuid.New().String()

	credentialKeyMap := struct {
		CredentialKey string `json:"credential_key"`
	}{
		CredentialKey: credentialKey,
	}

	_ = s.RedisService.SaveDataExpired(utils.CredentialKey, data.ClientID, 10, credentialKeyMap)

	return credentialKey, nil
}
