package http

import (
	authhttp "disability_system_backend/internal/modules/auth/adapters/http"
	"disability_system_backend/internal/modules/incapacidades/adapters/postgres"
	"disability_system_backend/internal/modules/incapacidades/usecase"
	"disability_system_backend/internal/shared/auth"
	"disability_system_backend/internal/shared/router"

	"gorm.io/gorm"
)

func Register(v1 *router.APIVersion, db *gorm.DB, jwtService *auth.JWTService) {
	incapacidadRepo := postgres.NewIncapacidadRepository(db)
	permissionRepo := postgres.NewPermissionRepository(db)
	incapacidadUseCase := usecase.NewIncapacidadUseCase(incapacidadRepo)
	incapacidadHandler := NewIncapacidadHandler(incapacidadUseCase)

	jwtMiddleware := authhttp.NewJWTMiddleware(jwtService)
	permissionMiddleware := NewPermissionMiddleware(permissionRepo)

	group := v1.Group("/incapacidades", jwtMiddleware.Authenticate(), permissionMiddleware.LoadActor())
	{
		group.GET("", incapacidadHandler.Listar)
		group.POST("", incapacidadHandler.Crear)
		group.GET("/estados", incapacidadHandler.ListarEstados)
		group.GET("/tipos", incapacidadHandler.ListarTipos)
		group.GET("/entidades", incapacidadHandler.ListarEntidades)
		group.GET("/:id", incapacidadHandler.Obtener)
		group.PUT("/:id", incapacidadHandler.Actualizar)
		group.PATCH("/:id/estado", incapacidadHandler.CambiarEstado)
		group.DELETE("/:id", incapacidadHandler.Archivar)
	}
}
