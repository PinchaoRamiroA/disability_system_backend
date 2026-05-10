package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Success bool        `json:"success"`
	Data    any         `json:"data,omitempty"`
	Message string      `json:"message"`
	Error   *ErrorInfo  `json:"error,omitempty"`
}

type ErrorInfo struct {
	Code    string `json:"code,omitempty"`
	Details any    `json:"details,omitempty"`
}

func Success(c *gin.Context, data any, message string) {
	c.JSON(http.StatusOK, Response{
		Success: true,
		Data:    data,
		Message: message,
	})
}

func Created(c *gin.Context, data any, message string) {
	c.JSON(http.StatusCreated, Response{
		Success: true,
		Data:    data,
		Message: message,
	})
}

func NoContent(c *gin.Context) {
	c.Status(http.StatusNoContent)
}

func Error(c *gin.Context, statusCode int, message string, errCode string, details any) {
	c.JSON(statusCode, Response{
		Success: false,
		Message: message,
		Error: &ErrorInfo{
			Code:    errCode,
			Details: details,
		},
	})
}

func BadRequest(c *gin.Context, message string, errCode string, details any) {
	Error(c, http.StatusBadRequest, message, errCode, details)
}

func Unauthorized(c *gin.Context, message string, errCode string) {
	Error(c, http.StatusUnauthorized, message, errCode, nil)
}

func Forbidden(c *gin.Context, message string, errCode string) {
	Error(c, http.StatusForbidden, message, errCode, nil)
}

func NotFound(c *gin.Context, message string, errCode string) {
	Error(c, http.StatusNotFound, message, errCode, nil)
}

func Conflict(c *gin.Context, message string, errCode string) {
	Error(c, http.StatusConflict, message, errCode, nil)
}

func InternalError(c *gin.Context, message string, errCode string) {
	Error(c, http.StatusInternalServerError, message, errCode, nil)
}

func ValidationError(c *gin.Context, message string, details any) {
	Error(c, http.StatusUnprocessableEntity, message, "VALIDATION_ERROR", details)
}

func Paginated(c *gin.Context, data any, total, page, limit int64) {
	c.JSON(http.StatusOK, Response{
		Success: true,
		Data: map[string]any{
			"items":      data,
			"total":      total,
			"page":       page,
			"limit":      limit,
			"total_pages": (total + limit - 1) / limit,
		},
		Message: "success",
	})
}