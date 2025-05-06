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
	nt "authentication/internal/utils/nats"
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

	server.initNats()
	server.initAesEncrypt()
	server.initRepository()
	server.initTransactional()
	server.initServices()
	server.initController()
	server.initMiddleware()
	server.initCron()
	return server, nil
}

// initRepository initializes database access objects (Repository)
func (s *ServerConfig) initRepository() {
	s.Repository = Repository{
		AuthRepository:         repository.NewAuthRepository(*s.DB),
		UserRepository:         repository.NewUserRepository(*s.DB),
		UserKeyRepository:      repository.NewUserKeyRepository(*s.DB),
		UserSettingRepository:  repository.NewUserSettingRepository(*s.DB),
		ResourceRepository:     repository.NewResourceRepository(*s.DB),
		RoleRepository:         repository.NewRoleRepository(*s.DB),
		UserRoleRepository:     repository.NewUserRoleRepository(*s.DB),
		UserSessionRepository:  repository.NewUserSessionRepository(*s.DB),
		UserResourceRepository: repository.NewUserResourceRepository(*s.DB),
	}
}

// initTransactional initializes transactional repository
func (s *ServerConfig) initTransactional() {
	s.Transactional = Transactional{
		UserTransactionalRepository: repository.NewUserTransactionalRepository(*s.DB),
	}
}

// initServices initializes the application services
func (s *ServerConfig) initServices() {
	s.Services = Services{
		AuthService: services.NewAuthService(s.Repository.AuthRepository,
			s.Repository.ResourceRepository,
			s.Repository.RoleRepository,
			s.Repository.UserResourceRepository,
			s.Repository.UserRepository,
			s.Repository.UserKeyRepository,
			s.Repository.UserRoleRepository,
			s.Repository.UserSessionRepository,
			s.Transactional.UserTransactionalRepository,
			s.Repository.UserSettingRepository,
			s.Redis,
			s.JWTService,
			s.Encryption.EncryptionService,
			s.Nats.NatsService),
		UserService:        services.NewUserService(s.Repository.UserRepository, s.Repository.UserKeyRepository, s.Repository.UserSettingRepository, s.Redis, s.JWTService, s.Encryption.EncryptionService),
		ResourceService:    services.NewResourceService(s.Repository.ResourceRepository, s.Repository.UserResourceRepository, s.Repository.RoleRepository, s.Repository.UserRepository, s.Nats.NatsService),
		RoleService:        services.NewRoleService(s.Repository.RoleRepository, s.Repository.UserRepository),
		UserSessionService: services.NewUsersSessionService(s.Repository.UserSessionRepository, s.Repository.UserRepository, s.JWTService, s.Redis),
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
		UserController:     controller.NewUserController(s.Services.UserService, s.JWTService, s.Config.CdnUrl),
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
		EncryptionService: utils.NewEncryption(s.Config.AesEncrypt, s.Config.AesFixedIV),
	}
}

// initNats initializes the application services
func (s *ServerConfig) initNats() {
	s.Nats = Nats{
		NatsService: nt.NewNatsService(s.Config.NatsUrl),
	}
}
