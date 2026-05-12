package domain

import "time"

type TipoReporte string

const (
	TipoReporteIncapacidades  TipoReporte = "incapacidades"
	TipoReporteAusentismo     TipoReporte = "ausentismo"
	TipoReporteCartera        TipoReporte = "cartera"
	TipoReporteEntidad        TipoReporte = "entidad"
	TipoReporteDocumentos     TipoReporte = "documentos"
	TipoReporteVencimientos  TipoReporte = "vencimientos"
	TipoReporteProductividad TipoReporte = "productividad"
)

type PeriodoReporte string

const (
	PeriodoDiario     PeriodoReporte = "diario"
	PeriodoSemanal    PeriodoReporte = "semanal"
	PeriodoMensual    PeriodoReporte = "mensual"
	PeriodoTrimestral PeriodoReporte = "trimestral"
	PeriodoAnual      PeriodoReporte = "anual"
)

type ReporteIncapacidades struct {
	IDReporte             uint64
	Tipo                  TipoReporte
	FechaGeneracion       time.Time
	Periodo               PeriodoReporte
	FechaInicio           time.Time
	FechaFin              time.Time
	TotalIncapacidades   int64
	IncapacidadesActivas  int64
	IncapacidadesCerradas int64
	Datos                []ReporteIncapacidadDetalle
}

type ReporteIncapacidadDetalle struct {
	IDIncapacidad   uint64
	NombreEmpleado  string
	TipoIncapacidad string
	Estado          string
	Entidad         string
	FechaInicio     time.Time
	FechaFin        *time.Time
	DiasDuracion    int
	Observaciones   string
}

type ReporteAusentismo struct {
	IDReporte           uint64
	FechaGeneracion     time.Time
	FechaInicio         time.Time
	FechaFin            time.Time
	TotalDiasPerdidos   int64
	TotalEmpleados      int64
	DiasPorTipo         map[string]int64
	DiasPorEntidad      map[string]int64
	TopIncapacidades    []ReporteTopIncapacidad
	Datos               []ReporteAusentismoDetalle
}

type ReporteTopIncapacidad struct {
	IDEmpleado      uint64
	NombreEmpleado  string
	TotalDias       int64
	CantidadINC     int64
}

type ReporteAusentismoDetalle struct {
	IDEmpleado       uint64
	NombreEmpleado   string
	FechaInicio      time.Time
	FechaFin         *time.Time
	TipoIncapacidad  string
	DiasPerdidos     int64
	Entidad         string
}

type ReporteCartera struct {
	IDReporte            uint64
	FechaGeneracion      time.Time
	FechaInicio          time.Time
	FechaFin            time.Time
	TotalValorCartera   string
	TotalValorCobrado   string
	TotalValorPendiente string
	PagosPendientes     int64
	PagosVencidos       int64
	PorEntidad          []ReporteCarteraEntidad
}

type ReporteCarteraEntidad struct {
	NombreEntidad      string
	ValorTotal         string
	ValorCobrado       string
	ValorPendiente    string
	CantidadINC       int64
	PagosPendientes   int64
	PagosVencidos     int64
}

type ReporteEntidad struct {
	IDReporte        uint64
	FechaGeneracion  time.Time
	EntidadID        uint64
	NombreEntidad    string
	Tipo             string
	Estadisticas     EntidadEstadisticas
}

type EntidadEstadisticas struct {
	TotalIncapacidades      int64
	IncapacidadesActivas    int64
	IncapacidadesPagadas    int64
	IncapacidadesRechazadas int64
	TotalValorEsperado      string
	TotalValorCobrado       string
	TiempoPromedioPago      int
	TasaRechazo             float64
}

type ReporteVencimientos struct {
	IDReporte              uint64
	FechaGeneracion        time.Time
	AlertasDocumentos     []ReporteVencimientoDoc
	AlertasPagos          []ReporteVencimientoPago
	AlertasIncapacidades  []ReporteVencimientoINC
}

type ReporteVencimientoDoc struct {
	IDIncapacidad  uint64
	NombreEmpleado string
	TipoDocumento  string
	FechaLimite    time.Time
	DiasRestantes  int
	Estado        string
}

type ReporteVencimientoPago struct {
	IDIncapacidad  uint64
	IDPago         uint64
	NombreEntidad  string
	Valor          string
	FechaLimitePago time.Time
	DiasVencido    int
	Estado        string
}

type ReporteVencimientoINC struct {
	IDIncapacidad    uint64
	NombreEmpleado   string
	DiasTranscurridos int
	FechaInicio      time.Time
	Estado           string
	Alerta           string
}

type FiltrosReporte struct {
	FechaInicio  *time.Time
	FechaFin     *time.Time
	IDEntidad   *uint64
	IDTipo      *uint64
	IDEstado    *uint64
	IDEmpleado  *uint64
	Origen      string
	Estado      string
	Periodo     PeriodoReporte
}
