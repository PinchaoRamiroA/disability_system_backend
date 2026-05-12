package postgres

import (
	"context"
	"errors"

	incmodels "disability_system_backend/internal/modules/incapacidades/adapters/postgres/models"
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
	var model incmodels.DocumentoModel
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
	var models []incmodels.DocumentoModel
	query := r.db.WithContext(ctx).Model(&incmodels.DocumentoModel{}).Where("is_deleted = false")

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
	result := r.db.WithContext(ctx).
		Model(&incmodels.DocumentoModel{}).
		Where("id_documento = ? AND is_deleted = false", documento.IDDocumento).
		Updates(map[string]interface{}{
			"estado_documento":  documento.Estado,
			"comentario":        documento.Comentario,
			"validado_por":      documento.ValidadoPor,
			"fecha_validacion":  documento.FechaValidacion,
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
		Model(&incmodels.DocumentoModel{}).
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
		Model(&incmodels.IncapacidadModel{}).
		Where("id_incapacidad = ? AND is_deleted = false", id).
		Count(&count).Error
	return count > 0, err
}

func toDocumentoModel(d *domain.Documento) *incmodels.DocumentoModel {
	return &incmodels.DocumentoModel{
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

func toDocumentoDomain(m *incmodels.DocumentoModel) *domain.Documento {
	return &domain.Documento{
		IDDocumento:     m.IDDocumento,
		IDIncapacidad:   m.IDIncapacidad,
		Nombre:          m.Nombre,
		Tipo:            m.Tipo,
		URL:             m.URL,
		Formato:         m.Formato,
		Estado:          m.Estado,
		Comentario:      m.Comentario,
		FechaCarga:      m.FechaCarga,
		ValidadoPor:     m.ValidadoPor,
		FechaValidacion: m.FechaValidacion,
		IsDeleted:       m.IsDeleted,
		CreatedAt:       m.CreatedAt,
		UpdatedAt:       m.UpdatedAt,
	}
}

var _ incports.DocumentoRepository = (*DocumentoRepository)(nil)