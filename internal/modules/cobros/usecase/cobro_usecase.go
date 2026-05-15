package usecase

import (
	"context"
	"strings"
	"time"

	"disability_system_backend/internal/modules/cobros/domain"
	"disability_system_backend/internal/modules/cobros/ports"
	shareddomain "disability_system_backend/internal/shared/domain"
	apperrors "disability_system_backend/internal/shared/errors"

	"github.com/shopspring/decimal"
)

const dateLayout = "2006-01-02"

type CrearPagoInput struct {
	IDIncapacidad   uint64
	IDEntidad       uint64
	TipoPago        string
	EstadoPago      string
	Descripcion     *string
	Valor           string
	FechaPago       string
	PeriodoContable *string
}

type ActualizarPagoInput struct {
	IDEntidad       *uint64
	TipoPago        *string
	EstadoPago      *string
	Descripcion     *string
	Valor           *string
	FechaPago       *string
	PeriodoContable *string
}

type CrearSeguimientoInput struct {
	IDIncapacidad   uint64
	TipoSeguimiento string
	Descripcion     *string
	Resultado       *string
}

type ActualizarSeguimientoInput struct {
	TipoSeguimiento *string
	Descripcion     *string
	Resultado       *string
}

type CobroUseCase struct {
	pagoRepo        ports.PagoRepository
	seguimientoRepo ports.SeguimientoRepository
}

func NewCobroUseCase(repo interface {
	ports.PagoRepository
	ports.SeguimientoRepository
}) *CobroUseCase {
	return &CobroUseCase{
		pagoRepo:        repo,
		seguimientoRepo: repo,
	}
}

func (uc *CobroUseCase) CrearPago(ctx context.Context, actor ports.Actor, input CrearPagoInput) (*domain.Pago, error) {
	if !actor.HasPermission("registrar_pago") {
		return nil, apperrors.ErrForbidden.WithMessage("no tienes permiso para registrar pagos")
	}
	if err := uc.ensureIncapacidadAndEntidad(ctx, input.IDIncapacidad, input.IDEntidad); err != nil {
		return nil, err
	}
	tipoPago, err := normalizeTipoPago(input.TipoPago)
	if err != nil {
		return nil, err
	}
	estadoPago := input.EstadoPago
	if strings.TrimSpace(estadoPago) == "" {
		estadoPago = string(shareddomain.EstadoPagoPendiente)
	}
	estadoPago, err = normalizeEstadoPago(estadoPago)
	if err != nil {
		return nil, err
	}
	valor, err := parseValor(input.Valor)
	if err != nil {
		return nil, err
	}
	fechaPago, err := parseDate(input.FechaPago, "fecha_pago")
	if err != nil {
		return nil, err
	}

	pago := &domain.Pago{
		IDIncapacidad:   input.IDIncapacidad,
		IDEntidad:       input.IDEntidad,
		TipoPago:        tipoPago,
		EstadoPago:      estadoPago,
		Descripcion:     input.Descripcion,
		Valor:           valor,
		FechaPago:       fechaPago,
		PeriodoContable: input.PeriodoContable,
		Conciliado:      strings.EqualFold(estadoPago, string(shareddomain.EstadoPagoConciliado)),
		RegistradoPor:   &actor.UserID,
	}
	if err := uc.pagoRepo.CreatePago(ctx, pago); err != nil {
		return nil, err
	}
	return pago, nil
}

func (uc *CobroUseCase) ObtenerPago(ctx context.Context, actor ports.Actor, id uint64) (*domain.Pago, error) {
	if !canReadCobros(actor) {
		return nil, apperrors.ErrForbidden.WithMessage("no tienes permiso para consultar pagos")
	}
	return uc.pagoRepo.FindPagoByID(ctx, id)
}

func (uc *CobroUseCase) ListarPagos(ctx context.Context, actor ports.Actor, filters ports.PagoFilters) ([]domain.Pago, int64, error) {
	if !canReadCobros(actor) {
		return nil, 0, apperrors.ErrForbidden.WithMessage("no tienes permiso para consultar pagos")
	}
	if !actor.HasPermission("registrar_pago") && !actor.HasPermission("gestionar_cobro_persuasivo") && !actor.HasPermission("gestionar_cobro_juridico") {
		filters.UserID = &actor.UserID
	}
	return uc.pagoRepo.ListPagos(ctx, filters)
}

