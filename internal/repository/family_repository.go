package repository

import (
	"authentication/internal/dto/out"
	"authentication/internal/models/family"
	"authentication/internal/utils"
	"gorm.io/gorm"
)

type FamilyRepository interface {
	CreateFamily(family *family.Family, permission *family.FamilyMemberPermission, member *family.FamilyMember) error
	GetFamilyByFamilyIdAndOwnerID(familyID, userID uint) (*family.Family, error)
	GetFamilyByOwnerID(userID uint) (*family.Family, error)
	UpdateFamily(family *family.Family) error
	DeleteFamily(family *family.Family) error
	GetAllFamilies() ([]family.Family, error)
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

func (r *familyRepository) CreateFamily(family *family.Family, permission *family.FamilyMemberPermission, member *family.FamilyMember) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Table(utils.TableFamilyName).Create(family).Error; err != nil {
			return err
		}
		permission.FamilyID = family.FamilyID
		if err := tx.Table(utils.TableFamilyMemberPermissionName).Create(permission).Error; err != nil {
			return err
		}
		member.FamilyID = family.FamilyID
		if err := tx.Table(utils.TableFamilyMemberName).Create(member).Error; err != nil {
			return err
		}
		return nil
	})
}

func (r *familyRepository) GetFamilyByFamilyIdAndOwnerID(familyID, userID uint) (*family.Family, error) {
	var f family.Family
	if err := r.db.Table(utils.TableFamilyName).Where("family_id = ? AND owner_id = ?", familyID, userID).First(&f).Error; err != nil {
		return nil, err
	}
	return &f, nil
}

func (r *familyRepository) GetFamilyByOwnerID(userID uint) (*family.Family, error) {
	var f family.Family
	if err := r.db.Table(utils.TableFamilyName).Where("owner_id = ?", userID).First(&f).Error; err != nil {
		return nil, err
	}
	return &f, nil
}

func (r *familyRepository) UpdateFamily(family *family.Family) error {
	return r.db.Table(utils.TableFamilyName).Save(family).Error
}

func (r *familyRepository) DeleteFamily(family *family.Family) error {
	return r.db.Table(utils.TableFamilyName).Delete(family).Error
}

func (r *familyRepository) GetAllFamilies() ([]family.Family, error) {
	var families []family.Family
	if err := r.db.Table(utils.TableFamilyName).Find(&families).Error; err != nil {
		return nil, err
	}
	return families, nil
}

func (r *familyRepository) GetFamilyResponseByClientID(clientID string) (*out.FamilyResponse, error) {
	var response out.FamilyResponse
	if err := r.db.Table(utils.TableFamilyName).Model(&family.Family{}).
		Where("owner_id = ?", clientID).
		Select("family_id, family_name, owner_id").
		Scan(&response).Error; err != nil {
		return nil, err
	}
	return &response, nil
}

func (r *familyRepository) DeleteFamilyByID(id uint) error {
	return r.db.Table(utils.TableFamilyName).Delete(&family.Family{}, id).Error
}

func (r *familyRepository) ChangeFamilyOwner(familyID uint, newOwnerID uint) error {
	return r.db.Table(utils.TableFamilyName).Model(&family.Family{}).
		Where("family_id = ?", familyID).
		Update("owner_id", newOwnerID).Error
}
