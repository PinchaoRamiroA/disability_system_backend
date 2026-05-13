package http

import (
	"errors"
	"strconv"

	"disability_system_backend/internal/modules/cobros/dto"
	"disability_system_backend/internal/modules/cobros/mapper"
	"disability_system_backend/internal/modules/cobros/ports"
	"disability_system_backend/internal/modules/cobros/usecase"
	apperrors "disability_system_backend/internal/shared/errors"
	"disability_system_backend/internal/shared/response"

	"github.com/gin-gonic/gin"
)

type CobroHandler struct {
	useCase *usecase.CobroUseCase
}

func NewCobroHandler(useCase *usecase.CobroUseCase) *CobroHandler {
	return &CobroHandler{useCase: useCase}
}

// CrearPago godoc
// @Summary Registrar pago
// @Description Registra un nuevo pago de EPS/ARL
// @Tags cobros
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.CrearPagoRequest true "Datos del pago"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /cobros/pagos [post]
func (h *CobroHandler) CrearPago(c *gin.Context) {
	actor, err := actorFromGin(c)
	if err != nil {
		handleError(c, err)
		return
	}
	var req dto.CrearPagoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, "datos inválidos", err.Error())
		return
	}
	pago, err := h.useCase.CrearPago(c.Request.Context(), actor, usecase.CrearPagoInput{
		IDIncapacidad:   req.IDIncapacidad,
		IDEntidad:       req.IDEntidad,
		TipoPago:        req.TipoPago,
		EstadoPago:      req.EstadoPago,
		Descripcion:     req.Descripcion,
		Valor:           req.Valor,
		FechaPago:       req.FechaPago,
		PeriodoContable: req.PeriodoContable,
	})
	if err != nil {
		handleError(c, err)
		return
	}
	response.Created(c, mapper.ToPagoResponse(pago), "pago registrado")
}

// ObtenerPago godoc
// @Summary Obtener pago por ID
// @Description Obtiene los detalles de un pago específico
// @Tags cobros
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID del pago"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /cobros/pagos/{id} [get]
func (h *CobroHandler) ObtenerPago(c *gin.Context) {
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
	pago, err := h.useCase.ObtenerPago(c.Request.Context(), actor, id)
	if err != nil {
		handleError(c, err)
		return
	}
	response.Success(c, mapper.ToPagoResponse(pago), "pago encontrado")
}

// ListarPagos godoc
// @Summary Listar pagos
// @Description Lista los pagos con filtros
// @Tags cobros
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id_incapacidad query int false "Filtrar por incapacidad"
// @Param id_entidad query int false "Filtrar por entidad"
// @Param tipo_pago query string false "Filtrar por tipo"
// @Param estado_pago query string false "Filtrar por estado"
// @Param conciliado query bool false "Filtrar por conciliación"
// @Param page query int false "Página" default(1)
// @Param limit query int false "Límite" default(20)
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /cobros/pagos [get]
func (h *CobroHandler) ListarPagos(c *gin.Context) {
	actor, err := actorFromGin(c)
	if err != nil {
		handleError(c, err)
		return
	}
	var query dto.ListarPagosQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.ValidationError(c, "filtros inválidos", err.Error())
		return
	}
	items, total, err := h.useCase.ListarPagos(c.Request.Context(), actor, ports.PagoFilters{
		IDIncapacidad: query.IDIncapacidad,
		IDEntidad:     query.IDEntidad,
		TipoPago:      query.TipoPago,
		EstadoPago:    query.EstadoPago,
		Conciliado:    query.Conciliado,
		Page:          query.Page,
		Limit:         query.Limit,
	})
	if err != nil {
		handleError(c, err)
		return
	}
	page, limit := normalizePagination(query.Page, query.Limit)
	response.Paginated(c, mapper.ToPagoResponses(items), total, int64(page), int64(limit))
}

// ActualizarPago godoc
// @Summary Actualizar pago
// @Description Actualiza los datos de un pago
// @Tags cobros
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID del pago"
// @Param request body dto.ActualizarPagoRequest true "Datos del pago"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /cobros/pagos/{id} [put]
func (h *CobroHandler) ActualizarPago(c *gin.Context) {
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
	var req dto.ActualizarPagoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, "datos inválidos", err.Error())
		return
	}
	pago, err := h.useCase.ActualizarPago(c.Request.Context(), actor, id, usecase.ActualizarPagoInput{
		IDEntidad:       req.IDEntidad,
		TipoPago:        req.TipoPago,
		EstadoPago:      req.EstadoPago,
		Descripcion:     req.Descripcion,
		Valor:           req.Valor,
		FechaPago:       req.FechaPago,
		PeriodoContable: req.PeriodoContable,
	})
	if err != nil {
		handleError(c, err)
		return
	}
	response.Success(c, mapper.ToPagoResponse(pago), "pago actualizado")
}

