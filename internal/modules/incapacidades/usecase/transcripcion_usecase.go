package usecase

import (
	"context"
	"errors"
	"time"

	"disability_system_backend/internal/modules/incapacidades/domain"
	"disability_system_backend/internal/modules/incapacidades/ports"
)

var (
	ErrIncapacidadNotFound     = errors.New("incapacidad no encontrada")
	ErrTranscripcionNotAllowed = errors.New("transcripcion no permitida para este estado")
	ErrTranscripcionVencida    = errors.New("la fecha limite de transcripcion ha vencido")
)

type TranscripcionUseCase struct {
	repo ports.IncapacidadRepository
}

func NewTranscripcionUseCase(repo ports.IncapacidadRepository) *TranscripcionUseCase {
	return &TranscripcionUseCase{repo: repo}
}

func (uc *TranscripcionUseCase) Transcribir(ctx context.Context, idIncapacidad uint64, actor ports.Actor, obs *string) (*domain.Incapacidad, error) {
	incapacidad, err := uc.repo.FindByID(ctx, idIncapacidad)
	if err != nil {
		return nil, ErrIncapacidadNotFound
	}

	if incapacidad.EstadoTranscripcion == "completado" {
		return nil, ErrTranscripcionNotAllowed
	}

	if incapacidad.FechaLimiteTranscripcion != nil && time.Now().After(*incapacidad.FechaLimiteTranscripcion) {
		incapacidad.EstadoTranscripcion = "vencida"
		uc.repo.Update(ctx, incapacidad)
		return nil, ErrTranscripcionVencida
	}

	now := time.Now()
	incapacidad.FechaTranscripcion = &now
	incapacidad.TranscritoPor = &actor.UserID
	incapacidad.EstadoTranscripcion = "completado"
	if obs != nil {
		incapacidad.ObservacionesTranscripcion = obs
	}

	if err := uc.repo.Update(ctx, incapacidad); err != nil {
		return nil, err
	}

	return incapacidad, nil
}

func (uc *TranscripcionUseCase) MarcarEnProceso(ctx context.Context, idIncapacidad uint64, actor ports.Actor) (*domain.Incapacidad, error) {
	incapacidad, err := uc.repo.FindByID(ctx, idIncapacidad)
	if err != nil {
		return nil, ErrIncapacidadNotFound
	}

	if incapacidad.EstadoTranscripcion == "completado" {
		return nil, ErrTranscripcionNotAllowed
	}

	incapacidad.EstadoTranscripcion = "en_proceso"
	incapacidad.TranscritoPor = &actor.UserID

	if err := uc.repo.Update(ctx, incapacidad); err != nil {
		return nil, err
	}

	return incapacidad, nil
}

func (uc *TranscripcionUseCase) ListarPendientes(ctx context.Context, estado string, page, limit int) ([]domain.Incapacidad, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	filters := ports.IncapacidadFilters{
		IncludeDeleted: false,
		Page:           page,
		Limit:          limit,
	}

	return uc.repo.List(ctx, filters)
}

func (uc *TranscripcionUseCase) CalcularFechaLimite(fechaCreacion time.Time, diasHabiles int) time.Time {
	fecha := fechaCreacion
	for diasHabiles > 0 {
		fecha = fecha.AddDate(0, 0, 1)
		if fecha.Weekday() != time.Saturday && fecha.Weekday() != time.Sunday {
			diasHabiles--
		}
	}
	return fecha
}

func (uc *TranscripcionUseCase) ObtenerAlertaVencimiento(incapacidad *domain.Incapacidad) *string {
	if incapacidad.EstadoTranscripcion == "completado" {
		return nil
	}

	if incapacidad.FechaLimiteTranscripcion == nil {
		return nil
	}

	diasRestantes := time.Until(*incapacidad.FechaLimiteTranscripcion).Hours() / 24

	var alerta string
	if diasRestantes < 0 {
		alerta = "VENCIDA"
	} else if diasRestantes <= 1 {
		alerta = "URGENTE - Vence mañana"
	} else if diasRestantes <= 3 {
		alerta = "ALERTA - Vence en " + formatDias(int(diasRestantes))
	}

	if alerta == "" {
		return nil
	}
	return &alerta
}

func formatDias(d int) string {
	if d == 1 {
		return "1 día"
	}
	return string(rune('0'+d%10)) + " días"
}
