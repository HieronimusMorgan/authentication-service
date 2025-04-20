package config

import (
	"authentication/internal/controller"
	"authentication/internal/middleware"
	"authentication/internal/repository"
	"authentication/internal/services"
	"authentication/internal/utils"
	controllercron "authentication/internal/utils/cron/controller"
	repositorycron "authentication/internal/utils/cron/repository"
	servicescron "authentication/internal/utils/cron/service"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// ServerConfig holds all initialized components
type ServerConfig struct {
	Gin           *gin.Engine
	Config        *Config
	DB            *gorm.DB
	Redis         utils.RedisService
	JWTService    utils.JWTService
	Controller    Controller
	Services      Services
	Repository    Repository
	Transactional Transactional
	Middleware    Middleware
	Cron          Cron
	Encryption    Encryption
}

// Services holds all service dependencies
type Services struct {
	AuthService        services.AuthService
	UserService        services.UserService
	UserSessionService services.UsersSessionService
	ResourceService    services.ResourceService
	RoleService        services.RoleService
}

// Repository contains repository (database access objects)
type Repository struct {
	AuthRepository         repository.AuthRepository
	UserRepository         repository.UserRepository
	UserSettingRepository  repository.UserSettingRepository
	ResourceRepository     repository.ResourceRepository
	UserResourceRepository repository.UserResourceRepository
	RoleRepository         repository.RoleRepository
	UserRoleRepository     repository.UserRoleRepository
	UserSessionRepository  repository.UserSessionRepository
}

type Controller struct {
	AuthController     controller.AuthController
	UserController     controller.UserController
	ResourceController controller.ResourceController
	RoleController     controller.RoleController
}

type Middleware struct {
	AuthMiddleware  middleware.AuthMiddleware
	AdminMiddleware middleware.AdminMiddleware
}

type Transactional struct {
	UserTransactionalRepository repository.UserTransactionalRepository
}

type Cron struct {
	CronService    servicescron.CronService
	CronRepository repositorycron.CronRepository
	CronController controllercron.CronJobController
}

type Encryption struct {
	EncryptionService utils.Encryption
}
