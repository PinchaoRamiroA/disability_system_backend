package http

import (
	"errors"
	"strconv"

	"disability_system_backend/internal/modules/notificaciones/dto"
	"disability_system_backend/internal/modules/notificaciones/mapper"
	"disability_system_backend/internal/modules/notificaciones/ports"
	"disability_system_backend/internal/modules/notificaciones/usecase"
	apperrors "disability_system_backend/internal/shared/errors"
	"disability_system_backend/internal/shared/response"

	"github.com/gin-gonic/gin"
)

type NotificacionHandler struct {
	useCase *usecase.NotificacionUseCase
}

func NewNotificacionHandler(useCase *usecase.NotificacionUseCase) *NotificacionHandler {
	return &NotificacionHandler{useCase: useCase}
}

func (h *NotificacionHandler) Crear(c *gin.Context) {
	actor, err := actorFromGin(c)
	if err != nil {
		handleError(c, err)
		return
	}
	var req dto.CrearNotificacionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, "datos inválidos", err.Error())
		return
	}
	notificacion, err := h.useCase.Crear(c.Request.Context(), actor, usecase.CrearNotificacionInput{
		IDUsuario:        req.IDUsuario,
		IDIncapacidad:    req.IDIncapacidad,
		TipoNotificacion: req.TipoNotificacion,
		Mensaje:          req.Mensaje,
	})
	if err != nil {
		handleError(c, err)
		return
	}
	response.Created(c, mapper.ToNotificacionResponse(notificacion), "notificación creada")
}

func (h *NotificacionHandler) Listar(c *gin.Context) {
	actor, err := actorFromGin(c)
	if err != nil {
		handleError(c, err)
		return
	}
	var query dto.ListarNotificacionesQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.ValidationError(c, "filtros inválidos", err.Error())
		return
	}
	items, total, err := h.useCase.Listar(c.Request.Context(), actor, ports.NotificacionFilters{
		IDUsuario:        query.IDUsuario,
		IDIncapacidad:    query.IDIncapacidad,
		TipoNotificacion: query.TipoNotificacion,
		Leida:            query.Leida,
		Page:             query.Page,
		Limit:            query.Limit,
	})
	if err != nil {
		handleError(c, err)
		return
	}
	page, limit := normalizePagination(query.Page, query.Limit)
	response.Paginated(c, mapper.ToNotificacionResponses(items), total, int64(page), int64(limit))
}

func (h *NotificacionHandler) Obtener(c *gin.Context) {
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
	notificacion, err := h.useCase.Obtener(c.Request.Context(), actor, id)
	if err != nil {
		handleError(c, err)
		return
	}
	response.Success(c, mapper.ToNotificacionResponse(notificacion), "notificación encontrada")
}

func (h *NotificacionHandler) ContarNoLeidas(c *gin.Context) {
	actor, err := actorFromGin(c)
	if err != nil {
		handleError(c, err)
		return
	}
	total, err := h.useCase.ContarNoLeidas(c.Request.Context(), actor)
	if err != nil {
		handleError(c, err)
		return
	}
	response.Success(c, dto.ConteoNoLeidasResponse{Total: total}, "notificaciones no leídas")
}

func (h *NotificacionHandler) MarcarLeida(c *gin.Context) {
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
	notificacion, err := h.useCase.MarcarLeida(c.Request.Context(), actor, id)
	if err != nil {
		handleError(c, err)
		return
	}
	response.Success(c, mapper.ToNotificacionResponse(notificacion), "notificación marcada como leída")
}

func (h *NotificacionHandler) MarcarTodasLeidas(c *gin.Context) {
	actor, err := actorFromGin(c)
	if err != nil {
		handleError(c, err)
		return
	}
	if err := h.useCase.MarcarTodasLeidas(c.Request.Context(), actor); err != nil {
		handleError(c, err)
		return
	}
	response.Success(c, gin.H{"leida": true}, "notificaciones marcadas como leídas")
}

func (h *NotificacionHandler) Eliminar(c *gin.Context) {
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
	if err := h.useCase.Eliminar(c.Request.Context(), actor, id); err != nil {
		handleError(c, err)
		return
	}
	response.NoContent(c)
}

func actorFromGin(c *gin.Context) (ports.Actor, error) {
	actorValue, exists := c.Get("actor")
	if !exists {
		return ports.Actor{}, apperrors.ErrUnauthorized.WithMessage("actor no encontrado en contexto")
	}
	actor, ok := actorValue.(ports.Actor)
	if !ok {
		return ports.Actor{}, apperrors.ErrInternal.WithMessage("actor inválido")
	}
	return actor, nil
}

func parseIDParam(c *gin.Context, name string) (uint64, error) {
	id, err := strconv.ParseUint(c.Param(name), 10, 64)
	if err != nil || id == 0 {
		return 0, apperrors.ErrBadRequest.WithMessage("id inválido")
	}
	return id, nil
}

func normalizePagination(page, limit int) (int, int) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}
	return page, limit
}

func handleError(c *gin.Context, err error) {
	var appErr *apperrors.AppError
	if errors.As(err, &appErr) {
		response.Error(c, appErr.HTTPStatus, appErr.Message, appErr.Code, appErr.Details)
		return
	}
	response.InternalError(c, "error interno", "INTERNAL_ERROR")
}
