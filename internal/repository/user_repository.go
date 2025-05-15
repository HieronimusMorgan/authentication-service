package repository

import (
	"authentication/internal/dto/out"
	"authentication/internal/models"
	"authentication/internal/utils"
	"errors"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type UserRepository interface {
	RegisterUser(user **models.Users) error
	CheckClientID(clientID string) bool
	GetUserByUsername(username string) (*models.Users, error)
	GetUserByEmail(email string) (*models.Users, error)
	GetUserByID(id uint) (*models.Users, error)
	UpdateUser(user *models.Users) error
	DeleteUser(user *models.Users) error
	GetAllUsers() (*[]models.Users, error)
	GetUsers() (*[]models.Users, error)
	GetUserByRole(role uint) (*[]models.Users, error)
	GetUserByRolePagination(role uint, index, size int) (*[]out.UserRoleResponse, error)
	GetUserByPhoneNumber(number string) (*models.Users, error)
	GetUserByClientID(clientID string) (*models.Users, error)
	GetUserByPinCodeAndClientID(pinCode, clientID string) (*models.Users, error)
	GetUserByClientAndRole(clientID, roleID uint) (*[]models.Users, error)
	GetUserResponseByClientID(clientID string) (*out.UserResponse, error)
	DeleteUserByID(id uint) error
	UpdateRole(user *models.Users) error
	GetListUser() (*[]models.Users, error)
	GetListUserResponse() (*[]out.UserRoleResourceSettingResponse, error)
	GetListUserByUserIDResponse(userID uint) (*[]out.UserRoleResourceSettingResponse, error)
	GetUserByResourceID(resourceID uint) (*[]models.Users, error)
	ChangePassword(user *models.Users) error
	UpdatePinAttempts(clientID string) error
	ResetPinAttempts(user *models.Users) error
	GetAllUsersByResourceId(resources *models.Resource) (*[]models.Users, error)
	GetUserRedisByClientID(clientID string) (*models.UserRedis, error)
	SaveUserKey(keys *models.UserKey) error
	GetUserKey(userID uint) (*models.UserKey, error)
	GetCountUserByRole(roleID uint) (int64, error)
}

type userRepository struct {
	db gorm.DB
}

func NewUserRepository(db gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r userRepository) RegisterUser(user **models.Users) error {
	if err := r.db.Create(user).Error; err != nil {
		return err
	}
	return nil
}

func (r userRepository) CheckClientID(clientID string) bool {
	var user models.Users
	if err := r.db.Where("client_id = ?", clientID).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false
		}
		return false
	}
	return true
}

func (r userRepository) GetUserByUsername(username string) (*models.Users, error) {
	var user models.Users
	err := r.db.Where("username = ?", username).First(&user).Error
	return &user, err
}

func (r userRepository) GetUserByEmail(email string) (*models.Users, error) {
	var user models.Users
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r userRepository) GetUserByID(id uint) (*models.Users, error) {
	var user models.Users
	if err := r.db.Where("user_id = ?", id).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r userRepository) UpdateUser(user *models.Users) error {
	if err := r.db.Save(&user).Error; err != nil {
		return err
	}
	return nil
}

func (r userRepository) DeleteUser(user *models.Users) error {
	if err := r.db.Model(&user).
		Update("deleted_by", user.DeletedBy).
		Delete(&user).Error; err != nil {
		return err
	}
	return nil
}

func (r userRepository) GetAllUsers() (*[]models.Users, error) {
	var users []models.Users
	if err := r.db.Find(&users).Error; err != nil {
		return nil, err
	}
	return &users, nil
}

func (r userRepository) GetUsers() (*[]models.Users, error) {
	var users []models.Users
	if err := r.db.Where("deleted_at IS NOT NULL").Find(&users).Error; err != nil {
		return nil, err
	}
	return &users, nil
}

