package usecase

import (
	"context"
	"time"

	"disability_system_backend/internal/modules/cobros/domain"
	"disability_system_backend/internal/modules/cobros/ports"
	apperrors "disability_system_backend/internal/shared/errors"

	"github.com/shopspring/decimal"
)

type EstadisticasCartera struct {
	TotalIncapacidades    int64
	IncapacidadesActivas  int64
	TotalValorCartera     string
	TotalValorCobrado     string
	TotalValorPendiente   string
	PagosPendientes       int64
	PagosVencidos         int64
	SeguimientosPendientes int64
}

type ResumenEntidad struct {
	IDEntidad   uint64
	Nombre      string
	Tipo        string
	CantidadINC int64
	ValorTotal  string
	ValorCobrado string
	ValorPendiente string
	PagosPendientes int64
	PagosVencidos int64
}

type AlertaVencimiento struct {
	IDIncapacidad   uint64
	IDEntidad       uint64
	NombreEntidad   string
	DiasVencido     int
	Estado         string
	FechaLimitePago *time.Time
	TipoAlerta     string
}

type CobroWorkflowService struct {
	pagoRepo        ports.PagoRepository
	seguimientoRepo ports.SeguimientoRepository
}

func NewCobroWorkflowService(pagoRepo ports.PagoRepository, seguimientoRepo ports.SeguimientoRepository) *CobroWorkflowService {
	return &CobroWorkflowService{
		pagoRepo:        pagoRepo,
		seguimientoRepo: seguimientoRepo,
	}
}

func (s *CobroWorkflowService) ObtenerEstadisticasGenerales(ctx context.Context) (*EstadisticasCartera, error) {
	pagos, total, err := s.pagoRepo.ListPagos(ctx, ports.PagoFilters{Limit: 1000})
	if err != nil {
		return nil, err
	}

	var totalValor, totalCobrado, totalPendiente float64
	var pagosPendientes, pagosVencidos int64

	for _, pago := range pagos {
		valor, _ := pago.Valor.Float64()
		totalValor += valor

		switch pago.EstadoPago {
		case "Pagado", "Conciliado":
			totalCobrado += valor
		default:
			totalPendiente += valor
			pagosPendientes++
			if pago.FechaPago.Before(time.Now()) && pago.EstadoPago != "Anulado" {
				pagosVencidos++
			}
		}
	}

	seguimientos, _, err := s.seguimientoRepo.ListSeguimientos(ctx, ports.SeguimientoFilters{
		Limit: 1000,
	})
	if err != nil {
		return nil, err
	}

	var seguimientosPendientes int64
	for _, seg := range seguimientos {
		if seg.Resultado == nil || *seg.Resultado == "Pendiente respuesta" || *seg.Resultado == "En revisión" {
			seguimientosPendientes++
		}
	}

	return &EstadisticasCartera{
		TotalIncapacidades:     total,
		IncapacidadesActivas:    total - int64(len(pagos)),
		TotalValorCartera:      formatCurrency(totalValor),
		TotalValorCobrado:      formatCurrency(totalCobrado),
		TotalValorPendiente:    formatCurrency(totalPendiente),
		PagosPendientes:        pagosPendientes,
		PagosVencidos:          pagosVencidos,
		SeguimientosPendientes: seguimientosPendientes,
	}, nil
}

