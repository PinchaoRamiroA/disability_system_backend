package mapper

import (
	"time"

	"disability_system_backend/internal/modules/incapacidades/domain"
	"disability_system_backend/internal/modules/incapacidades/dto"
)

const dateLayout = "2006-01-02"

func ToIncapacidadResponse(i *domain.Incapacidad) dto.IncapacidadResponse {
	resp := dto.IncapacidadResponse{
		IDIncapacidad:   i.IDIncapacidad,
		IDUsuario:       i.IDUsuario,
		CanalRecepcion:  i.CanalRecepcion,
		Titulo:          i.Titulo,
		FechaInicio:     i.FechaInicio.Format(dateLayout),
		FechaFin:        formatDatePtr(i.FechaFin),
		Origen:          i.Origen,
		FechaRadicacion: formatDatePtr(i.FechaRadicacion),
		FechaPago:       formatDatePtr(i.FechaPago),
		Observaciones:   i.Observaciones,
		CreatedBy:       i.CreatedBy,
		CreatedAt:       i.CreatedAt.Format(time.RFC3339),
		UpdatedAt:       i.UpdatedAt.Format(time.RFC3339),
	}
	if i.Estado != nil {
		resp.Estado = &dto.EstadoIncapacidadResponse{
			IDEstado:          i.Estado.IDEstado,
			Nombre:            i.Estado.Nombre,
			Descripcion:       i.Estado.Descripcion,
			PermiteTransicion: i.Estado.PermiteTransicion,
		}
	}
	if i.Tipo != nil {
		resp.Tipo = &dto.TipoIncapacidadResponse{
			IDTipo:               i.Tipo.IDTipo,
			Nombre:               i.Tipo.Nombre,
			DocumentosRequeridos: i.Tipo.DocumentosRequeridos,
		}
	}
	if i.Entidad != nil {
		resp.Entidad = &dto.EntidadResponse{
			IDEntidad:              i.Entidad.IDEntidad,
			Nombre:                 i.Entidad.Nombre,
			Tipo:                   i.Entidad.Tipo,
			PlazoTranscripcionDias: i.Entidad.PlazoTranscripcionDias,
			TiempoMaximoPagoDias:   i.Entidad.TiempoMaximoPagoDias,
			CanalAtencion:          i.Entidad.CanalAtencion,
			CanalesAtencion:        i.Entidad.CanalesAtencion,
			RequiereTranscripcion:  i.Entidad.RequiereTranscripcion,
		}
	}
	return resp
}

func ToIncapacidadResponses(items []domain.Incapacidad) []dto.IncapacidadResponse {
	responses := make([]dto.IncapacidadResponse, 0, len(items))
	for i := range items {
		responses = append(responses, ToIncapacidadResponse(&items[i]))
	}
	return responses
}

func ToEstadoResponses(items []domain.EstadoIncapacidad) []dto.EstadoIncapacidadResponse {
	responses := make([]dto.EstadoIncapacidadResponse, 0, len(items))
	for _, item := range items {
		responses = append(responses, dto.EstadoIncapacidadResponse{
			IDEstado:          item.IDEstado,
			Nombre:            item.Nombre,
			Descripcion:       item.Descripcion,
			PermiteTransicion: item.PermiteTransicion,
		})
	}
	return responses
}

func ToTipoResponses(items []domain.TipoIncapacidad) []dto.TipoIncapacidadResponse {
	responses := make([]dto.TipoIncapacidadResponse, 0, len(items))
	for _, item := range items {
		responses = append(responses, dto.TipoIncapacidadResponse{
			IDTipo:               item.IDTipo,
			Nombre:               item.Nombre,
			DocumentosRequeridos: item.DocumentosRequeridos,
		})
	}
	return responses
}

func ToEntidadResponses(items []domain.Entidad) []dto.EntidadResponse {
	responses := make([]dto.EntidadResponse, 0, len(items))
	for _, item := range items {
		responses = append(responses, dto.EntidadResponse{
			IDEntidad:              item.IDEntidad,
			Nombre:                 item.Nombre,
			Tipo:                   item.Tipo,
			PlazoTranscripcionDias: item.PlazoTranscripcionDias,
			TiempoMaximoPagoDias:   item.TiempoMaximoPagoDias,
			CanalAtencion:          item.CanalAtencion,
			CanalesAtencion:        item.CanalesAtencion,
			RequiereTranscripcion:  item.RequiereTranscripcion,
		})
	}
	return responses
}

func formatDatePtr(t *time.Time) *string {
	if t == nil {
		return nil
	}
	formatted := t.Format(dateLayout)
	return &formatted
}
