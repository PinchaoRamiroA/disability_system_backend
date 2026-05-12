package http

import (
	"context"

	authhttp "disability_system_backend/internal/modules/auth/adapters/http"
	historialpostgres "disability_system_backend/internal/modules/historial/adapters/postgres"
	historialdomain "disability_system_backend/internal/modules/historial/domain"
	historialuc "disability_system_backend/internal/modules/historial/usecase"
	inicapapostgres "disability_system_backend/internal/modules/incapacidades/adapters/postgres"
	"disability_system_backend/internal/modules/incapacidades/usecase"
	"disability_system_backend/internal/shared/auth"
	"disability_system_backend/internal/shared/router"
	"disability_system_backend/internal/shared/storage"

	"gorm.io/gorm"
)

func Register(v1 *router.APIVersion, db *gorm.DB, jwtService *auth.JWTService, storageService *storage.StorageService) {
	incapacidadRepo := inicapapostgres.NewIncapacidadRepository(db)
	documentoRepo := inicapapostgres.NewDocumentoRepository(db)
	permissionRepo := inicapapostgres.NewPermissionRepository(db)
	historialRepo := historialpostgres.NewHistorialRepository(db)

	incapacidadUseCase := usecase.NewIncapacidadUseCase(incapacidadRepo)

	historialService := historialuc.NewHistorialService(historialRepo)
	incapacidadUseCase.SetHistorialService(func(ctx context.Context, incapacidadID uint64, tipoID uint64, descripcion string, gestorID *uint64) error {
		return historialService.CreateEntry(ctx, incapacidadID, tipoID, descripcion, gestorID)
	})

	documentoUseCase := usecase.NewDocumentoUseCase(documentoRepo, historialService)

	incapacidadHandler := NewIncapacidadHandler(incapacidadUseCase)
	documentoHandler := NewDocumentoHandler(documentoUseCase, func(incapacidadID uint64, tipoID *uint64, page, limit int) ([]historialdomain.Historial, int64, error) {
		return historialService.List(context.Background(), incapacidadID, tipoID, page, limit)
	}, storageService)

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
		group.GET("/:id/documentos", documentoHandler.Listar)
		group.POST("/:id/documentos", documentoHandler.Subir)
		group.POST("/:id/documentos/url", documentoHandler.GenerarURLPrefirmada)
		group.GET("/:id/historial", documentoHandler.ListarHistorial)
		group.GET("/:id/plazos", incapacidadHandler.ObtenerPlazos)

	group.GET("/tipos/:tipo_id/documentos-requeridos", incapacidadHandler.ObtenerDocumentosRequeridos)
	}

	docGroup := v1.Group("/documentos", jwtMiddleware.Authenticate(), permissionMiddleware.LoadActor())
	{
		docGroup.PATCH("/:id/validar", documentoHandler.Validar)
		docGroup.DELETE("/:id", documentoHandler.Eliminar)
	}

	catalogosGroup := v1.Group("/catalogos", jwtMiddleware.Authenticate(), permissionMiddleware.LoadActor())
	{
		catalogosGroup.GET("/estados-documento", incapacidadHandler.ListarEstadosDocumento)
		catalogosGroup.GET("/tipos-documento", incapacidadHandler.ListarTiposDocumento)
		catalogosGroup.GET("/tipos-pago", incapacidadHandler.ListarTiposPago)
	}
}
