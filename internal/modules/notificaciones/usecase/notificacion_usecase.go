package usecase

import (
	"context"
	"strings"

	"disability_system_backend/internal/modules/notificaciones/domain"
	"disability_system_backend/internal/modules/notificaciones/ports"
	shareddomain "disability_system_backend/internal/shared/domain"
	apperrors "disability_system_backend/internal/shared/errors"
)

type CrearNotificacionInput struct {
	IDUsuario        uint64
	IDIncapacidad    *uint64
	TipoNotificacion string
	Mensaje          string
}

type NotificacionUseCase struct {
	repo ports.NotificacionRepository
}

func NewNotificacionUseCase(repo ports.NotificacionRepository) *NotificacionUseCase {
	return &NotificacionUseCase{repo: repo}
}

func (uc *NotificacionUseCase) Crear(ctx context.Context, actor ports.Actor, input CrearNotificacionInput) (*domain.Notificacion, error) {
	if !actor.CanManageNotifications() {
		return nil, apperrors.ErrForbidden.WithMessage("no tienes permiso para generar notificaciones")
	}
	if strings.TrimSpace(input.Mensaje) == "" {
		return nil, apperrors.ErrValidation.WithMessage("mensaje es obligatorio")
	}
	tipo, err := normalizeTipoNotificacion(input.TipoNotificacion)
	if err != nil {
		return nil, err
	}
	ok, err := uc.repo.UserExists(ctx, input.IDUsuario)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, apperrors.ErrUserNotFound
	}
	if input.IDIncapacidad != nil {
		ok, err := uc.repo.IncapacidadExists(ctx, *input.IDIncapacidad)
		if err != nil {
			return nil, err
		}
		if !ok {
			return nil, apperrors.ErrIncapacidadNotFound
		}
	}

	notificacion := &domain.Notificacion{
		IDUsuario:        input.IDUsuario,
		IDIncapacidad:    input.IDIncapacidad,
		TipoNotificacion: tipo,
		Mensaje:          strings.TrimSpace(input.Mensaje),
	}
	if err := uc.repo.Create(ctx, notificacion); err != nil {
		return nil, err
	}
	return notificacion, nil
}

func (uc *NotificacionUseCase) Obtener(ctx context.Context, actor ports.Actor, id uint64) (*domain.Notificacion, error) {
	notificacion, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if err := ensureCanAccess(actor, notificacion); err != nil {
		return nil, err
	}
	return notificacion, nil
}

func (uc *NotificacionUseCase) Listar(ctx context.Context, actor ports.Actor, filters ports.NotificacionFilters) ([]domain.Notificacion, int64, error) {
	if !actor.CanManageNotifications() {
		filters.IDUsuario = &actor.UserID
	}
	if filters.TipoNotificacion != "" {
		tipo, err := normalizeTipoNotificacion(filters.TipoNotificacion)
		if err != nil {
			return nil, 0, err
		}
		filters.TipoNotificacion = tipo
	}
	return uc.repo.List(ctx, filters)
}

func (uc *NotificacionUseCase) ContarNoLeidas(ctx context.Context, actor ports.Actor) (int64, error) {
	return uc.repo.CountUnread(ctx, actor.UserID)
}

func (uc *NotificacionUseCase) MarcarLeida(ctx context.Context, actor ports.Actor, id uint64) (*domain.Notificacion, error) {
	notificacion, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if err := ensureCanAccess(actor, notificacion); err != nil {
		return nil, err
	}
	if err := uc.repo.MarkAsRead(ctx, id); err != nil {
		return nil, err
	}
	notificacion.Leida = true
	return notificacion, nil
}

func (uc *NotificacionUseCase) MarcarTodasLeidas(ctx context.Context, actor ports.Actor) error {
	return uc.repo.MarkAllAsRead(ctx, actor.UserID)
}

func (uc *NotificacionUseCase) Eliminar(ctx context.Context, actor ports.Actor, id uint64) error {
	notificacion, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if err := ensureCanAccess(actor, notificacion); err != nil {
		return err
	}
	return uc.repo.Delete(ctx, id)
}

func ensureCanAccess(actor ports.Actor, notificacion *domain.Notificacion) error {
	if actor.CanManageNotifications() || notificacion.IDUsuario == actor.UserID {
		return nil
	}
	return apperrors.ErrForbidden.WithMessage("solo puedes acceder a tus propias notificaciones")
}

func normalizeTipoNotificacion(value string) (string, error) {
	tipo := shareddomain.TipoNotificacion(strings.TrimSpace(value))
	if !tipo.IsValid() {
		return "", apperrors.ErrValidation.WithMessage("tipo_notificacion inválido")
	}
	return string(tipo), nil
}
