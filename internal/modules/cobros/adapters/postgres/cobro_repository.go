package postgres

import (
	"context"
	"errors"

	cobrosmodels "disability_system_backend/internal/modules/cobros/adapters/postgres/models"
	"disability_system_backend/internal/modules/cobros/domain"
	"disability_system_backend/internal/modules/cobros/ports"
	incmodels "disability_system_backend/internal/modules/incapacidades/adapters/postgres/models"
	apperrors "disability_system_backend/internal/shared/errors"

	"gorm.io/gorm"
)

type CobroRepository struct {
	db *gorm.DB
}

func NewCobroRepository(db *gorm.DB) *CobroRepository {
	return &CobroRepository{db: db}
}

func (r *CobroRepository) CreatePago(ctx context.Context, pago *domain.Pago) error {
	model := toPagoModel(pago)
	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		return err
	}
	created, err := r.FindPagoByID(ctx, model.IDPago)
	if err != nil {
		return err
	}
	*pago = *created
	return nil
}

func (r *CobroRepository) FindPagoByID(ctx context.Context, id uint64) (*domain.Pago, error) {
	var model cobrosmodels.PagoModel
	err := r.db.WithContext(ctx).
		Table("pago").
		Select("pago.*, entidad.nombre as nombre_entidad").
		Joins("left join entidad on entidad.id_entidad = pago.id_entidad").
		Where("pago.id_pago = ? AND pago.is_deleted = false", id).
		First(&model).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.ErrNotFound.WithMessage("pago no encontrado").WithError(err)
		}
		return nil, err
	}
	return toPagoDomain(&model), nil
}

func (r *CobroRepository) ListPagos(ctx context.Context, filters ports.PagoFilters) ([]domain.Pago, int64, error) {
	var models []cobrosmodels.PagoModel
	query := r.db.WithContext(ctx).
		Table("pago").
		Select("pago.*, entidad.nombre as nombre_entidad").
		Joins("left join entidad on entidad.id_entidad = pago.id_entidad").
		Where("pago.is_deleted = false")

	if filters.UserID != nil {
		query = query.Joins("INNER JOIN incapacidad ON incapacidad.id_incapacidad = pago.id_incapacidad").
			Where("incapacidad.id_usuario = ?", *filters.UserID)
	}

	if filters.IDIncapacidad != nil {
		query = query.Where("pago.id_incapacidad = ?", *filters.IDIncapacidad)
	}
	if filters.IDEntidad != nil {
		query = query.Where("pago.id_entidad = ?", *filters.IDEntidad)
	}
	if filters.TipoPago != "" {
		query = query.Where("pago.tipo_pago = ?", filters.TipoPago)
	}
	if filters.EstadoPago != "" {
		query = query.Where("pago.estado_pago = ?", filters.EstadoPago)
	}
	if filters.Conciliado != nil {
		query = query.Where("pago.conciliado = ?", *filters.Conciliado)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	page, limit := normalizePagination(filters.Page, filters.Limit)
	err := query.Order("pago.created_at DESC").
		Limit(limit).
		Offset((page - 1) * limit).
		Find(&models).Error
	if err != nil {
		return nil, 0, err
	}

	items := make([]domain.Pago, 0, len(models))
	for i := range models {
		items = append(items, *toPagoDomain(&models[i]))
	}
	return items, total, nil
}

func (r *CobroRepository) UpdatePago(ctx context.Context, pago *domain.Pago) error {
	model := toPagoModel(pago)
	err := r.db.WithContext(ctx).
		Model(&cobrosmodels.PagoModel{}).
		Where("id_pago = ? AND is_deleted = false", pago.IDPago).
		Updates(model).Error
	if err != nil {
		return err
	}
	updated, err := r.FindPagoByID(ctx, pago.IDPago)
	if err != nil {
		return err
	}
	*pago = *updated
	return nil
}

func (r *CobroRepository) SoftDeletePago(ctx context.Context, id uint64) error {
	result := r.db.WithContext(ctx).
		Model(&cobrosmodels.PagoModel{}).
		Where("id_pago = ? AND is_deleted = false", id).
		Update("is_deleted", true)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return apperrors.ErrNotFound.WithMessage("pago no encontrado")
	}
	return nil
}

func (r *CobroRepository) CreateSeguimiento(ctx context.Context, seguimiento *domain.SeguimientoCobro) error {
	model := toSeguimientoModel(seguimiento)
	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		return err
	}
	created, err := r.FindSeguimientoByID(ctx, model.IDSeguimiento)
	if err != nil {
		return err
	}
	*seguimiento = *created
	return nil
}

func (r *CobroRepository) FindSeguimientoByID(ctx context.Context, id uint64) (*domain.SeguimientoCobro, error) {
	var model cobrosmodels.SeguimientoCobroModel
	err := r.db.WithContext(ctx).Where("id_seguimiento = ?", id).First(&model).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.ErrNotFound.WithMessage("seguimiento de cobro no encontrado").WithError(err)
		}
		return nil, err
	}
	return toSeguimientoDomain(&model), nil
}

