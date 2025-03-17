package repository

import (
	"authentication/internal/dto/out"
	"authentication/internal/models/family"
	"authentication/internal/utils"
	"gorm.io/gorm"
)

type FamilyMemberRepository interface {
	CreateFamilyMember(family *family.FamilyMember, familyPermission *family.FamilyMemberPermission) error
	GetFamilyMemberByID(id uint) (*family.FamilyMember, error)
	UpdateFamilyMember(family *family.FamilyMember) error
	DeleteFamilyMember(family *family.FamilyMember) error
	GetAllFamilyMembers() ([]family.FamilyMember, error)
	GetFamilyMembersByFamilyID(familyID uint) ([]family.FamilyMember, error)
	GetFamilyMembersByMemberID(memberID uint) ([]family.FamilyMember, error)
	GetFamilyMembersByUserID(userID uint) (family.FamilyMember, error)
	GetFamilyMembersByFamilyIDAndMemberID(familyID uint, memberID uint) (*family.FamilyMember, error)
	GetAllFamilyMemberResponseByFamilyID(familyID uint) ([]out.FamilyMemberResponse, error)
	GetAllFamilyMemberResponseByMemberID(memberID uint) ([]out.FamilyMemberResponse, error)
	GetAllFamilyMemberResponseByFamilyIDAndMemberID(familyID uint, memberID uint) (*out.FamilyMemberResponse, error)
}

type familyMemberRepository struct {
	db gorm.DB
}

func NewFamilyMemberRepository(db gorm.DB) FamilyMemberRepository {
	return &familyMemberRepository{db: db}
}

func (r *familyMemberRepository) CreateFamilyMember(family *family.FamilyMember, familyPermission *family.FamilyMemberPermission) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Table(utils.TableFamilyMemberName).Create(family).Error; err != nil {
			return err
		}
		if err := tx.Table(utils.TableFamilyMemberPermissionName).Create(familyPermission).Error; err != nil {
			return err
		}
		return nil
	})
}

func (r *familyMemberRepository) GetFamilyMemberByID(id uint) (*family.FamilyMember, error) {
	var familyMember family.FamilyMember
	if err := r.db.Table(utils.TableFamilyMemberName).First(&familyMember, id).Error; err != nil {
		return nil, err
	}
	return &familyMember, nil
}

func (r *familyMemberRepository) UpdateFamilyMember(family *family.FamilyMember) error {
	return r.db.Table(utils.TableFamilyMemberName).Save(family).Error
}

func (r *familyMemberRepository) DeleteFamilyMember(family *family.FamilyMember) error {
	return r.db.Table(utils.TableFamilyMemberName).Delete(family).Error
}

func (r *familyMemberRepository) GetAllFamilyMembers() ([]family.FamilyMember, error) {
	var familyMembers []family.FamilyMember
	if err := r.db.Table(utils.TableFamilyMemberName).Find(&familyMembers).Error; err != nil {
		return nil, err
	}
	return familyMembers, nil
}

func (r *familyMemberRepository) GetFamilyMembersByFamilyID(familyID uint) ([]family.FamilyMember, error) {
	var familyMembers []family.FamilyMember
	if err := r.db.Table(utils.TableFamilyMemberName).Where("family_id = ?", familyID).Find(&familyMembers).Error; err != nil {
		return nil, err
	}
	return familyMembers, nil
}

func (r *familyMemberRepository) GetFamilyMembersByMemberID(memberID uint) ([]family.FamilyMember, error) {
	var familyMembers []family.FamilyMember
	if err := r.db.Table(utils.TableFamilyMemberName).Where("user_id = ?", memberID).Find(&familyMembers).Error; err != nil {
		return nil, err
	}
	return familyMembers, nil
}

func (r *familyMemberRepository) GetFamilyMembersByUserID(userID uint) (family.FamilyMember, error) {
	var familyMember family.FamilyMember
	if err := r.db.Table(utils.TableFamilyMemberName).Where("user_id = ?", userID).First(&familyMember).Error; err != nil {
		return familyMember, err
	}
	return familyMember, nil
}

func (r *familyMemberRepository) GetFamilyMembersByFamilyIDAndMemberID(familyID uint, memberID uint) (*family.FamilyMember, error) {
	var familyMember family.FamilyMember
	if err := r.db.Table(utils.TableFamilyMemberName).Where("family_id = ? AND user_id = ?", familyID, memberID).First(&familyMember).Error; err != nil {
		return nil, err
	}
	return &familyMember, nil
}

func (r *familyMemberRepository) GetAllFamilyMemberResponseByFamilyID(familyID uint) ([]out.FamilyMemberResponse, error) {
	var responses []out.FamilyMemberResponse
	if err := r.db.Table(utils.TableFamilyMemberName).Where("family_id = ?", familyID).Joins("JOIN users ON family_members.user_id = users.user_id").Scan(&responses).Error; err != nil {
		return nil, err
	}
	return responses, nil
}

func (r *familyMemberRepository) GetAllFamilyMemberResponseByMemberID(memberID uint) ([]out.FamilyMemberResponse, error) {
	var responses []out.FamilyMemberResponse
	if err := r.db.Table(utils.TableFamilyMemberName).Where("user_id = ?", memberID).Joins("JOIN family ON family_members.family_id = family.family_id").Scan(&responses).Error; err != nil {
		return nil, err
	}
	return responses, nil
}

func (r *familyMemberRepository) GetAllFamilyMemberResponseByFamilyIDAndMemberID(familyID uint, memberID uint) (*out.FamilyMemberResponse, error) {
	var response out.FamilyMemberResponse
	if err := r.db.Table(utils.TableFamilyMemberName).Where("family_id = ? AND user_id = ?", familyID, memberID).Joins("JOIN users ON family_members.user_id = users.user_id").Scan(&response).Error; err != nil {
		return nil, err
	}
	return &response, nil
}
