package app

import (
	"context"

	authhttp "disability_system_backend/internal/modules/auth/adapters/http"
	"disability_system_backend/internal/modules/auth/adapters/postgres"
	"disability_system_backend/internal/modules/auth/usecase"
	cobroshttp "disability_system_backend/internal/modules/cobros/adapters/http"
	incapacidadeshttp "disability_system_backend/internal/modules/incapacidades/adapters/http"
	notificacioneshttp "disability_system_backend/internal/modules/notificaciones/adapters/http"
	usuarioshttp "disability_system_backend/internal/modules/usuarios/adapters/http"
	"disability_system_backend/internal/shared/auth"
	"disability_system_backend/internal/shared/database"
	"disability_system_backend/internal/shared/middleware"
	"disability_system_backend/internal/shared/router"
	"disability_system_backend/internal/shared/storage"
)

func (a *App) InitRouter() error {
	if a.DB == nil {
		a.Logger.Error("database not initialized, cannot create router")
		return ErrDatabaseNotInitialized
	}

	r := router.New(a.Config, a.Logger)
	r.Use(middleware.Recovery(a.Logger))
	r.Use(middleware.RequestID())
	r.Use(middleware.CORS())
	r.Use(middleware.Logger(a.Logger))

	r.SetupHealth(a.DB)

	a.Router = r
	return nil
}

func (a *App) InitDependencies() error {
	db, err := database.NewConnection(a.Config)
	if err != nil {
		a.Logger.Error("failed to connect to database", "error", err)
		return err
	}

	if err := database.HealthCheck(db); err != nil {
		a.Logger.Error("database health check failed", "error", err)
		return err
	}

	a.DB = db
	return nil
}

func (a *App) InitAuth() *auth.JWTService {
	jwtService := auth.NewJWTService(
		a.Config.JWT.Secret,
		a.Config.JWT.Expiration,
		a.Config.JWT.RefreshExpiry,
	)

	// Initialize Repositories
	userRepo := postgres.NewUserRepository(a.DB)
	roleRepo := postgres.NewRoleRepository(a.DB)
	tokenService := postgres.NewTokenService(jwtService)
	passwordHasher := postgres.NewPasswordHasher()

	// Initialize UseCases
	loginUseCase := usecase.NewLoginUseCase(userRepo, roleRepo, tokenService, passwordHasher, a.Config.JWT.Expiration)
	registerUseCase := usecase.NewRegisterUseCase(userRepo, passwordHasher)
	refreshUseCase := usecase.NewRefreshTokenUseCase(tokenService)

	// Initialize Handler
	a.AuthHandler = authhttp.NewAuthHandler(loginUseCase, registerUseCase, refreshUseCase)

	// Register Routes
	a.InitAuthRoutes()
	incapacidadeshttp.Register(a.Router.V1(), a.DB, jwtService, a.StorageService)
	cobroshttp.Register(a.Router.V1(), a.DB, jwtService)
	notificacioneshttp.Register(a.Router.V1(), a.DB, jwtService)
	usuarioshttp.Register(a.Router.V1(), a.DB, jwtService)

	if a.Config.App.Env != "test" {
		storageService, err := storage.NewStorageService(context.Background(), storage.LoadR2Config())
		if err != nil {
			a.Logger.Warn("storage service not configured", "error", err)
		} else {
			a.StorageService = storageService
		}
	}

	return jwtService
}
