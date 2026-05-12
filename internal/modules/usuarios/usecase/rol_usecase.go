package usecase

import (
	"context"

	"disability_system_backend/internal/modules/usuarios/domain"
	apperrors "disability_system_backend/internal/shared/errors"
)

type RolUseCase struct {
	repo interface {
		FindByID(ctx context.Context, id uint64) (*domain.Rol, error)
		FindByName(ctx context.Context, name string) (*domain.Rol, error)
		FindAll(ctx context.Context, page, limit int) ([]domain.Rol, int64, error)
		Create(ctx context.Context, rol *domain.Rol) error
		Update(ctx context.Context, rol *domain.Rol) error
		Delete(ctx context.Context, id uint64) error
	}
}

func NewRolUseCase(repo interface {
	FindByID(ctx context.Context, id uint64) (*domain.Rol, error)
	FindByName(ctx context.Context, name string) (*domain.Rol, error)
	FindAll(ctx context.Context, page, limit int) ([]domain.Rol, int64, error)
	Create(ctx context.Context, rol *domain.Rol) error
	Update(ctx context.Context, rol *domain.Rol) error
	Delete(ctx context.Context, id uint64) error
}) *RolUseCase {
	return &RolUseCase{repo: repo}
}

func (uc *RolUseCase) Listar(ctx context.Context, page, limit int) ([]domain.Rol, int64, error) {
	return uc.repo.FindAll(ctx, page, limit)
}

func (uc *RolUseCase) Obtener(ctx context.Context, id uint64) (*domain.Rol, error) {
	return uc.repo.FindByID(ctx, id)
}

func (uc *RolUseCase) Crear(ctx context.Context, nombre string, permisos []string) (*domain.Rol, error) {
	_, err := uc.repo.FindByName(ctx, nombre)
	if err == nil {
		return nil, apperrors.ErrConflict.WithMessage("El rol ya existe")
	}

	rol := &domain.Rol{
		Nombre:   nombre,
		Permisos: permisos,
	}

	if err := uc.repo.Create(ctx, rol); err != nil {
		return nil, err
	}

	return rol, nil
}

func (uc *RolUseCase) Actualizar(ctx context.Context, id uint64, nombre string, permisos []string) (*domain.Rol, error) {
	rol, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if nombre != rol.Nombre {
		_, err := uc.repo.FindByName(ctx, nombre)
		if err == nil {
			return nil, apperrors.ErrConflict.WithMessage("El nombre del rol ya existe")
		}
	}

	rol.Nombre = nombre
	rol.Permisos = permisos

	if err := uc.repo.Update(ctx, rol); err != nil {
		return nil, err
	}

	return rol, nil
}

func (uc *RolUseCase) Eliminar(ctx context.Context, id uint64) error {
	_, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	return uc.repo.Delete(ctx, id)
}
