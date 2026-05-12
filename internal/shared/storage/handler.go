package storage

import (
	"net/http"
	"time"

	"disability_system_backend/internal/shared/response"

	"github.com/gin-gonic/gin"
)

type StorageHandler struct {
	service *StorageService
}

func NewStorageHandler(service *StorageService) *StorageHandler {
	return &StorageHandler{service: service}
}

func (h *StorageHandler) GenerarURLSubida(c *gin.Context) {
	var req struct {
		IDIncapacidad uint64 `json:"id_incapacidad" binding:"required"`
		Nombre        string `json:"nombre" binding:"required"`
		Tipo          string `json:"tipo" binding:"required"`
		Formato       string `json:"formato" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Datos inválidos", "INVALID_REQUEST", err.Error())
		return
	}

	ext := "." + req.Formato
	contentType := h.service.GetValidator().GetMimeTypeFromExtension(req.Nombre)
	if contentType == "" {
		contentType = getContentTypeFromExtension(ext)
	}

	if !h.service.GetValidator().IsAllowedType(contentType) {
		response.Error(c, http.StatusUnsupportedMediaType, "Tipo de archivo no permitido", "INVALID_FILE_TYPE", nil)
		return
	}

	filename := req.Nombre + ext
	result, err := h.service.GenerateUploadURL(c.Request.Context(), filename, contentType, req.IDIncapacidad)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Error al generar URL", "PRESIGN_FAILED", err.Error())
		return
	}

	response.Success(c, gin.H{
		"upload_url":   result.URL,
		"key":          result.Key,
		"expires_at":   result.ExpiresAt.Format(time.RFC3339),
		"id_incapacidad": req.IDIncapacidad,
		"nombre":       filename,
		"tipo":         req.Tipo,
		"formato":       req.Formato,
		"max_size":     h.service.GetMaxFileSize(),
	}, "URL generada correctamente")
}

func (h *StorageHandler) ObtenerInfo(c *gin.Context) {
	allowedTypes := h.service.GetValidator().GetAllowedMimeTypes()
	extensions := h.service.GetValidator().GetAllowedExtensions()

	response.Success(c, gin.H{
		"max_file_size":    h.service.GetMaxFileSize(),
		"allowed_types":     allowedTypes,
		"allowed_extensions": extensions,
		"configured":        h.service.IsConfigured(),
	}, "Información de storage")
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
