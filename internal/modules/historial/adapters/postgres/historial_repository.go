package postgres

import (
	"context"
	"errors"

	histmodels "disability_system_backend/internal/modules/historial/adapters/postgres/models"
	"disability_system_backend/internal/modules/historial/domain"
	"disability_system_backend/internal/modules/historial/ports"
	apperrors "disability_system_backend/internal/shared/errors"

	"gorm.io/gorm"
)

type HistorialRepository struct {
	db *gorm.DB
}

func NewHistorialRepository(db *gorm.DB) *HistorialRepository {
	return &HistorialRepository{db: db}
}

func (r *HistorialRepository) Create(ctx context.Context, historial *domain.Historial) error {
	model := &histmodels.HistorialModel{
		IDIncapacidad:   historial.IDIncapacidad,
		IDTipoHistorial: historial.IDTipoHistorial,
		Descripcion:     historial.Descripcion,
		Fecha:           historial.Fecha,
		GestorID:        historial.GestorID,
	}
	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		return err
	}
	historial.IDHistorial = model.IDHistorial
	return nil
}

func (r *HistorialRepository) List(ctx context.Context, incapacidadID uint64, tipoID *uint64, page, limit int) ([]domain.Historial, int64, error) {
	var models []histmodels.HistorialModel
	query := r.db.WithContext(ctx).Model(&histmodels.HistorialModel{}).Where("id_incapacidad = ?", incapacidadID)

	if tipoID != nil {
		query = query.Where("id_tipo_historial = ?", *tipoID)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	err := query.
		Order("fecha DESC").
		Limit(limit).
		Offset((page - 1) * limit).
		Find(&models).Error
	if err != nil {
		return nil, 0, err
	}

	items := make([]domain.Historial, 0, len(models))
	for i := range models {
		items = append(items, domain.Historial{
			IDHistorial:     models[i].IDHistorial,
			IDIncapacidad:   models[i].IDIncapacidad,
			IDTipoHistorial: models[i].IDTipoHistorial,
			Descripcion:     models[i].Descripcion,
			Fecha:           models[i].Fecha,
			GestorID:        models[i].GestorID,
			CreatedAt:       models[i].CreatedAt,
		})
	}
	return items, total, nil
}

func (r *HistorialRepository) FindByID(ctx context.Context, id uint64) (*domain.Historial, error) {
	var model histmodels.HistorialModel
	err := r.db.WithContext(ctx).Where("id_historial = ?", id).First(&model).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.ErrNotFound.WithMessage("historial no encontrado").WithError(err)
		}
		return nil, err
	}
	return &domain.Historial{
		IDHistorial:     model.IDHistorial,
		IDIncapacidad:   model.IDIncapacidad,
		IDTipoHistorial: model.IDTipoHistorial,
		Descripcion:     model.Descripcion,
		Fecha:           model.Fecha,
		GestorID:        model.GestorID,
		CreatedAt:       model.CreatedAt,
	}, nil
}

func (r *HistorialRepository) FindTipoByID(ctx context.Context, id uint64) (*domain.TipoHistorial, error) {
	var model histmodels.TipoHistorialModel
	err := r.db.WithContext(ctx).Where("id_tipo_historial = ?", id).First(&model).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.ErrNotFound.WithMessage("tipo de historial no encontrado").WithError(err)
		}
		return nil, err
	}
	return &domain.TipoHistorial{
		IDTipoHistorial: model.IDTipoHistorial,
		Nombre:          model.Nombre,
		Descripcion:     model.Descripcion,
		CreatedAt:       model.CreatedAt,
		UpdatedAt:       model.UpdatedAt,
	}, nil
}

var _ ports.HistorialRepository = (*HistorialRepository)(nil)