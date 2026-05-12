package postgres

import (
	"context"
	"errors"

	incmodels "disability_system_backend/internal/modules/incapacidades/adapters/postgres/models"
	notifmodels "disability_system_backend/internal/modules/notificaciones/adapters/postgres/models"
	"disability_system_backend/internal/modules/notificaciones/domain"
	"disability_system_backend/internal/modules/notificaciones/ports"
	usuariosmodels "disability_system_backend/internal/modules/usuarios/adapters/postgres/models"
	apperrors "disability_system_backend/internal/shared/errors"

	"gorm.io/gorm"
)

type NotificacionRepository struct {
	db *gorm.DB
}

func NewNotificacionRepository(db *gorm.DB) *NotificacionRepository {
	return &NotificacionRepository{db: db}
}

func (r *NotificacionRepository) Create(ctx context.Context, notificacion *domain.Notificacion) error {
	model := toNotificacionModel(notificacion)
	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		return err
	}
	created, err := r.FindByID(ctx, model.IDNotificacion)
	if err != nil {
		return err
	}
	*notificacion = *created
	return nil
}

func (r *NotificacionRepository) FindByID(ctx context.Context, id uint64) (*domain.Notificacion, error) {
	var model notifmodels.NotificacionModel
	err := r.db.WithContext(ctx).Where("id_notificacion = ?", id).First(&model).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.ErrNotFound.WithMessage("notificación no encontrada").WithError(err)
		}
		return nil, err
	}
	return toNotificacionDomain(&model), nil
}

func (r *NotificacionRepository) List(ctx context.Context, filters ports.NotificacionFilters) ([]domain.Notificacion, int64, error) {
	var models []notifmodels.NotificacionModel
	query := r.db.WithContext(ctx).Model(&notifmodels.NotificacionModel{})

	if filters.IDUsuario != nil {
		query = query.Where("id_usuario = ?", *filters.IDUsuario)
	}
	if filters.IDIncapacidad != nil {
		query = query.Where("id_incapacidad = ?", *filters.IDIncapacidad)
	}
	if filters.TipoNotificacion != "" {
		query = query.Where("tipo_notificacion = ?", filters.TipoNotificacion)
	}
	if filters.Leida != nil {
		query = query.Where("leida = ?", *filters.Leida)
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

	items := make([]domain.Notificacion, 0, len(models))
	for i := range models {
		items = append(items, *toNotificacionDomain(&models[i]))
	}
	return items, total, nil
}

func (r *NotificacionRepository) CountUnread(ctx context.Context, userID uint64) (int64, error) {
	var total int64
	err := r.db.WithContext(ctx).
		Model(&notifmodels.NotificacionModel{}).
		Where("id_usuario = ? AND leida = false", userID).
		Count(&total).Error
	return total, err
}

func (r *NotificacionRepository) MarkAsRead(ctx context.Context, id uint64) error {
	result := r.db.WithContext(ctx).
		Model(&notifmodels.NotificacionModel{}).
		Where("id_notificacion = ?", id).
		Update("leida", true)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return apperrors.ErrNotFound.WithMessage("notificación no encontrada")
	}
	return nil
}

func (r *NotificacionRepository) MarkAllAsRead(ctx context.Context, userID uint64) error {
	return r.db.WithContext(ctx).
		Model(&notifmodels.NotificacionModel{}).
		Where("id_usuario = ? AND leida = false", userID).
		Update("leida", true).Error
}

func (r *NotificacionRepository) Delete(ctx context.Context, id uint64) error {
	result := r.db.WithContext(ctx).Delete(&notifmodels.NotificacionModel{}, "id_notificacion = ?", id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return apperrors.ErrNotFound.WithMessage("notificación no encontrada")
	}
	return nil
}

func (r *NotificacionRepository) UserExists(ctx context.Context, id uint64) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&usuariosmodels.UsuarioModel{}).
		Where("id_usuario = ? AND is_deleted = false", id).
		Count(&count).Error
	return count > 0, err
}

func (r *NotificacionRepository) IncapacidadExists(ctx context.Context, id uint64) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&incmodels.IncapacidadModel{}).
		Where("id_incapacidad = ? AND is_deleted = false", id).
		Count(&count).Error
	return count > 0, err
}

func toNotificacionModel(n *domain.Notificacion) *notifmodels.NotificacionModel {
	return &notifmodels.NotificacionModel{
		IDNotificacion:   n.IDNotificacion,
		IDUsuario:        n.IDUsuario,
		IDIncapacidad:    n.IDIncapacidad,
		TipoNotificacion: n.TipoNotificacion,
		Mensaje:          n.Mensaje,
		Fecha:            n.Fecha,
		Leida:            n.Leida,
	}
}

func toNotificacionDomain(m *notifmodels.NotificacionModel) *domain.Notificacion {
	return &domain.Notificacion{
		IDNotificacion:   m.IDNotificacion,
		IDUsuario:        m.IDUsuario,
		IDIncapacidad:    m.IDIncapacidad,
		TipoNotificacion: m.TipoNotificacion,
		Mensaje:          m.Mensaje,
		Fecha:            m.Fecha,
		Leida:            m.Leida,
		CreatedAt:        m.CreatedAt,
		UpdatedAt:        m.UpdatedAt,
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
