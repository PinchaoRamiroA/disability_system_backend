package postgres

import (
	"context"
	"errors"

	authdomain "disability_system_backend/internal/modules/auth/domain"
	usuariosmodels "disability_system_backend/internal/modules/usuarios/adapters/postgres/models"
	apperrors "disability_system_backend/internal/shared/errors"
	"gorm.io/gorm"
)

type RoleRepository struct {
	db *gorm.DB
}

func NewRoleRepository(db *gorm.DB) *RoleRepository {
	return &RoleRepository{db: db}
}

func (r *RoleRepository) FindByID(ctx context.Context, id uint64) (*authdomain.Role, error) {
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

func (r *RoleRepository) FindByName(ctx context.Context, name string) (*authdomain.Role, error) {
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

func toDomainRole(m *usuariosmodels.RolModel) *authdomain.Role {
	return &authdomain.Role{
		ID:       m.IDRol,
		Nombre:   m.Nombre,
		Permisos: m.GetPermisos(),
	}
}