package usecase

import (
	"context"
	"strings"
	"time"

	"disability_system_backend/internal/modules/incapacidades/domain"
	"disability_system_backend/internal/modules/incapacidades/ports"
	apperrors "disability_system_backend/internal/shared/errors"
)

const dateLayout = "2006-01-02"

type CrearIncapacidadInput struct {
	IDUsuario       uint64
	IDEstado        *uint64
	IDTipo          uint64
	IDEntidad       uint64
	CanalRecepcion  string
	Titulo          string
	FechaInicio     string
	FechaFin        *string
	Origen          string
	FechaRadicacion *string
	FechaPago       *string
	Observaciones   *string
}

type ActualizarIncapacidadInput struct {
	IDUsuario       *uint64
	IDTipo          *uint64
	IDEntidad       *uint64
	CanalRecepcion  *string
	Titulo          *string
	FechaInicio     *string
	FechaFin        *string
	Origen          *string
	FechaRadicacion *string
	FechaPago       *string
	Observaciones   *string
}

type IncapacidadUseCase struct {
	repo             ports.IncapacidadRepository
	docService       *IncapacidadDocumentosService
	historialFn      func(ctx context.Context, incapacidadID uint64, tipoID uint64, descripcion string, gestorID *uint64) error
	faltanteNotifier ports.DocumentoFaltanteNotifier
}

func NewIncapacidadUseCase(repo ports.IncapacidadRepository) *IncapacidadUseCase {
	uc := &IncapacidadUseCase{
		repo:       repo,
		docService: NewIncapacidadDocumentosService(repo),
	}
	return uc
}

func (uc *IncapacidadUseCase) SetHistorialService(fn func(ctx context.Context, incapacidadID uint64, tipoID uint64, descripcion string, gestorID *uint64) error) {
	uc.historialFn = fn
}

func (uc *IncapacidadUseCase) SetDocumentoFaltanteNotifier(notifier ports.DocumentoFaltanteNotifier) {
	uc.faltanteNotifier = notifier
}

func (uc *IncapacidadUseCase) Crear(ctx context.Context, actor ports.Actor, input CrearIncapacidadInput) (*domain.Incapacidad, error) {
	if !actor.HasPermission("crear_incapacidad") {
		return nil, apperrors.ErrForbidden.WithMessage("no tienes permiso para crear incapacidades")
	}
	if strings.TrimSpace(input.Titulo) == "" || strings.TrimSpace(input.Origen) == "" {
		return nil, apperrors.ErrValidation.WithMessage("título y origen son obligatorios")
	}
	if input.IDUsuario == 0 {
		input.IDUsuario = actor.UserID
	}
	if err := uc.ensureReferences(ctx, input.IDUsuario, input.IDTipo, input.IDEntidad); err != nil {
		return nil, err
	}

	estadoID := input.IDEstado
	if estadoID == nil {
		estado, err := uc.repo.FindEstadoByName(ctx, "Recibida")
		if err != nil {
			return nil, err
		}
		estadoID = &estado.IDEstado
	} else if _, err := uc.repo.FindEstadoByID(ctx, *estadoID); err != nil {
		return nil, err
	}

	fechaInicio, err := parseRequiredDate(input.FechaInicio, "fecha_inicio")
	if err != nil {
		return nil, err
	}
	fechaFin, err := parseOptionalDate(input.FechaFin, "fecha_fin")
	if err != nil {
		return nil, err
	}
	fechaRadicacion, err := parseOptionalDate(input.FechaRadicacion, "fecha_radicacion")
	if err != nil {
		return nil, err
	}
	fechaPago, err := parseOptionalDate(input.FechaPago, "fecha_pago")
	if err != nil {
		return nil, err
	}
	if fechaFin != nil && fechaFin.Before(fechaInicio) {
		return nil, apperrors.ErrValidation.WithMessage("fecha_fin no puede ser anterior a fecha_inicio")
	}

	incapacidad := &domain.Incapacidad{
		IDUsuario:       input.IDUsuario,
		IDEstado:        *estadoID,
		IDTipo:          input.IDTipo,
		IDEntidad:       input.IDEntidad,
		CanalRecepcion:  strings.TrimSpace(input.CanalRecepcion),
		Titulo:          strings.TrimSpace(input.Titulo),
		FechaInicio:     fechaInicio,
		FechaFin:        fechaFin,
		Origen:          strings.TrimSpace(input.Origen),
		FechaRadicacion: fechaRadicacion,
		FechaPago:       fechaPago,
		Observaciones:   input.Observaciones,
		CreatedBy:       &actor.UserID,
	}
	if err := uc.repo.Create(ctx, incapacidad); err != nil {
		return nil, err
	}

	if uc.historialFn != nil {
		tipoIncapacidad, _ := uc.repo.FindTipoByID(ctx, input.IDTipo)
		tipoNombre := ""
		if tipoIncapacidad != nil {
			tipoNombre = tipoIncapacidad.Nombre
		}
		descripcion := "Incapacidad creada - Tipo: " + tipoNombre
		_ = uc.historialFn(ctx, incapacidad.IDIncapacidad, 1, descripcion, &actor.UserID)
	}
	uc.notificarDocumentosFaltantesIniciales(ctx, incapacidad)

	return incapacidad, nil
}

