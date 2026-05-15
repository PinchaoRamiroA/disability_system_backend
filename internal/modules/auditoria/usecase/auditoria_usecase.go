package usecase

import (
	"context"

	"disability_system_backend/internal/modules/auditoria/domain"
	"disability_system_backend/internal/modules/auditoria/dto"
	"disability_system_backend/internal/modules/auditoria/ports"
	apperrors "disability_system_backend/internal/shared/errors"
)

type AuditoriaUseCase struct {
	repo ports.AuditoriaRepository
}

func NewAuditoriaUseCase(repo ports.AuditoriaRepository) *AuditoriaUseCase {
	return &AuditoriaUseCase{repo: repo}
}

func (uc *AuditoriaUseCase) Crear(ctx context.Context, actor ports.Actor, req dto.CrearAuditoriaRequest) (*domain.Auditoria, error) {
	// Solo admins o el propio sistema deberían crear auditoría
	// Asumiremos que si el sistema la llama desde otro módulo, el context u actor ya está autorizado.
	// Por ahora lo dejaremos crear sin permisos restrictivos, o puedes restringirlo según la necesidad.
	
	idUsuario := req.IDUsuario
	if idUsuario == nil && actor != nil {
		userID := actor.GetUserID()
		idUsuario = &userID
	}

	auditoria := &domain.Auditoria{
		IDUsuario:      idUsuario,
		IDIncapacidad:  req.IDIncapacidad,
		TipoAccion:     req.TipoAccion,
		Modulo:         req.Modulo,
		Descripcion:    req.Descripcion,
		CambioAnterior: req.CambioAnterior,
		CambioNuevo:    req.CambioNuevo,
	}

	if err := uc.repo.Create(ctx, auditoria); err != nil {
		return nil, err
	}

	return auditoria, nil
}

func (uc *AuditoriaUseCase) Listar(ctx context.Context, actor ports.Actor, query dto.ListarAuditoriaQuery) ([]domain.Auditoria, int64, error) {
	// Verificar permisos (por ejemplo, "consultar_auditoria")
	if actor != nil && !actor.HasPermission("consultar_auditoria") {
		return nil, 0, apperrors.ErrForbidden.WithMessage("No tienes permiso para consultar auditoría")
	}

	return uc.repo.List(ctx, query)
}