// EliminarPago godoc
// @Summary Eliminar pago
// @Description Elimina un pago del sistema
// @Tags cobros
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID del pago"
// @Success 204
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /cobros/pagos/{id} [delete]
func (h *CobroHandler) EliminarPago(c *gin.Context) {
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
	if err := h.useCase.EliminarPago(c.Request.Context(), actor, id); err != nil {
		handleError(c, err)
		return
	}
	response.NoContent(c)
}

// ConciliarPago godoc
// @Summary Conciliar pago
// @Description Concilia un pago para contabilidad
// @Tags cobros
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID del pago"
// @Param request body dto.ConciliarPagoRequest true "Datos de conciliación"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /cobros/pagos/{id}/conciliar [patch]
func (h *CobroHandler) ConciliarPago(c *gin.Context) {
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
	var req dto.ConciliarPagoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, "datos inválidos", err.Error())
		return
	}
	pago, err := h.useCase.ConciliarPago(c.Request.Context(), actor, id, req.Conciliado, req.EstadoPago, req.Descripcion)
	if err != nil {
		handleError(c, err)
		return
	}
	response.Success(c, mapper.ToPagoResponse(pago), "pago conciliado")
}

// CrearSeguimiento godoc
// @Summary Registrar seguimiento de cobro
// @Description Registra un seguimiento de cobro persuasivo o jurídico
// @Tags cobros
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.CrearSeguimientoRequest true "Datos del seguimiento"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /cobros/seguimientos [post]
func (h *CobroHandler) CrearSeguimiento(c *gin.Context) {
	actor, err := actorFromGin(c)
	if err != nil {
		handleError(c, err)
		return
	}
	var req dto.CrearSeguimientoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, "datos inválidos", err.Error())
		return
	}
	seguimiento, err := h.useCase.CrearSeguimiento(c.Request.Context(), actor, usecase.CrearSeguimientoInput{
		IDIncapacidad:   req.IDIncapacidad,
		TipoSeguimiento: req.TipoSeguimiento,
		Descripcion:     req.Descripcion,
		Resultado:       req.Resultado,
	})
	if err != nil {
		handleError(c, err)
		return
	}
	response.Created(c, mapper.ToSeguimientoResponse(seguimiento), "seguimiento registrado")
}

// ObtenerSeguimiento godoc
// @Summary Obtener seguimiento por ID
// @Description Obtiene los detalles de un seguimiento
// @Tags cobros
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID del seguimiento"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /cobros/seguimientos/{id} [get]
func (h *CobroHandler) ObtenerSeguimiento(c *gin.Context) {
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
	seguimiento, err := h.useCase.ObtenerSeguimiento(c.Request.Context(), actor, id)
	if err != nil {
		handleError(c, err)
		return
	}
	response.Success(c, mapper.ToSeguimientoResponse(seguimiento), "seguimiento encontrado")
}

// ListarSeguimientos godoc
// @Summary Listar seguimientos
// @Description Lista los seguimientos de cobro
// @Tags cobros
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id_incapacidad query int false "Filtrar por incapacidad"
// @Param tipo_seguimiento query string false "Filtrar por tipo"
// @Param page query int false "Página" default(1)
// @Param limit query int false "Límite" default(20)
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /cobros/seguimientos [get]
func (h *CobroHandler) ListarSeguimientos(c *gin.Context) {
	actor, err := actorFromGin(c)
	if err != nil {
		handleError(c, err)
		return
	}
	var query dto.ListarSeguimientosQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.ValidationError(c, "filtros inválidos", err.Error())
		return
	}
	items, total, err := h.useCase.ListarSeguimientos(c.Request.Context(), actor, ports.SeguimientoFilters{
		IDIncapacidad:   query.IDIncapacidad,
		TipoSeguimiento: query.TipoSeguimiento,
		Page:            query.Page,
		Limit:           query.Limit,
	})
	if err != nil {
		handleError(c, err)
		return
	}
	page, limit := normalizePagination(query.Page, query.Limit)
	response.Paginated(c, mapper.ToSeguimientoResponses(items), total, int64(page), int64(limit))
}