func (uc *IncapacidadUseCase) notificarDocumentosFaltantesIniciales(ctx context.Context, incapacidad *domain.Incapacidad) {
	if uc.faltanteNotifier == nil || uc.docService == nil {
		return
	}
	requeridos, err := uc.docService.ObtenerDocumentosRequeridos(ctx, incapacidad.IDTipo)
	if err != nil {
		return
	}
	faltantes := make([]domain.TipoDocumento, 0, len(requeridos))
	for _, documento := range requeridos {
		if documento.Requerido {
			faltantes = append(faltantes, documento)
		}
	}
	if len(faltantes) == 0 {
		return
	}
	_ = uc.faltanteNotifier.NotificarDocumentosFaltantes(ctx, incapacidad.IDUsuario, incapacidad.IDIncapacidad, faltantes)
}

func (uc *IncapacidadUseCase) Obtener(ctx context.Context, actor ports.Actor, id uint64) (*domain.Incapacidad, error) {
	if !actor.HasPermission("consultar_incapacidad") {
		return nil, apperrors.ErrForbidden.WithMessage("no tienes permiso para consultar incapacidades")
	}
	incapacidad, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if err := ensureCanRead(actor, incapacidad); err != nil {
		return nil, err
	}
	return incapacidad, nil
}

func (uc *IncapacidadUseCase) Listar(ctx context.Context, actor ports.Actor, filters ports.IncapacidadFilters) ([]domain.Incapacidad, int64, error) {
	if !actor.HasPermission("consultar_incapacidad") {
		return nil, 0, apperrors.ErrForbidden.WithMessage("no tienes permiso para consultar incapacidades")
	}
	if !actor.CanManageIncapacidades() {
		filters.UserID = &actor.UserID
	}
	return uc.repo.List(ctx, filters)
}

