package services

import (
	"Authentication/internal/dto/in"
	"Authentication/internal/dto/out"
	"Authentication/internal/models"
	"Authentication/internal/repository"
	"Authentication/internal/utils"
	"encoding/json"
	"gorm.io/gorm"
)

type AuthService struct {
	AuthRepository     *repository.AuthRepository
	ResourceRepository *repository.ResourceRepository
}

func NewAuthService(db *gorm.DB) *AuthService {
	userRepo := repository.NewAuthRepository(db)
	resourceRepo := repository.NewResourceRepository(db)
	return &AuthService{AuthRepository: userRepo, ResourceRepository: resourceRepo}
}

func (s AuthService) Register(i *in.RegisterRequest) (interface{}, interface{}) {
	if err := utils.ValidateUsername(i.Username); err != nil {
		return nil, err
	}
	pass, err := utils.HashPassword(i.Password)
	if err != nil {
		return nil, err
	}

	firstName := utils.ValidationTrimSpace(i.FirstName)
	lastName := utils.ValidationTrimSpace(i.LastName)
	fullName := firstName + " " + lastName

	user := models.User{
		ClientID:       utils.GenerateClientID(),
		Username:       i.Username,
		Password:       pass,
		FirstName:      firstName,
		LastName:       lastName,
		FullName:       fullName,
		PhoneNumber:    i.PhoneNumber,
		ProfilePicture: i.ProfilePicture,
		RoleID:         2,
	}
	err = s.AuthRepository.CreateUser(&user)
	if err != nil {
		return nil, err
	}

	assignResource, err := s.AuthRepository.AssignUserResource(user.UserID, 1)
	if err != nil {
		return nil, err
	}

	token, err := utils.GenerateToken(user)
	if err != nil {
		return nil, err
	}

	response := out.RegisterResponse{
		UserID:         user.UserID,
		Username:       user.Username,
		FirstName:      user.FirstName,
		LastName:       user.LastName,
		PhoneNumber:    user.PhoneNumber,
		ProfilePicture: user.ProfilePicture,
		Role:           assignResource.Role,
		Resource:       assignResource.Resource,
		Token:          token.AccessToken,
		RefreshToken:   token.RefreshToken,
	}

	utils.SaveDataToRedis("token", user.ClientID, token)
	utils.SaveDataToRedis("user", user.ClientID, user)
	return response, nil
}

func (s AuthService) Login(loginRequest *struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}) (interface{}, error) {
	user, err := s.AuthRepository.GetUserByUsername(loginRequest.Username)
	if err != nil {
		return nil, err
	}

	err = utils.CheckPassword(user.(models.User).Password, loginRequest.Password)
	if err != nil {
		return nil, err
	}

	token, err := utils.GenerateToken(user.(models.User))
	if err != nil {
		return nil, err
	}

	response := out.LoginResponse{
		UserID:         user.(models.User).UserID,
		Username:       user.(models.User).Username,
		FirstName:      user.(models.User).FirstName,
		LastName:       user.(models.User).LastName,
		PhoneNumber:    user.(models.User).PhoneNumber,
		ProfilePicture: user.(models.User).ProfilePicture,
		Token:          token.AccessToken,
		RefreshToken:   token.RefreshToken,
	}
	jsonData, err := json.Marshal(token)
	if err != nil {
		return nil, err
	}
	utils.SaveDataToRedis("token", user.(models.User).ClientID, jsonData)
	utils.SaveDataToRedis("user", user.(models.User).ClientID, user)
	return response, nil
}

func (s AuthService) GetProfile(token string) (*models.User, error) {
	tokenClaims, err := utils.ExtractClaims(token)
	if err != nil {
		return nil, err
	}
	user, err := s.AuthRepository.GetUserByClientID(tokenClaims.ClientID)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s AuthService) RegisterInternalToken(i *struct {
	ResourceName string `json:"resource_name" binding:"required"`
}) (interface{}, error) {
	resource, err := s.ResourceRepository.GetResourceByName(i.ResourceName)
	if err != nil {
		return nil, err
	}

	token, err := utils.GenerateInternalToken(resource.Name)
	if err != nil {
		return nil, err
	}

	err = s.ResourceRepository.CreateInternalToken(resource.ResourceID, token)
	if err != nil {
		return nil, err
	}

	return token, nil
}
