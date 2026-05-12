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

// Crear godoc
// @Summary Crear notificación
// @Description Crea una nueva notificación
// @Tags notificaciones
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.CrearNotificacionRequest true "Datos de la notificación"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /notificaciones [post]
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

// Listar godoc
// @Summary Listar notificaciones
// @Description Lista las notificaciones del usuario
// @Tags notificaciones
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id_usuario query int false "Filtrar por usuario"
// @Param id_incapacidad query int false "Filtrar por incapacidad"
// @Param tipo_notificacion query string false "Filtrar por tipo"
// @Param leida query bool false "Filtrar por estado de lectura"
// @Param page query int false "Página" default(1)
// @Param limit query int false "Límite" default(20)
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /notificaciones [get]
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

// Obtener godoc
// @Summary Obtener notificación por ID
// @Description Obtiene los detalles de una notificación
// @Tags notificaciones
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID de la notificación"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /notificaciones/{id} [get]
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

// ContarNoLeidas godoc
// @Summary Contar notificaciones no leídas
// @Description Obtiene el conteo de notificaciones no leídas
// @Tags notificaciones
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /notificaciones/no-leidas/count [get]
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

// MarcarLeida godoc
// @Summary Marcar notificación como leída
// @Description Marca una notificación específica como leída
// @Tags notificaciones
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID de la notificación"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /notificaciones/{id}/leida [patch]
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

// MarcarTodasLeidas godoc
// @Summary Marcar todas como leídas
// @Description Marca todas las notificaciones como leídas
// @Tags notificaciones
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /notificaciones/marcar-todas-leidas [patch]
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

// Eliminar godoc
// @Summary Eliminar notificación
// @Description Elimina una notificación
// @Tags notificaciones
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID de la notificación"
// @Success 204
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /notificaciones/{id} [delete]
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
