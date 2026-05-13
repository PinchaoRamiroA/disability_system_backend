package dto

import "time"

type GenerarReporteRequest struct {
	Tipo        string  `json:"tipo"`
	FechaInicio string  `json:"fecha_inicio"`
	FechaFin    string  `json:"fecha_fin"`
	IDEntidad   *uint64 `json:"id_entidad,omitempty"`
	IDTipo      *uint64 `json:"id_tipo,omitempty"`
	IDEstado    *uint64 `json:"id_estado,omitempty"`
	IDEmpleado  *uint64 `json:"id_empleado,omitempty"`
	Origen      string  `json:"origen,omitempty"`
	Periodo     string  `json:"periodo,omitempty"`
}

type ResumenEjecutivoResponse struct {
	FechaGeneracion      time.Time `json:"fecha_generacion"`
	TotalIncapacidades   int64     `json:"total_incapacidades"`
	IncapacidadesActivas int64     `json:"incapacidades_activas"`
	TotalDiasPerdidos    int64     `json:"total_dias_perdidos"`
	TotalValorCartera    string    `json:"total_valor_cartera"`
	TotalValorCobrado    string    `json:"total_valor_cobrado"`
	TotalValorPendiente   string    `json:"total_valor_pendiente"`
	PagosPendientes      int64     `json:"pagos_pendientes"`
	PagosVencidos        int64     `json:"pagos_vencidos"`
}
