package ports

import (
	"context"

	"disability_system_backend/internal/modules/usuarios/domain"
)

type UsuarioRepository interface {
	FindByID(ctx context.Context, id uint64) (*domain.Usuario, error)
	FindByEmail(ctx context.Context, email string) (*domain.Usuario, error)
	FindByDocumentNumber(ctx context.Context, docNumber string) (*domain.Usuario, error)
	FindAll(ctx context.Context, page, limit int, estado *bool, idRol *uint64, search string) ([]domain.Usuario, int64, error)
	Create(ctx context.Context, usuario *domain.Usuario) error
	Update(ctx context.Context, usuario *domain.Usuario) error
	SoftDelete(ctx context.Context, id uint64) error
	EmailExists(ctx context.Context, email string, excludeID *uint64) (bool, error)
	DocumentExists(ctx context.Context, docNumber string, excludeID *uint64) (bool, error)
}

type RolRepository interface {
	FindByID(ctx context.Context, id uint64) (*domain.Rol, error)
	FindByName(ctx context.Context, name string) (*domain.Rol, error)
	FindAll(ctx context.Context, page, limit int) ([]domain.Rol, int64, error)
	Create(ctx context.Context, rol *domain.Rol) error
	Update(ctx context.Context, rol *domain.Rol) error
	Delete(ctx context.Context, id uint64) error
}

type PermisoRepository interface {
	FindAll(ctx context.Context) ([]domain.Permiso, error)
}
