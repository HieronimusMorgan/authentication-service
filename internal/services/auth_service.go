package services

import (
	"authentication/internal/dto/in"
	"authentication/internal/dto/out"
	"authentication/internal/models"
	"authentication/internal/repository"
	"authentication/internal/utils"
	"gorm.io/gorm"
)

type AuthService struct {
	AuthRepository         *repository.AuthRepository
	ResourceRepository     *repository.ResourceRepository
	RoleRepository         *repository.RoleRepository
	RoleResourceRepository *repository.RoleResourceRepository
	UserRepository         *repository.UserRepository
	UserRoleRepository     *repository.UserRoleRepository
}

func NewAuthService(db *gorm.DB) *AuthService {
	return &AuthService{
		AuthRepository:         repository.NewAuthRepository(db),
		ResourceRepository:     repository.NewResourceRepository(db),
		RoleRepository:         repository.NewRoleRepository(db),
		RoleResourceRepository: repository.NewRoleResourceRepository(db),
		UserRepository:         repository.NewUserRepository(db),
		UserRoleRepository:     repository.NewUserRoleRepository(db),
	}
}

func (s AuthService) Register(req *in.RegisterRequest) (interface{}, error) {
	if err := utils.ValidateUsername(req.Username); err != nil {
		return nil, err
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
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
		RoleID:         2,
	}

	if err := s.UserRepository.RegisterUser(&user); err != nil {
		return nil, err
	}

	userRole := &models.UserRole{
		UserID:    user.UserID,
		RoleID:    user.Role.RoleID,
		CreatedBy: "system",
		UpdatedBy: "system",
	}
	if err := s.UserRoleRepository.RegisterUserRole(&userRole); err != nil {
		return nil, err
	}

	token, err := utils.GenerateToken(*user)
	if err != nil {
		return nil, err
	}

	_ = utils.SaveDataToRedis("token", user.ClientID, token)
	_ = utils.SaveDataToRedis("user", user.ClientID, user)

	// Construct response
	response := out.RegisterResponse{
		UserID:         user.UserID,
		Username:       user.Username,
		FirstName:      user.FirstName,
		LastName:       user.LastName,
		PhoneNumber:    user.PhoneNumber,
		ProfilePicture: user.ProfilePicture,
		Token:          token.AccessToken,
		RefreshToken:   token.RefreshToken,
	}
	return response, nil
}

func (s AuthService) Login(req *struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}) (interface{}, error) {
	user, err := s.AuthRepository.GetUserByUsername(req.Username)
	if err != nil {
		return nil, err
	}
	u := user.(models.Users)
	if err := utils.CheckPassword(u.Password, req.Password); err != nil {
		return nil, err
	}

	role, err := s.RoleRepository.GetRoleByID(u.RoleID)
	u.Role = *role

	token, err := utils.GenerateToken(u)
	if err != nil {
		return nil, err
	}

	_ = utils.SaveDataToRedis("token", u.ClientID, token)
	_ = utils.SaveDataToRedis("user", u.ClientID, user)

	response := out.LoginResponse{
		UserID:         u.UserID,
		Username:       u.Username,
		FirstName:      u.FirstName,
		LastName:       u.LastName,
		PhoneNumber:    u.PhoneNumber,
		ProfilePicture: u.ProfilePicture,
		Token:          token.AccessToken,
		RefreshToken:   token.RefreshToken,
	}
	return response, nil
}

func (s AuthService) GetProfile(token string) (*models.Users, error) {
	claims, err := utils.ExtractClaims(token)
	if err != nil {
		return nil, err
	}

	user, err := s.AuthRepository.GetUserByClientID(claims.ClientID)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s AuthService) RegisterInternalToken(req *struct {
	ResourceName string `json:"resource_name" binding:"required"`
}) (interface{}, error) {
	resource, err := s.ResourceRepository.GetResourceByName(req.ResourceName)
	if err != nil {
		return nil, err
	}

	token, err := utils.GenerateInternalToken(resource.Name)
	if err != nil {
		return nil, err
	}

	if err := s.ResourceRepository.CreateInternalToken(resource.ResourceID, token); err != nil {
		return nil, err
	}

	return token, nil
}
