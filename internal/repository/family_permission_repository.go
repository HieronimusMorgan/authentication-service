package repository

import (
	"authentication/internal/dto/out"
	"authentication/internal/models/family"
	"authentication/internal/utils"
	"gorm.io/gorm"
)

type FamilyPermissionRepository interface {
	CreateFamilyPermission(permission *family.FamilyPermission) error
	GetFamilyPermissionByName(name string) (*family.FamilyPermission, error)
	GetFamilyPermissionByID(id uint) (*family.FamilyPermission, error)
	UpdateFamilyPermission(permission *family.FamilyPermission) error
	DeleteFamilyPermission(permission *family.FamilyPermission) error
	GetAllFamilyPermissions() ([]family.FamilyPermission, error)
	GetAllFamilyPermissionsResponse() ([]out.FamilyPermissionResponse, error)
	GetPermissionsByUserID(userID uint) ([]family.FamilyPermission, error)
}

type familyPermissionRepository struct {
	db gorm.DB
}

func NewFamilyPermissionRepository(db gorm.DB) FamilyPermissionRepository {
	return &familyPermissionRepository{db: db}
}

func (r *familyPermissionRepository) CreateFamilyPermission(permission *family.FamilyPermission) error {
	return r.db.Table(utils.TableFamilyPermissionName).Create(permission).Error
}

func (r *familyPermissionRepository) GetFamilyPermissionByName(name string) (*family.FamilyPermission, error) {
	var permission family.FamilyPermission
	if err := r.db.Table(utils.TableFamilyPermissionName).Where("permission_name = ?", name).First(&permission).Error; err != nil {
		return nil, err
	}
	return &permission, nil
}

func (r *familyPermissionRepository) GetFamilyPermissionByID(id uint) (*family.FamilyPermission, error) {
	var permission family.FamilyPermission
	if err := r.db.Table(utils.TableFamilyPermissionName).First(&permission, id).Error; err != nil {
		return nil, err
	}
	return &permission, nil
}

func (r *familyPermissionRepository) UpdateFamilyPermission(permission *family.FamilyPermission) error {
	return r.db.Table(utils.TableFamilyPermissionName).Save(permission).Error
}

func (r *familyPermissionRepository) DeleteFamilyPermission(permission *family.FamilyPermission) error {
	return r.db.Table(utils.TableFamilyPermissionName).Delete(permission).Error
}

func (r *familyPermissionRepository) GetAllFamilyPermissions() ([]family.FamilyPermission, error) {
	var permissions []family.FamilyPermission
	if err := r.db.Table(utils.TableFamilyPermissionName).Find(&permissions).Error; err != nil {
		return nil, err
	}
	return permissions, nil
}

func (r *familyPermissionRepository) GetAllFamilyPermissionsResponse() ([]out.FamilyPermissionResponse, error) {
	var responses []out.FamilyPermissionResponse
	if err := r.db.Table(utils.TableFamilyPermissionName).Model(&family.FamilyPermission{}).
		Select("permission_id, permission_name, description").
		Scan(&responses).Error; err != nil {
		return nil, err
	}
	return responses, nil
}

func (r *familyPermissionRepository) GetPermissionsByUserID(userID uint) ([]family.FamilyPermission, error) {
	var permissions []family.FamilyPermission
	if err := r.db.Table(utils.TableFamilyPermissionName).Joins("JOIN family_member_permission fmp ON family_permission.permission_id = fmp.permission_id").
		Where("fmp.user_id = ?", userID).
		Find(&permissions).Error; err != nil {
		return nil, err
	}
	return permissions, nil
}