func (r *CobroRepository) ListSeguimientos(ctx context.Context, filters ports.SeguimientoFilters) ([]domain.SeguimientoCobro, int64, error) {
	var models []cobrosmodels.SeguimientoCobroModel
	query := r.db.WithContext(ctx).Model(&cobrosmodels.SeguimientoCobroModel{})

	if filters.IDIncapacidad != nil {
		query = query.Where("id_incapacidad = ?", *filters.IDIncapacidad)
	}
	if filters.TipoSeguimiento != "" {
		query = query.Where("tipo_seguimiento = ?", filters.TipoSeguimiento)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	page, limit := normalizePagination(filters.Page, filters.Limit)
	err := query.Order("fecha DESC").
		Limit(limit).
		Offset((page - 1) * limit).
		Find(&models).Error
	if err != nil {
		return nil, 0, err
	}

	items := make([]domain.SeguimientoCobro, 0, len(models))
	for i := range models {
		items = append(items, *toSeguimientoDomain(&models[i]))
	}
	return items, total, nil
}

func (r *CobroRepository) UpdateSeguimiento(ctx context.Context, seguimiento *domain.SeguimientoCobro) error {
	model := toSeguimientoModel(seguimiento)
	err := r.db.WithContext(ctx).
		Model(&cobrosmodels.SeguimientoCobroModel{}).
		Where("id_seguimiento = ?", seguimiento.IDSeguimiento).
		Updates(model).Error
	if err != nil {
		return err
	}
	updated, err := r.FindSeguimientoByID(ctx, seguimiento.IDSeguimiento)
	if err != nil {
		return err
	}
	*seguimiento = *updated
	return nil
}

func (r *CobroRepository) IncapacidadExists(ctx context.Context, id uint64) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&incmodels.IncapacidadModel{}).
		Where("id_incapacidad = ? AND is_deleted = false", id).
		Count(&count).Error
	return count > 0, err
}

func (r *CobroRepository) EntidadExists(ctx context.Context, id uint64) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&incmodels.EntidadModel{}).
		Where("id_entidad = ?", id).
		Count(&count).Error
	return count > 0, err
}

func (r *CobroRepository) GetEntidadInfo(ctx context.Context) (map[uint64]struct{ Nombre, Tipo string }, error) {
	var entities []incmodels.EntidadModel
	if err := r.db.WithContext(ctx).Select("id_entidad, nombre, tipo").Find(&entities).Error; err != nil {
		return nil, err
	}
	info := make(map[uint64]struct{ Nombre, Tipo string }, len(entities))
	for _, e := range entities {
		info[e.IDEntidad] = struct{ Nombre, Tipo string }{
			Nombre: e.Nombre,
			Tipo:   e.Tipo,
		}
	}
	return info, nil
}

func toPagoModel(p *domain.Pago) *cobrosmodels.PagoModel {
	return &cobrosmodels.PagoModel{
		IDPago:          p.IDPago,
		IDIncapacidad:   p.IDIncapacidad,
		IDEntidad:       p.IDEntidad,
		TipoPago:        p.TipoPago,
		EstadoPago:      p.EstadoPago,
		Descripcion:     p.Descripcion,
		Valor:           p.Valor,
		FechaPago:       p.FechaPago,
		PeriodoContable: p.PeriodoContable,
		Conciliado:      p.Conciliado,
		RegistradoPor:   p.RegistradoPor,
		IsDeleted:       p.IsDeleted,
	}
}

func toPagoDomain(m *cobrosmodels.PagoModel) *domain.Pago {
	return &domain.Pago{
		IDPago:          m.IDPago,
		IDIncapacidad:   m.IDIncapacidad,
		IDEntidad:       m.IDEntidad,
		TipoPago:        m.TipoPago,
		EstadoPago:      m.EstadoPago,
		Descripcion:     m.Descripcion,
		Valor:           m.Valor,
		FechaPago:       m.FechaPago,
		PeriodoContable: m.PeriodoContable,
		Conciliado:      m.Conciliado,
		RegistradoPor:   m.RegistradoPor,
		CreatedAt:       m.CreatedAt,
		UpdatedAt:       m.UpdatedAt,
		IsDeleted:       m.IsDeleted,
		NombreEntidad:   m.NombreEntidad,
	}
}

func toSeguimientoModel(s *domain.SeguimientoCobro) *cobrosmodels.SeguimientoCobroModel {
	return &cobrosmodels.SeguimientoCobroModel{
		IDSeguimiento:   s.IDSeguimiento,
		IDIncapacidad:   s.IDIncapacidad,
		TipoSeguimiento: s.TipoSeguimiento,
		Descripcion:     s.Descripcion,
		Fecha:           s.Fecha,
		Resultado:       s.Resultado,
		GestionadoPor:   s.GestionadoPor,
	}
}

func toSeguimientoDomain(m *cobrosmodels.SeguimientoCobroModel) *domain.SeguimientoCobro {
	return &domain.SeguimientoCobro{
		IDSeguimiento:   m.IDSeguimiento,
		IDIncapacidad:   m.IDIncapacidad,
		TipoSeguimiento: m.TipoSeguimiento,
		Descripcion:     m.Descripcion,
		Fecha:           m.Fecha,
		Resultado:       m.Resultado,
		GestionadoPor:   m.GestionadoPor,
		CreatedAt:       m.CreatedAt,
		UpdatedAt:       m.UpdatedAt,
	}
}

func normalizePagination(page, limit int) (int, int) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}
	return page, limit
}

func (r *CobroRepository) GetIncapacidadesDetailed(ctx context.Context, ids []uint64) (map[uint64]ports.IncapacidadInfo, error) {
	var models []incmodels.IncapacidadModel
	err := r.db.WithContext(ctx).
		Select("id_incapacidad, titulo").
		Where("id_incapacidad IN ?", ids).
		Find(&models).Error
	if err != nil {
		return nil, err
	}

	result := make(map[uint64]ports.IncapacidadInfo, len(models))
	for _, m := range models {
		result[m.IDIncapacidad] = ports.IncapacidadInfo{
			ID:     m.IDIncapacidad,
			Titulo: m.Titulo,
		}
	}
	return result, nil
}
