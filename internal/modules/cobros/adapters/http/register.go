package http

import (
	authhttp "disability_system_backend/internal/modules/auth/adapters/http"
	"disability_system_backend/internal/modules/cobros/adapters/postgres"
	"disability_system_backend/internal/modules/cobros/usecase"
	"disability_system_backend/internal/shared/auth"
	"disability_system_backend/internal/shared/router"

	"gorm.io/gorm"
)

func Register(v1 *router.APIVersion, db *gorm.DB, jwtService *auth.JWTService) {
	cobroRepo := postgres.NewCobroRepository(db)
	permissionRepo := postgres.NewPermissionRepository(db)
	cobroUseCase := usecase.NewCobroUseCase(cobroRepo)
	cobroHandler := NewCobroHandler(cobroUseCase)

	jwtMiddleware := authhttp.NewJWTMiddleware(jwtService)
	permissionMiddleware := NewPermissionMiddleware(permissionRepo)

	group := v1.Group("/cobros", jwtMiddleware.Authenticate(), permissionMiddleware.LoadActor())
	{
		group.GET("/pagos", cobroHandler.ListarPagos)
		group.POST("/pagos", cobroHandler.CrearPago)
		group.GET("/pagos/:id", cobroHandler.ObtenerPago)
		group.PUT("/pagos/:id", cobroHandler.ActualizarPago)
		group.DELETE("/pagos/:id", cobroHandler.EliminarPago)
		group.PATCH("/pagos/:id/conciliar", cobroHandler.ConciliarPago)

		group.GET("/seguimientos", cobroHandler.ListarSeguimientos)
		group.POST("/seguimientos", cobroHandler.CrearSeguimiento)
		group.GET("/seguimientos/:id", cobroHandler.ObtenerSeguimiento)
		group.PUT("/seguimientos/:id", cobroHandler.ActualizarSeguimiento)
	}
}
