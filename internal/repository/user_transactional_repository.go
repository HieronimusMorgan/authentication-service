package repository

import (
	"authentication/internal/models"
	"authentication/internal/utils"
	"gorm.io/gorm"
)

type UserTransactionalRepository interface {
	RegistrationUser(user *models.Users) error
	DeleteUser(user models.Users) error
}

type userTransactionalRepository struct {
	db gorm.DB
}

func NewUserTransactionalRepository(db gorm.DB) UserTransactionalRepository {
	return &userTransactionalRepository{db: db}
}

func (r *userTransactionalRepository) RegistrationUser(user *models.Users) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Table(utils.TableUsersName).Create(&user).Error; err != nil {
			return err
		}
		var role models.Role
		if err := r.db.Table(utils.TableRolesName).Where("name = ?", "User").First(&role).Error; err != nil {
			return err
		}

		userRole := &models.UserRole{
			UserID:    user.UserID,
			RoleID:    role.RoleID,
			CreatedBy: "system",
			UpdatedBy: "system",
		}

		if err := tx.Table(utils.TableUserRolesName).Create(userRole).Error; err != nil {
			return err
		}

		return nil
	})
}

func (r *userTransactionalRepository) DeleteUser(user models.Users) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(&user).Error; err != nil {
			return err
		}
		return nil
	})
}