// ActualizarSeguimiento godoc
// @Summary Actualizar seguimiento
// @Description Actualiza un seguimiento de cobro
// @Tags cobros
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID del seguimiento"
// @Param request body dto.ActualizarSeguimientoRequest true "Datos del seguimiento"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /cobros/seguimientos/{id} [put]
func (h *CobroHandler) ActualizarSeguimiento(c *gin.Context) {
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
	var req dto.ActualizarSeguimientoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, "datos inválidos", err.Error())
		return
	}
	seguimiento, err := h.useCase.ActualizarSeguimiento(c.Request.Context(), actor, id, usecase.ActualizarSeguimientoInput{
		TipoSeguimiento: req.TipoSeguimiento,
		Descripcion:     req.Descripcion,
		Resultado:       req.Resultado,
	})
	if err != nil {
		handleError(c, err)
		return
	}
	response.Success(c, mapper.ToSeguimientoResponse(seguimiento), "seguimiento actualizado")
}

// ObtenerEstadisticas godoc
// @Summary Obtener estadísticas de cartera
// @Description Obtiene estadísticas generales de la cartera de cobros
// @Tags cartera
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /cartera/estadisticas [get]
func (h *CobroHandler) ObtenerEstadisticas(c *gin.Context) {
	actor, err := actorFromGin(c)
	if err != nil {
		handleError(c, err)
		return
	}
	estadisticas, err := h.useCase.ObtenerEstadisticasGenerales(c.Request.Context(), actor)
	if err != nil {
		handleError(c, err)
		return
	}
	response.Success(c, estadisticas, "estadísticas de cartera")
}

// ObtenerResumenEntidad godoc
// @Summary Obtener resumen por entidad
// @Description Obtiene el resumen de cobros agrupado por entidad
// @Tags cartera
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /cartera/resumen-entidad [get]
func (h *CobroHandler) ObtenerResumenEntidad(c *gin.Context) {
	actor, err := actorFromGin(c)
	if err != nil {
		handleError(c, err)
		return
	}
	resumen, err := h.useCase.ObtenerResumenPorEntidad(c.Request.Context(), actor)
	if err != nil {
		handleError(c, err)
		return
	}
	response.Success(c, resumen, "resumen por entidad")
}

// ObtenerAlertasVencimiento godoc
// @Summary Obtener alertas de vencimiento
// @Description Lista las alertas de vencimiento de pagos
// @Tags cartera
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param dias_minimos query int false "Días mínimos de vencimiento" default(0)
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /cartera/alertas-vencimiento [get]
func (h *CobroHandler) ObtenerAlertasVencimiento(c *gin.Context) {
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
	alertas, err := h.useCase.ObtenerAlertasVencimiento(c.Request.Context(), actor, query.DiasMinimos)
	if err != nil {
		handleError(c, err)
		return
	}
	response.Success(c, alertas, "alertas de vencimiento")
}

// ObtenerCarteraVencida godoc
// @Summary Obtener cartera vencida
// @Description Lista los pagos vencidos
// @Tags cartera
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /cartera/vencida [get]
func (h *CobroHandler) ObtenerCarteraVencida(c *gin.Context) {
	actor, err := actorFromGin(c)
	if err != nil {
		handleError(c, err)
		return
	}
	pagos, err := h.useCase.ObtenerCarteraVencida(c.Request.Context(), actor)
	if err != nil {
		handleError(c, err)
		return
	}
	response.Success(c, mapper.ToPagoResponses(pagos), "cartera vencida")
}

// ObtenerProximoEstado godoc
// @Summary Obtener el próximo estado de una incapacidad
// @Description Retorna el próximo estado sugerido basado en la acción a realizar
// @Tags Cartera
// @Accept json
// @Produce json
// @Param id path int true "ID de la incapacidad"
// @Param accion query string true "Acción a realizar"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /cartera/incapacidades/{id}/proximo-estado [get]
func (h *CobroHandler) ObtenerProximoEstado(c *gin.Context) {
	actor, err := actorFromGin(c)
	if err != nil {
		handleError(c, err)
		return
	}
	incapacidadID, err := parseIDParam(c, "id")
	if err != nil {
		handleError(c, err)
		return
	}
	accion := c.Query("accion")
	if accion == "" {
		response.BadRequest(c, "parámetro 'accion' requerido", "MISSING_PARAM", nil)
		return
	}
	estado, err := h.useCase.ObtenerProximoEstadoIncapacidad(c.Request.Context(), actor, incapacidadID, accion)
	if err != nil {
		handleError(c, err)
		return
	}
	response.Success(c, gin.H{"proximo_estado": estado}, "próximo estado sugerido")
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
