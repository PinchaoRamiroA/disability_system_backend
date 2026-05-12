package http

import (
	"errors"
	"strconv"
	"time"

	"disability_system_backend/internal/modules/historial/domain"
	incapacidaddomain "disability_system_backend/internal/modules/incapacidades/domain"
	"disability_system_backend/internal/modules/incapacidades/dto"
	"disability_system_backend/internal/modules/incapacidades/usecase"
	apperrors "disability_system_backend/internal/shared/errors"
	"disability_system_backend/internal/shared/response"
	"disability_system_backend/internal/shared/storage"

	"github.com/gin-gonic/gin"
)

type DocumentoHandler struct {
	useCase          *usecase.DocumentoUseCase
	historialListFn  func(incapacidadID uint64, tipoID *uint64, page, limit int) ([]domain.Historial, int64, error)
	storageService   *storage.StorageService
}

func NewDocumentoHandler(useCase *usecase.DocumentoUseCase, historialListFn func(uint64, *uint64, int, int) ([]domain.Historial, int64, error), storageService *storage.StorageService) *DocumentoHandler {
	return &DocumentoHandler{
		useCase:         useCase,
		historialListFn: historialListFn,
		storageService:  storageService,
	}
}

// Subir godoc
// @Summary Subir documento
// @Description Registra un documento asociado a una incapacidad
// @Tags documentos
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.SubirDocumentoRequest true "Datos del documento"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /incapacidades/{id}/documentos [post]
func (h *DocumentoHandler) Subir(c *gin.Context) {
	actor, err := actorFromGin(c)
	if err != nil {
		handleError(c, err)
		return
	}

	var req dto.SubirDocumentoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, "datos inválidos", err.Error())
		return
	}

	documento, err := h.useCase.Subir(c.Request.Context(), actor, struct {
		IDIncapacidad uint64
		Nombre        string
		Tipo          string
		URL           string
		Formato       string
	}{
		IDIncapacidad: req.IDIncapacidad,
		Nombre:        req.Nombre,
		Tipo:          req.Tipo,
		URL:           req.URL,
		Formato:       req.Formato,
	})
	if err != nil {
		handleError(c, err)
		return
	}

	response.Created(c, toDocumentoResponse(documento), "documento subido exitosamente")
}

// Validar godoc
// @Summary Validar documento
// @Description Valida o rechaza un documento
// @Tags documentos
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID del documento"
// @Param request body dto.ValidarDocumentoRequest true "Estado de validación"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /documentos/{id}/validar [patch]
func (h *DocumentoHandler) Validar(c *gin.Context) {
	actor, err := actorFromGin(c)
	if err != nil {
		handleError(c, err)
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		response.BadRequest(c, "id de documento inválido", "BAD_REQUEST", nil)
		return
	}

	var req dto.ValidarDocumentoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, "datos inválidos", err.Error())
		return
	}

	validStates := map[string]bool{"Validado": true, "Rechazado": true, "Incompleto": true}
	if !validStates[req.Estado] {
		response.ValidationError(c, "estado inválido", "estado debe ser: Validado, Rechazado o Incompleto")
		return
	}

	comentario := ""
	if req.Comentario != nil {
		comentario = *req.Comentario
	}

	documento, err := h.useCase.Validar(c.Request.Context(), actor, id, req.Estado, comentario)
	if err != nil {
		handleError(c, err)
		return
	}

	response.Success(c, toDocumentoResponse(documento), "documento validado exitosamente")
}

// Listar godoc
// @Summary Listar documentos
// @Description Lista los documentos de una incapacidad
// @Tags documentos
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID de la incapacidad"
// @Param id_incapacidad query int true "ID de la incapacidad"
// @Param estado query string false "Filtrar por estado"
// @Param tipo query string false "Filtrar por tipo"
// @Param page query int false "Página" default(1)
// @Param limit query int false "Límite" default(20)
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /incapacidades/{id}/documentos [get]
func (h *DocumentoHandler) Listar(c *gin.Context) {
	actor, err := actorFromGin(c)
	if err != nil {
		handleError(c, err)
		return
	}

	var query dto.ListarDocumentosQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.ValidationError(c, "filtros inválidos", err.Error())
		return
	}

	if query.IDIncapacidad == 0 {
		response.BadRequest(c, "id_incapacidad es requerido", "BAD_REQUEST", nil)
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

	items, total, err := h.useCase.Listar(c.Request.Context(), actor, query.IDIncapacidad, query.Estado, query.Tipo, page, limit)
	if err != nil {
		handleError(c, err)
		return
	}

	response.Paginated(c, toDocumentoResponses(items), total, int64(page), int64(limit))
}

