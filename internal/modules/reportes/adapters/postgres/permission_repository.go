package postgres

import (
	"context"
	"errors"

	usuariosmodels "disability_system_backend/internal/modules/usuarios/adapters/postgres/models"
	apperrors "disability_system_backend/internal/shared/errors"

	"gorm.io/gorm"
)

type PermissionRepository struct {
	db *gorm.DB
}

func NewPermissionRepository(db *gorm.DB) *PermissionRepository {
	return &PermissionRepository{db: db}
}

func (r *PermissionRepository) FindPermissionsByRoleName(ctx context.Context, role string) ([]string, error) {
	var model usuariosmodels.RolModel
	err := r.db.WithContext(ctx).Where("nombre = ? AND is_deleted = false", role).First(&model).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.ErrForbidden.WithMessage("rol sin permisos configurados").WithError(err)
		}
		return nil, err
	}
	return model.GetPermisos(), nil
}