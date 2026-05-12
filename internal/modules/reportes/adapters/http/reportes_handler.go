package http

import (
	"errors"
	"strconv"
	"time"

	reportesdomain "disability_system_backend/internal/modules/reportes/domain"
	"disability_system_backend/internal/modules/reportes/dto"
	"disability_system_backend/internal/modules/reportes/ports"
	"disability_system_backend/internal/modules/reportes/usecase"
	apperrors "disability_system_backend/internal/shared/errors"
	"disability_system_backend/internal/shared/response"

	"github.com/gin-gonic/gin"
)

type ReportesHandler struct {
	useCase *usecase.ReportesUseCase
}

func NewReportesHandler(useCase *usecase.ReportesUseCase) *ReportesHandler {
	return &ReportesHandler{useCase: useCase}
}

func (h *ReportesHandler) GenerarReporte(c *gin.Context) {
	actor, err := actorFromGin(c)
	if err != nil {
		handleError(c, err)
		return
	}

	var req dto.GenerarReporteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, "datos inválidos", err.Error())
		return
	}

	filtros := buildFiltrosReporte(req)

	var result interface{}
	switch req.Tipo {
	case "incapacidades":
		result, err = h.useCase.GenerarReporteIncapacidades(c.Request.Context(), actor, filtros)
	case "ausentismo":
		result, err = h.useCase.GenerarReporteAusentismo(c.Request.Context(), actor, filtros)
	case "cartera":
		result, err = h.useCase.GenerarReporteCartera(c.Request.Context(), actor, filtros)
	default:
		response.BadRequest(c, "tipo de reporte no válido", "INVALID_REPORT_TYPE", nil)
		return
	}

	if err != nil {
		handleError(c, err)
		return
	}

	response.Success(c, result, "reporte generado")
}

func (h *ReportesHandler) GenerarReporteEntidad(c *gin.Context) {
	actor, err := actorFromGin(c)
	if err != nil {
		handleError(c, err)
		return
	}

	entidadID, err := strconv.ParseUint(c.Param("entidad_id"), 10, 64)
	if err != nil || entidadID == 0 {
		response.BadRequest(c, "id de entidad inválido", "INVALID_ID", nil)
		return
	}

	var req dto.GenerarReporteRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.ValidationError(c, "filtros inválidos", err.Error())
		return
	}

	filtros := buildFiltrosReporte(req)

	result, err := h.useCase.GenerarReporteEntidad(c.Request.Context(), actor, entidadID, filtros)
	if err != nil {
		handleError(c, err)
		return
	}

	response.Success(c, result, "reporte de entidad generado")
}

func (h *ReportesHandler) GenerarReporteVencimientos(c *gin.Context) {
	actor, err := actorFromGin(c)
	if err != nil {
		handleError(c, err)
		return
	}

	var query struct {
		DiasMinimos int `form:"dias_minimos,default=0"`
	}
	if err := c.ShouldBindQuery(&query); err != nil {
		response.ValidationError(c, "parámetros inválidos", err.Error())
		return
	}

	result, err := h.useCase.GenerarReporteVencimientos(c.Request.Context(), actor, query.DiasMinimos)
	if err != nil {
		handleError(c, err)
		return
	}

	response.Success(c, result, "reporte de vencimientos generado")
}

func (h *ReportesHandler) ObtenerResumenEjecutivo(c *gin.Context) {
	actor, err := actorFromGin(c)
	if err != nil {
		handleError(c, err)
		return
	}

	resumen, err := h.useCase.ObtenerResumenEjecutivo(c.Request.Context(), actor)
	if err != nil {
		handleError(c, err)
		return
	}

	response.Success(c, gin.H{
		"fecha_generacion":       resumen.FechaGeneracion.Format(time.RFC3339),
		"total_incapacidades":    resumen.TotalIncapacidades,
		"incapacidades_activas":  resumen.IncapacidadesActivas,
		"total_dias_perdidos":    resumen.TotalDiasPerdidos,
		"total_valor_cartera":   resumen.TotalValorCartera,
		"total_valor_cobrado":    resumen.TotalValorCobrado,
		"total_valor_pendiente": resumen.TotalValorPendiente,
		"pagos_pendientes":      resumen.PagosPendientes,
		"pagos_vencidos":        resumen.PagosVencidos,
	}, "resumen ejecutivo")
}

func buildFiltrosReporte(req dto.GenerarReporteRequest) reportesdomain.FiltrosReporte {
	filtros := reportesdomain.FiltrosReporte{
		IDEntidad: req.IDEntidad,
		IDTipo:    req.IDTipo,
		IDEstado:  req.IDEstado,
		IDEmpleado: req.IDEmpleado,
		Origen:    req.Origen,
		Periodo:   reportesdomain.PeriodoReporte(req.Periodo),
	}

	if req.FechaInicio != "" {
		if t, err := time.Parse("2006-01-02", req.FechaInicio); err == nil {
			filtros.FechaInicio = &t
		}
	}
	if req.FechaFin != "" {
		if t, err := time.Parse("2006-01-02", req.FechaFin); err == nil {
			filtros.FechaFin = &t
		}
	}

	return filtros
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

func handleError(c *gin.Context, err error) {
	var appErr *apperrors.AppError
	if errors.As(err, &appErr) {
		response.Error(c, appErr.HTTPStatus, appErr.Message, appErr.Code, appErr.Details)
		return
	}
	response.InternalError(c, "error interno", "INTERNAL_ERROR")
}
