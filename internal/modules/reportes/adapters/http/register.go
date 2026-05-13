package http

import (
	authhttp "disability_system_backend/internal/modules/auth/adapters/http"
	"disability_system_backend/internal/modules/reportes/adapters/postgres"
	reportesports "disability_system_backend/internal/modules/reportes/ports"
	"disability_system_backend/internal/modules/reportes/usecase"
	"disability_system_backend/internal/shared/auth"
	"disability_system_backend/internal/shared/router"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Register(v1 *router.APIVersion, db *gorm.DB, jwtService *auth.JWTService) {
	reportesRepo := postgres.NewReportesRepository(db)
	permissionRepo := postgres.NewPermissionRepository(db)
	reportesUseCase := usecase.NewReportesUseCase(reportesRepo)
	reportesHandler := NewReportesHandler(reportesUseCase)

	jwtMiddleware := authhttp.NewJWTMiddleware(jwtService)
	permissionMiddleware := NewPermissionMiddleware(permissionRepo)

	reportesGroup := v1.Group("/reportes", jwtMiddleware.Authenticate(), permissionMiddleware.LoadActor())
	{
		reportesGroup.POST("", reportesHandler.GenerarReporte)
		reportesGroup.GET("/incapacidades", reportesHandler.GenerarReporte)
		reportesGroup.GET("/ausentismo", reportesHandler.GenerarReporte)
		reportesGroup.GET("/cartera", reportesHandler.GenerarReporte)
		reportesGroup.GET("/entidades/:entidad_id", reportesHandler.GenerarReporteEntidad)
		reportesGroup.GET("/vencimientos", reportesHandler.GenerarReporteVencimientos)
		reportesGroup.GET("/resumen-ejecutivo", reportesHandler.ObtenerResumenEjecutivo)
	}
}

type PermissionMiddleware struct {
	permissionRepo reportesports.PermissionRepository
}

func NewPermissionMiddleware(permissionRepo reportesports.PermissionRepository) *PermissionMiddleware {
	return &PermissionMiddleware{permissionRepo: permissionRepo}
}

func (m *PermissionMiddleware) LoadActor() gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDValue, ok := c.Get("user_id")
		if !ok {
			c.AbortWithStatusJSON(401, gin.H{"success": false, "message": "usuario no autenticado"})
			return
		}
		roleValue, ok := c.Get("user_role")
		if !ok {
			c.AbortWithStatusJSON(401, gin.H{"success": false, "message": "rol no encontrado"})
			return
		}
		userID, ok := userIDValue.(uint64)
		if !ok {
			c.AbortWithStatusJSON(500, gin.H{"success": false, "message": "usuario inválido en contexto"})
			return
		}
		role, ok := roleValue.(string)
		if !ok {
			c.AbortWithStatusJSON(500, gin.H{"success": false, "message": "rol inválido en contexto"})
			return
		}

		permisos, err := m.permissionRepo.FindPermissionsByRoleName(c.Request.Context(), role)
		if err != nil {
			c.AbortWithStatusJSON(403, gin.H{"success": false, "message": "rol sin permisos"})
			return
		}

		actor := reportesports.Actor{
			UserID:   userID,
			Role:     role,
			Permisos: permisos,
		}
		c.Set("actor", actor)
		c.Next()
	}
}