func (uc *CobroUseCase) ActualizarPago(ctx context.Context, actor ports.Actor, id uint64, input ActualizarPagoInput) (*domain.Pago, error) {
	if !actor.HasPermission("registrar_pago") {
		return nil, apperrors.ErrForbidden.WithMessage("no tienes permiso para actualizar pagos")
	}
	pago, err := uc.pagoRepo.FindPagoByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if input.IDEntidad != nil {
		ok, err := uc.pagoRepo.EntidadExists(ctx, *input.IDEntidad)
		if err != nil {
			return nil, err
		}
		if !ok {
			return nil, apperrors.ErrNotFound.WithMessage("entidad no encontrada")
		}
		pago.IDEntidad = *input.IDEntidad
	}
	if input.TipoPago != nil {
		tipo, err := normalizeTipoPago(*input.TipoPago)
		if err != nil {
			return nil, err
		}
		pago.TipoPago = tipo
	}
	if input.EstadoPago != nil {
		estado, err := normalizeEstadoPago(*input.EstadoPago)
		if err != nil {
			return nil, err
		}
		pago.EstadoPago = estado
		pago.Conciliado = strings.EqualFold(estado, string(shareddomain.EstadoPagoConciliado))
	}
	if input.Descripcion != nil {
		pago.Descripcion = input.Descripcion
	}
	if input.Valor != nil {
		valor, err := parseValor(*input.Valor)
		if err != nil {
			return nil, err
		}
		pago.Valor = valor
	}
	if input.FechaPago != nil {
		fecha, err := parseDate(*input.FechaPago, "fecha_pago")
		if err != nil {
			return nil, err
		}
		pago.FechaPago = fecha
	}
	if input.PeriodoContable != nil {
		pago.PeriodoContable = input.PeriodoContable
	}
	if err := uc.pagoRepo.UpdatePago(ctx, pago); err != nil {
		return nil, err
	}
	return pago, nil
}

func (uc *CobroUseCase) EliminarPago(ctx context.Context, actor ports.Actor, id uint64) error {
	if !actor.HasPermission("registrar_pago") {
		return apperrors.ErrForbidden.WithMessage("no tienes permiso para eliminar pagos")
	}
	return uc.pagoRepo.SoftDeletePago(ctx, id)
}

func (uc *CobroUseCase) ConciliarPago(ctx context.Context, actor ports.Actor, id uint64, conciliado bool, estadoPago *string, descripcion *string) (*domain.Pago, error) {
	if !actor.HasPermission("realizar_conciliacion") {
		return nil, apperrors.ErrForbidden.WithMessage("no tienes permiso para conciliar pagos")
	}
	pago, err := uc.pagoRepo.FindPagoByID(ctx, id)
	if err != nil {
		return nil, err
	}
	pago.Conciliado = conciliado
	if estadoPago != nil && strings.TrimSpace(*estadoPago) != "" {
		estado, err := normalizeEstadoPago(*estadoPago)
		if err != nil {
			return nil, err
		}
		pago.EstadoPago = estado
	} else if conciliado {
		pago.EstadoPago = string(shareddomain.EstadoPagoConciliado)
	} else {
		pago.EstadoPago = string(shareddomain.EstadoPagoEnProceso)
	}
	if descripcion != nil {
		pago.Descripcion = descripcion
	}
	if err := uc.pagoRepo.UpdatePago(ctx, pago); err != nil {
		return nil, err
	}
	return pago, nil
}

func (uc *CobroUseCase) CrearSeguimiento(ctx context.Context, actor ports.Actor, input CrearSeguimientoInput) (*domain.SeguimientoCobro, error) {
	tipo, err := normalizeTipoSeguimiento(input.TipoSeguimiento)
	if err != nil {
		return nil, err
	}
	if err := ensureCanManageSeguimiento(actor, tipo); err != nil {
		return nil, err
	}
	ok, err := uc.pagoRepo.IncapacidadExists(ctx, input.IDIncapacidad)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, apperrors.ErrIncapacidadNotFound
	}
	if input.Resultado != nil {
		if _, err := normalizeResultadoSeguimiento(*input.Resultado); err != nil {
			return nil, err
		}
	}
	seguimiento := &domain.SeguimientoCobro{
		IDIncapacidad:   input.IDIncapacidad,
		TipoSeguimiento: tipo,
		Descripcion:     input.Descripcion,
		Resultado:       input.Resultado,
		GestionadoPor:   &actor.UserID,
	}
	if err := uc.seguimientoRepo.CreateSeguimiento(ctx, seguimiento); err != nil {
		return nil, err
	}
	return seguimiento, nil
}

func (uc *CobroUseCase) ObtenerSeguimiento(ctx context.Context, actor ports.Actor, id uint64) (*domain.SeguimientoCobro, error) {
	if !canReadCobros(actor) && !actor.HasPermission("consultar_historial") {
		return nil, apperrors.ErrForbidden.WithMessage("no tienes permiso para consultar seguimientos")
	}
	return uc.seguimientoRepo.FindSeguimientoByID(ctx, id)
}

func (uc *CobroUseCase) ListarSeguimientos(ctx context.Context, actor ports.Actor, filters ports.SeguimientoFilters) ([]domain.SeguimientoCobro, int64, error) {
	if !canReadCobros(actor) && !actor.HasPermission("consultar_historial") {
		return nil, 0, apperrors.ErrForbidden.WithMessage("no tienes permiso para consultar seguimientos")
	}
	return uc.seguimientoRepo.ListSeguimientos(ctx, filters)
}

