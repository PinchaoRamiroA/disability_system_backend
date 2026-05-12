package mapper

import (
	"time"

	"disability_system_backend/internal/modules/notificaciones/domain"
	"disability_system_backend/internal/modules/notificaciones/dto"
)

func ToNotificacionResponse(n *domain.Notificacion) dto.NotificacionResponse {
	return dto.NotificacionResponse{
		IDNotificacion:   n.IDNotificacion,
		IDUsuario:        n.IDUsuario,
		IDIncapacidad:    n.IDIncapacidad,
		TipoNotificacion: n.TipoNotificacion,
		Mensaje:          n.Mensaje,
		Fecha:            n.Fecha.Format(time.RFC3339),
		Leida:            n.Leida,
		CreatedAt:        n.CreatedAt.Format(time.RFC3339),
		UpdatedAt:        n.UpdatedAt.Format(time.RFC3339),
	}
}

func ToNotificacionResponses(items []domain.Notificacion) []dto.NotificacionResponse {
	responses := make([]dto.NotificacionResponse, 0, len(items))
	for i := range items {
		responses = append(responses, ToNotificacionResponse(&items[i]))
	}
	return responses
}
