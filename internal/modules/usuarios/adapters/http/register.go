package http

import (
	authhttp "disability_system_backend/internal/modules/auth/adapters/http"
	"disability_system_backend/internal/modules/usuarios/adapters/postgres"
	"disability_system_backend/internal/modules/usuarios/usecase"
	"disability_system_backend/internal/shared/auth"
	"disability_system_backend/internal/shared/router"

	"gorm.io/gorm"
)

func Register(v1 *router.APIVersion, db *gorm.DB, jwtService *auth.JWTService) {
	rolRepo := postgres.NewRolRepository(db)

	usuarioUseCase := usecase.NewUsuarioUseCase(db)
	rolUseCase := usecase.NewRolUseCase(rolRepo)

	usuarioHandler := NewUsuarioHandler(usuarioUseCase)
	rolHandler := NewRolHandler(rolUseCase)

	jwtMiddleware := authhttp.NewJWTMiddleware(jwtService)

	usuariosGroup := v1.Group("/usuarios", jwtMiddleware.Authenticate())
	{
		usuariosGroup.GET("", usuarioHandler.Listar)
		usuariosGroup.GET("/:id", usuarioHandler.Obtener)
		usuariosGroup.POST("", usuarioHandler.Crear)
		usuariosGroup.PUT("/:id", usuarioHandler.Actualizar)
		usuariosGroup.PATCH("/:id/estado", usuarioHandler.CambiarEstado)
		usuariosGroup.POST("/:id/rol", usuarioHandler.AsignarRol)
		usuariosGroup.POST("/:id/password", usuarioHandler.CambiarPassword)
		usuariosGroup.DELETE("/:id", usuarioHandler.Eliminar)
	}

	rolesGroup := v1.Group("/roles", jwtMiddleware.Authenticate())
	{
		rolesGroup.GET("", rolHandler.Listar)
		rolesGroup.GET("/:id", rolHandler.Obtener)
		rolesGroup.POST("", rolHandler.Crear)
		rolesGroup.PUT("/:id", rolHandler.Actualizar)
		rolesGroup.DELETE("/:id", rolHandler.Eliminar)
	}
}
