package usecase

import (
	"context"
	"time"

	"disability_system_backend/internal/modules/auth/domain"
	"disability_system_backend/internal/modules/auth/ports"
	apperrors "disability_system_backend/internal/shared/errors"
)

const DefaultRoleName = "Empleado"

type RegisterUseCase struct {
	userRepo       ports.UserRepository
	passwordHasher ports.PasswordHasher
	roleRepo       ports.RoleRepository
}

func NewRegisterUseCase(
	userRepo ports.UserRepository,
	passwordHasher ports.PasswordHasher,
	roleRepo ...ports.RoleRepository,
) *RegisterUseCase {
	var rRepo ports.RoleRepository
	if len(roleRepo) > 0 {
		rRepo = roleRepo[0]
	}
	return &RegisterUseCase{
		userRepo:       userRepo,
		passwordHasher: passwordHasher,
		roleRepo:       rRepo,
	}
}

func (uc *RegisterUseCase) Execute(
	ctx context.Context,
	nombre, email, password, numeroDocumento string,
) (*domain.User, error) {
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

	var roleID uint64 = 4
	if uc.roleRepo != nil {
		role, err := uc.roleRepo.FindByName(ctx, DefaultRoleName)
		if err != nil {
			return nil, apperrors.ErrRolNotFound.WithMessage("rol por defecto no encontrado").WithError(err)
		}
		if role == nil {
			return nil, apperrors.ErrRolNotFound.WithMessage("rol por defecto no encontrado")
		}
		roleID = role.ID
	}

	hashedPassword, err := uc.passwordHasher.Hash(password)
	if err != nil {
		return nil, apperrors.ErrHashPassword.WithError(err)
	}

	now := time.Now()
	user := &domain.User{
		IDRol:           roleID,
		Nombre:          nombre,
		Correo:          email,
		PasswordHash:    hashedPassword,
		NumeroDocumento: numeroDocumento,
		Estado:          true,
		IsDeleted:       false,
		CreatedAt:       now,
	}

	if err := uc.userRepo.Create(ctx, user); err != nil {
		return nil, apperrors.ErrDatabase.WithError(err)
	}

	return user, nil
}