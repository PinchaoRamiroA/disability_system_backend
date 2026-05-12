package postgres

import (
	"context"
	"encoding/json"
	"errors"

	incmodels "disability_system_backend/internal/modules/incapacidades/adapters/postgres/models"
	"disability_system_backend/internal/modules/incapacidades/domain"
	"disability_system_backend/internal/modules/incapacidades/ports"
	usuariosmodels "disability_system_backend/internal/modules/usuarios/adapters/postgres/models"
	apperrors "disability_system_backend/internal/shared/errors"

	"gorm.io/gorm"
)

type IncapacidadRepository struct {
	db *gorm.DB
}

func NewIncapacidadRepository(db *gorm.DB) *IncapacidadRepository {
	return &IncapacidadRepository{db: db}
}

func (r *IncapacidadRepository) Create(ctx context.Context, incapacidad *domain.Incapacidad) error {
	model := toIncapacidadModel(incapacidad)
	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		return err
	}
	created, err := r.FindByID(ctx, model.IDIncapacidad)
	if err != nil {
		return err
	}
	*incapacidad = *created
	return nil
}

func (r *IncapacidadRepository) FindByID(ctx context.Context, id uint64) (*domain.Incapacidad, error) {
	var model incmodels.IncapacidadModel
	err := r.db.WithContext(ctx).
		Preload("Estado").
		Preload("Tipo").
		Preload("Entidad").
		Where("id_incapacidad = ? AND is_deleted = false", id).
		First(&model).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.ErrIncapacidadNotFound.WithError(err)
		}
		return nil, err
	}
	return toIncapacidadDomain(&model), nil
}

func (r *IncapacidadRepository) List(ctx context.Context, filters ports.IncapacidadFilters) ([]domain.Incapacidad, int64, error) {
	var models []incmodels.IncapacidadModel
	query := r.db.WithContext(ctx).Model(&incmodels.IncapacidadModel{})

	if !filters.IncludeDeleted {
		query = query.Where("incapacidad.is_deleted = false")
	}
	if filters.UserID != nil {
		query = query.Where("incapacidad.id_usuario = ?", *filters.UserID)
	}
	if filters.EstadoID != nil {
		query = query.Where("incapacidad.id_estado = ?", *filters.EstadoID)
	}
	if filters.TipoID != nil {
		query = query.Where("incapacidad.id_tipo = ?", *filters.TipoID)
	}
	if filters.EntidadID != nil {
		query = query.Where("incapacidad.id_entidad = ?", *filters.EntidadID)
	}
	if filters.Origen != "" {
		query = query.Where("incapacidad.origen = ?", filters.Origen)
	}
	if filters.CanalRecepcion != "" {
		query = query.Where("incapacidad.canal_recepcion = ?", filters.CanalRecepcion)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	page := filters.Page
	if page < 1 {
		page = 1
	}
	limit := filters.Limit
	if limit < 1 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	err := query.
		Preload("Estado").
		Preload("Tipo").
		Preload("Entidad").
		Order("incapacidad.created_at DESC").
		Limit(limit).
		Offset((page - 1) * limit).
		Find(&models).Error
	if err != nil {
		return nil, 0, err
	}

	items := make([]domain.Incapacidad, 0, len(models))
	for i := range models {
		items = append(items, *toIncapacidadDomain(&models[i]))
	}
	return items, total, nil
}

func (r *IncapacidadRepository) Update(ctx context.Context, incapacidad *domain.Incapacidad) error {
	model := toIncapacidadModel(incapacidad)
	err := r.db.WithContext(ctx).
		Model(&incmodels.IncapacidadModel{}).
		Where("id_incapacidad = ? AND is_deleted = false", incapacidad.IDIncapacidad).
		Updates(model).Error
	if err != nil {
		return err
	}
	updated, err := r.FindByID(ctx, incapacidad.IDIncapacidad)
	if err != nil {
		return err
	}
	*incapacidad = *updated
	return nil
}

func (r *IncapacidadRepository) SoftDelete(ctx context.Context, id uint64) error {
	result := r.db.WithContext(ctx).
		Model(&incmodels.IncapacidadModel{}).
		Where("id_incapacidad = ? AND is_deleted = false", id).
		Update("is_deleted", true)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return apperrors.ErrIncapacidadNotFound
	}
	return nil
}

func (r *IncapacidadRepository) ExistsUsuario(ctx context.Context, id uint64) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&usuariosmodels.UsuarioModel{}).
		Where("id_usuario = ? AND is_deleted = false", id).
		Count(&count).Error
	return count > 0, err
}

func (r *IncapacidadRepository) FindEstadoByID(ctx context.Context, id uint64) (*domain.EstadoIncapacidad, error) {
	var model incmodels.EstadoIncapacidadModel
	err := r.db.WithContext(ctx).Where("id_estado = ?", id).First(&model).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.ErrNotFound.WithMessage("estado de incapacidad no encontrado").WithError(err)
		}
		return nil, err
	}
	return toEstadoDomain(&model), nil
}

