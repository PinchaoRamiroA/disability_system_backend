package authhttp

import (
	"time"

	"disability_system_backend/internal/modules/auth/adapters/postgres"
	"disability_system_backend/internal/modules/auth/usecase"
	"disability_system_backend/internal/shared/auth"
	"disability_system_backend/internal/shared/router"

	"gorm.io/gorm"
)

func Register(v1 *router.APIVersion, db *gorm.DB, jwtService *auth.JWTService, tokenTTL time.Duration) {
	userRepo := postgres.NewUserRepository(db)
	roleRepo := postgres.NewRoleRepository(db)
	tokenService := postgres.NewTokenService(jwtService)
	passwordHasher := postgres.NewPasswordHasher()

	loginUseCase := usecase.NewLoginUseCase(userRepo, roleRepo, tokenService, passwordHasher, tokenTTL)
	registerUseCase := usecase.NewRegisterUseCase(userRepo, passwordHasher)
	refreshUseCase := usecase.NewRefreshTokenUseCase(tokenService)

	authHandler := NewAuthHandler(loginUseCase, registerUseCase, refreshUseCase)

	auth := v1.Group("/auth")
	{
		auth.POST("/login", authHandler.Login)
		auth.POST("/register", authHandler.Register)
		auth.POST("/refresh", authHandler.RefreshToken)
	}
}
