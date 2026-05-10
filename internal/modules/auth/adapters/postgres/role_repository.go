package postgres

import (
	"context"
	"errors"

	usuariosmodels "disability_system_backend/internal/modules/usuarios/adapters/postgres/models"
	usuariosdomain "disability_system_backend/internal/modules/usuarios/domain"
	apperrors "disability_system_backend/internal/shared/errors"
	"gorm.io/gorm"
)

type RoleRepository struct {
	db *gorm.DB
}

func NewRoleRepository(db *gorm.DB) *RoleRepository {
	return &RoleRepository{db: db}
}

func (r *RoleRepository) FindByID(ctx context.Context, id uint64) (*usuariosdomain.Rol, error) {
	var model usuariosmodels.RolModel
	err := r.db.WithContext(ctx).Where("id_rol = ?", id).First(&model).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.ErrNotFound.WithMessage("rol no encontrado").WithError(err)
		}
		return nil, err
	}
	return toDomainRole(&model), nil
}

func (r *RoleRepository) FindByName(ctx context.Context, name string) (*usuariosdomain.Rol, error) {
	var model usuariosmodels.RolModel
	err := r.db.WithContext(ctx).Where("nombre = ?", name).First(&model).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.ErrNotFound.WithMessage("rol no encontrado").WithError(err)
		}
		return nil, err
	}
	return toDomainRole(&model), nil
}

func toDomainRole(m *usuariosmodels.RolModel) *usuariosdomain.Rol {
	return &usuariosdomain.Rol{
		ID:        m.IDRol,
		Nombre:    m.Nombre,
		Permisos:  m.GetPermisos(),
		IsDeleted: m.IsDeleted,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}

var _ RoleRepositoryI = (*RoleRepository)(nil)

type RoleRepositoryI interface {
	FindByID(ctx context.Context, id uint64) (*usuariosdomain.Rol, error)
	FindByName(ctx context.Context, name string) (*usuariosdomain.Rol, error)
}