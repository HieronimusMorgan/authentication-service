package repository

import (
	"authentication/internal/dto/out"
	"authentication/internal/models"
	"authentication/internal/utils"
	"gorm.io/gorm"
)

type FamilyMemberPermissionRepository interface {
	GetFamilyMemberPermissionByUserID(userID uint) ([]models.FamilyMemberPermission, error)
	AddFamilyMemberPermission(family *models.FamilyMemberPermission) error
	RemoveFamilyMemberPermission(family *models.FamilyMemberPermission) error
	GetListFamilyMemberPermissionByFamilyIDAndUserID(familyID uint, userID uint) ([]out.FamilyPermissionResponse, error)

	CreateFamilyMemberPermission(family *models.FamilyMemberPermission) error
	GetFamilyMemberPermissionByID(id uint) (*models.FamilyMemberPermission, error)
	DeleteFamilyMemberPermissionByFamilyAndUserAndPermission(familyID, userID, permissionID uint) error
	DeleteFamilyMemberPermission(family *models.FamilyMemberPermission) error
	DeleteFamilyMemberPermissionByFamilyAndMember(familyID uint, memberID uint) error
	GetAllFamilyMemberPermissions() ([]models.FamilyMemberPermission, error)
	GetFamilyMemberPermissionsByFamilyID(familyID uint) ([]models.FamilyMemberPermission, error)
	GetFamilyMemberPermissionsByMemberID(memberID uint) ([]models.FamilyMemberPermission, error)
	GetFamilyMemberPermissionsByFamilyIDAndMemberID(familyID uint, memberID uint) (*models.FamilyMemberPermission, error)
	GetAllFamilyMemberPermissionResponseByFamilyID(familyID uint) ([]out.FamilyMemberPermissionResponse, error)
	GetAllFamilyMemberPermissionResponseByMemberID(memberID uint) ([]out.FamilyMemberPermissionResponse, error)
	GetAllFamilyMemberPermissionResponseByFamilyIDAndMemberID(familyID uint, memberID uint) (*out.FamilyMemberPermissionResponse, error)
}

type familyMemberPermissionRepository struct {
	db gorm.DB
}

func NewFamilyMemberPermissionRepository(db gorm.DB) FamilyMemberPermissionRepository {
	return &familyMemberPermissionRepository{db: db}
}

func (r *familyMemberPermissionRepository) CreateFamilyMemberPermission(family *models.FamilyMemberPermission) error {
	return r.db.Table(utils.TableFamilyMemberPermissionName).Create(family).Error
}

func (r *familyMemberPermissionRepository) GetFamilyMemberPermissionByID(id uint) (*models.FamilyMemberPermission, error) {
	var familyMemberPermission models.FamilyMemberPermission
	if err := r.db.Table(utils.TableFamilyMemberPermissionName).First(&familyMemberPermission, id).Error; err != nil {
		return nil, err
	}
	return &familyMemberPermission, nil
}

func (r *familyMemberPermissionRepository) GetFamilyMemberPermissionByUserID(userID uint) ([]models.FamilyMemberPermission, error) {
	var familyMemberPermissions []models.FamilyMemberPermission
	if err := r.db.Table(utils.TableFamilyMemberPermissionName).Where("user_id = ?", userID).Find(&familyMemberPermissions).Error; err != nil {
		return nil, err
	}
	return familyMemberPermissions, nil
}

func (r *familyMemberPermissionRepository) AddFamilyMemberPermission(familyMemberPermission *models.FamilyMemberPermission) error {
	return r.db.Table(utils.TableFamilyMemberPermissionName).Create(familyMemberPermission).Error
}

func (r *familyMemberPermissionRepository) RemoveFamilyMemberPermission(familyMemberPermission *models.FamilyMemberPermission) error {
	return r.db.Table(utils.TableFamilyMemberPermissionName).Delete(familyMemberPermission).Error
}

func (r *familyMemberPermissionRepository) DeleteFamilyMemberPermissionByFamilyAndUserAndPermission(familyID, userID, permissionID uint) error {
	return r.db.Unscoped().Table(utils.TableFamilyMemberPermissionName).Where("family_id = ? AND user_id = ? AND permission_id = ?", familyID, userID, permissionID).Delete(&models.FamilyMemberPermission{}).Error
}

func (r *familyMemberPermissionRepository) DeleteFamilyMemberPermission(family *models.FamilyMemberPermission) error {
	return r.db.Table(utils.TableFamilyMemberPermissionName).Delete(family).Error
}

