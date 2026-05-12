package postgres

import (
	"context"
	"errors"
	"time"

	"disability_system_backend/internal/modules/incapacidades/adapters/postgres/models"
	"disability_system_backend/internal/modules/incapacidades/domain"
	incports "disability_system_backend/internal/modules/incapacidades/ports"
	apperrors "disability_system_backend/internal/shared/errors"

	"gorm.io/gorm"
)

type DocumentoRepository struct {
	db *gorm.DB
}

func NewDocumentoRepository(db *gorm.DB) *DocumentoRepository {
	return &DocumentoRepository{db: db}
}

func (r *DocumentoRepository) Create(ctx context.Context, documento *domain.Documento) error {
	model := toDocumentoModel(documento)
	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		return err
	}
	created, err := r.FindByID(ctx, model.IDDocumento)
	if err != nil {
		return err
	}
	*documento = *created
	return nil
}

func (r *DocumentoRepository) FindByID(ctx context.Context, id uint64) (*domain.Documento, error) {
	var model models.DocumentoModel
	err := r.db.WithContext(ctx).
		Where("id_documento = ? AND is_deleted = false", id).
		First(&model).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.ErrNotFound.WithMessage("documento no encontrado").WithError(err)
		}
		return nil, err
	}
	return toDocumentoDomain(&model), nil
}

func (r *DocumentoRepository) List(ctx context.Context, incapacidadID uint64, estado, tipo string, page, limit int) ([]domain.Documento, int64, error) {
	var models []models.DocumentoModel
	query := r.db.WithContext(ctx).Model(&models.DocumentoModel{}).Where("is_deleted = false")

	if incapacidadID > 0 {
		query = query.Where("id_incapacidad = ?", incapacidadID)
	}
	if estado != "" {
		query = query.Where("estado_documento = ?", estado)
	}
	if tipo != "" {
		query = query.Where("tipo_documento = ?", tipo)
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
		Order("fecha_carga DESC").
		Limit(limit).
		Offset((page - 1) * limit).
		Find(&models).Error
	if err != nil {
		return nil, 0, err
	}

	items := make([]domain.Documento, 0, len(models))
	for i := range models {
		items = append(items, *toDocumentoDomain(&models[i]))
	}
	return items, total, nil
}

func (r *DocumentoRepository) Update(ctx context.Context, documento *domain.Documento) error {
	model := toDocumentoModel(documento)
	result := r.db.WithContext(ctx).
		Model(&models.DocumentoModel{}).
		Where("id_documento = ? AND is_deleted = false", documento.IDDocumento).
		Updates(map[string]interface{}{
			"estado_documento": documento.Estado,
			"comentario":       documento.Comentario,
			"validado_por":     documento.ValidadoPor,
			"fecha_validacion": documento.FechaValidacion,
		})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return apperrors.ErrNotFound.WithMessage("documento no encontrado")
	}
	updated, err := r.FindByID(ctx, documento.IDDocumento)
	if err != nil {
		return err
	}
	*documento = *updated
	return nil
}

func (r *DocumentoRepository) Delete(ctx context.Context, id uint64) error {
	result := r.db.WithContext(ctx).
		Model(&models.DocumentoModel{}).
		Where("id_documento = ?", id).
		Update("is_deleted", true)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return apperrors.ErrNotFound.WithMessage("documento no encontrado")
	}
	return nil
}

func (r *DocumentoRepository) ExistsIncapacidad(ctx context.Context, id uint64) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&models.IncapacidadModel{}).
		Where("id_incapacidad = ? AND is_deleted = false", id).
		Count(&count).Error
	return count > 0, err
}

func toDocumentoModel(d *domain.Documento) *models.DocumentoModel {
	return &models.DocumentoModel{
		IDDocumento:     d.IDDocumento,
		IDIncapacidad:   d.IDIncapacidad,
		Nombre:          d.Nombre,
		Tipo:            d.Tipo,
		URL:             d.URL,
		Formato:         d.Formato,
		Estado:          d.Estado,
		Comentario:      d.Comentario,
		FechaCarga:      d.FechaCarga,
		ValidadoPor:     d.ValidadoPor,
		FechaValidacion: d.FechaValidacion,
		IsDeleted:       d.IsDeleted,
	}
}

func toDocumentoDomain(m *models.DocumentoModel) *domain.Documento {
	return &domain.Documento{
		IDDocumento:     m.IDDocumento,
		IDIncapacidad:   m.IDIncapacidad,
		Nombre:           m.Nombre,
		Tipo:             m.Tipo,
		URL:              m.URL,
		Formato:          m.Formato,
		Estado:           m.Estado,
		Comentario:       m.Comentario,
		FechaCarga:       m.FechaCarga,
		ValidadoPor:      m.ValidadoPor,
		FechaValidacion:  m.FechaValidacion,
		IsDeleted:        m.IsDeleted,
		CreatedAt:        m.CreatedAt,
		UpdatedAt:        m.UpdatedAt,
	}
}

type HistorialRepository struct {
	db *gorm.DB
}

func NewHistorialRepository(db *gorm.DB) *HistorialRepository {
	return &HistorialRepository{db: db}
}

func (r *HistorialRepository) Create(ctx context.Context, historial *domain.Historial) error {
	model := &models.HistorialModel{
		IDIncapacidad:   historial.IDIncapacidad,
		IDTipoHistorial: historial.IDTipoHistorial,
		Descripcion:     historial.Descripcion,
		Fecha:           time.Now(),
		GestorID:        historial.GestorID,
	}
	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		return err
	}
	historial.IDHistorial = model.IDHistorial
	historial.Fecha = model.Fecha
	return nil
}

func (r *HistorialRepository) List(ctx context.Context, incapacidadID uint64, tipoID *uint64, page, limit int) ([]domain.Historial, int64, error) {
	var models []models.HistorialModel
	query := r.db.WithContext(ctx).Model(&models.HistorialModel{}).Where("id_incapacidad = ?", incapacidadID)

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
			IDHistorial:    models[i].IDHistorial,
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
	var model models.HistorialModel
	err := r.db.WithContext(ctx).Where("id_historial = ?", id).First(&model).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.ErrNotFound.WithMessage("historial no encontrado").WithError(err)
		}
		return nil, err
	}
	return &domain.Historial{
		IDHistorial:    model.IDHistorial,
		IDIncapacidad:  model.IDIncapacidad,
		IDTipoHistorial: model.IDTipoHistorial,
		Descripcion:    model.Descripcion,
		Fecha:          model.Fecha,
		GestorID:       model.GestorID,
		CreatedAt:      model.CreatedAt,
	}, nil
}

func (r *HistorialRepository) FindTipoByID(ctx context.Context, id uint64) (*domain.TipoHistorial, error) {
	var model models.TipoHistorialModel
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
		CreatedAt:      model.CreatedAt,
		UpdatedAt:      model.UpdatedAt,
	}, nil
}

var _ incports.DocumentoRepository = (*DocumentoRepository)(nil)
var _ incports.HistorialRepository = (*HistorialRepository)(nil)