func (uc *IncapacidadUseCase) Actualizar(ctx context.Context, actor ports.Actor, id uint64, input ActualizarIncapacidadInput) (*domain.Incapacidad, error) {
	if !actor.HasPermission("editar_incapacidad") {
		return nil, apperrors.ErrForbidden.WithMessage("no tienes permiso para editar incapacidades")
	}
	incapacidad, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if input.IDUsuario != nil {
		ok, err := uc.repo.ExistsUsuario(ctx, *input.IDUsuario)
		if err != nil {
			return nil, err
		}
		if !ok {
			return nil, apperrors.ErrUserNotFound
		}
		incapacidad.IDUsuario = *input.IDUsuario
	}
	if input.IDTipo != nil {
		if _, err := uc.repo.FindTipoByID(ctx, *input.IDTipo); err != nil {
			return nil, err
		}
		incapacidad.IDTipo = *input.IDTipo
	}
	if input.IDEntidad != nil {
		if _, err := uc.repo.FindEntidadByID(ctx, *input.IDEntidad); err != nil {
			return nil, err
		}
		incapacidad.IDEntidad = *input.IDEntidad
	}
	if input.CanalRecepcion != nil {
		incapacidad.CanalRecepcion = strings.TrimSpace(*input.CanalRecepcion)
	}
	if input.Titulo != nil {
		if strings.TrimSpace(*input.Titulo) == "" {
			return nil, apperrors.ErrValidation.WithMessage("título no puede estar vacío")
		}
		incapacidad.Titulo = strings.TrimSpace(*input.Titulo)
	}
	if input.Origen != nil {
		if strings.TrimSpace(*input.Origen) == "" {
			return nil, apperrors.ErrValidation.WithMessage("origen no puede estar vacío")
		}
		incapacidad.Origen = strings.TrimSpace(*input.Origen)
	}
	if input.FechaInicio != nil {
		fechaInicio, err := parseRequiredDate(*input.FechaInicio, "fecha_inicio")
		if err != nil {
			return nil, err
		}
		incapacidad.FechaInicio = fechaInicio
	}
	if input.FechaFin != nil {
		fechaFin, err := parseOptionalDate(input.FechaFin, "fecha_fin")
		if err != nil {
			return nil, err
		}
		incapacidad.FechaFin = fechaFin
	}
	if input.FechaRadicacion != nil {
		fechaRadicacion, err := parseOptionalDate(input.FechaRadicacion, "fecha_radicacion")
		if err != nil {
			return nil, err
		}
		incapacidad.FechaRadicacion = fechaRadicacion
	}
	if input.FechaPago != nil {
		fechaPago, err := parseOptionalDate(input.FechaPago, "fecha_pago")
		if err != nil {
			return nil, err
		}
		incapacidad.FechaPago = fechaPago
	}
	if incapacidad.FechaFin != nil && incapacidad.FechaFin.Before(incapacidad.FechaInicio) {
		return nil, apperrors.ErrValidation.WithMessage("fecha_fin no puede ser anterior a fecha_inicio")
	}
	if input.Observaciones != nil {
		incapacidad.Observaciones = input.Observaciones
	}

	if err := uc.repo.Update(ctx, incapacidad); err != nil {
		return nil, err
	}
	return incapacidad, nil
}

func (uc *IncapacidadUseCase) CambiarEstado(ctx context.Context, actor ports.Actor, id, estadoID uint64, observaciones *string) (*domain.Incapacidad, error) {
	if !actor.HasPermission("editar_incapacidad") && !actor.HasPermission("archivar_incapacidad") {
		return nil, apperrors.ErrForbidden.WithMessage("no tienes permiso para cambiar estados de incapacidades")
	}
	incapacidad, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if incapacidad.Estado != nil && !incapacidad.Estado.PermiteTransicion {
		return nil, apperrors.ErrConflict.WithMessage("el estado actual no permite transición")
	}
	estado, err := uc.repo.FindEstadoByID(ctx, estadoID)
	if err != nil {
		return nil, err
	}
	if strings.EqualFold(estado.Nombre, "Archivada") && !actor.HasPermission("archivar_incapacidad") {
		return nil, apperrors.ErrForbidden.WithMessage("no tienes permiso para archivar incapacidades")
	}
	incapacidad.IDEstado = estadoID
	if observaciones != nil {
		incapacidad.Observaciones = observaciones
	}
	if err := uc.repo.Update(ctx, incapacidad); err != nil {
		return nil, err
	}
	return incapacidad, nil
}

