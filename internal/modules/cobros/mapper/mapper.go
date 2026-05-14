package mapper

import (
	"time"

	"disability_system_backend/internal/modules/cobros/domain"
	"disability_system_backend/internal/modules/cobros/dto"
)

const dateLayout = "2006-01-02"

func ToPagoResponse(p *domain.Pago) dto.PagoResponse {
	return dto.PagoResponse{
		IDPago:          p.IDPago,
		IDIncapacidad:   p.IDIncapacidad,
		IDEntidad:       p.IDEntidad,
		TipoPago:        p.TipoPago,
		EstadoPago:      p.EstadoPago,
		Descripcion:     p.Descripcion,
		Valor:           p.Valor.StringFixed(2),
		FechaPago:       p.FechaPago.Format(dateLayout),
		PeriodoContable: p.PeriodoContable,
		Conciliado:      p.Conciliado,
		RegistradoPor:   p.RegistradoPor,
		NombreEntidad:   p.NombreEntidad,
		CreatedAt:       p.CreatedAt.Format(time.RFC3339),
		UpdatedAt:       p.UpdatedAt.Format(time.RFC3339),
	}
}

func ToPagoResponses(items []domain.Pago) []dto.PagoResponse {
	responses := make([]dto.PagoResponse, 0, len(items))
	for i := range items {
		responses = append(responses, ToPagoResponse(&items[i]))
	}
	return responses
}

func ToSeguimientoResponse(s *domain.SeguimientoCobro) dto.SeguimientoCobroResponse {
	return dto.SeguimientoCobroResponse{
		IDSeguimiento:   s.IDSeguimiento,
		IDIncapacidad:   s.IDIncapacidad,
		TipoSeguimiento: s.TipoSeguimiento,
		Descripcion:     s.Descripcion,
		Fecha:           s.Fecha.Format(time.RFC3339),
		Resultado:       s.Resultado,
		GestionadoPor:   s.GestionadoPor,
		CreatedAt:       s.CreatedAt.Format(time.RFC3339),
		UpdatedAt:       s.UpdatedAt.Format(time.RFC3339),
	}
}

func ToSeguimientoResponses(items []domain.SeguimientoCobro) []dto.SeguimientoCobroResponse {
	responses := make([]dto.SeguimientoCobroResponse, 0, len(items))
	for i := range items {
		responses = append(responses, ToSeguimientoResponse(&items[i]))
	}
	return responses
}
