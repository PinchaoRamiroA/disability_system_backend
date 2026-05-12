package usecase

import (
	"context"
	"time"

	"disability_system_backend/internal/modules/incapacidades/domain"
	"disability_system_backend/internal/modules/incapacidades/ports"
	apperrors "disability_system_backend/internal/shared/errors"
)

type DocumentoUseCase struct {
	repo          ports.DocumentoRepository
	historialSvc  ports.HistorialService
}

func NewDocumentoUseCase(repo ports.DocumentoRepository, historialSvc ports.HistorialService) *DocumentoUseCase {
	return &DocumentoUseCase{repo: repo, historialSvc: historialSvc}
}

func (uc *DocumentoUseCase) Subir(ctx context.Context, actor ports.Actor, input struct {
	IDIncapacidad uint64
	Nombre        string
	Tipo          string
	URL           string
	Formato       string
}) (*domain.Documento, error) {
	if !actor.HasPermission("crear_incapacidad") && !actor.HasPermission("editar_incapacidad") {
		return nil, apperrors.ErrForbidden.WithMessage("no tienes permiso para subir documentos")
	}
	exists, err := uc.repo.ExistsIncapacidad(ctx, input.IDIncapacidad)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, apperrors.ErrIncapacidadNotFound
	}

	documento := &domain.Documento{
		IDIncapacidad: input.IDIncapacidad,
		Nombre:        input.Nombre,
		Tipo:          input.Tipo,
		URL:           input.URL,
		Formato:       input.Formato,
		Estado:        "Pendiente",
		FechaCarga:    time.Now(),
	}

	if err := uc.repo.Create(ctx, documento); err != nil {
		return nil, err
	}

	descripcion := "Documento '" + input.Nombre + "' subido al sistema"
	uc.historialSvc.CreateEntry(ctx, input.IDIncapacidad, 1, descripcion, &actor.UserID)

	return documento, nil
}

func (uc *DocumentoUseCase) Validar(ctx context.Context, actor ports.Actor, id uint64, estado, comentario string) (*domain.Documento, error) {
	if !actor.HasPermission("validar_documentos") && !actor.HasPermission("editar_incapacidad") {
		return nil, apperrors.ErrForbidden.WithMessage("no tienes permiso para validar documentos")
	}

	documento, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	validStates := map[string]bool{"Validado": true, "Rechazado": true, "Incompleto": true}
	if !validStates[estado] {
		return nil, apperrors.ErrValidation.WithMessage("estado debe ser: Validado, Rechazado o Incompleto")
	}

	documento.Estado = estado
	documento.Comentario = &comentario
	documento.ValidadoPor = &actor.UserID
	now := time.Now()
	documento.FechaValidacion = &now

	if err := uc.repo.Update(ctx, documento); err != nil {
		return nil, err
	}

	descripcion := "Documento '" + documento.Nombre + "' " + estado
	if comentario != "" {
		descripcion += ": " + comentario
	}
	uc.historialSvc.CreateEntry(ctx, documento.IDIncapacidad, 1, descripcion, &actor.UserID)

	return documento, nil
}

func (uc *DocumentoUseCase) Listar(ctx context.Context, actor ports.Actor, incapacidadID uint64, estado, tipo string, page, limit int) ([]domain.Documento, int64, error) {
	if !actor.HasPermission("consultar_incapacidad") {
		return nil, 0, apperrors.ErrForbidden.WithMessage("no tienes permiso para consultar documentos")
	}
	return uc.repo.List(ctx, incapacidadID, estado, tipo, page, limit)
}

func (uc *DocumentoUseCase) Eliminar(ctx context.Context, actor ports.Actor, id uint64) error {
	if !actor.HasPermission("editar_incapacidad") {
		return apperrors.ErrForbidden.WithMessage("no tienes permiso para eliminar documentos")
	}
	return uc.repo.Delete(ctx, id)
}