package http

import (
	"errors"
	"log/slog"
	"strconv"
	"time"

	"disability_system_backend/internal/modules/incapacidades/dto"
	"disability_system_backend/internal/modules/incapacidades/usecase"
	"disability_system_backend/internal/shared/response"

	"github.com/gin-gonic/gin"
)

type TranscripcionHandler struct {
	useCase *usecase.TranscripcionUseCase
}

func NewTranscripcionHandler(useCase *usecase.TranscripcionUseCase) *TranscripcionHandler {
	return &TranscripcionHandler{useCase: useCase}
}

// Transcribir godoc
// @Summary Registrar transcripción EPS/ARL
// @Description Registra la transcripción de una incapacidad ante la EPS/ARL
// @Tags incapacidades
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID de la incapacidad"
// @Param request body dto.TranscribirIncapacidadRequest true "Datos de transcripción"
// @Success 200 {object} response.Response "transcripción registrada exitosamente"
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 409 {object} response.Response
// @Router /incapacidades/{id}/transcribir [post]
func (h *TranscripcionHandler) Transcribir(c *gin.Context) {
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

	var req dto.TranscribirIncapacidadRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, "datos inválidos", err.Error())
		return
	}

	obs := req.ObservacionesTranscripcion
	if obs == nil && req.Observaciones != nil {
		obs = req.Observaciones
	}
	if req.NumeroRadicado != nil && *req.NumeroRadicado != "" {
		radicadoNote := "Radicado: " + *req.NumeroRadicado
		if obs != nil && *obs != "" {
			combined := radicadoNote + " | " + *obs
			obs = &combined
		} else {
			obs = &radicadoNote
		}
	}

	incapacidad, err := h.useCase.Transcribir(c.Request.Context(), id, actor, obs)
	if err != nil {
		handleError(c, err)
		return
	}

	response.Success(c, gin.H{
		"id_incapacidad":       incapacidad.IDIncapacidad,
		"estado_transcripcion": incapacidad.EstadoTranscripcion,
		"fecha_transcripcion":  incapacidad.FechaTranscripcion,
		"transcrito_por":       incapacidad.TranscritoPor,
		"observaciones":        incapacidad.ObservacionesTranscripcion,
	}, "transcripción registrada exitosamente")
}

// MarcarEnProceso godoc
// @Summary Marcar transcripción en proceso
// @Description Marca el estado de transcripción de una incapacidad como 'en_proceso'
// @Tags incapacidades
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID de la incapacidad"
// @Param request body dto.MarcarTranscripcionRequest true "Estado de transcripción"
// @Success 200 {object} response.Response "estado actualizado"
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /incapacidades/{id}/transcripcion [patch]
func (h *TranscripcionHandler) MarcarEnProceso(c *gin.Context) {
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

	var req dto.MarcarTranscripcionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, "datos inválidos", err.Error())
		return
	}

	incapacidad, err := h.useCase.MarcarEnProceso(c.Request.Context(), id, actor)
	if err != nil {
		slog.ErrorContext(c.Request.Context(), "MarcarEnProceso error", "error", err, "id", id)
		if errors.Is(err, usecase.ErrTranscripcionNotAllowed) {
			response.BadRequest(c, "no se puede cambiar el estado de transcripción cuando ya está completada", "TRANSCRIPCION_NOT_ALLOWED", nil)
			return
		}
		handleError(c, err)
		return
	}

	response.Success(c, gin.H{
		"id_incapacidad":       incapacidad.IDIncapacidad,
		"estado_transcripcion": incapacidad.EstadoTranscripcion,
	}, "estado de transcripción actualizado")
}

// ListarPendientes godoc
// @Summary Listar transcripciones pendientes
// @Description Lista las incapacidades con transcripción pendiente
// @Tags incapacidades
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param estado query string false "Estado de transcripción (pendiente, en_proceso, completado, vencida)"
// @Param page query int false "Número de página" default(1)
// @Param limit query int false "Límite de resultados" default(20)
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /incapacidades/transcripciones/pendientes [get]
func (h *TranscripcionHandler) ListarPendientes(c *gin.Context) {
	var query dto.ListarTranscripcionesQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.ValidationError(c, "filtros inválidos", err.Error())
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

	items, total, err := h.useCase.ListarPendientes(c.Request.Context(), query.Estado, page, limit)
	if err != nil {
		handleError(c, err)
		return
	}

	var results []gin.H
	for _, item := range items {
		alerta := h.useCase.ObtenerAlertaVencimiento(&item)
		fechaLimite := ""
		if item.FechaLimiteTranscripcion != nil {
			fechaLimite = item.FechaLimiteTranscripcion.Format("2006-01-02")
		}

		var entidadData gin.H
		if item.Entidad != nil {
			entidadData = gin.H{
				"id_entidad":               item.Entidad.IDEntidad,
				"nombre":                   item.Entidad.Nombre,
				"tipo":                     item.Entidad.Tipo,
				"plazo_transcripcion_dias": item.Entidad.PlazoTranscripcionDias,
				"tiempo_maximo_pago_dias":  item.Entidad.TiempoMaximoPagoDias,
				"requiere_transcripcion":   item.Entidad.RequiereTranscripcion,
			}
		}
		var tipoData gin.H
		if item.Tipo != nil {
			tipoData = gin.H{
				"id_tipo": item.Tipo.IDTipo,
				"nombre":  item.Tipo.Nombre,
				"origen":  item.Tipo.Origen,
			}
		}
		var estadoData gin.H
		if item.Estado != nil {
			estadoData = gin.H{
				"id_estado": item.Estado.IDEstado,
				"nombre":    item.Estado.Nombre,
			}
		}

		var diasRestantes *int
		if item.FechaLimiteTranscripcion != nil {
			dr := int(time.Until(*item.FechaLimiteTranscripcion).Hours() / 24)
			diasRestantes = &dr
		}
		fechaFinStr := ""
		if item.FechaFin != nil {
			fechaFinStr = item.FechaFin.Format("2006-01-02")
		}
		fechaRadicacionStr := ""
		if item.FechaRadicacion != nil {
			fechaRadicacionStr = item.FechaRadicacion.Format("2006-01-02")
		}

		results = append(results, gin.H{
			"id_incapacidad":              item.IDIncapacidad,
			"id_usuario":                  item.IDUsuario,
			"titulo":                      item.Titulo,
			"origen":                      item.Origen,
			"fecha_inicio":                item.FechaInicio.Format("2006-01-02"),
			"fecha_fin":                   fechaFinStr,
			"fecha_radicacion":            fechaRadicacionStr,
			"estado_transcripcion":        item.EstadoTranscripcion,
			"fecha_limite_transcripcion":  fechaLimite,
			"dias_restantes":              diasRestantes,
			"alerta_vencimiento":          alerta,
			"observaciones_transcripcion": item.ObservacionesTranscripcion,
			"entidad":                     entidadData,
			"tipo":                        tipoData,
			"estado":                      estadoData,
		})
	}

	response.Paginated(c, results, total, int64(page), int64(limit))
}

func parsePageLimit(c *gin.Context) (int, int) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	return page, limit
}
