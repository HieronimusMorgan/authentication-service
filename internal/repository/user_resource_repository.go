package repository

import (
	"authentication/internal/models"
	"authentication/internal/utils"
	"gorm.io/gorm"
)

type UserResourceRepository interface {
	RegisterUserResource(userResource models.UserResource) error
	GetUserResourceByUserID(userID uint) (*models.UserResource, error)
	GetUserResourceByUserIDAndResourceID(userID, resourceID uint) (*models.UserResource, error)
	UpdateUserResource(userResource models.UserResource) error
	DeleteUserResource(userResource *models.UserResource) error
	GetAllUserResource() (*[]models.UserResource, error)
	GetUserResourceByID(id uint) (*models.UserResource, error)
	GetUserResourceByResourceID(roleID uint) (*models.UserResource, error)
}

type userResourceRepository struct {
	db gorm.DB
}

func NewUserResourceRepository(db gorm.DB) UserResourceRepository {
	return &userResourceRepository{db: db}
}

func (r userResourceRepository) RegisterUserResource(userResource models.UserResource) error {
	err := r.db.Table(utils.TableUserResourceName).Create(&userResource).Error
	if err != nil {
		return err
	}
	return nil
}

func (r userResourceRepository) GetUserResourceByUserID(userID uint) (*models.UserResource, error) {
	var userResource models.UserResource
	err := r.db.Table(utils.TableUserResourceName).Where("user_id = ?", userID).First(&userResource).Error
	if err != nil {
		return nil, err
	}
	return &userResource, nil
}

func (r userResourceRepository) GetUserResourceByUserIDAndResourceID(userID, resourceID uint) (*models.UserResource, error) {
	var userResource models.UserResource
	err := r.db.Table(utils.TableUserResourceName).Where("user_id = ? AND resource_id = ?", userID, resourceID).First(&userResource).Error
	if err != nil {
		return nil, err
	}
	return &userResource, nil
}

func (r userResourceRepository) UpdateUserResource(userResource models.UserResource) error {
	err := r.db.Table(utils.TableUserResourceName).Save(userResource).Error
	if err != nil {
		return err
	}
	return nil
}

func (r userResourceRepository) DeleteUserResource(userResource *models.UserResource) error {
	err := r.db.Unscoped().Table(utils.TableUserResourceName).Model(userResource).
		Update("deleted_by", userResource.DeletedBy).
		Delete(userResource).Error
	if err != nil {
		return err
	}
	return nil
}

func (r userResourceRepository) GetAllUserResource() (*[]models.UserResource, error) {
	var userResources []models.UserResource
	err := r.db.Table(utils.TableUserResourceName).Find(&userResources).Error
	if err != nil {
		return nil, err
	}
	return &userResources, nil
}

func (r userResourceRepository) GetUserResourceByID(id uint) (*models.UserResource, error) {
	var userResource models.UserResource
	err := r.db.Table(utils.TableUserResourceName).Where("id = ?", id).First(&userResource).Error
	if err != nil {
		return nil, err
	}
	return &userResource, nil
}

func (r userResourceRepository) GetUserResourceByResourceID(roleID uint) (*models.UserResource, error) {
	var userResource models.UserResource
	err := r.db.Table(utils.TableUserResourceName).Where("role_id = ?", roleID).First(&userResource).Error
	if err != nil {
		return nil, err
	}
	return &userResource, nil
}
