package usecase

import (
	"context"
	"strconv"
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
	IDIncapacidad    uint64            `json:"id_incapacidad"`
	Incapacidad      IncapacidadResumen `json:"incapacidad"`
	NombreEntidad    string            `json:"nombre_entidad"`
	TipoAlerta       string            `json:"tipo_alerta"`
	FechaVencimiento string            `json:"fecha_vencimiento"`
	DiasRestantes    int               `json:"dias_restantes"`
	Prioridad        string            `json:"prioridad"`
	Mensaje          string            `json:"mensaje"`
}

type IncapacidadResumen struct {
	ID     uint64 `json:"id"`
	Titulo string `json:"titulo"`
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

type resumenEntidadAccumulator struct {
	IDEntidad       uint64
	Nombre          string
	Tipo            string
	CantidadINC     int64
	ValorTotal      decimal.Decimal
	ValorCobrado    decimal.Decimal
	ValorPendiente  decimal.Decimal
	PagosPendientes int64
	PagosVencidos   int64
}

func (s *CobroWorkflowService) ObtenerEstadisticasGenerales(ctx context.Context) (*EstadisticasCartera, error) {
	pagos, total, err := s.pagoRepo.ListPagos(ctx, ports.PagoFilters{Limit: 10000})
	if err != nil {
		return nil, err
	}

	var totalValor, totalCobrado, totalPendiente decimal.Decimal
	var pagosPendientes, pagosVencidos int64
	incapacidadesTotal := make(map[uint64]bool)
	incapacidadesActivas := make(map[uint64]bool)

	for _, pago := range pagos {
		totalValor = totalValor.Add(pago.Valor)

		if pago.IDIncapacidad != 0 {
			incapacidadesTotal[pago.IDIncapacidad] = true
		}

		switch pago.EstadoPago {
		case "Pagado", "Conciliado":
			totalCobrado = totalCobrado.Add(pago.Valor)
		default:
			totalPendiente = totalPendiente.Add(pago.Valor)
			pagosPendientes++
			if pago.FechaPago.Before(time.Now()) && pago.EstadoPago != "Anulado" {
				pagosVencidos++
			}
			if pago.EstadoPago != "Anulado" && pago.IDIncapacidad != 0 {
				incapacidadesActivas[pago.IDIncapacidad] = true
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

	totalIncapacidades := int64(len(incapacidadesTotal))
	if totalIncapacidades == 0 && total > 0 {
		totalIncapacidades = total
	}

	totalActivas := int64(len(incapacidadesActivas))
	if totalActivas == 0 && len(incapacidadesTotal) == 0 && pagosPendientes > 0 {
		totalActivas = pagosPendientes
	}

	return &EstadisticasCartera{
		TotalIncapacidades:     totalIncapacidades,
		IncapacidadesActivas:   totalActivas,
		TotalValorCartera:      formatDecimal(totalValor),
		TotalValorCobrado:      formatDecimal(totalCobrado),
		TotalValorPendiente:    formatDecimal(totalPendiente),
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

	entidadInfo, err := s.pagoRepo.GetEntidadInfo(ctx)
	if err != nil {
		return nil, err
	}

	resumenPorEntidad := make(map[uint64]*resumenEntidadAccumulator)

	for _, pago := range pagos {
		acc, ok := resumenPorEntidad[pago.IDEntidad]
		if !ok {
			info := entidadInfo[pago.IDEntidad]
			acc = &resumenEntidadAccumulator{
				IDEntidad: pago.IDEntidad,
				Nombre:    info.Nombre,
				Tipo:      info.Tipo,
			}
			resumenPorEntidad[pago.IDEntidad] = acc
		}

		acc.CantidadINC++
		acc.ValorTotal = acc.ValorTotal.Add(pago.Valor)

		switch pago.EstadoPago {
		case "Pagado", "Conciliado":
			acc.ValorCobrado = acc.ValorCobrado.Add(pago.Valor)
		default:
			acc.ValorPendiente = acc.ValorPendiente.Add(pago.Valor)
			acc.PagosPendientes++
			if pago.FechaPago.Before(time.Now()) && pago.EstadoPago != "Anulado" {
				acc.PagosVencidos++
			}
		}
	}

	result := make([]ResumenEntidad, 0, len(resumenPorEntidad))
	for _, acc := range resumenPorEntidad {
		result = append(result, ResumenEntidad{
			IDEntidad:       acc.IDEntidad,
			Nombre:          acc.Nombre,
			Tipo:            acc.Tipo,
			CantidadINC:     acc.CantidadINC,
			ValorTotal:      formatDecimal(acc.ValorTotal),
			ValorCobrado:    formatDecimal(acc.ValorCobrado),
			ValorPendiente:  formatDecimal(acc.ValorPendiente),
			PagosPendientes: acc.PagosPendientes,
			PagosVencidos:   acc.PagosVencidos,
		})
	}
	return result, nil
}

func (s *CobroWorkflowService) ObtenerAlertasVencimiento(ctx context.Context, diasMinimos int, actor ports.Actor) ([]AlertaVencimiento, error) {
	filters := ports.PagoFilters{Limit: 1000}
	if !actor.HasPermission("registrar_pago") && !actor.HasPermission("gestionar_cobro_persuasivo") && !actor.HasPermission("gestionar_cobro_juridico") && !actor.HasPermission("generar_alertas") {
		filters.UserID = &actor.UserID
	}
	pagos, _, err := s.pagoRepo.ListPagos(ctx, filters)
	if err != nil {
		return nil, err
	}

	incapacidadIDs := make([]uint64, 0)
	for _, p := range pagos {
		incapacidadIDs = append(incapacidadIDs, p.IDIncapacidad)
	}

	incapacidades, err := s.pagoRepo.GetIncapacidadesDetailed(ctx, incapacidadIDs)
	if err != nil {
		return nil, err
	}

	alertas := make([]AlertaVencimiento, 0)
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	diasVentana := diasMinimos
	if diasVentana < 0 {
		diasVentana = -diasVentana
	}
	fechaLimite := today.AddDate(0, 0, diasVentana+1)

	for _, pago := range pagos {
		pagoDate := time.Date(pago.FechaPago.Year(), pago.FechaPago.Month(), pago.FechaPago.Day(), 0, 0, 0, 0, time.UTC)
		if pagoDate.Before(fechaLimite) && pago.EstadoPago != "Pagado" && pago.EstadoPago != "Anulado" && pago.EstadoPago != "Conciliado" {
			diasVencido := int(today.Sub(pagoDate).Hours() / 24)
			tipoAlerta := getTipoAlerta(diasVencido)

			var mensaje string
			if diasVencido > 0 {
				mensaje = "Vencida hace " + strconv.Itoa(diasVencido) + " días"
			} else if diasVencido == 0 {
				mensaje = "Vence hoy"
			} else {
				mensaje = "Vence en " + strconv.Itoa(-diasVencido) + " días"
			}

			incInfo := IncapacidadResumen{
				ID: pago.IDIncapacidad,
			}
			if inc, ok := incapacidades[pago.IDIncapacidad]; ok {
				incInfo.Titulo = inc.Titulo
			}

			alertas = append(alertas, AlertaVencimiento{
				IDIncapacidad:    pago.IDIncapacidad,
				Incapacidad:      incInfo,
				NombreEntidad:    pago.NombreEntidad,
				TipoAlerta:       tipoAlerta,
				FechaVencimiento: pago.FechaPago.Format("2006-01-02"),
				DiasRestantes:    -diasVencido,
				Prioridad:        tipoAlerta,
				Mensaje:          mensaje,
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

func formatDecimal(d decimal.Decimal) string {
	return d.Round(2).String()
}

func formatCurrency(value float64) string {
	return formatFloat(value)
}

func sumCurrency(current string, addition float64) float64 {
	var currentDec decimal.Decimal
	if current != "" {
		if d, err := decimal.NewFromString(current); err == nil {
			currentDec = d
		}
	}
	res, _ := currentDec.Add(decimal.NewFromFloat(addition)).Float64()
	return res
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
	return workflowSvc.ObtenerAlertasVencimiento(ctx, diasMinimos, actor)
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
