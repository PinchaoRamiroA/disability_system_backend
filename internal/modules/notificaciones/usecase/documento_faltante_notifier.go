package usecase

import (
	"context"
	"strings"

	incdomain "disability_system_backend/internal/modules/incapacidades/domain"
	"disability_system_backend/internal/modules/notificaciones/domain"
	"disability_system_backend/internal/modules/notificaciones/ports"
	shareddomain "disability_system_backend/internal/shared/domain"
)

type DocumentoFaltanteNotifier struct {
	repo ports.NotificacionRepository
}

func NewDocumentoFaltanteNotifier(repo ports.NotificacionRepository) *DocumentoFaltanteNotifier {
	return &DocumentoFaltanteNotifier{repo: repo}
}

func (n *DocumentoFaltanteNotifier) NotificarDocumentosFaltantes(ctx context.Context, userID, incapacidadID uint64, documentos []incdomain.TipoDocumento) error {
	if len(documentos) == 0 {
		return nil
	}
	mensaje := buildDocumentoFaltanteMessage(documentos)
	if exists, err := n.existsUnread(ctx, userID, incapacidadID, mensaje); err != nil {
		return err
	} else if exists {
		return nil
	}

	notificacion := &domain.Notificacion{
		IDUsuario:        userID,
		IDIncapacidad:    &incapacidadID,
		TipoNotificacion: string(shareddomain.TipoNotifDocFaltante),
		Mensaje:          mensaje,
	}
	return n.repo.Create(ctx, notificacion)
}

func (n *DocumentoFaltanteNotifier) existsUnread(ctx context.Context, userID, incapacidadID uint64, mensaje string) (bool, error) {
	leida := false
	items, _, err := n.repo.List(ctx, ports.NotificacionFilters{
		IDUsuario:        &userID,
		IDIncapacidad:    &incapacidadID,
		TipoNotificacion: string(shareddomain.TipoNotifDocFaltante),
		Leida:            &leida,
		Page:             1,
		Limit:            100,
	})
	if err != nil {
		return false, err
	}
	for _, item := range items {
		if item.Mensaje == mensaje {
			return true, nil
		}
	}
	return false, nil
}

func buildDocumentoFaltanteMessage(documentos []incdomain.TipoDocumento) string {
	nombres := make([]string, 0, len(documentos))
	for _, documento := range documentos {
		nombres = append(nombres, documento.Nombre)
	}
	if len(nombres) == 1 {
		return "Documento faltante: " + nombres[0]
	}
	return "Documentos faltantes: " + strings.Join(nombres, ", ")
}
