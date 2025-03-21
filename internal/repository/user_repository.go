package repository

import (
	"authentication/internal/dto/out"
	"authentication/internal/models"
	"gorm.io/gorm"
)

type UserRepository interface {
	RegisterUser(user **models.Users) error
	GetUserByUsername(username string) (*models.Users, error)
	GetUserByEmail(email string) (*models.Users, error)
	GetUserByID(id uint) (*models.Users, error)
	UpdateUser(user *models.Users) error
	DeleteUser(user *models.Users) error
	GetAllUsers() (*[]models.Users, error)
	GetUsers() (*[]models.Users, error)
	GetUserByRole(role uint) (*[]models.Users, error)
	GetUserByPhoneNumber(number string) (*models.Users, error)
	GetUserByClientID(clientID string) (*models.Users, error)
	GetUserByPinCodeAndClientID(pinCode, clientID string) (*models.Users, error)
	GetUserByClientAndRole(clientID, roleID uint) (*[]models.Users, error)
	GetUserResponseByClientID(clientID string) (*out.UserResponse, error)
	DeleteUserByID(id uint) error
	UpdateRole(user *models.Users) error
	GetListUser() (*[]models.Users, error)
	GetUserByResourceID(resourceID uint) (*[]models.Users, error)
	ChangePassword(user *models.Users) error
	UpdatePinAttempts(clientID string) error
	ResetPinAttempts(user *models.Users) error
}

type userRepository struct {
	db gorm.DB
}

func NewUserRepository(db gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r userRepository) RegisterUser(user **models.Users) error {
	if err := r.db.Create(user).Error; err != nil {
		return err
	}
	return nil
}

func (r userRepository) GetUserByUsername(username string) (*models.Users, error) {
	var user models.Users
	err := r.db.Where("username = ?", username).First(&user).Error
	return &user, err
}

func (r userRepository) GetUserByEmail(email string) (*models.Users, error) {
	var user models.Users
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r userRepository) GetUserByID(id uint) (*models.Users, error) {
	var user models.Users
	if err := r.db.Where("user_id = ?", id).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r userRepository) UpdateUser(user *models.Users) error {
	if err := r.db.Save(&user).Error; err != nil {
		return err
	}
	return nil
}

func (r userRepository) DeleteUser(user *models.Users) error {
	if err := r.db.Model(&user).
		Update("deleted_by", user.DeletedBy).
		Delete(&user).Error; err != nil {
		return err
	}
	return nil
}

func (r userRepository) GetAllUsers() (*[]models.Users, error) {
	var users []models.Users
	if err := r.db.Find(&users).Error; err != nil {
		return nil, err
	}
	return &users, nil
}

func (r userRepository) GetUsers() (*[]models.Users, error) {
	var users []models.Users
	if err := r.db.Where("deleted_at IS NOT NULL").Find(&users).Error; err != nil {
		return nil, err
	}
	return &users, nil
}

func (r userRepository) GetUserByPhoneNumber(number string) (*models.Users, error) {
	var user models.Users
	if err := r.db.Where("phone_number = ?", number).Find(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r userRepository) GetUserByRole(role uint) (*[]models.Users, error) {
	var users []models.Users
	if err := r.db.Where("role_id = ?", role).Find(&users).Error; err != nil {
		return nil, err
	}
	return &users, nil
}

func (r userRepository) GetUserByClientID(clientID string) (*models.Users, error) {
	var users models.Users
	if err := r.db.Where("client_id = ?", clientID).Find(&users).Error; err != nil {
		return nil, err
	}
	return &users, nil
}

func (r userRepository) GetUserByPinCodeAndClientID(pinCode, clientID string) (*models.Users, error) {
	var user models.Users
	if err := r.db.Where("pin_code = ? AND client_id = ?", pinCode, clientID).Find(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r userRepository) GetUserByClientAndRole(clientID, roleID uint) (*[]models.Users, error) {
	var users []models.Users
	if err := r.db.Where("client_id = ? AND role_id = ?", clientID, roleID).Find(&users).Error; err != nil {
		return nil, err
	}
	return &users, nil
}

func (r userRepository) GetUserResponseByClientID(clientID string) (*out.UserResponse, error) {
	var user out.UserResponse
	if err := r.db.Table("users").Where("client_id = ?", clientID).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r userRepository) DeleteUserByID(id uint) error {
	if err := r.db.Where("user_id = ?", id).Delete(&models.Users{}).Error; err != nil {
		return err
	}
	return nil
}

func (r userRepository) UpdateRole(user *models.Users) error {
	if err := r.db.Model(&user).
		Update("role_id", user.RoleID).
		Update("updated_by", user.UpdatedBy).
		Error; err != nil {
		return err
	}
	return nil
}

func (r userRepository) GetListUser() (*[]models.Users, error) {
	var users []models.Users
	if err := r.db.Preload("Role").Find(&users).Where("delete_at IS NULL").Error; err != nil {
		return nil, err
	}
	return &users, nil
}

func (r userRepository) GetUserByResourceID(resourceID uint) (*[]models.Users, error) {
	var users []models.Users
	if err := r.db.Preload("Role").Joins("JOIN role_resources rr ON rr.role_id = users.role_id").
		Joins("JOIN resources r ON r.resource_id = rr.resource_id").
		Where("r.resource_id = ?", resourceID).
		Find(&users).Error; err != nil {
		return nil, err
	}
	return &users, nil
}

func (r userRepository) ChangePassword(user *models.Users) error {
	if err := r.db.Model(&user).
		Update("password", user.Password).
		Update("updated_by", user.UpdatedBy).
		Error; err != nil {
		return err
	}
	return nil
}

func (r userRepository) UpdatePinAttempts(clientID string) error {
	if err := r.db.Model(&models.Users{}).
		Where("client_id = ?", clientID).
		Update("pin_attempts", gorm.Expr("pin_attempts + 1")).
		Error; err != nil {
		return err
	}
	return nil
}

func (r userRepository) ResetPinAttempts(user *models.Users) error {
	if err := r.db.Model(&user).
		Update("pin_attempts", 0).
		Error; err != nil {
		return err
	}
	return nil
}
