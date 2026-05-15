package http

import (
	"errors"
	"net/http"
	"strconv"

	"disability_system_backend/internal/modules/auditoria/dto"
	"disability_system_backend/internal/modules/auditoria/mapper"
	"disability_system_backend/internal/modules/auditoria/ports"
	"disability_system_backend/internal/modules/auditoria/usecase"
	apperrors "disability_system_backend/internal/shared/errors"
	"disability_system_backend/internal/shared/response"

	"github.com/gin-gonic/gin"
)

type actorWrapper struct {
	userID   uint64
	permisos []string
}

func (a *actorWrapper) GetUserID() uint64 {
	return a.userID
}

func (a *actorWrapper) HasPermission(p string) bool {
	for _, perm := range a.permisos {
		if perm == p {
			return true
		}
	}
	return false
}

func actorFromGin(c *gin.Context) (ports.Actor, error) {
	userID := c.GetUint64("user_id")
	if userID == 0 {
		return nil, apperrors.ErrUnauthorized
	}

	rawPermisos, exists := c.Get("permisos")
	var permisos []string
	if exists {
		if p, ok := rawPermisos.([]string); ok {
			permisos = p
		}
	}

	return &actorWrapper{
		userID:   userID,
		permisos: permisos,
	}, nil
}

type AuditoriaHandler struct {
	useCase *usecase.AuditoriaUseCase
}

func NewAuditoriaHandler(uc *usecase.AuditoriaUseCase) *AuditoriaHandler {
	return &AuditoriaHandler{useCase: uc}
}

// Listar godoc
// @Summary Listar entradas de auditoría
// @Description Lista el historial de acciones y auditoría en el sistema
// @Tags auditoria
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id_usuario query int false "Filtrar por ID de usuario"
// @Param id_incapacidad query int false "Filtrar por ID de incapacidad"
// @Param tipo_accion query string false "Filtrar por tipo de acción"
// @Param modulo query string false "Filtrar por módulo"
// @Param fecha_inicio query string false "Fecha de inicio (YYYY-MM-DD)"
// @Param fecha_fin query string false "Fecha de fin (YYYY-MM-DD)"
// @Param page query int false "Página" default(1)
// @Param limit query int false "Límite" default(20)
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Router /auditoria [get]
func (h *AuditoriaHandler) Listar(c *gin.Context) {
	actor, err := actorFromGin(c)
	if err != nil {
		handleError(c, err)
		return
	}

	var query dto.ListarAuditoriaQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.ValidationError(c, "Parámetros de consulta inválidos", err.Error())
		return
	}

	if query.Page < 1 {
		query.Page = 1
	}
	if query.Limit < 1 {
		query.Limit = 20
	}
	if query.Limit > 100 {
		query.Limit = 100
	}

	auditorias, total, err := h.useCase.Listar(c.Request.Context(), actor, query)
	if err != nil {
		handleError(c, err)
		return
	}

	items := make([]dto.AuditoriaResponse, len(auditorias))
	for i, a := range auditorias {
		items[i] = mapper.ToResponse(a)
	}

	response.Paginated(c, items, total, int64(query.Page), int64(query.Limit))
}

// ListarPorUsuario godoc
// @Summary Listar auditoría de un usuario
// @Description Lista el historial de acciones realizadas por un usuario específico
// @Tags auditoria
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID del usuario"
// @Param page query int false "Página" default(1)
// @Param limit query int false "Límite" default(20)
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Router /auditoria/usuario/{id} [get]
func (h *AuditoriaHandler) ListarPorUsuario(c *gin.Context) {
	actor, err := actorFromGin(c)
	if err != nil {
		handleError(c, err)
		return
	}

	idParam := c.Param("id")
	idUsuario, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "ID inválido", "INVALID_ID", err.Error())
		return
	}

	var query dto.ListarAuditoriaQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.ValidationError(c, "Parámetros de consulta inválidos", err.Error())
		return
	}
	query.IDUsuario = &idUsuario

	if query.Page < 1 {
		query.Page = 1
	}
	if query.Limit < 1 {
		query.Limit = 20
	}
	if query.Limit > 100 {
		query.Limit = 100
	}

	auditorias, total, err := h.useCase.Listar(c.Request.Context(), actor, query)
	if err != nil {
		handleError(c, err)
		return
	}

	items := make([]dto.AuditoriaResponse, len(auditorias))
	for i, a := range auditorias {
		items[i] = mapper.ToResponse(a)
	}

	response.Paginated(c, items, total, int64(query.Page), int64(query.Limit))
}

func handleError(c *gin.Context, err error) {
	var appErr *apperrors.AppError
	if errors.As(err, &appErr) {
		response.Error(c, appErr.HTTPStatus, appErr.Message, appErr.Code, appErr.Details)
		return
	}
	response.InternalError(c, "Error interno", "INTERNAL_ERROR")
}
