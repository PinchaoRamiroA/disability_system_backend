package usecase

import (
	"context"
	"time"

	usuariosdomain "disability_system_backend/internal/modules/usuarios/domain"
	"disability_system_backend/internal/modules/auth/ports"
	apperrors "disability_system_backend/internal/shared/errors"
)

type RegisterUseCase struct {
	userRepo       ports.UserRepository
	passwordHasher ports.PasswordHasher
}

func NewRegisterUseCase(
	userRepo ports.UserRepository,
	passwordHasher ports.PasswordHasher,
) *RegisterUseCase {
	return &RegisterUseCase{
		userRepo:       userRepo,
		passwordHasher: passwordHasher,
	}
}

func (uc *RegisterUseCase) Execute(
	ctx context.Context,
	nombre, email, password, numeroDocumento string,
) (*usuariosdomain.Usuario, error) {
	exists, err := uc.userRepo.EmailExists(ctx, email)
	if err != nil {
		return nil, apperrors.ErrDatabase.WithError(err)
	}
	if exists {
		return nil, apperrors.ErrEmailAlreadyExists.WithMessage("el email ya está registrado")
	}

	exists, err = uc.userRepo.DocumentExists(ctx, numeroDocumento)
	if err != nil {
		return nil, apperrors.ErrDatabase.WithError(err)
	}
	if exists {
		return nil, apperrors.ErrConflict.WithMessage("el número de documento ya está registrado")
	}

	hashedPassword, err := uc.passwordHasher.Hash(password)
	if err != nil {
		return nil, apperrors.ErrHashPassword.WithError(err)
	}

	now := time.Now()
	user := &usuariosdomain.Usuario{
		IDRol:           4,
		Nombre:          nombre,
		Correo:          email,
		PasswordHash:    hashedPassword,
		NumeroDocumento: numeroDocumento,
		Estado:          true,
		IsDeleted:       false,
		CreatedAt:       now,
		UpdatedAt:       now,
	}

	if err := uc.userRepo.Create(ctx, user); err != nil {
		return nil, apperrors.ErrDatabase.WithError(err)
	}

	return user, nil
}