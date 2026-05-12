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
	historialRepo ports.HistorialRepository
}

func NewDocumentoUseCase(repo ports.DocumentoRepository, historialRepo ports.HistorialRepository) *DocumentoUseCase {
	return &DocumentoUseCase{repo: repo, historialRepo: historialRepo}
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

	uc.registrarHistorial(ctx, input.IDIncapacidad, actor.UserID, "documento_subido", "Documento '"+input.Nombre+"' subido al sistema")

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
	uc.registrarHistorial(ctx, documento.IDIncapacidad, actor.UserID, "documento_validado", descripcion)

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

func (uc *DocumentoUseCase) registrarHistorial(ctx context.Context, incapacidadID uint64, gestorID uint64, tipoNombre, descripcion string) {
	tipoID := uint64(1)
	uc.historialRepo.Create(ctx, &domain.Historial{
		IDIncapacidad:   incapacidadID,
		IDTipoHistorial: tipoID,
		Descripcion:     descripcion,
		GestorID:        &gestorID,
	})
}

type HistorialUseCase struct {
	repo ports.HistorialRepository
}

func NewHistorialUseCase(repo ports.HistorialRepository) *HistorialUseCase {
	return &HistorialUseCase{repo: repo}
}

func (uc *HistorialUseCase) Listar(ctx context.Context, actor ports.Actor, incapacidadID uint64, tipoID *uint64, page, limit int) ([]domain.Historial, int64, error) {
	if !actor.HasPermission("consultar_incapacidad") && !actor.HasPermission("consultar_historial") {
		return nil, 0, apperrors.ErrForbidden.WithMessage("no tienes permiso para consultar historial")
	}
	return uc.repo.List(ctx, incapacidadID, tipoID, page, limit)
}

func (uc *HistorialUseCase) Crear(ctx context.Context, actor ports.Actor, input struct {
	IDIncapacidad  uint64
	IDTipoHistorial uint64
	Descripcion    string
}) (*domain.Historial, error) {
	if !actor.HasPermission("crear_incapacidad") && !actor.HasPermission("editar_incapacidad") {
		return nil, apperrors.ErrForbidden.WithMessage("no tienes permiso para crear entradas de historial")
	}

	tipo, err := uc.repo.FindTipoByID(ctx, input.IDTipoHistorial)
	if err != nil {
		return nil, err
	}

	historial := &domain.Historial{
		IDIncapacidad:   input.IDIncapacidad,
		IDTipoHistorial: input.IDTipoHistorial,
		Descripcion:     input.Descripcion,
		GestorID:        &actor.UserID,
	}

	if err := uc.repo.Create(ctx, historial); err != nil {
		return nil, err
	}

	_ = tipo
	return historial, nil
}