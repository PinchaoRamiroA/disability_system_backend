package mapper

import (
	"disability_system_backend/internal/modules/auditoria/domain"
	"disability_system_backend/internal/modules/auditoria/dto"
)

func ToResponse(a domain.Auditoria) dto.AuditoriaResponse {
	return dto.AuditoriaResponse{
		IDAuditoria:    a.ID,
		IDUsuario:      a.IDUsuario,
		UsuarioNombre:  a.UsuarioNombre,
		IDIncapacidad:  a.IDIncapacidad,
		TipoAccion:     a.TipoAccion,
		Modulo:         a.Modulo,
		Descripcion:    a.Descripcion,
		CambioAnterior: a.CambioAnterior,
		CambioNuevo:    a.CambioNuevo,
		CreatedAt:      a.CreatedAt,
	}
}
