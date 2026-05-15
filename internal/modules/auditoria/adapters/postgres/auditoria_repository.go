package postgres

import (
	"context"

	"disability_system_backend/internal/modules/auditoria/adapters/postgres/models"
	"disability_system_backend/internal/modules/auditoria/domain"
	"disability_system_backend/internal/modules/auditoria/dto"

	"gorm.io/gorm"
)

type AuditoriaRepository struct {
	db *gorm.DB
}

func NewAuditoriaRepository(db *gorm.DB) *AuditoriaRepository {
	return &AuditoriaRepository{db: db}
}

func (r *AuditoriaRepository) Create(ctx context.Context, auditoria *domain.Auditoria) error {
	model := toModel(auditoria)
	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		return err
	}
	*auditoria = *toDomain(model)
	return nil
}

func (r *AuditoriaRepository) List(ctx context.Context, filters dto.ListarAuditoriaQuery) ([]domain.Auditoria, int64, error) {
	var modelsList []models.AuditoriaModel
	var total int64

	query := r.db.WithContext(ctx).
		Table("auditoria").
		Select("auditoria.*, usuario.nombre as usuario_nombre").
		Joins("LEFT JOIN usuario ON usuario.id_usuario = auditoria.id_usuario")

	if filters.IDUsuario != nil {
		query = query.Where("auditoria.id_usuario = ?", *filters.IDUsuario)
	}
	if filters.IDIncapacidad != nil {
		query = query.Where("auditoria.id_incapacidad = ?", *filters.IDIncapacidad)
	}
	if filters.TipoAccion != "" {
		query = query.Where("auditoria.tipo_accion = ?", filters.TipoAccion)
	}
	if filters.Modulo != "" {
		query = query.Where("auditoria.modulo = ?", filters.Modulo)
	}
	if filters.FechaInicio != "" {
		query = query.Where("auditoria.created_at >= ?", filters.FechaInicio)
	}
	if filters.FechaFin != "" {
		query = query.Where("auditoria.created_at <= ?", filters.FechaFin)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (filters.Page - 1) * filters.Limit
	err := query.Order("auditoria.created_at DESC").
		Limit(filters.Limit).
		Offset(offset).
		Find(&modelsList).Error

	if err != nil {
		return nil, 0, err
	}

	items := make([]domain.Auditoria, len(modelsList))
	for i, m := range modelsList {
		items[i] = *toDomain(&m)
	}

	return items, total, nil
}

func toModel(a *domain.Auditoria) *models.AuditoriaModel {
	return &models.AuditoriaModel{
		IDAuditoria:    a.ID,
		IDUsuario:      a.IDUsuario,
		IDIncapacidad:  a.IDIncapacidad,
		TipoAccion:     a.TipoAccion,
		Modulo:         a.Modulo,
		Descripcion:    a.Descripcion,
		CambioAnterior: a.CambioAnterior,
		CambioNuevo:    a.CambioNuevo,
		CreatedAt:      a.CreatedAt,
	}
}

func toDomain(m *models.AuditoriaModel) *domain.Auditoria {
	return &domain.Auditoria{
		ID:             m.IDAuditoria,
		IDUsuario:      m.IDUsuario,
		UsuarioNombre:  m.UsuarioNombre,
		IDIncapacidad:  m.IDIncapacidad,
		TipoAccion:     m.TipoAccion,
		Modulo:         m.Modulo,
		Descripcion:    m.Descripcion,
		CambioAnterior: m.CambioAnterior,
		CambioNuevo:    m.CambioNuevo,
		CreatedAt:      m.CreatedAt,
	}
}
