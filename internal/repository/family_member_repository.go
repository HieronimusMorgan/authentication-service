package repository

import (
	"authentication/internal/dto/out"
	"authentication/internal/models"
	"authentication/internal/utils"
	"gorm.io/gorm"
)

type FamilyMemberRepository interface {
	CreateFamilyMember(family *models.FamilyMember, familyPermission *models.FamilyMemberPermission) error
	GetFamilyMemberByID(id uint) (*models.FamilyMember, error)
	UpdateFamilyMember(family *models.FamilyMember) error
	DeleteFamilyMember(family *models.FamilyMember) error
	GetAllFamilyMembers() ([]models.FamilyMember, error)
	GetFamilyMembersByFamilyID(familyID uint) ([]out.FamilyMembersResponse, error)
	GetFamilyMembersByMemberID(memberID uint) ([]models.FamilyMember, error)
	GetFamilyMembersByUserID(userID uint) (models.FamilyMember, error)
	GetFamilyMembersByFamilyIDAndMemberID(familyID uint, memberID uint) (*models.FamilyMember, error)
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

func (r *familyMemberRepository) CreateFamilyMember(family *models.FamilyMember, familyPermission *models.FamilyMemberPermission) error {
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

func (r *familyMemberRepository) GetFamilyMemberByID(id uint) (*models.FamilyMember, error) {
	var familyMember models.FamilyMember
	if err := r.db.Table(utils.TableFamilyMemberName).First(&familyMember, id).Error; err != nil {
		return nil, err
	}
	return &familyMember, nil
}

func (r *familyMemberRepository) UpdateFamilyMember(family *models.FamilyMember) error {
	return r.db.Table(utils.TableFamilyMemberName).Save(family).Error
}

func (r *familyMemberRepository) DeleteFamilyMember(f *models.FamilyMember) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Table(utils.TableFamilyMemberPermissionName).Where("user_id = ?", f.UserID).
			Update("deleted_by", f.DeletedBy).
			Delete(f).Error; err != nil {
			return err
		}
		if err := tx.Table(utils.TableFamilyMemberName).Where("user_id = ?", f.UserID).
			Update("deleted_by", f.DeletedBy).
			Delete(f).Error; err != nil {
			return err
		}
		return nil
	})
}

func (r *familyMemberRepository) GetAllFamilyMembers() ([]models.FamilyMember, error) {
	var familyMembers []models.FamilyMember
	if err := r.db.Table(utils.TableFamilyMemberName).Find(&familyMembers).Error; err != nil {
		return nil, err
	}
	return familyMembers, nil
}

func (r *familyMemberRepository) GetFamilyMembersByFamilyID(familyID uint) ([]out.FamilyMembersResponse, error) {
	var responses []out.FamilyMembersResponse

	err := r.db.Table(utils.TableFamilyMemberName+" AS fm").
		Select("u.user_id, u.username, u.first_name, u.last_name, u.phone_number, u.profile_picture").
		Joins("JOIN users AS u ON fm.user_id = u.user_id").
		Where("fm.family_id = ?", familyID).
		Scan(&responses).Error

	if err != nil {
		return nil, err
	}

	return responses, nil
}

func (r *familyMemberRepository) GetFamilyMembersByMemberID(memberID uint) ([]models.FamilyMember, error) {
	var familyMembers []models.FamilyMember
	if err := r.db.Table(utils.TableFamilyMemberName).Where("user_id = ?", memberID).Find(&familyMembers).Error; err != nil {
		return nil, err
	}
	return familyMembers, nil
}

func (r *familyMemberRepository) GetFamilyMembersByUserID(userID uint) (models.FamilyMember, error) {
	var familyMember models.FamilyMember
	if err := r.db.Table(utils.TableFamilyMemberName).Where("user_id = ?", userID).First(&familyMember).Error; err != nil {
		return familyMember, err
	}
	return familyMember, nil
}

func (r *familyMemberRepository) GetFamilyMembersByFamilyIDAndMemberID(familyID uint, memberID uint) (*models.FamilyMember, error) {
	var familyMember models.FamilyMember
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