func (r *IncapacidadRepository) FindEstadoByName(ctx context.Context, name string) (*domain.EstadoIncapacidad, error) {
	var model incmodels.EstadoIncapacidadModel
	err := r.db.WithContext(ctx).Where("nombre = ?", name).First(&model).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.ErrNotFound.WithMessage("estado de incapacidad no encontrado").WithError(err)
		}
		return nil, err
	}
	return toEstadoDomain(&model), nil
}

func (r *IncapacidadRepository) FindTipoByID(ctx context.Context, id uint64) (*domain.TipoIncapacidad, error) {
	var model incmodels.TipoIncapacidadModel
	err := r.db.WithContext(ctx).Where("id_tipo = ?", id).First(&model).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.ErrNotFound.WithMessage("tipo de incapacidad no encontrado").WithError(err)
		}
		return nil, err
	}
	return toTipoDomain(&model), nil
}

func (r *IncapacidadRepository) FindEntidadByID(ctx context.Context, id uint64) (*domain.Entidad, error) {
	var model incmodels.EntidadModel
	err := r.db.WithContext(ctx).Where("id_entidad = ?", id).First(&model).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.ErrNotFound.WithMessage("entidad no encontrada").WithError(err)
		}
		return nil, err
	}
	return toEntidadDomain(&model), nil
}

func (r *IncapacidadRepository) ListEstados(ctx context.Context) ([]domain.EstadoIncapacidad, error) {
	var models []incmodels.EstadoIncapacidadModel
	if err := r.db.WithContext(ctx).Order("nombre ASC").Find(&models).Error; err != nil {
		return nil, err
	}
	items := make([]domain.EstadoIncapacidad, 0, len(models))
	for i := range models {
		items = append(items, *toEstadoDomain(&models[i]))
	}
	return items, nil
}

func (r *IncapacidadRepository) ListTipos(ctx context.Context) ([]domain.TipoIncapacidad, error) {
	var models []incmodels.TipoIncapacidadModel
	if err := r.db.WithContext(ctx).Order("nombre ASC").Find(&models).Error; err != nil {
		return nil, err
	}
	items := make([]domain.TipoIncapacidad, 0, len(models))
	for i := range models {
		items = append(items, *toTipoDomain(&models[i]))
	}
	return items, nil
}

func (r *IncapacidadRepository) ListEntidades(ctx context.Context) ([]domain.Entidad, error) {
	var models []incmodels.EntidadModel
	if err := r.db.WithContext(ctx).Order("nombre ASC").Find(&models).Error; err != nil {
		return nil, err
	}
	items := make([]domain.Entidad, 0, len(models))
	for i := range models {
		items = append(items, *toEntidadDomain(&models[i]))
	}
	return items, nil
}

func (r *IncapacidadRepository) ListEstadosDocumento(ctx context.Context) ([]domain.EstadoDocumento, error) {
	var models []incmodels.EstadoDocumentoModel
	if err := r.db.WithContext(ctx).Order("nombre ASC").Find(&models).Error; err != nil {
		return nil, err
	}
	items := make([]domain.EstadoDocumento, 0, len(models))
	for i := range models {
		items = append(items, *toEstadoDocumentoDomain(&models[i]))
	}
	return items, nil
}

func (r *IncapacidadRepository) ListTiposDocumento(ctx context.Context) ([]domain.TipoDocumento, error) {
	var models []incmodels.TipoDocumentoModel
	if err := r.db.WithContext(ctx).Order("nombre ASC").Find(&models).Error; err != nil {
		return nil, err
	}
	items := make([]domain.TipoDocumento, 0, len(models))
	for i := range models {
		items = append(items, *toTipoDocumentoDomain(&models[i]))
	}
	return items, nil
}

func (r *IncapacidadRepository) ListTiposPago(ctx context.Context) ([]domain.TipoPago, error) {
	var models []incmodels.TipoPagoModel
	if err := r.db.WithContext(ctx).Order("nombre ASC").Find(&models).Error; err != nil {
		return nil, err
	}
	items := make([]domain.TipoPago, 0, len(models))
	for i := range models {
		items = append(items, *toTipoPagoDomain(&models[i]))
	}
	return items, nil
}

func toIncapacidadModel(i *domain.Incapacidad) *incmodels.IncapacidadModel {
	return &incmodels.IncapacidadModel{
		IDIncapacidad:   i.IDIncapacidad,
		IDUsuario:       i.IDUsuario,
		IDEstado:        i.IDEstado,
		IDTipo:          i.IDTipo,
		IDEntidad:       i.IDEntidad,
		CanalRecepcion:  i.CanalRecepcion,
		Titulo:          i.Titulo,
		FechaInicio:     i.FechaInicio,
		FechaFin:        i.FechaFin,
		Origen:          i.Origen,
		FechaRadicacion: i.FechaRadicacion,
		FechaPago:       i.FechaPago,
		Observaciones:   i.Observaciones,
		CreatedBy:       i.CreatedBy,
		IsDeleted:       i.IsDeleted,
	}
}

