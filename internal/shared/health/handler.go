package health

import (
	"errors"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var (
	ErrDatabaseNotInitialized = errors.New("database not initialized")
)

type Handler struct {
	db *gorm.DB
}

func NewHandler(db *gorm.DB) *Handler {
	return &Handler{db: db}
}

func (h *Handler) Liveness(c *gin.Context) {
	c.JSON(200, gin.H{
		"status": "alive",
	})
}

func (h *Handler) Readiness(c *gin.Context) {
	if h.db == nil {
		c.JSON(503, gin.H{
			"status":   "unhealthy",
			"database": "not initialized",
		})
		return
	}

	sqlDB, err := h.db.DB()
	if err != nil {
		c.JSON(503, gin.H{
			"status":   "unhealthy",
			"database": "error",
		})
		return
	}

	if err := sqlDB.Ping(); err != nil {
		c.JSON(503, gin.H{
			"status":   "unhealthy",
			"database": "down",
		})
		return
	}

	c.JSON(200, gin.H{
		"status":   "ready",
		"database": "connected",
	})
}