func (s *CobroWorkflowService) ObtenerResumenPorEntidad(ctx context.Context) ([]ResumenEntidad, error) {
	pagos, _, err := s.pagoRepo.ListPagos(ctx, ports.PagoFilters{Limit: 10000})
	if err != nil {
		return nil, err
	}

	resumenPorEntidad := make(map[uint64]*ResumenEntidad)

	for _, pago := range pagos {
		if _, ok := resumenPorEntidad[pago.IDEntidad]; !ok {
			resumenPorEntidad[pago.IDEntidad] = &ResumenEntidad{
				IDEntidad: pago.IDEntidad,
			}
		}

		valor, _ := pago.Valor.Float64()
		r := resumenPorEntidad[pago.IDEntidad]
		r.CantidadINC++

		switch pago.EstadoPago {
		case "Pagado", "Conciliado":
			r.ValorCobrado = formatCurrency(sumCurrency(r.ValorCobrado, valor))
		default:
			r.ValorPendiente = formatCurrency(sumCurrency(r.ValorPendiente, valor))
			r.PagosPendientes++
			if pago.FechaPago.Before(time.Now()) && pago.EstadoPago != "Anulado" {
				r.PagosVencidos++
			}
		}
		r.ValorTotal = formatCurrency(sumCurrency(r.ValorTotal, valor))
	}

	result := make([]ResumenEntidad, 0, len(resumenPorEntidad))
	for _, r := range resumenPorEntidad {
		result = append(result, *r)
	}
	return result, nil
}

func (s *CobroWorkflowService) ObtenerAlertasVencimiento(ctx context.Context, diasMinimos int) ([]AlertaVencimiento, error) {
	pagos, _, err := s.pagoRepo.ListPagos(ctx, ports.PagoFilters{Limit: 1000})
	if err != nil {
		return nil, err
	}

	var alertas []AlertaVencimiento
	fechaLimite := time.Now().AddDate(0, 0, -diasMinimos)

	for _, pago := range pagos {
		if pago.FechaPago.Before(fechaLimite) && pago.EstadoPago != "Pagado" && pago.EstadoPago != "Anulado" && pago.EstadoPago != "Conciliado" {
			diasVencido := int(time.Since(pago.FechaPago).Hours() / 24)
			alertas = append(alertas, AlertaVencimiento{
				IDIncapacidad:   pago.IDIncapacidad,
				IDEntidad:       pago.IDEntidad,
				DiasVencido:     diasVencido,
				Estado:          pago.EstadoPago,
				FechaLimitePago: &pago.FechaPago,
				TipoAlerta:      getTipoAlerta(diasVencido),
			})
		}
	}

	return alertas, nil
}

func (s *CobroWorkflowService) DeterminarProximoEstadoIncapacidad(ctx context.Context, incapacidadID uint64, accion string) (string, error) {
	switch accion {
	case "registrar_pago":
		return "Pagada", nil
	case "conciliar":
		return "Conciliada", nil
	case "iniciar_cobro_persuasivo":
		return "Cobro persuasivo", nil
	case "iniciar_cobro_juridico":
		return "Cobro jurídico", nil
	case "completar_cobro":
		return "Pagada", nil
	case "aprobar_pago":
		return "Pagada", nil
	default:
		return "", apperrors.ErrValidation.WithMessage("acción no reconocida para determinar estado")
	}
}

func (s *CobroWorkflowService) ValidarTransicionEstado(estadoActual, nuevoEstado string) error {
	validTransitions := map[string][]string{
		"Cobrada":             {"Pendiente pago", "Cobro persuasivo", "Cobro jurídico", "Rechazada", "Archivada"},
		"Pendiente pago":      {"Pagada", "Cobro persuasivo", "Cobro jurídico", "Rechazada", "Archivada"},
		"Pagada":              {"En conciliación", "Archivada"},
		"En conciliación":     {"Conciliada", "Archivada"},
		"Conciliada":          {"Archivada", "Cerrada"},
		"Cobro persuasivo":   {"Cobro jurídico", "Pagada", "Archivada"},
		"Cobro jurídico":     {"Pagada", "Rechazada", "Archivada"},
	}

	if transitions, ok := validTransitions[estadoActual]; ok {
		for _, t := range transitions {
			if t == nuevoEstado {
				return nil
			}
		}
	}

	return apperrors.ErrConflict.WithMessage("transición de estado no válida: de " + estadoActual + " a " + nuevoEstado)
}

