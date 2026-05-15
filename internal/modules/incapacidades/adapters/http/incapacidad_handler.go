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

// Crear godoc
// @Summary Crear incapacidad
// @Description Registra una nueva incapacidad médica
// @Tags incapacidades
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.CrearIncapacidadRequest true "Datos de la incapacidad"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Router /incapacidades [post]
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

	var userID uint64
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

// Obtener godoc
// @Summary Obtener incapacidad por ID
// @Description Obtiene los detalles de una incapacidad específica
// @Tags incapacidades
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID de la incapacidad"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /incapacidades/{id} [get]
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

// Listar godoc
// @Summary Listar incapacidades
// @Description Lista las incapacidades con filtros y paginación
// @Tags incapacidades
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id_usuario query int false "Filtrar por usuario"
// @Param id_estado query int false "Filtrar por estado"
// @Param id_tipo query int false "Filtrar por tipo"
// @Param id_entidad query int false "Filtrar por entidad"
// @Param origen query string false "Filtrar por origen"
// @Param canal_recepcion query string false "Filtrar por canal"
// @Param search query string false "Buscar por título, observaciones, nombre o documento"
// @Param fecha_desde query string false "Fecha inicio desde (YYYY-MM-DD)"
// @Param fecha_hasta query string false "Fecha inicio hasta (YYYY-MM-DD)"
// @Param page query int false "Página" default(1)
// @Param limit query int false "Límite" default(20)
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /incapacidades [get]
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
		Search:         query.Search,
		FechaDesde:     query.FechaDesde,
		FechaHasta:     query.FechaHasta,
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

// Actualizar godoc
// @Summary Actualizar incapacidad
// @Description Actualiza los datos de una incapacidad
// @Tags incapacidades
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID de la incapacidad"
// @Param request body dto.ActualizarIncapacidadRequest true "Datos a actualizar"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /incapacidades/{id} [put]
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

// CambiarEstado godoc
// @Summary Cambiar estado de incapacidad
// @Description Actualiza el estado de una incapacidad
// @Tags incapacidades
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID de la incapacidad"
// @Param request body dto.CambiarEstadoRequest true "Nuevo estado"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /incapacidades/{id}/estado [patch]
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

// Archivar godoc
// @Summary Archivar incapacidad
// @Description Archiva una incapacidad (soft delete)
// @Tags incapacidades
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID de la incapacidad"
// @Success 204
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /incapacidades/{id} [delete]
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

// ListarEstados godoc
// @Summary Listar estados de incapacidad
// @Description Lista todos los estados disponibles
// @Tags catalogos
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Router /incapacidades/estados [get]
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

// ListarTipos godoc
// @Summary Listar tipos de incapacidad
// @Description Lista todos los tipos de incapacidad
// @Tags catalogos
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Router /incapacidades/tipos [get]
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

// ListarEntidades godoc
// @Summary Listar entidades
// @Description Lista todas las EPS y ARL registradas
// @Tags catalogos
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Router /incapacidades/entidades [get]
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

// ListarEstadosDocumento godoc
// @Summary Listar estados de documento
// @Description Lista todos los estados de documento
// @Tags catalogos
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Router /catalogos/estados-documento [get]
func (h *IncapacidadHandler) ListarEstadosDocumento(c *gin.Context) {
	actor, err := actorFromGin(c)
	if err != nil {
		handleError(c, err)
		return
	}
	items, err := h.useCase.ListarEstadosDocumento(c.Request.Context(), actor)
	if err != nil {
		handleError(c, err)
		return
	}
	response.Success(c, mapper.ToEstadoDocumentoResponses(items), "estados de documento")
}

