package http

import (
	authhttp "disability_system_backend/internal/modules/auth/adapters/http"
	"disability_system_backend/internal/modules/reportes/adapters/postgres"
	"disability_system_backend/internal/modules/reportes/usecase"
	"disability_system_backend/internal/shared/auth"
	"disability_system_backend/internal/shared/router"

	"gorm.io/gorm"
)

func Register(v1 *router.APIVersion, db *gorm.DB, jwtService *auth.JWTService) {
	reportesRepo := postgres.NewReportesRepository(db)
	reportesUseCase := usecase.NewReportesUseCase(reportesRepo)
	reportesHandler := NewReportesHandler(reportesUseCase)

	jwtMiddleware := authhttp.NewJWTMiddleware(jwtService)

	reportesGroup := v1.Group("/reportes", jwtMiddleware.Authenticate())
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
