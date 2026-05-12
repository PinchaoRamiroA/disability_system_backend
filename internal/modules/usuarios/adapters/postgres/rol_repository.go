package postgres

import (
	"context"
	"errors"

	"disability_system_backend/internal/modules/usuarios/adapters/postgres/models"
	"disability_system_backend/internal/modules/usuarios/domain"
	apperrors "disability_system_backend/internal/shared/errors"

	"gorm.io/gorm"
)

type RolRepository struct {
	db *gorm.DB
}

func NewRolRepository(db *gorm.DB) *RolRepository {
	return &RolRepository{db: db}
}

func (r *RolRepository) FindByID(ctx context.Context, id uint64) (*domain.Rol, error) {
	var model models.RolModel
	err := r.db.WithContext(ctx).
		Where("id_rol = ? AND is_deleted = false", id).
		First(&model).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.ErrRolNotFound.WithError(err)
		}
		return nil, err
	}
	return toDomainRol(&model), nil
}

func (r *RolRepository) FindByName(ctx context.Context, name string) (*domain.Rol, error) {
	var model models.RolModel
	err := r.db.WithContext(ctx).
		Where("nombre = ? AND is_deleted = false", name).
		First(&model).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.ErrRolNotFound.WithError(err)
		}
		return nil, err
	}
	return toDomainRol(&model), nil
}

func (r *RolRepository) FindAll(ctx context.Context, page, limit int) ([]domain.Rol, int64, error) {
	var modelList []models.RolModel
	var total int64

	query := r.db.WithContext(ctx).Model(&models.RolModel{}).Where("is_deleted = false")

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	err = query.Offset(offset).Limit(limit).Order("nombre ASC").Find(&modelList).Error
	if err != nil {
		return nil, 0, err
	}

	roles := make([]domain.Rol, len(modelList))
	for i, m := range modelList {
		roles[i] = *toDomainRol(&m)
	}

	return roles, total, nil
}

func (r *RolRepository) Create(ctx context.Context, rol *domain.Rol) error {
	model := toModelRol(rol)
	return r.db.WithContext(ctx).Create(model).Error
}

func (r *RolRepository) Update(ctx context.Context, rol *domain.Rol) error {
	model := toModelRol(rol)
	return r.db.WithContext(ctx).Save(model).Error
}

func (r *RolRepository) Delete(ctx context.Context, id uint64) error {
	return r.db.WithContext(ctx).
		Model(&models.RolModel{}).
		Where("id_rol = ?", id).
		Update("is_deleted", true).Error
}

func toDomainRol(m *models.RolModel) *domain.Rol {
	return &domain.Rol{
		ID:        m.IDRol,
		Nombre:    m.Nombre,
		Permisos:  m.GetPermisos(),
		IsDeleted: m.IsDeleted,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}

func toModelRol(u *domain.Rol) *models.RolModel {
	model := &models.RolModel{
		IDRol:     u.ID,
		Nombre:    u.Nombre,
		IsDeleted: u.IsDeleted,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
	model.SetPermisos(u.Permisos)
	return model
}