// ListarTiposDocumento godoc
// @Summary Listar tipos de documento
// @Description Lista todos los tipos de documento
// @Tags catalogos
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Router /catalogos/tipos-documento [get]
func (h *IncapacidadHandler) ListarTiposDocumento(c *gin.Context) {
	actor, err := actorFromGin(c)
	if err != nil {
		handleError(c, err)
		return
	}
	items, err := h.useCase.ListarTiposDocumento(c.Request.Context(), actor)
	if err != nil {
		handleError(c, err)
		return
	}
	response.Success(c, mapper.ToTipoDocumentoResponses(items), "tipos de documento")
}

// ListarTiposPago godoc
// @Summary Listar tipos de pago
// @Description Lista todos los tipos de pago
// @Tags catalogos
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Router /catalogos/tipos-pago [get]
func (h *IncapacidadHandler) ListarTiposPago(c *gin.Context) {
	actor, err := actorFromGin(c)
	if err != nil {
		handleError(c, err)
		return
	}
	items, err := h.useCase.ListarTiposPago(c.Request.Context(), actor)
	if err != nil {
		handleError(c, err)
		return
	}
	response.Success(c, mapper.ToTipoPagoResponses(items), "tipos de pago")
}

// ObtenerDocumentosRequeridos godoc
// @Summary Obtener documentos requeridos
// @Description Lista los documentos requeridos para un tipo de incapacidad
// @Tags incapacidades
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param tipo_id path int true "ID del tipo de incapacidad"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /incapacidades/tipos/{tipo_id}/documentos-requeridos [get]
func (h *IncapacidadHandler) ObtenerDocumentosRequeridos(c *gin.Context) {
	actor, err := actorFromGin(c)
	if err != nil {
		handleError(c, err)
		return
	}
	tipoID, err := parseIDParam(c, "tipo_id")
	if err != nil {
		handleError(c, err)
		return
	}

	items, err := h.useCase.ObtenerDocumentosRequeridos(c.Request.Context(), actor, tipoID)
	if err != nil {
		handleError(c, err)
		return
	}
	response.Success(c, mapper.ToTipoDocumentoResponses(items), "documentos requeridos")
}

// ObtenerPlazos godoc
// @Summary Obtener información de plazos
// @Description Obtiene los plazos y alertas de vencimiento para una incapacidad
// @Tags incapacidades
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID de la incapacidad"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /incapacidades/{id}/plazos [get]
func (h *IncapacidadHandler) ObtenerPlazos(c *gin.Context) {
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

	info, err := h.useCase.ObtenerInfoPlazos(c.Request.Context(), actor, incapacidadID)
	if err != nil {
		handleError(c, err)
		return
	}

	var docsRequeridos []dto.TipoDocumentoResponse
	for _, doc := range info.DocumentosRequeridos {
		docsRequeridos = append(docsRequeridos, dto.TipoDocumentoResponse{
			IDTipoDocumento: doc.IDTipoDocumento,
			Nombre:          doc.Nombre,
			Descripcion:     doc.Descripcion,
			Requerido:       doc.Requerido,
		})
	}

	var alertas []string
	if len(info.AlertasVencimientos) > 0 {
		alertas = info.AlertasVencimientos
	}

	response.Success(c, gin.H{
		"id_incapacidad":               info.IncapacidadID,
		"tipo_incapacidad":             info.TipoIncapacidad,
		"documentos_requeridos":        docsRequeridos,
		"plazo_entrega_dias":          info.PlazoEntregaDias,
		"fecha_limite_entrega":        info.FechaLimiteEntrega.Format("2006-01-02"),
		"plazo_transcripcion_dias":    info.PlazoTranscripcionDias,
		"fecha_limite_transcripcion":  info.FechaLimiteTranscripcion.Format("2006-01-02"),
		"tiempo_maximo_pago_dias":     info.TiempoMaximoPagoDias,
		"fecha_limite_pago":           info.FechaLimitePago.Format("2006-01-02"),
		"dias_transcurridos":          info.DiasTranscurridos,
		"alertas_vencimiento":          alertas,
	}, "información de plazos obtenida")
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
