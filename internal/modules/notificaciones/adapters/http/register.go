package http

import (
	authhttp "disability_system_backend/internal/modules/auth/adapters/http"
	"disability_system_backend/internal/modules/notificaciones/adapters/postgres"
	"disability_system_backend/internal/modules/notificaciones/usecase"
	"disability_system_backend/internal/shared/auth"
	"disability_system_backend/internal/shared/router"

	"gorm.io/gorm"
)

func Register(v1 *router.APIVersion, db *gorm.DB, jwtService *auth.JWTService) {
	notificacionRepo := postgres.NewNotificacionRepository(db)
	permissionRepo := postgres.NewPermissionRepository(db)
	notificacionUseCase := usecase.NewNotificacionUseCase(notificacionRepo)
	notificacionHandler := NewNotificacionHandler(notificacionUseCase)

	jwtMiddleware := authhttp.NewJWTMiddleware(jwtService)
	permissionMiddleware := NewPermissionMiddleware(permissionRepo)

	group := v1.Group("/notificaciones", jwtMiddleware.Authenticate(), permissionMiddleware.LoadActor())
	{
		group.GET("", notificacionHandler.Listar)
		group.POST("", notificacionHandler.Crear)
		group.GET("/no-leidas/count", notificacionHandler.ContarNoLeidas)
		group.PATCH("/marcar-todas-leidas", notificacionHandler.MarcarTodasLeidas)
		group.GET("/:id", notificacionHandler.Obtener)
		group.PATCH("/:id/leida", notificacionHandler.MarcarLeida)
		group.DELETE("/:id", notificacionHandler.Eliminar)
	}
}
