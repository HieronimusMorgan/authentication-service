package config

import (
	"authentication/internal/handler"
	"authentication/internal/middleware"
	"authentication/internal/repository"
	"authentication/internal/services"
	"authentication/internal/utils"
	"gorm.io/gorm"
)

// ServerConfig holds all initialized components
type ServerConfig struct {
	Config     *Config
	DB         *gorm.DB
	Redis      utils.RedisService
	JWTService utils.JWTService
	Handler    Handler
	Services   Services
	Repository Repository
	Middleware Middleware
}

// Services holds all service dependencies
type Services struct {
	AuthService        services.AuthService
	UserSessionService services.UsersSessionService
	ResourceService    services.ResourceService
	RoleService        services.RoleService
}

// Repository contains repository (database access objects)
type Repository struct {
	AuthRepo         repository.AuthRepository
	UserRepo         repository.UserRepository
	ResourceRepo     repository.ResourceRepository
	RoleResourceRepo repository.RoleResourceRepository
	RoleRepo         repository.RoleRepository
	UserRoleRepo     repository.UserRoleRepository
	UserSessionRepo  repository.UserSessionRepository
}

type Handler struct {
	AuthHandler     handler.AuthHandler
	ResourceHandler handler.ResourceHandler
	RoleHandler     handler.RoleHandler
}

type Middleware struct {
	AuthMiddleware  middleware.AuthMiddleware
	AdminMiddleware middleware.AdminMiddleware
}
