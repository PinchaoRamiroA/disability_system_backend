package usecase

import (
	"context"
	"time"

	"disability_system_backend/internal/modules/auth/domain"
	"disability_system_backend/internal/modules/auth/ports"
	apperrors "disability_system_backend/internal/shared/errors"
)

type LoginUseCase struct {
	userRepo       ports.UserRepository
	roleRepo       ports.RoleRepository
	tokenService   ports.TokenService
	passwordHasher ports.PasswordHasher
	tokenExpiry    time.Duration
}

func NewLoginUseCase(
	userRepo ports.UserRepository,
	roleRepo ports.RoleRepository,
	tokenService ports.TokenService,
	passwordHasher ports.PasswordHasher,
	tokenExpiry time.Duration,
) *LoginUseCase {
	return &LoginUseCase{
		userRepo:       userRepo,
		roleRepo:       roleRepo,
		tokenService:   tokenService,
		passwordHasher: passwordHasher,
		tokenExpiry:    tokenExpiry,
	}
}

func (uc *LoginUseCase) Execute(ctx context.Context, email, password string) (*ports.TokenPair, *domain.User, *domain.Role, error) {
	user, err := uc.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return nil, nil, nil, apperrors.ErrInvalidCredentials.WithError(err)
	}

	if !user.Estado || user.IsDeleted {
		return nil, nil, nil, apperrors.ErrUnauthorized.WithMessage("usuario inactivo")
	}

	if !uc.passwordHasher.Check(password, user.PasswordHash) {
		return nil, nil, nil, apperrors.ErrInvalidCredentials.WithMessage("contraseña incorrecta")
	}

	role, err := uc.roleRepo.FindByID(ctx, user.IDRol)
	if err != nil {
		return nil, nil, nil, apperrors.ErrInternal.WithError(err)
	}

	tokens, err := uc.tokenService.GenerateTokenPair(user.ID, user.Correo, role.Nombre)
	if err != nil {
		return nil, nil, nil, apperrors.ErrInternal.WithError(err)
	}

	return tokens, user, role, nil
}

func (uc *LoginUseCase) GetExpirationSeconds() int64 {
	return int64(uc.tokenExpiry.Seconds())
}