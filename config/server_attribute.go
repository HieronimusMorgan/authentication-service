package config

import (
	"authentication/internal/handler"
	"authentication/internal/middleware"
	"authentication/internal/repository"
	"authentication/internal/services"
	"authentication/internal/utils"
	"log"

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

func NewServerConfig() (*ServerConfig, error) {
	cfg := LoadConfig()
	redisClient := InitRedis(cfg)
	redisService := utils.NewRedisService(*redisClient)
	db := InitDatabase(cfg)

	server := &ServerConfig{
		Config:     cfg,
		DB:         db,
		Redis:      redisService,
		JWTService: utils.NewJWTService(cfg.JWTSecret),
	}

	server.initRepository()
	server.initServices()
	server.initHandler()
	server.initMiddleware()
	return server, nil
}

// initRepository initializes database access objects (Repository)
func (s *ServerConfig) initRepository() {
	s.Repository = Repository{
		AuthRepo:         repository.NewAuthRepository(s.DB),
		UserRepo:         repository.NewUserRepository(s.DB),
		ResourceRepo:     repository.NewResourceRepository(s.DB),
		RoleRepo:         repository.NewRoleRepository(s.DB),
		UserRoleRepo:     repository.NewUserRoleRepository(s.DB),
		UserSessionRepo:  repository.NewUserSessionRepository(s.DB),
		RoleResourceRepo: repository.NewRoleResourceRepository(s.DB),
	}
}

// initServices initializes the application services
func (s *ServerConfig) initServices() {
	s.Services = Services{
		AuthService: services.NewAuthService(s.Repository.AuthRepo,
			s.Repository.ResourceRepo,
			s.Repository.RoleRepo,
			s.Repository.RoleResourceRepo,
			s.Repository.UserRepo,
			s.Repository.UserRoleRepo,
			s.Repository.UserSessionRepo,
			s.Redis,
			s.JWTService),
		ResourceService:    services.NewResourceService(s.Repository.ResourceRepo, s.Repository.RoleResourceRepo, s.Repository.RoleRepo, s.Repository.UserRepo),
		RoleService:        services.NewRoleService(s.Repository.RoleRepo, s.Repository.UserRepo),
		UserSessionService: services.NewUsersSessionService(s.Repository.UserSessionRepo, s.Repository.UserRepo, s.JWTService, s.Redis),
	}
}

// Start initializes everything and returns an error if something fails
func (s *ServerConfig) Start() error {
	log.Println("✅ Server configuration initialized successfully!")
	return nil
}

func (s *ServerConfig) initHandler() {
	s.Handler = Handler{
		AuthHandler:     handler.NewAuthHandler(s.Services.AuthService, s.Services.UserSessionService, s.JWTService),
		ResourceHandler: handler.NewResourceHandler(s.Services.ResourceService, s.JWTService),
		RoleHandler:     handler.NewRoleHandler(s.Services.RoleService, s.JWTService),
	}
}

func (s *ServerConfig) initMiddleware() {
	s.Middleware = Middleware{
		AuthMiddleware:  middleware.NewAuthMiddleware(s.JWTService),
		AdminMiddleware: middleware.NewAdminMiddleware(s.JWTService),
	}
}