func (r userRepository) GetUserByPhoneNumber(number string) (*models.Users, error) {
	var user models.Users
	if err := r.db.Where("phone_number = ?", number).Find(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r userRepository) GetUserByRole(role uint) (*[]models.Users, error) {
	var users []models.Users
	if err := r.db.Where("role_id = ?", role).Find(&users).Error; err != nil {
		return nil, err
	}
	return &users, nil
}

func (r userRepository) GetUserByRolePagination(role uint, index, size int) (*[]out.UserRoleResponse, error) {

	query := `
		SELECT 
			u.user_id, u.client_id, u.username, u.first_name, u.last_name, u.phone_number, u.profile_picture
		FROM users u
		LEFT JOIN user_roles ur ON ur.user_id = u.user_id
		LEFT JOIN roles r ON r.role_id = ur.role_id
		WHERE r.role_id = ? AND u.deleted_at IS NULL
		ORDER BY u.user_id ASC
		LIMIT ? OFFSET ?;
	`

	var rows []out.UserRoleResponse
	if err := r.db.Raw(query, role, size, (index-1)*size).Scan(&rows).Error; err != nil {
		return nil, err
	}

	return &rows, nil
}

func (r userRepository) GetUserByClientID(clientID string) (*models.Users, error) {
	var users models.Users
	if err := r.db.Where("client_id = ?", clientID).Find(&users).Error; err != nil {
		return nil, err
	}
	return &users, nil
}

func (r userRepository) GetUserByPinCodeAndClientID(pinCode, clientID string) (*models.Users, error) {
	var user models.Users
	if err := r.db.Where("pin_code = ? AND client_id = ?", pinCode, clientID).Find(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r userRepository) GetUserByClientAndRole(clientID, roleID uint) (*[]models.Users, error) {
	var users []models.Users
	if err := r.db.Where("client_id = ? AND role_id = ?", clientID, roleID).Find(&users).Error; err != nil {
		return nil, err
	}
	return &users, nil
}

func (r userRepository) GetUserResponseByClientID(clientID string) (*out.UserResponse, error) {
	var user out.UserResponse
	if err := r.db.Table("users").Where("client_id = ?", clientID).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r userRepository) DeleteUserByID(id uint) error {
	if err := r.db.Where("user_id = ?", id).Delete(&models.Users{}).Error; err != nil {
		return err
	}
	return nil
}

func (r userRepository) UpdateRole(user *models.Users) error {
	if err := r.db.Model(&user).
		Update("role_id", user.RoleID).
		Update("updated_by", user.UpdatedBy).
		Error; err != nil {
		return err
	}
	return nil
}

func (r userRepository) GetListUser() (*[]models.Users, error) {
	var users []models.Users
	if err := r.db.Where("deleted_at IS NULL").Find(&users).Error; err != nil {
		return nil, err
	}
	return &users, nil
}

func (r userRepository) GetListUserResponse() (*[]out.UserRoleResourceSettingResponse, error) {
	type userRow struct {
		UserID                uint
		ClientID              string
		Username              string
		FirstName             string
		LastName              string
		PhoneNumber           string
		ProfilePicture        string
		RoleID                *uint
		RoleName              *string
		RoleDescription       *string
		ResourceID            *uint
		ResourceName          *string
		ResourceDescription   *string
		SettingID             *uint
		GroupInviteType       *int
		GroupInviteDisallowed pq.Int32Array `gorm:"type:integer[]"`
	}

	query := `
		SELECT 
			u.user_id, u.client_id, u.username, u.first_name, u.last_name, u.phone_number, u.profile_picture,
			r.role_id, r.name AS role_name, r.description AS role_description,
			res.resource_id, res.name AS resource_name, res.description AS resource_description,
			us.setting_id, us.group_invite_type, us.group_invite_disallowed
		FROM users u
		LEFT JOIN user_roles ur ON ur.user_id = u.user_id
		LEFT JOIN roles r ON r.role_id = ur.role_id
		LEFT JOIN user_resources ure ON ure.user_id = u.user_id
		LEFT JOIN resources res ON res.resource_id = ure.resource_id
		LEFT JOIN user_settings us ON us.user_id = u.user_id
		WHERE u.deleted_at IS NULL
		ORDER BY u.user_id ASC, r.role_id ASC, res.resource_id ASC;
	`

	var rows []userRow
	if err := r.db.Raw(query).Scan(&rows).Error; err != nil {
		return nil, err
	}

	userMap := make(map[uint]*out.UserRoleResourceSettingResponse)

	for _, row := range rows {
		user, exists := userMap[row.UserID]
		if !exists {
			user = &out.UserRoleResourceSettingResponse{
				UserID:         row.UserID,
				ClientID:       row.ClientID,
				Username:       row.Username,
				FirstName:      row.FirstName,
				LastName:       row.LastName,
				PhoneNumber:    row.PhoneNumber,
				ProfilePicture: &row.ProfilePicture,
				Role:           []out.RoleResponse{},
				Resource:       []out.ResourceResponse{},
			}
			if row.SettingID != nil {
				user.UserSetting = out.UserSettingResponse{
					SettingID:             *row.SettingID,
					GroupInviteType:       utils.DerefInt(row.GroupInviteType),
					GroupInviteDisallowed: row.GroupInviteDisallowed,
				}
			}
			userMap[row.UserID] = user
		}

		if row.RoleID != nil && !utils.ContainsRole(user.Role, *row.RoleID) {
			user.Role = append(user.Role, out.RoleResponse{
				RoleID:      *row.RoleID,
				Name:        utils.DerefStr(row.RoleName),
				Description: utils.DerefStr(row.RoleDescription),
			})
		}

		if row.ResourceID != nil && !utils.ContainsResource(user.Resource, *row.ResourceID) {
			user.Resource = append(user.Resource, out.ResourceResponse{
				ResourceID:  *row.ResourceID,
				Name:        utils.DerefStr(row.ResourceName),
				Description: utils.DerefStr(row.ResourceDescription),
			})
		}
	}

	result := make([]out.UserRoleResourceSettingResponse, 0, len(userMap))
	for _, u := range userMap {
		result = append(result, *u)
	}

	return &result, nil
}
func (r userRepository) GetListUserByUserIDResponse(userID uint) (*[]out.UserRoleResourceSettingResponse, error) {
	type userRow struct {
		UserID                uint
		ClientID              string
		Username              string
		FirstName             string
		LastName              string
		PhoneNumber           string
		ProfilePicture        string
		RoleID                *uint
		RoleName              *string
		RoleDescription       *string
		ResourceID            *uint
		ResourceName          *string
		ResourceDescription   *string
		SettingID             *uint
		GroupInviteType       *int
		GroupInviteDisallowed pq.Int32Array `gorm:"type:integer[]"`
	}

	query := `
		SELECT 
			u.user_id, u.client_id, u.username, u.first_name, u.last_name, u.phone_number, u.profile_picture,
			r.role_id, r.name AS role_name, r.description AS role_description,
			res.resource_id, res.name AS resource_name, res.description AS resource_description,
			us.setting_id, us.group_invite_type, us.group_invite_disallowed
		FROM users u
		LEFT JOIN user_roles ur ON ur.user_id = u.user_id
		LEFT JOIN roles r ON r.role_id = ur.role_id
		LEFT JOIN user_resources ure ON ure.user_id = u.user_id
		LEFT JOIN resources res ON res.resource_id = ure.resource_id
		LEFT JOIN user_settings us ON us.user_id = u.user_id
		WHERE u.user_id = ? AND u.deleted_at IS NULL
		ORDER BY u.user_id ASC, r.role_id ASC, res.resource_id ASC;
	`

	var rows []userRow
	if err := r.db.Raw(query, userID).Scan(&rows).Error; err != nil {
		return nil, err
	}

	userMap := make(map[uint]*out.UserRoleResourceSettingResponse)

	for _, row := range rows {
		user, exists := userMap[row.UserID]
		if !exists {
			user = &out.UserRoleResourceSettingResponse{
				UserID:         row.UserID,
				ClientID:       row.ClientID,
				Username:       row.Username,
				FirstName:      row.FirstName,
				LastName:       row.LastName,
				PhoneNumber:    row.PhoneNumber,
				ProfilePicture: &row.ProfilePicture,
				Role:           []out.RoleResponse{},
				Resource:       []out.ResourceResponse{},
			}
			if row.SettingID != nil {
				user.UserSetting = out.UserSettingResponse{
					SettingID:             *row.SettingID,
					GroupInviteType:       utils.DerefInt(row.GroupInviteType),
					GroupInviteDisallowed: row.GroupInviteDisallowed,
				}
			}
			userMap[row.UserID] = user
		}

		if row.RoleID != nil && !utils.ContainsRole(user.Role, *row.RoleID) {
			user.Role = append(user.Role, out.RoleResponse{
				RoleID:      *row.RoleID,
				Name:        utils.DerefStr(row.RoleName),
				Description: utils.DerefStr(row.RoleDescription),
			})
		}

		if row.ResourceID != nil && !utils.ContainsResource(user.Resource, *row.ResourceID) {
			user.Resource = append(user.Resource, out.ResourceResponse{
				ResourceID:  *row.ResourceID,
				Name:        utils.DerefStr(row.ResourceName),
				Description: utils.DerefStr(row.ResourceDescription),
			})
		}
	}

	result := make([]out.UserRoleResourceSettingResponse, 0, len(userMap))
	for _, u := range userMap {
		result = append(result, *u)
	}

	return &result, nil
}

func (r userRepository) GetUserByResourceID(resourceID uint) (*[]models.Users, error) {
	var users []models.Users
	if err := r.db.Preload("Roles").Joins("JOIN user_resource ur ON ur.user_id = users.user_id").
		Joins("JOIN resources r ON r.resource_id = ur.resource_id").
		Where("r.resource_id = ?", resourceID).
		Find(&users).Error; err != nil {
		return nil, err
	}
	return &users, nil
}

func (r userRepository) ChangePassword(user *models.Users) error {
	if err := r.db.Model(&user).
		Update("password", user.Password).
		Update("updated_by", user.UpdatedBy).
		Error; err != nil {
		return err
	}
	return nil
}

func (r userRepository) UpdatePinAttempts(clientID string) error {
	if err := r.db.Model(&models.Users{}).
		Where("client_id = ?", clientID).
		Update("pin_attempts", gorm.Expr("pin_attempts + 1")).
		Error; err != nil {
		return err
	}
	return nil
}

func (r userRepository) ResetPinAttempts(user *models.Users) error {
	if err := r.db.Model(&user).
		Update("pin_attempts", 0).
		Error; err != nil {
		return err
	}
	return nil
}

func (r userRepository) GetAllUsersByResourceId(resources *models.Resource) (*[]models.Users, error) {
	var users []models.Users
	err := r.db.Table(utils.TableUsersName).
		Select("users.user_id, users.username, users.first_name, users.last_name").
		Joins("JOIN user_resources ur ON users.user_id = ur.user_id").
		Where("ur.resource_id = ?", resources.ResourceID).
		Find(&users).Error
	if err != nil {
		return nil, err
	}
	return &users, nil
}

func (r userRepository) GetUserRedisByClientID(clientID string) (*models.UserRedis, error) {
	// Fetch user and settings
	var row struct {
		UserID                uint
		ClientID              string
		Username              string
		Email                 string
		Password              string
		PinCode               string
		PinAttempts           int
		FirstName             string
		LastName              string
		FullName              string
		PhoneNumber           string
		ProfilePicture        string
		DeviceID              *string
		SettingID             uint
		GroupInviteType       int
		GroupInviteDisallowed pq.Int32Array
	}

	query := `
		SELECT 
			u.user_id, u.client_id, u.username, u.email, u.password, u.pin_code, u.pin_attempts,
			u.first_name, u.last_name, u.full_name, u.phone_number, u.profile_picture, u.device_id,
			us.setting_id, us.group_invite_type, us.group_invite_disallowed
		FROM users u
		LEFT JOIN user_settings us ON us.user_id = u.user_id
		WHERE u.client_id = ?
	`
	if err := r.db.Raw(query, clientID).Scan(&row).Error; err != nil {
		return nil, err
	}

	user := &models.UserRedis{
		UserID:         row.UserID,
		ClientID:       row.ClientID,
		Username:       row.Username,
		Email:          row.Email,
		Password:       row.Password,
		PinCode:        row.PinCode,
		PinAttempts:    row.PinAttempts,
		FirstName:      row.FirstName,
		LastName:       row.LastName,
		FullName:       row.FullName,
		PhoneNumber:    row.PhoneNumber,
		ProfilePicture: row.ProfilePicture,
		DeviceID:       row.DeviceID,
		UserSetting: models.UserSettingRedis{
			SettingID:             row.SettingID,
			GroupInviteType:       row.GroupInviteType,
			GroupInviteDisallowed: row.GroupInviteDisallowed,
		},
	}

	var roles []models.RoleRedis
	roleQuery := `
		SELECT r.role_id, r.name, r.description
		FROM roles r
		JOIN user_roles ur ON ur.role_id = r.role_id
		WHERE ur.user_id = ?
	`
	if err := r.db.Raw(roleQuery, user.UserID).Scan(&roles).Error; err != nil {
		return nil, err
	}
	user.Role = roles

	var resources []models.ResourceRedis
	resourceQuery := `
		SELECT r.resource_id, r.name, r.description
		FROM resources r
		JOIN user_resources ur ON ur.resource_id = r.resource_id
		WHERE ur.user_id = ?
	`
	if err := r.db.Raw(resourceQuery, user.UserID).Scan(&resources).Error; err != nil {
		return nil, err
	}
	user.Resource = resources

	return user, nil
}

func (r userRepository) SaveUserKey(keys *models.UserKey) error {
	if err := r.db.Table(utils.TableUserKeysName).Create(keys).Error; err != nil {
		return err
	}
	return nil
}

func (r userRepository) GetUserKey(userID uint) (*models.UserKey, error) {
	var userKey models.UserKey
	if err := r.db.Table(utils.TableUserKeysName).Where("user_id = ?", userID).First(&userKey).Error; err != nil {
		return nil, err
	}
	return &userKey, nil
}

func (r userRepository) GetCountUserByRole(roleID uint) (int64, error) {
	var count int64
	query := `
		SELECT COUNT(ur.*)
		FROM roles r
		JOIN user_roles ur ON ur.role_id = r.role_id
		WHERE r.role_id = ?
	`
	if err := r.db.Raw(query, roleID).Scan(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}
