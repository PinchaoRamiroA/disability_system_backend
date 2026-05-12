package router

import (
	"disability_system_backend/internal/shared/config"
	"disability_system_backend/internal/shared/logger"
	"disability_system_backend/internal/shared/response"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"
)

// @title Sistema de Gestión de Incapacidades
// @version 1.0
// @description API REST para la gestión de incapacidades médicas, cobro a EPS/ARL y seguimiento documental
// @BasePath /api/v1
// @schemes http https

type Router struct {
	engine *gin.Engine
	config *config.Config
	logger *logger.Logger
}

type APIVersion struct {
	group *gin.RouterGroup
}

func New(cfg *config.Config, log *logger.Logger) *Router {
	if cfg.App.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	return &Router{
		engine: gin.New(),
		config: cfg,
		logger: log,
	}
}

func (r *Router) Engine() *gin.Engine {
	return r.engine
}

func (r *Router) V1() *APIVersion {
	return &APIVersion{
		group: r.engine.Group("/api/v1"),
	}
}

func (r *Router) Public() *APIVersion {
	return &APIVersion{
		group: r.engine.Group(""),
	}
}

func (r *Router) Use(middlewares ...gin.HandlerFunc) {
	r.engine.Use(middlewares...)
}

func (r *Router) SetupHealth(db *gorm.DB) {
	handler := healthHandler{db: db}
	r.engine.GET("/health/live", handler.liveness)
	r.engine.GET("/health/ready", handler.readiness)
	r.engine.GET("/health", handler.readiness)
	r.engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}

type healthHandler struct {
	db *gorm.DB
}

func (h *healthHandler) liveness(c *gin.Context) {
	response.Success(c, gin.H{"status": "alive"}, "service is alive")
}

func (h *healthHandler) readiness(c *gin.Context) {
	if h.db == nil {
		response.Error(c, 503, "database not initialized", "DB_NOT_INITIALIZED", nil)
		return
	}

	sqlDB, err := h.db.DB()
	if err != nil {
		response.Error(c, 503, "database error", "DB_ERROR", nil)
		return
	}

	if err := sqlDB.Ping(); err != nil {
		response.Error(c, 503, "database down", "DB_DOWN", nil)
		return
	}

	response.Success(c, gin.H{
		"status":   "ready",
		"database": "connected",
	}, "service is ready")
}

func (r *Router) Run(port string) error {
	return r.engine.Run(port)
}

func (v *APIVersion) GET(path string, handlers ...gin.HandlerFunc) {
	v.group.GET(path, handlers...)
}

func (v *APIVersion) POST(path string, handlers ...gin.HandlerFunc) {
	v.group.POST(path, handlers...)
}

func (v *APIVersion) PUT(path string, handlers ...gin.HandlerFunc) {
	v.group.PUT(path, handlers...)
}

func (v *APIVersion) PATCH(path string, handlers ...gin.HandlerFunc) {
	v.group.PATCH(path, handlers...)
}

func (v *APIVersion) DELETE(path string, handlers ...gin.HandlerFunc) {
	v.group.DELETE(path, handlers...)
}

func (v *APIVersion) Group(prefix string, handlers ...gin.HandlerFunc) *APIVersion {
	return &APIVersion{
		group: v.group.Group(prefix, handlers...),
	}
}

func (r *Router) RegisterModule(fn func(*APIVersion), v *APIVersion) {
	fn(v)
}

type ModuleRegistrar func(*APIVersion)

func (r *Router) Register(modules ...ModuleRegistrar) {
	for _, m := range modules {
		m(r.V1())
	}
}
