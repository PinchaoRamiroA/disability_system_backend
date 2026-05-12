package usecase

import (
	"context"
	"fmt"
	"time"

	"disability_system_backend/internal/modules/reportes/domain"
	"disability_system_backend/internal/modules/reportes/ports"
	apperrors "disability_system_backend/internal/shared/errors"
)

type ReportesUseCase struct {
	repo ports.ReportesRepository
}

func NewReportesUseCase(repo ports.ReportesRepository) *ReportesUseCase {
	return &ReportesUseCase{repo: repo}
}

func (uc *ReportesUseCase) GenerarReporteIncapacidades(ctx context.Context, actor ports.Actor, filtros domain.FiltrosReporte) (*domain.ReporteIncapacidades, error) {
	if !actor.HasPermission("consultar_reportes") {
		return nil, apperrors.ErrForbidden.WithMessage("no tienes permiso para generar reportes")
	}

	reporte := &domain.ReporteIncapacidades{
		FechaGeneracion:     time.Now(),
		FechaInicio:         timePtrToTime(filtros.FechaInicio),
		FechaFin:            timePtrToTime(filtros.FechaFin),
		Periodo:             filtros.Periodo,
	}

	total, activas, cerradas, err := uc.repo.GetTotalIncapacidades(ctx, filtros)
	if err != nil {
		return nil, err
	}
	reporte.TotalIncapacidades = total
	reporte.IncapacidadesActivas = activas
	reporte.IncapacidadesCerradas = cerradas

	datos, err := uc.repo.GetIncapacidadesReport(ctx, filtros)
	if err != nil {
		return nil, err
	}
	reporte.Datos = datos

	return reporte, nil
}

func (uc *ReportesUseCase) GenerarReporteAusentismo(ctx context.Context, actor ports.Actor, filtros domain.FiltrosReporte) (*domain.ReporteAusentismo, error) {
	if !actor.HasPermission("consultar_reportes") {
		return nil, apperrors.ErrForbidden.WithMessage("no tienes permiso para generar reportes")
	}

	reporte := &domain.ReporteAusentismo{
		FechaGeneracion: time.Now(),
		FechaInicio:     timePtrToTime(filtros.FechaInicio),
		FechaFin:        timePtrToTime(filtros.FechaFin),
	}

	diasPerdidos, err := uc.repo.GetTotalDiasPerdidos(ctx, filtros)
	if err == nil {
		reporte.TotalDiasPerdidos = diasPerdidos
	}

	diasPorTipo, err := uc.repo.GetDiasPorTipo(ctx, filtros)
	if err == nil {
		reporte.DiasPorTipo = diasPorTipo
	}

	diasPorEntidad, err := uc.repo.GetDiasPorEntidad(ctx, filtros)
	if err == nil {
		reporte.DiasPorEntidad = diasPorEntidad
	}

	datos, top, err := uc.repo.GetAusentismoReport(ctx, filtros)
	if err == nil {
		reporte.Datos = datos
		reporte.TopIncapacidades = top
		reporte.TotalEmpleados = int64(len(datos))
	}

	return reporte, nil
}

func (uc *ReportesUseCase) GenerarReporteCartera(ctx context.Context, actor ports.Actor, filtros domain.FiltrosReporte) (*domain.ReporteCartera, error) {
	if !actor.HasPermission("consultar_reportes") && !actor.HasPermission("realizar_conciliacion") {
		return nil, apperrors.ErrForbidden.WithMessage("no tienes permiso para generar reportes de cartera")
	}

	return uc.repo.GetCarteraReport(ctx, filtros)
}

func (uc *ReportesUseCase) GenerarReporteEntidad(ctx context.Context, actor ports.Actor, entidadID uint64, filtros domain.FiltrosReporte) (*domain.ReporteEntidad, error) {
	if !actor.HasPermission("consultar_reportes") {
		return nil, apperrors.ErrForbidden.WithMessage("no tienes permiso para generar reportes")
	}

	reporte, err := uc.repo.GetEntidadReport(ctx, entidadID, filtros)
	if err != nil {
		return nil, err
	}
	reporte.FechaGeneracion = time.Now()

	return reporte, nil
}

func (uc *ReportesUseCase) GenerarReporteVencimientos(ctx context.Context, actor ports.Actor, diasMinimos int) (*domain.ReporteVencimientos, error) {
	if !actor.HasPermission("consultar_reportes") && !actor.HasPermission("generar_alertas") {
		return nil, apperrors.ErrForbidden.WithMessage("no tienes permiso para generar reportes de vencimientos")
	}

	reporte, err := uc.repo.GetVencimientosReport(ctx, diasMinimos)
	if err != nil {
		return nil, err
	}
	reporte.FechaGeneracion = time.Now()

	return reporte, nil
}

func (uc *ReportesUseCase) ObtenerResumenEjecutivo(ctx context.Context, actor ports.Actor) (*ResumenEjecutivo, error) {
	if !actor.HasPermission("consultar_reportes") {
		return nil, apperrors.ErrForbidden.WithMessage("no tienes permiso para consultar el resumen ejecutivo")
	}

	resumen := &ResumenEjecutivo{FechaGeneracion: time.Now()}

	total, activas, _, err := uc.repo.GetTotalIncapacidades(ctx, domain.FiltrosReporte{})
	if err == nil {
		resumen.TotalIncapacidades = total
		resumen.IncapacidadesActivas = activas
	}

	diasPerdidos, err := uc.repo.GetTotalDiasPerdidos(ctx, domain.FiltrosReporte{})
	if err == nil {
		resumen.TotalDiasPerdidos = diasPerdidos
	}

	pagados, pendientes, err := uc.repo.SumValorCartera(ctx)
	if err == nil {
		resumen.TotalValorCartera = formatCurrency(pagados + pendientes)
		resumen.TotalValorCobrado = formatCurrency(pagados)
		resumen.TotalValorPendiente = formatCurrency(pendientes)
	}

	pagosPendientes, err := uc.repo.CountPagosPendientes(ctx)
	if err == nil {
		resumen.PagosPendientes = pagosPendientes
	}

	pagosVencidos, err := uc.repo.CountPagosVencidos(ctx, 0)
	if err == nil {
		resumen.PagosVencidos = pagosVencidos
	}

	incActivas, err := uc.repo.CountIncapacidadesActivas(ctx)
	if err == nil {
		resumen.IncapacidadesActivas = incActivas
	}

	return resumen, nil
}

type ResumenEjecutivo struct {
	FechaGeneracion      time.Time
	TotalIncapacidades   int64
	IncapacidadesActivas int64
	TotalDiasPerdidos    int64
	TotalValorCartera    string
	TotalValorCobrado    string
	TotalValorPendiente   string
	PagosPendientes      int64
	PagosVencidos        int64
}

func timePtrToTime(t *time.Time) time.Time {
	if t == nil {
		return time.Time{}
	}
	return *t
}

func formatCurrency(value float64) string {
	return formatMoney(value)
}

func formatMoney(f float64) string {
	return fmt.Sprintf("%.2f", f)
}
