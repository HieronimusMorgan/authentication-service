package repository

import (
	"authentication/internal/dto/out"
	"authentication/internal/models"
	"authentication/internal/utils"
	"fmt"
	"gorm.io/gorm"
)

type FamilyRepository interface {
	CreateFamily(family *models.Family, familyPermission []models.FamilyPermission, member *models.FamilyMember) error
	GetFamilyByID(id uint) (*models.Family, error)
	GetFamilyByFamilyIdAndOwnerID(familyID, userID uint) (*models.Family, error)
	GetFamilyByOwnerID(userID uint) (*models.Family, error)
	UpdateFamily(family *models.Family) error
	GetAllFamilies() ([]models.Family, error)
	GetFamilyResponseByClientID(clientID string) (*out.FamilyResponse, error)
	DeleteFamilyByID(id uint) error
	ChangeFamilyOwner(familyID uint, newOwnerID uint) error
}

type familyRepository struct {
	db gorm.DB
}

func NewFamilyRepository(db gorm.DB) FamilyRepository {
	return &familyRepository{db: db}
}

func (r *familyRepository) CreateFamily(family *models.Family, familyPermissions []models.FamilyPermission, member *models.FamilyMember) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// Insert Family
		if err := tx.Table(utils.TableFamilyName).Create(family).Error; err != nil {
			return fmt.Errorf("failed to create family: %w", err)
		}

		// Insert Family Member Permissions
		for _, p := range familyPermissions {
			permissionRecord := models.FamilyMemberPermission{
				FamilyID:     family.FamilyID,
				UserID:       family.OwnerID,
				CreatedBy:    family.CreatedBy,
				PermissionID: p.PermissionID,
			}
			if err := tx.Table(utils.TableFamilyMemberPermissionName).Create(&permissionRecord).Error; err != nil {
				return fmt.Errorf("failed to insert permission (ID: %d) for user (ID: %d): %w", p.PermissionID, family.OwnerID, err)
			}
		}

		// Assign Family ID to Member
		member.FamilyID = family.FamilyID

		// Insert Family Member
		if err := tx.Table(utils.TableFamilyMemberName).Create(member).Error; err != nil {
			return fmt.Errorf("failed to create family member: %w", err)
		}

		return nil
	})
}

func (r *familyRepository) GetFamilyByID(id uint) (*models.Family, error) {
	var f models.Family
	if err := r.db.Table(utils.TableFamilyName).Where("family_id = ?", id).First(&f).Error; err != nil {
		return nil, err
	}
	return &f, nil
}

func (r *familyRepository) GetFamilyByFamilyIdAndOwnerID(familyID, userID uint) (*models.Family, error) {
	var f models.Family
	if err := r.db.Table(utils.TableFamilyName).Where("family_id = ? AND owner_id = ?", familyID, userID).First(&f).Error; err != nil {
		return nil, err
	}
	return &f, nil
}

func (r *familyRepository) GetFamilyByOwnerID(userID uint) (*models.Family, error) {
	var f models.Family
	if err := r.db.Table(utils.TableFamilyName).Where("owner_id = ?", userID).First(&f).Error; err != nil {
		return nil, err
	}
	return &f, nil
}

func (r *familyRepository) UpdateFamily(family *models.Family) error {
	return r.db.Table(utils.TableFamilyName).Save(family).Error
}

func (r *familyRepository) GetAllFamilies() ([]models.Family, error) {
	var families []models.Family
	if err := r.db.Table(utils.TableFamilyName).Find(&families).Error; err != nil {
		return nil, err
	}
	return families, nil
}

func (r *familyRepository) GetFamilyResponseByClientID(clientID string) (*out.FamilyResponse, error) {
	var response out.FamilyResponse
	if err := r.db.Table(utils.TableFamilyName).Model(&models.Family{}).
		Where("owner_id = ?", clientID).
		Select("family_id, family_name, owner_id").
		Scan(&response).Error; err != nil {
		return nil, err
	}
	return &response, nil
}

func (r *familyRepository) DeleteFamilyByID(id uint) error {
	return r.db.Table(utils.TableFamilyName).Delete(&models.Family{}, id).Error
}

func (r *familyRepository) ChangeFamilyOwner(familyID uint, newOwnerID uint) error {
	return r.db.Table(utils.TableFamilyName).Model(&models.Family{}).
		Where("family_id = ?", familyID).
		Update("owner_id", newOwnerID).Error
}