// Eliminar godoc
// @Summary Eliminar documento
// @Description Elimina un documento del sistema
// @Tags documentos
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID del documento"
// @Success 204
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /documentos/{id} [delete]
func (h *DocumentoHandler) Eliminar(c *gin.Context) {
	actor, err := actorFromGin(c)
	if err != nil {
		handleError(c, err)
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		response.BadRequest(c, "id de documento inválido", "BAD_REQUEST", nil)
		return
	}

	if err := h.useCase.Eliminar(c.Request.Context(), actor, id); err != nil {
		handleError(c, err)
		return
	}

	response.NoContent(c)
}

// ListarHistorial godoc
// @Summary Listar historial de incapacidad
// @Description Lista el historial de cambios de una incapacidad
// @Tags incapacidades
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID de la incapacidad"
// @Param id_incapacidad query int false "ID de la incapacidad"
// @Param id_tipo_historial query int false "Filtrar por tipo de evento"
// @Param page query int false "Página" default(1)
// @Param limit query int false "Límite" default(20)
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /incapacidades/{id}/historial [get]
func (h *DocumentoHandler) ListarHistorial(c *gin.Context) {
	var query dto.ListarHistorialQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.ValidationError(c, "filtros inválidos", err.Error())
		return
	}

	if query.IDIncapacidad == 0 {
		response.BadRequest(c, "id_incapacidad es requerido", "BAD_REQUEST", nil)
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

	var tipoID *uint64
	if query.IDTipoHistorial > 0 {
		tipoID = &query.IDTipoHistorial
	}

	items, total, err := h.historialListFn(query.IDIncapacidad, tipoID, page, limit)
	if err != nil {
		var appErr *apperrors.AppError
		if errors.As(err, &appErr) {
			response.Error(c, appErr.HTTPStatus, appErr.Message, appErr.Code, appErr.Details)
			return
		}
		response.InternalError(c, "error interno", "INTERNAL_ERROR")
		return
	}

	response.Paginated(c, toHistorialResponses(items), total, int64(page), int64(limit))
}

func toDocumentoResponse(d *incapacidaddomain.Documento) dto.DocumentoResponse {
	resp := dto.DocumentoResponse{
		IDDocumento:   d.IDDocumento,
		IDIncapacidad: d.IDIncapacidad,
		Nombre:        d.Nombre,
		Tipo:          d.Tipo,
		URL:           d.URL,
		Formato:       d.Formato,
		Estado:        d.Estado,
		Comentario:    d.Comentario,
		ValidadoPor:   d.ValidadoPor,
		CreatedAt:     d.CreatedAt.Format(time.RFC3339),
	}
	if !d.FechaCarga.IsZero() {
		resp.FechaCarga = d.FechaCarga.Format("2006-01-02 15:04:05")
	}
	if d.FechaValidacion != nil {
		resp.FechaValidacion = stringPtr(d.FechaValidacion.Format("2006-01-02 15:04:05"))
	}
	return resp
}

func toDocumentoResponses(items []incapacidaddomain.Documento) []dto.DocumentoResponse {
	result := make([]dto.DocumentoResponse, 0, len(items))
	for _, item := range items {
		result = append(result, toDocumentoResponse(&item))
	}
	return result
}

func toHistorialResponses(items []domain.Historial) []dto.HistorialResponse {
	result := make([]dto.HistorialResponse, 0, len(items))
	for _, item := range items {
		result = append(result, dto.HistorialResponse{
			IDHistorial:    item.IDHistorial,
			IDIncapacidad:  item.IDIncapacidad,
			IDTipoHistorial: item.IDTipoHistorial,
			Descripcion:    item.Descripcion,
			Fecha:          item.Fecha.Format("2006-01-02 15:04:05"),
			GestorID:       item.GestorID,
		})
	}
	return result
}

func stringPtr(s string) *string {
	return &s
}

