package repository

import (
	"authentication/internal/models"
	"gorm.io/gorm"
)

type UserSettingRepository interface {
	// Create
	AddUserSetting(userSetting *models.UserSetting) error

	// Read
	GetUserSettingBySettingID(settingID uint) (*models.UserSetting, error)
	GetUserSettingByUserID(userID uint) (*models.UserSetting, error)
	GetUserSettingByUserIDAndSettingID(userID, settingID uint) (*models.UserSetting, error)
	GetAllUserSettings() ([]models.UserSetting, error)

	// Update
	UpdateUserSetting(userSetting *models.UserSetting) error
	UpdateGroupInviteTypeSettings(userID uint, inviteType int) error
	UpdateGroupInviteSettings(userID uint, inviteType int, disallowed []string) error

	// Delete
	DeleteUserSettingByUserID(userID uint) error

	// Exists
	UserSettingExists(userID uint) (bool, error)
}

type userSettingRepository struct {
	db gorm.DB
}

func NewUserSettingRepository(db gorm.DB) UserSettingRepository {
	return &userSettingRepository{db: db}
}

func (r *userSettingRepository) AddUserSetting(userSetting *models.UserSetting) error {
	return r.db.Create(&userSetting).Error
}

func (r *userSettingRepository) GetUserSettingBySettingID(settingID uint) (*models.UserSetting, error) {
	var setting models.UserSetting
	if err := r.db.First(&setting, "setting_id = ?", settingID).Error; err != nil {
		return nil, err
	}
	return &setting, nil
}

func (r *userSettingRepository) GetUserSettingByUserID(userID uint) (*models.UserSetting, error) {
	var setting models.UserSetting
	if err := r.db.First(&setting, "user_id = ?", userID).Error; err != nil {
		return nil, err
	}
	return &setting, nil
}

func (r *userSettingRepository) GetUserSettingByUserIDAndSettingID(userID, settingID uint) (*models.UserSetting, error) {
	var setting models.UserSetting
	if err := r.db.First(&setting, "user_id = ? AND setting_id = ?", userID, settingID).Error; err != nil {
		return nil, err
	}
	return &setting, nil
}

func (r *userSettingRepository) GetAllUserSettings() ([]models.UserSetting, error) {
	var settings []models.UserSetting
	if err := r.db.Find(&settings).Error; err != nil {
		return nil, err
	}
	return settings, nil
}

func (r *userSettingRepository) UpdateUserSetting(userSetting *models.UserSetting) error {
	return r.db.Save(userSetting).Error
}

func (r *userSettingRepository) UpdateGroupInviteTypeSettings(userID uint, inviteType int) error {
	return r.db.Model(&models.UserSetting{}).
		Where("user_id = ?", userID).
		Update("group_invite_type", inviteType).Error
}

func (r *userSettingRepository) UpdateGroupInviteSettings(userID uint, inviteType int, disallowed []string) error {
	return r.db.Model(&models.UserSetting{}).
		Where("user_id = ?", userID).
		Updates(map[string]interface{}{
			"group_invite_type":       inviteType,
			"group_invite_disallowed": disallowed,
		}).Error
}

func (r *userSettingRepository) DeleteUserSettingByUserID(userID uint) error {
	return r.db.Where("user_id = ?", userID).Delete(&models.UserSetting{}).Error
}

func (r *userSettingRepository) UserSettingExists(userID uint) (bool, error) {
	var count int64
	err := r.db.Model(&models.UserSetting{}).Where("user_id = ?", userID).Count(&count).Error
	return count > 0, err
}