func (uc *IncapacidadUseCase) Archivar(ctx context.Context, actor ports.Actor, id uint64) error {
	if !actor.HasPermission("archivar_incapacidad") {
		return apperrors.ErrForbidden.WithMessage("no tienes permiso para archivar incapacidades")
	}
	return uc.repo.SoftDelete(ctx, id)
}

func (uc *IncapacidadUseCase) ListarEstados(ctx context.Context, actor ports.Actor) ([]domain.EstadoIncapacidad, error) {
	if !actor.HasPermission("consultar_incapacidad") {
		return nil, apperrors.ErrForbidden.WithMessage("no tienes permiso para consultar catálogos")
	}
	return uc.repo.ListEstados(ctx)
}

func (uc *IncapacidadUseCase) ListarTipos(ctx context.Context, actor ports.Actor) ([]domain.TipoIncapacidad, error) {
	if !actor.HasPermission("consultar_incapacidad") {
		return nil, apperrors.ErrForbidden.WithMessage("no tienes permiso para consultar catálogos")
	}
	return uc.repo.ListTipos(ctx)
}

func (uc *IncapacidadUseCase) ListarEntidades(ctx context.Context, actor ports.Actor) ([]domain.Entidad, error) {
	if !actor.HasPermission("consultar_incapacidad") {
		return nil, apperrors.ErrForbidden.WithMessage("no tienes permiso para consultar catálogos")
	}
	return uc.repo.ListEntidades(ctx)
}

func (uc *IncapacidadUseCase) ListarEstadosDocumento(ctx context.Context, actor ports.Actor) ([]domain.EstadoDocumento, error) {
	if !actor.HasPermission("consultar_incapacidad") {
		return nil, apperrors.ErrForbidden.WithMessage("no tienes permiso para consultar catálogos")
	}
	return uc.repo.ListEstadosDocumento(ctx)
}

func (uc *IncapacidadUseCase) ListarTiposDocumento(ctx context.Context, actor ports.Actor) ([]domain.TipoDocumento, error) {
	if !actor.HasPermission("consultar_incapacidad") {
		return nil, apperrors.ErrForbidden.WithMessage("no tienes permiso para consultar catálogos")
	}
	return uc.repo.ListTiposDocumento(ctx)
}

func (uc *IncapacidadUseCase) ListarTiposPago(ctx context.Context, actor ports.Actor) ([]domain.TipoPago, error) {
	if !actor.HasPermission("consultar_incapacidad") && !actor.HasPermission("registrar_pago") {
		return nil, apperrors.ErrForbidden.WithMessage("no tienes permiso para consultar catálogos de pago")
	}
	return uc.repo.ListTiposPago(ctx)
}

func (uc *IncapacidadUseCase) ObtenerDocumentosRequeridos(ctx context.Context, actor ports.Actor, tipoID uint64) ([]domain.TipoDocumento, error) {
	if !actor.HasPermission("consultar_incapacidad") {
		return nil, apperrors.ErrForbidden.WithMessage("no tienes permiso para consultar catálogos")
	}
	return uc.docService.ObtenerDocumentosRequeridos(ctx, tipoID)
}