// GenerarURLPrefirmada godoc
// @Summary Generar URL prefirmada para upload
// @Description Genera una URL prefirmada para subir documentos a R2
// @Tags documentos
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID de la incapacidad"
// @Param request body map[string]interface{} true "Datos del archivo (nombre, formato, tipo)"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /incapacidades/{id}/documentos/url [post]
func (h *DocumentoHandler) GenerarURLPrefirmada(c *gin.Context) {
	if h.storageService == nil {
		response.InternalError(c, "servicio de almacenamiento no disponible", "STORAGE_NOT_CONFIGURED")
		return
	}

	idStr := c.Param("id")
	incapacidadID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil || incapacidadID == 0 {
		response.BadRequest(c, "id de incapacidad inválido", "BAD_REQUEST", nil)
		return
	}

	var req struct {
		Nombre  string `json:"nombre" binding:"required"`
		Formato string `json:"formato" binding:"required"`
		Tipo    string `json:"tipo" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, "datos inválidos", err.Error())
		return
	}

	contentType := h.storageService.GetValidator().GetMimeTypeFromExtension(req.Nombre)
	if contentType == "" {
		ext := "." + req.Formato
		contentType = getContentTypeFromExtension(ext)
	}

	if err := h.storageService.Validate(contentType, h.storageService.GetMaxFileSize()); err != nil {
		handleStorageError(c, err)
		return
	}

	filename := req.Nombre + "." + req.Formato
	result, err := h.storageService.GenerateUploadURL(c.Request.Context(), filename, contentType, incapacidadID)
	if err != nil {
		handleStorageError(c, err)
		return
	}

	response.Success(c, gin.H{
		"upload_url": result.URL,
		"key":       result.Key,
		"expires_at": result.ExpiresAt.Format(time.RFC3339),
	}, "URL prefirmada generada")
}

// SubirBinario godoc
// @Summary Subir documento (binario)
// @Description Sube un archivo directamente al bucket R2
// @Tags documentos
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID de la incapacidad"
// @Param file formData file true "Archivo PDF/JPG/PNG (max 10MB)"
// @Param tipo formData string true "Tipo de documento"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 413 {object} map[string]interface{}
// @Router /incapacidades/{id}/documentos/upload [post]
func (h *DocumentoHandler) SubirBinario(c *gin.Context) {
	if h.storageService == nil {
		response.InternalError(c, "servicio de almacenamiento no disponible", "STORAGE_NOT_CONFIGURED")
		return
	}

	actor, err := actorFromGin(c)
	if err != nil {
		handleError(c, err)
		return
	}

	idStr := c.Param("id")
	incapacidadID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil || incapacidadID == 0 {
		response.BadRequest(c, "id de incapacidad inválido", "BAD_REQUEST", nil)
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		response.BadRequest(c, "archivo requerido", "FILE_REQUIRED", nil)
		return
	}

	tipo := c.PostForm("tipo")
	if tipo == "" {
		response.BadRequest(c, "tipo de documento requerido", "TIPO_REQUIRED", nil)
		return
	}

	contentType := file.Header.Get("Content-Type")
	if err := h.storageService.Validate(contentType, h.storageService.GetMaxFileSize()); err != nil {
		handleStorageError(c, err)
		return
	}

	filename := file.Filename
	result, err := h.storageService.UploadFile(c.Request.Context(), file, filename, contentType, incapacidadID)
	if err != nil {
		handleStorageError(c, err)
		return
	}

	ext := getExtensionFromFilename(filename)
	documento, err := h.useCase.Subir(c.Request.Context(), actor, struct {
		IDIncapacidad uint64
		Nombre        string
		Tipo          string
		URL           string
		Formato       string
	}{
		IDIncapacidad: incapacidadID,
		Nombre:        filename,
		Tipo:          tipo,
		URL:           result.URL,
		Formato:       ext,
	})
	if err != nil {
		handleError(c, err)
		return
	}

	response.Created(c, toDocumentoResponse(documento), "documento subido exitosamente")
}

func getContentTypeFromExtension(ext string) string {
	switch ext {
	case ".pdf":
		return "application/pdf"
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".png":
		return "image/png"
	default:
		return "application/octet-stream"
	}
}

func getExtensionFromFilename(filename string) string {
	ext := ""
	if len(filename) > 4 {
		ext = filename[len(filename)-4:]
		if ext[0] != '.' {
			ext = filename[len(filename)-5:]
		}
	}
	return ext
}

func handleStorageError(c *gin.Context, err error) {
	var storageErr *storage.StorageError
	if errors.As(err, &storageErr) {
		response.Error(c, storageErr.HTTPStatus, storageErr.Message, storageErr.Code, nil)
		return
	}
	response.InternalError(c, "error interno", "INTERNAL_ERROR")
}

type DocumentoHandlerInterface interface {
	Subir(c *gin.Context)
	Validar(c *gin.Context)
	Listar(c *gin.Context)
	Eliminar(c *gin.Context)
	ListarHistorial(c *gin.Context)
	GenerarURLPrefirmada(c *gin.Context)
	SubirBinario(c *gin.Context)
}

var _ DocumentoHandlerInterface = (*DocumentoHandler)(nil)