func toIncapacidadDomain(m *incmodels.IncapacidadModel) *domain.Incapacidad {
	incapacidad := &domain.Incapacidad{
		IDIncapacidad:   m.IDIncapacidad,
		IDUsuario:       m.IDUsuario,
		IDEstado:        m.IDEstado,
		IDTipo:          m.IDTipo,
		IDEntidad:       m.IDEntidad,
		CanalRecepcion:  m.CanalRecepcion,
		Titulo:          m.Titulo,
		FechaInicio:     m.FechaInicio,
		FechaFin:        m.FechaFin,
		Origen:          m.Origen,
		FechaRadicacion: m.FechaRadicacion,
		FechaPago:       m.FechaPago,
		Observaciones:   m.Observaciones,
		CreatedBy:       m.CreatedBy,
		CreatedAt:       m.CreatedAt,
		UpdatedAt:       m.UpdatedAt,
		IsDeleted:       m.IsDeleted,
	}
	if m.Estado.IDEstado != 0 {
		incapacidad.Estado = toEstadoDomain(&m.Estado)
	}
	if m.Tipo.IDTipo != 0 {
		incapacidad.Tipo = toTipoDomain(&m.Tipo)
	}
	if m.Entidad.IDEntidad != 0 {
		incapacidad.Entidad = toEntidadDomain(&m.Entidad)
	}
	return incapacidad
}

func toEstadoDomain(m *incmodels.EstadoIncapacidadModel) *domain.EstadoIncapacidad {
	return &domain.EstadoIncapacidad{
		IDEstado:          m.IDEstado,
		Nombre:            m.Nombre,
		Descripcion:       m.Descripcion,
		PermiteTransicion: m.PermiteTransicion,
		CreatedAt:         m.CreatedAt,
		UpdatedAt:         m.UpdatedAt,
	}
}

func toTipoDomain(m *incmodels.TipoIncapacidadModel) *domain.TipoIncapacidad {
	var docs []string
	if m.DocumentosRequeridos != nil {
		_ = json.Unmarshal(m.DocumentosRequeridos, &docs)
	}
	return &domain.TipoIncapacidad{
		IDTipo:               m.IDTipo,
		Nombre:               m.Nombre,
		DocumentosRequeridos: docs,
		CreatedAt:            m.CreatedAt,
		UpdatedAt:            m.UpdatedAt,
	}
}

func toEntidadDomain(m *incmodels.EntidadModel) *domain.Entidad {
	return &domain.Entidad{
		IDEntidad:              m.IDEntidad,
		Nombre:                 m.Nombre,
		Tipo:                   m.Tipo,
		PlazoTranscripcionDias: m.PlazoTranscripcionDias,
		TiempoMaximoPagoDias:   m.TiempoMaximoPagoDias,
		CanalAtencion:          m.CanalAtencion,
		CanalesAtencion:        m.GetCanalesAtencion(),
		RequiereTranscripcion:  m.RequiereTranscripcion,
		CreatedAt:              m.CreatedAt,
		UpdatedAt:              m.UpdatedAt,
	}
}

func toEstadoDocumentoDomain(m *incmodels.EstadoDocumentoModel) *domain.EstadoDocumento {
	return &domain.EstadoDocumento{
		IDEstadoDocumento: m.IDEstadoDocumento,
		Nombre:            m.Nombre,
		Descripcion:       m.Descripcion,
		Color:             m.Color,
		CreatedAt:         m.CreatedAt,
		UpdatedAt:         m.UpdatedAt,
	}
}

func toTipoDocumentoDomain(m *incmodels.TipoDocumentoModel) *domain.TipoDocumento {
	return &domain.TipoDocumento{
		IDTipoDocumento: m.IDTipoDocumento,
		Nombre:          m.Nombre,
		Descripcion:     m.Descripcion,
		Requerido:       m.Requerido,
		CreatedAt:       m.CreatedAt,
		UpdatedAt:       m.UpdatedAt,
	}
}

func toTipoPagoDomain(m *incmodels.TipoPagoModel) *domain.TipoPago {
	return &domain.TipoPago{
		IDTipoPago:  m.IDTipoPago,
		Nombre:      m.Nombre,
		Descripcion: m.Descripcion,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.CreatedAt,
	}
}

func (r *IncapacidadRepository) FindTiposDocumentoByNombre(ctx context.Context, nombres []string) ([]domain.TipoDocumento, error) {
	var models []incmodels.TipoDocumentoModel
	query := r.db.WithContext(ctx).Where("nombre IN (?)", nombres)
	if err := query.Find(&models).Error; err != nil {
		return nil, err
	}
	items := make([]domain.TipoDocumento, 0, len(models))
	for i := range models {
		items = append(items, *toTipoDocumentoDomain(&models[i]))
	}
	return items, nil
}