func (s *CobroWorkflowService) ObtenerCarteraVencida(ctx context.Context) ([]domain.Pago, error) {
	pagos, _, err := s.pagoRepo.ListPagos(ctx, ports.PagoFilters{Limit: 10000})
	if err != nil {
		return nil, err
	}

	var carteraVencida []domain.Pago
	for _, pago := range pagos {
		if pago.FechaPago.Before(time.Now()) &&
			pago.EstadoPago != "Pagado" &&
			pago.EstadoPago != "Anulado" &&
			pago.EstadoPago != "Conciliado" {
			carteraVencida = append(carteraVencida, pago)
		}
	}

	return carteraVencida, nil
}

func (s *CobroWorkflowService) CalcularDiasVencido(fechaPago time.Time) int {
	if fechaPago.After(time.Now()) {
		return 0
	}
	return int(time.Since(fechaPago).Hours() / 24)
}

func formatCurrency(value float64) string {
	return formatFloat(value)
}

func sumCurrency(current string, addition float64) float64 {
	var currentValue float64
	if current != "" {
		d, err := decimal.NewFromString(current)
		if err == nil {
			currentValue, _ = d.Float64()
		}
	}
	return currentValue + addition
}

func formatFloat(f float64) string {
	return decimal.NewFromFloat(f).Round(2).String()
}

func getTipoAlerta(diasVencido int) string {
	if diasVencido > 60 {
		return "Crítico"
	}
	if diasVencido > 30 {
		return "Alto"
	}
	if diasVencido > 15 {
		return "Medio"
	}
	return "Bajo"
}

func (uc *CobroUseCase) ObtenerEstadisticasGenerales(ctx context.Context, actor ports.Actor) (*EstadisticasCartera, error) {
	if !canReadCobros(actor) {
		return nil, apperrors.ErrForbidden.WithMessage("no tienes permiso para consultar estadísticas")
	}

	workflowSvc := NewCobroWorkflowService(uc.pagoRepo, uc.seguimientoRepo)
	return workflowSvc.ObtenerEstadisticasGenerales(ctx)
}

func (uc *CobroUseCase) ObtenerResumenPorEntidad(ctx context.Context, actor ports.Actor) ([]ResumenEntidad, error) {
	if !canReadCobros(actor) {
		return nil, apperrors.ErrForbidden.WithMessage("no tienes permiso para consultar estadísticas")
	}

	workflowSvc := NewCobroWorkflowService(uc.pagoRepo, uc.seguimientoRepo)
	return workflowSvc.ObtenerResumenPorEntidad(ctx)
}

func (uc *CobroUseCase) ObtenerAlertasVencimiento(ctx context.Context, actor ports.Actor, diasMinimos int) ([]AlertaVencimiento, error) {
	if !canReadCobros(actor) && !actor.HasPermission("generar_alertas") {
		return nil, apperrors.ErrForbidden.WithMessage("no tienes permiso para consultar alertas")
	}

	workflowSvc := NewCobroWorkflowService(uc.pagoRepo, uc.seguimientoRepo)
	return workflowSvc.ObtenerAlertasVencimiento(ctx, diasMinimos)
}

func (uc *CobroUseCase) ObtenerCarteraVencida(ctx context.Context, actor ports.Actor) ([]domain.Pago, error) {
	if !canReadCobros(actor) {
		return nil, apperrors.ErrForbidden.WithMessage("no tienes permiso para consultar cartera vencida")
	}

	workflowSvc := NewCobroWorkflowService(uc.pagoRepo, uc.seguimientoRepo)
	return workflowSvc.ObtenerCarteraVencida(ctx)
}

func (uc *CobroUseCase) ObtenerProximoEstadoIncapacidad(ctx context.Context, actor ports.Actor, incapacidadID uint64, accion string) (string, error) {
	if !canReadCobros(actor) {
		return "", apperrors.ErrForbidden.WithMessage("no tienes permiso para consultar workflow")
	}

	workflowSvc := NewCobroWorkflowService(uc.pagoRepo, uc.seguimientoRepo)
	return workflowSvc.DeterminarProximoEstadoIncapacidad(ctx, incapacidadID, accion)
}