func (uc *IncapacidadUseCase) ObtenerInfoPlazos(ctx context.Context, actor ports.Actor, incapacidadID uint64) (*PlazosInfo, error) {
	if !actor.HasPermission("consultar_incapacidad") {
		return nil, apperrors.ErrForbidden.WithMessage("no tienes permiso para consultar incapacidades")
	}
	incapacidad, err := uc.repo.FindByID(ctx, incapacidadID)
	if err != nil {
		return nil, err
	}

	tipoIncapacidad, err := uc.repo.FindTipoByID(ctx, incapacidad.IDTipo)
	if err != nil {
		return nil, err
	}

	documentosRequeridos, err := uc.docService.ObtenerDocumentosRequeridos(ctx, incapacidad.IDTipo)
	if err != nil {
		return nil, err
	}

	plazoEntregaDias, _ := uc.docService.ObtenerPlazoEntrega(ctx, incapacidad.IDTipo, incapacidad.CanalRecepcion)
	fechaLimiteEntrega, _ := uc.docService.ObtenerFechaLimiteEntrega(ctx, incapacidad.CreatedAt, incapacidad.IDTipo, incapacidad.CanalRecepcion)

	plazoTranscripcionDias, _ := uc.docService.ObtenerPlazoTranscripcion(ctx, incapacidad.IDEntidad)
	fechaLimiteTranscripcion, _ := uc.docService.ObtenerFechaLimiteTranscripcion(ctx, incapacidad.CreatedAt, incapacidad.IDEntidad)

	tiempoMaximoPagoDias, _ := uc.docService.ObtenerTiempoMaximoPago(ctx, incapacidad.IDEntidad)
	fechaLimitePago := incapacidad.FechaInicio.AddDate(0, 0, tiempoMaximoPagoDias)

	alertas := uc.docService.ObtenerAlertasVencimientos(incapacidad.FechaInicio)

	return &PlazosInfo{
		IncapacidadID:            incapacidad.IDIncapacidad,
		TipoIncapacidad:          tipoIncapacidad.Nombre,
		DocumentosRequeridos:     documentosRequeridos,
		PlazoEntregaDias:         plazoEntregaDias,
		FechaLimiteEntrega:       fechaLimiteEntrega,
		PlazoTranscripcionDias:   plazoTranscripcionDias,
		FechaLimiteTranscripcion: fechaLimiteTranscripcion,
		TiempoMaximoPagoDias:     tiempoMaximoPagoDias,
		FechaLimitePago:          fechaLimitePago,
		DiasTranscurridos:        uc.docService.ObtenerDiasTranscurridos(incapacidad.FechaInicio),
		AlertasVencimientos:      alertas,
	}, nil
}

type PlazosInfo struct {
	IncapacidadID            uint64
	TipoIncapacidad          string
	DocumentosRequeridos     []domain.TipoDocumento
	PlazoEntregaDias         int
	FechaLimiteEntrega       time.Time
	PlazoTranscripcionDias   int
	FechaLimiteTranscripcion time.Time
	TiempoMaximoPagoDias     int
	FechaLimitePago          time.Time
	DiasTranscurridos        int
	AlertasVencimientos      []string
}

func (uc *IncapacidadUseCase) ensureReferences(ctx context.Context, userID, tipoID, entidadID uint64) error {
	ok, err := uc.repo.ExistsUsuario(ctx, userID)
	if err != nil {
		return err
	}
	if !ok {
		return apperrors.ErrUserNotFound
	}
	if _, err := uc.repo.FindTipoByID(ctx, tipoID); err != nil {
		return err
	}
	if _, err := uc.repo.FindEntidadByID(ctx, entidadID); err != nil {
		return err
	}
	return nil
}

func ensureCanRead(actor ports.Actor, incapacidad *domain.Incapacidad) error {
	if actor.CanManageIncapacidades() || incapacidad.IDUsuario == actor.UserID {
		return nil
	}
	return apperrors.ErrForbidden.WithMessage("solo puedes consultar tus propias incapacidades")
}

func parseRequiredDate(value, field string) (time.Time, error) {
	parsed, err := time.Parse(dateLayout, strings.TrimSpace(value))
	if err != nil {
		return time.Time{}, apperrors.ErrValidation.WithMessage(field + " debe tener formato YYYY-MM-DD")
	}
	return parsed, nil
}

func parseOptionalDate(value *string, field string) (*time.Time, error) {
	if value == nil || strings.TrimSpace(*value) == "" {
		return nil, nil
	}
	parsed, err := parseRequiredDate(*value, field)
	if err != nil {
		return nil, err
	}
	return &parsed, nil
}
