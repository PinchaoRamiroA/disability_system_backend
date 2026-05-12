package ports

import (
	"context"

	"disability_system_backend/internal/modules/notificaciones/domain"
)

type Actor struct {
	UserID   uint64
	Role     string
	Permisos []string
}

func (a Actor) HasPermission(permission string) bool {
	for _, p := range a.Permisos {
		if p == permission {
			return true
		}
	}
	return false
}

func (a Actor) CanManageNotifications() bool {
	return a.HasPermission("generar_alertas") || a.HasPermission("gestionar_usuarios")
}

type NotificacionFilters struct {
	IDUsuario        *uint64
	IDIncapacidad    *uint64
	TipoNotificacion string
	Leida            *bool
	Page             int
	Limit            int
}

type NotificacionRepository interface {
	Create(ctx context.Context, notificacion *domain.Notificacion) error
	FindByID(ctx context.Context, id uint64) (*domain.Notificacion, error)
	List(ctx context.Context, filters NotificacionFilters) ([]domain.Notificacion, int64, error)
	CountUnread(ctx context.Context, userID uint64) (int64, error)
	MarkAsRead(ctx context.Context, id uint64) error
	MarkAllAsRead(ctx context.Context, userID uint64) error
	Delete(ctx context.Context, id uint64) error
	UserExists(ctx context.Context, id uint64) (bool, error)
	IncapacidadExists(ctx context.Context, id uint64) (bool, error)
}

type PermissionRepository interface {
	FindPermissionsByRoleName(ctx context.Context, role string) ([]string, error)
}
