package http

import (
	"errors"
	"strconv"

	"disability_system_backend/internal/modules/incapacidades/dto"
	"disability_system_backend/internal/modules/incapacidades/mapper"
	"disability_system_backend/internal/modules/incapacidades/ports"
	"disability_system_backend/internal/modules/incapacidades/usecase"
	apperrors "disability_system_backend/internal/shared/errors"
	"disability_system_backend/internal/shared/response"

	"github.com/gin-gonic/gin"
)

type IncapacidadHandler struct {
	useCase *usecase.IncapacidadUseCase
}

func NewIncapacidadHandler(useCase *usecase.IncapacidadUseCase) *IncapacidadHandler {
	return &IncapacidadHandler{useCase: useCase}
}

func (h *IncapacidadHandler) Crear(c *gin.Context) {
	actor, err := actorFromGin(c)
	if err != nil {
		handleError(c, err)
		return
	}

	var req dto.CrearIncapacidadRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, "datos inválidos", err.Error())
		return
	}

	userID := actor.UserID
	if req.IDUsuario != nil {
		userID = *req.IDUsuario
	}

	incapacidad, err := h.useCase.Crear(c.Request.Context(), actor, usecase.CrearIncapacidadInput{
		IDUsuario:       userID,
		IDEstado:        req.IDEstado,
		IDTipo:          req.IDTipo,
		IDEntidad:       req.IDEntidad,
		CanalRecepcion:  req.CanalRecepcion,
		Titulo:          req.Titulo,
		FechaInicio:     req.FechaInicio,
		FechaFin:        req.FechaFin,
		Origen:          req.Origen,
		FechaRadicacion: req.FechaRadicacion,
		FechaPago:       req.FechaPago,
		Observaciones:   req.Observaciones,
	})
	if err != nil {
		handleError(c, err)
		return
	}

	response.Created(c, mapper.ToIncapacidadResponse(incapacidad), "incapacidad creada exitosamente")
}

func (h *IncapacidadHandler) Obtener(c *gin.Context) {
	actor, err := actorFromGin(c)
	if err != nil {
		handleError(c, err)
		return
	}
	id, err := parseIDParam(c, "id")
	if err != nil {
		handleError(c, err)
		return
	}

	incapacidad, err := h.useCase.Obtener(c.Request.Context(), actor, id)
	if err != nil {
		handleError(c, err)
		return
	}
	response.Success(c, mapper.ToIncapacidadResponse(incapacidad), "incapacidad encontrada")
}

func (h *IncapacidadHandler) Listar(c *gin.Context) {
	actor, err := actorFromGin(c)
	if err != nil {
		handleError(c, err)
		return
	}

	var query dto.ListarIncapacidadesQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.ValidationError(c, "filtros inválidos", err.Error())
		return
	}

	items, total, err := h.useCase.Listar(c.Request.Context(), actor, ports.IncapacidadFilters{
		UserID:         query.IDUsuario,
		EstadoID:       query.IDEstado,
		TipoID:         query.IDTipo,
		EntidadID:      query.IDEntidad,
		Origen:         query.Origen,
		CanalRecepcion: query.CanalRecepcion,
		Page:           query.Page,
		Limit:          query.Limit,
	})
	if err != nil {
		handleError(c, err)
		return
	}

	page := query.Page
	if page < 1 {
		page = 1
	}
	limit := query.Limit
	if limit < 1 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}
	response.Paginated(c, mapper.ToIncapacidadResponses(items), total, int64(page), int64(limit))
}

func (h *IncapacidadHandler) Actualizar(c *gin.Context) {
	actor, err := actorFromGin(c)
	if err != nil {
		handleError(c, err)
		return
	}
	id, err := parseIDParam(c, "id")
	if err != nil {
		handleError(c, err)
		return
	}

	var req dto.ActualizarIncapacidadRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, "datos inválidos", err.Error())
		return
	}

	incapacidad, err := h.useCase.Actualizar(c.Request.Context(), actor, id, usecase.ActualizarIncapacidadInput{
		IDUsuario:       req.IDUsuario,
		IDTipo:          req.IDTipo,
		IDEntidad:       req.IDEntidad,
		CanalRecepcion:  req.CanalRecepcion,
		Titulo:          req.Titulo,
		FechaInicio:     req.FechaInicio,
		FechaFin:        req.FechaFin,
		Origen:          req.Origen,
		FechaRadicacion: req.FechaRadicacion,
		FechaPago:       req.FechaPago,
		Observaciones:   req.Observaciones,
	})
	if err != nil {
		handleError(c, err)
		return
	}

	response.Success(c, mapper.ToIncapacidadResponse(incapacidad), "incapacidad actualizada")
}

func (h *IncapacidadHandler) CambiarEstado(c *gin.Context) {
	actor, err := actorFromGin(c)
	if err != nil {
		handleError(c, err)
		return
	}
	id, err := parseIDParam(c, "id")
	if err != nil {
		handleError(c, err)
		return
	}

	var req dto.CambiarEstadoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, "datos inválidos", err.Error())
		return
	}

	incapacidad, err := h.useCase.CambiarEstado(c.Request.Context(), actor, id, req.IDEstado, req.Observaciones)
	if err != nil {
		handleError(c, err)
		return
	}
	response.Success(c, mapper.ToIncapacidadResponse(incapacidad), "estado de incapacidad actualizado")
}

func (h *IncapacidadHandler) Archivar(c *gin.Context) {
	actor, err := actorFromGin(c)
	if err != nil {
		handleError(c, err)
		return
	}
	id, err := parseIDParam(c, "id")
	if err != nil {
		handleError(c, err)
		return
	}

	if err := h.useCase.Archivar(c.Request.Context(), actor, id); err != nil {
		handleError(c, err)
		return
	}
	response.NoContent(c)
}

func (h *IncapacidadHandler) ListarEstados(c *gin.Context) {
	actor, err := actorFromGin(c)
	if err != nil {
		handleError(c, err)
		return
	}
	items, err := h.useCase.ListarEstados(c.Request.Context(), actor)
	if err != nil {
		handleError(c, err)
		return
	}
	response.Success(c, mapper.ToEstadoResponses(items), "estados de incapacidad")
}

func (h *IncapacidadHandler) ListarTipos(c *gin.Context) {
	actor, err := actorFromGin(c)
	if err != nil {
		handleError(c, err)
		return
	}
	items, err := h.useCase.ListarTipos(c.Request.Context(), actor)
	if err != nil {
		handleError(c, err)
		return
	}
	response.Success(c, mapper.ToTipoResponses(items), "tipos de incapacidad")
}

func (h *IncapacidadHandler) ListarEntidades(c *gin.Context) {
	actor, err := actorFromGin(c)
	if err != nil {
		handleError(c, err)
		return
	}
	items, err := h.useCase.ListarEntidades(c.Request.Context(), actor)
	if err != nil {
		handleError(c, err)
		return
	}
	response.Success(c, mapper.ToEntidadResponses(items), "entidades")
}

func parseIDParam(c *gin.Context, name string) (uint64, error) {
	id, err := strconv.ParseUint(c.Param(name), 10, 64)
	if err != nil || id == 0 {
		return 0, apperrors.ErrBadRequest.WithMessage("id inválido")
	}
	return id, nil
}

func handleError(c *gin.Context, err error) {
	var appErr *apperrors.AppError
	if errors.As(err, &appErr) {
		response.Error(c, appErr.HTTPStatus, appErr.Message, appErr.Code, appErr.Details)
		return
	}
	response.InternalError(c, "error interno", "INTERNAL_ERROR")
}
