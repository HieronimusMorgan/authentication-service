package repository

import (
	"authentication/internal/dto/out"
	"authentication/internal/models"
	"authentication/internal/utils"
	"gorm.io/gorm"
)

type FamilyPermissionRepository interface {
	CreateFamilyPermission(permission *models.FamilyPermission) error
	GetFamilyPermissionByName(name string) (*models.FamilyPermission, error)
	GetListFamilyPermissionAccess(userID uint) ([]models.FamilyPermission, error)
	GetFamilyPermissionByID(id uint) (*models.FamilyPermission, error)
	UpdateFamilyPermission(permission *models.FamilyPermission) error
	DeleteFamilyPermission(permission *models.FamilyPermission) error
	GetAllFamilyPermissions() ([]models.FamilyPermission, error)
	GetAllFamilyPermissionsResponse() ([]out.FamilyPermissionResponse, error)
	GetPermissionsByUserID(userID uint) ([]models.FamilyPermission, error)
}

type familyPermissionRepository struct {
	db gorm.DB
}

func NewFamilyPermissionRepository(db gorm.DB) FamilyPermissionRepository {
	return &familyPermissionRepository{db: db}
}

func (r *familyPermissionRepository) CreateFamilyPermission(permission *models.FamilyPermission) error {
	return r.db.Table(utils.TableFamilyPermissionName).Create(permission).Error
}

func (r *familyPermissionRepository) GetFamilyPermissionByName(name string) (*models.FamilyPermission, error) {
	var permission models.FamilyPermission
	if err := r.db.Table(utils.TableFamilyPermissionName).Where("permission_name = ?", name).First(&permission).Error; err != nil {
		return nil, err
	}
	return &permission, nil
}

func (r *familyPermissionRepository) GetListFamilyPermissionAccess(userID uint) ([]models.FamilyPermission, error) {
	var permissions []models.FamilyPermission
	query := `
		SELECT fp.*
		FROM "family_permission" fp
		JOIN "family_member_permission" fmp ON fp.permission_id = fmp.permission_id
		WHERE fmp.user_id = ?
	`
	if err := r.db.Raw(query, userID).Scan(&permissions).Error; err != nil {
		return nil, err
	}
	return permissions, nil
}

func (r *familyPermissionRepository) GetFamilyPermissionByID(id uint) (*models.FamilyPermission, error) {
	var permission models.FamilyPermission
	if err := r.db.Table(utils.TableFamilyPermissionName).First(&permission, id).Error; err != nil {
		return nil, err
	}
	return &permission, nil
}

func (r *familyPermissionRepository) UpdateFamilyPermission(permission *models.FamilyPermission) error {
	return r.db.Table(utils.TableFamilyPermissionName).Save(permission).Error
}

func (r *familyPermissionRepository) DeleteFamilyPermission(permission *models.FamilyPermission) error {
	return r.db.Table(utils.TableFamilyPermissionName).Delete(permission).Error
}

func (r *familyPermissionRepository) GetAllFamilyPermissions() ([]models.FamilyPermission, error) {
	var permissions []models.FamilyPermission
	if err := r.db.Table(utils.TableFamilyPermissionName).Find(&permissions).Error; err != nil {
		return nil, err
	}
	return permissions, nil
}

func (r *familyPermissionRepository) GetAllFamilyPermissionsResponse() ([]out.FamilyPermissionResponse, error) {
	var responses []out.FamilyPermissionResponse
	if err := r.db.Table(utils.TableFamilyPermissionName).Model(&models.FamilyPermission{}).
		Select("permission_id, permission_name, description").
		Scan(&responses).Error; err != nil {
		return nil, err
	}
	return responses, nil
}

func (r *familyPermissionRepository) GetPermissionsByUserID(userID uint) ([]models.FamilyPermission, error) {
	var permissions []models.FamilyPermission
	if err := r.db.Table(utils.TableFamilyPermissionName).Joins("JOIN family_member_permission fmp ON family_permission.permission_id = fmp.permission_id").
		Where("fmp.user_id = ?", userID).
		Find(&permissions).Error; err != nil {
		return nil, err
	}
	return permissions, nil
}
