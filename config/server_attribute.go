package config

import (
	"authentication/internal/controller"
	"authentication/internal/middleware"
	"authentication/internal/repository"
	"authentication/internal/services"
	"authentication/internal/utils"
	controllercron "authentication/internal/utils/cron/controller"
	repositorycron "authentication/internal/utils/cron/repository"
	"authentication/internal/utils/cron/service"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func NewServerConfig() (*ServerConfig, error) {
	cfg := LoadConfig()
	redisClient := InitRedis(cfg)
	redisService := utils.NewRedisService(*redisClient)
	db := InitDatabase(cfg)
	engine := InitGin()

	// Graceful Shutdown Handling
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-quit
		log.Println("🛑 Shutting down gracefully...")

		// Close database and Redis before exiting
		CloseDatabase(db)
		CloseRedis(redisClient)

		os.Exit(0)
	}()

	server := &ServerConfig{
		Gin:        engine,
		Config:     cfg,
		DB:         db,
		Redis:      redisService,
		JWTService: utils.NewJWTService(cfg.JWTSecret),
	}

	server.initAesEncrypt()
	server.initRepository()
	server.initServices()
	server.initController()
	server.initMiddleware()
	server.initCron()
	return server, nil
}

// initRepository initializes database access objects (Repository)
func (s *ServerConfig) initRepository() {
	s.Repository = Repository{
		AuthRepo:         repository.NewAuthRepository(*s.DB),
		UserRepo:         repository.NewUserRepository(*s.DB),
		ResourceRepo:     repository.NewResourceRepository(*s.DB),
		RoleRepo:         repository.NewRoleRepository(*s.DB),
		UserRoleRepo:     repository.NewUserRoleRepository(*s.DB),
		UserSessionRepo:  repository.NewUserSessionRepository(*s.DB),
		RoleResourceRepo: repository.NewRoleResourceRepository(*s.DB),
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
			s.JWTService,
			s.Encryption.EncryptionService),
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

func (s *ServerConfig) initController() {
	s.Controller = Controller{
		AuthController:     controller.NewAuthController(s.Services.AuthService, s.Services.UserSessionService, s.JWTService),
		ResourceController: controller.NewResourceController(s.Services.ResourceService, s.JWTService),
		RoleController:     controller.NewRoleController(s.Services.RoleService, s.JWTService),
	}
}

func (s *ServerConfig) initMiddleware() {
	s.Middleware = Middleware{
		AuthMiddleware:  middleware.NewAuthMiddleware(s.JWTService),
		AdminMiddleware: middleware.NewAdminMiddleware(s.JWTService),
	}
}

func (s *ServerConfig) initCron() {
	s.Cron = Cron{
		CronRepository: repositorycron.NewCronRepository(*s.DB),
		CronService:    service.NewCronService(*s.DB, repositorycron.NewCronRepository(*s.DB), s.Services.UserSessionService, s.Services.AuthService),
		CronController: controllercron.NewCronJobController(service.NewCronService(*s.DB, repositorycron.NewCronRepository(*s.DB), s.Services.UserSessionService, s.Services.AuthService)),
	}
	s.Cron.CronService.Start()
}

func (s *ServerConfig) initAesEncrypt() {
	s.Encryption = Encryption{
		EncryptionService: utils.NewEncryption(s.Config.AesEncrypt),
	}
}