func (uc *CobroUseCase) ActualizarSeguimiento(ctx context.Context, actor ports.Actor, id uint64, input ActualizarSeguimientoInput) (*domain.SeguimientoCobro, error) {
	seguimiento, err := uc.seguimientoRepo.FindSeguimientoByID(ctx, id)
	if err != nil {
		return nil, err
	}
	tipo := seguimiento.TipoSeguimiento
	if input.TipoSeguimiento != nil {
		tipo, err = normalizeTipoSeguimiento(*input.TipoSeguimiento)
		if err != nil {
			return nil, err
		}
	}
	if err := ensureCanManageSeguimiento(actor, tipo); err != nil {
		return nil, err
	}
	seguimiento.TipoSeguimiento = tipo
	if input.Descripcion != nil {
		seguimiento.Descripcion = input.Descripcion
	}
	if input.Resultado != nil {
		resultado, err := normalizeResultadoSeguimiento(*input.Resultado)
		if err != nil {
			return nil, err
		}
		seguimiento.Resultado = &resultado
	}
	if err := uc.seguimientoRepo.UpdateSeguimiento(ctx, seguimiento); err != nil {
		return nil, err
	}
	return seguimiento, nil
}

func (uc *CobroUseCase) ensureIncapacidadAndEntidad(ctx context.Context, incapacidadID, entidadID uint64) error {
	ok, err := uc.pagoRepo.IncapacidadExists(ctx, incapacidadID)
	if err != nil {
		return err
	}
	if !ok {
		return apperrors.ErrIncapacidadNotFound
	}
	ok, err = uc.pagoRepo.EntidadExists(ctx, entidadID)
	if err != nil {
		return err
	}
	if !ok {
		return apperrors.ErrNotFound.WithMessage("entidad no encontrada")
	}
	return nil
}

func canReadCobros(actor ports.Actor) bool {
	return actor.HasPermission("consultar_incapacidad") ||
		actor.HasPermission("registrar_pago") ||
		actor.HasPermission("realizar_conciliacion") ||
		actor.HasPermission("gestionar_cobro_persuasivo") ||
		actor.HasPermission("gestionar_cobro_juridico")
}

func ensureCanManageSeguimiento(actor ports.Actor, tipo string) error {
	if strings.EqualFold(tipo, string(shareddomain.TipoSegJuridico)) {
		if actor.HasPermission("gestionar_cobro_juridico") {
			return nil
		}
		return apperrors.ErrForbidden.WithMessage("no tienes permiso para gestionar cobro jurídico")
	}
	if strings.EqualFold(tipo, string(shareddomain.TipoSegPersuasivo)) ||
		strings.EqualFold(tipo, string(shareddomain.TipoSegCobroAdministrativo)) ||
		strings.EqualFold(tipo, string(shareddomain.TipoSegNormal)) ||
		strings.EqualFold(tipo, string(shareddomain.TipoSegPreventivo)) {
		if actor.HasPermission("gestionar_cobro_persuasivo") || actor.HasPermission("gestionar_cobro_juridico") {
			return nil
		}
		return apperrors.ErrForbidden.WithMessage("no tienes permiso para gestionar seguimiento de cobro")
	}
	return apperrors.ErrValidation.WithMessage("tipo de seguimiento inválido")
}

func normalizeTipoPago(value string) (string, error) {
	tipo := shareddomain.TipoPago(strings.TrimSpace(value))
	if !tipo.IsValid() {
		return "", apperrors.ErrValidation.WithMessage("tipo_pago inválido")
	}
	return string(tipo), nil
}

func normalizeEstadoPago(value string) (string, error) {
	estado := shareddomain.EstadoPago(strings.TrimSpace(value))
	if !estado.IsValid() {
		return "", apperrors.ErrValidation.WithMessage("estado_pago inválido")
	}
	return string(estado), nil
}

func normalizeTipoSeguimiento(value string) (string, error) {
	tipo := shareddomain.TipoSeguimiento(strings.TrimSpace(value))
	if !tipo.IsValid() {
		return "", apperrors.ErrValidation.WithMessage("tipo_seguimiento inválido")
	}
	return string(tipo), nil
}

func normalizeResultadoSeguimiento(value string) (string, error) {
	resultado := shareddomain.ResultadoSeguimiento(strings.TrimSpace(value))
	if !resultado.IsValid() {
		return "", apperrors.ErrValidation.WithMessage("resultado de seguimiento inválido")
	}
	return string(resultado), nil
}

func parseValor(value string) (decimal.Decimal, error) {
	valor, err := decimal.NewFromString(strings.TrimSpace(value))
	if err != nil {
		return decimal.Zero, apperrors.ErrValidation.WithMessage("valor debe ser numérico")
	}
	if !valor.GreaterThanOrEqual(decimal.Zero) {
		return decimal.Zero, apperrors.ErrValidation.WithMessage("valor debe ser mayor o igual a cero")
	}
	return valor, nil
}

func parseDate(value, field string) (time.Time, error) {
	parsed, err := time.Parse(dateLayout, strings.TrimSpace(value))
	if err != nil {
		return time.Time{}, apperrors.ErrValidation.WithMessage(field + " debe tener formato YYYY-MM-DD")
	}
	return parsed, nil
}