func (r *familyMemberPermissionRepository) DeleteFamilyMemberPermissionByFamilyAndMember(familyID uint, memberID uint) error {
	return r.db.Table(utils.TableFamilyMemberPermissionName).Where("family_id = ? AND user_id = ?", familyID, memberID).Delete(&models.FamilyMemberPermission{}).Error
}

func (r *familyMemberPermissionRepository) GetAllFamilyMemberPermissions() ([]models.FamilyMemberPermission, error) {
	var familyMemberPermissions []models.FamilyMemberPermission
	if err := r.db.Table(utils.TableFamilyMemberPermissionName).Find(&familyMemberPermissions).Error; err != nil {
		return nil, err
	}
	return familyMemberPermissions, nil
}

func (r *familyMemberPermissionRepository) GetFamilyMemberPermissionsByFamilyID(familyID uint) ([]models.FamilyMemberPermission, error) {
	var permissions []models.FamilyMemberPermission
	if err := r.db.Table(utils.TableFamilyMemberPermissionName).Where("family_id = ?", familyID).Find(&permissions).Error; err != nil {
		return nil, err
	}
	return permissions, nil
}

func (r *familyMemberPermissionRepository) GetFamilyMemberPermissionsByMemberID(memberID uint) ([]models.FamilyMemberPermission, error) {
	var permissions []models.FamilyMemberPermission
	if err := r.db.Table(utils.TableFamilyMemberPermissionName).Where("user_id = ?", memberID).Find(&permissions).Error; err != nil {
		return nil, err
	}
	return permissions, nil
}

func (r *familyMemberPermissionRepository) GetFamilyMemberPermissionsByFamilyIDAndMemberID(familyID uint, memberID uint) (*models.FamilyMemberPermission, error) {
	var permission models.FamilyMemberPermission
	if err := r.db.Table(utils.TableFamilyMemberPermissionName).Where("family_id = ? AND user_id = ?", familyID, memberID).First(&permission).Error; err != nil {
		return nil, err
	}
	return &permission, nil
}

func (r *familyMemberPermissionRepository) GetAllFamilyMemberPermissionResponseByFamilyID(familyID uint) ([]out.FamilyMemberPermissionResponse, error) {
	var responses []out.FamilyMemberPermissionResponse
	if err := r.db.Table(utils.TableFamilyMemberPermissionName).Where("family_id = ?", familyID).Joins("JOIN family_permission ON family_member_permission.permission_id = family_permission.permission_id").Scan(&responses).Error; err != nil {
		return nil, err
	}
	return responses, nil
}

func (r *familyMemberPermissionRepository) GetAllFamilyMemberPermissionResponseByMemberID(memberID uint) ([]out.FamilyMemberPermissionResponse, error) {
	var responses []out.FamilyMemberPermissionResponse
	if err := r.db.Table(utils.TableFamilyMemberPermissionName).Where("user_id = ?", memberID).Joins("JOIN family_permission ON family_member_permission.permission_id = family_permission.permission_id").Scan(&responses).Error; err != nil {
		return nil, err
	}
	return responses, nil
}

func (r *familyMemberPermissionRepository) GetAllFamilyMemberPermissionResponseByFamilyIDAndMemberID(familyID uint, memberID uint) (*out.FamilyMemberPermissionResponse, error) {
	var response out.FamilyMemberPermissionResponse
	if err := r.db.Table(utils.TableFamilyMemberPermissionName).Where("family_id = ? AND user_id = ?", familyID, memberID).Joins("JOIN family_permission ON family_member_permission.permission_id = family_permission.permission_id").Scan(&response).Error; err != nil {
		return nil, err
	}
	return &response, nil
}

func (r *familyMemberPermissionRepository) GetListFamilyMemberPermissionByFamilyIDAndUserID(familyID uint, userID uint) ([]out.FamilyPermissionResponse, error) {
	var permissions []out.FamilyPermissionResponse
	query := `
		SELECT fp.*
		FROM "family_permission" fp
		JOIN "family_member_permission" fmp ON fp.permission_id = fmp.permission_id
		WHERE fmp.user_id = ? AND fmp.family_id = ?
	`
	if err := r.db.Raw(query, userID, familyID).Scan(&permissions).Error; err != nil {
		return nil, err
	}
	return permissions, nil
}
