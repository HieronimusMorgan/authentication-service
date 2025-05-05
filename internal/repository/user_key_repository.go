package repository

import (
	"authentication/internal/models"
	"authentication/internal/utils"
	"gorm.io/gorm"
)

type UserKeyRepository interface {
	AddUserKey(userKey *models.UserKey) error
	GetUserKeyByUserID(userID uint) (*models.UserKey, error)
}

type userKeyRepository struct {
	db gorm.DB
}

func NewUserKeyRepository(db gorm.DB) UserKeyRepository {
	return &userKeyRepository{
		db: db,
	}
}

func (r *userKeyRepository) AddUserKey(userKey *models.UserKey) error {
	if err := r.db.Table(utils.TableUserKeysName).Create(userKey).Error; err != nil {
		return err
	}
	return nil
}

func (r *userKeyRepository) GetUserKeyByUserID(userID uint) (*models.UserKey, error) {
	var userKey models.UserKey
	if err := r.db.Table(utils.TableUserKeysName).Where("user_id = ?", userID).First(&userKey).Error; err != nil {
		return nil, err
	}
	return &userKey, nil
}
