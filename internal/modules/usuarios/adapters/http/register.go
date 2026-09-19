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
	usuarioRepo := postgres.NewUsuarioRepository(db)
	rolRepo := postgres.NewRolRepository(db)

	usuarioUseCase := usecase.NewUsuarioUseCase(usuarioRepo, rolRepo)
	rolUseCase := usecase.NewRolUseCase(rolRepo)

	usuarioHandler := NewUsuarioHandler(usuarioUseCase)
	rolHandler := NewRolHandler(rolUseCase)

	jwtMiddleware := authhttp.NewJWTMiddleware(jwtService)

	usuariosGroup := v1.Group("/usuarios", jwtMiddleware.Authenticate())
	{
		usuariosGroup.GET("", jwtMiddleware.RequireRole("Administrador", "admin", "Gestión Humana", "SG-SST", "Recepcionista"), usuarioHandler.Listar)
		usuariosGroup.GET("/:id", usuarioHandler.Obtener)
		usuariosGroup.POST("", jwtMiddleware.RequireRole("Administrador", "admin", "Gestión Humana"), usuarioHandler.Crear)
		usuariosGroup.PUT("/:id", jwtMiddleware.RequireRole("Administrador", "admin", "Gestión Humana"), usuarioHandler.Actualizar)
		usuariosGroup.PATCH("/:id/estado", jwtMiddleware.RequireRole("Administrador", "admin", "Gestión Humana"), usuarioHandler.CambiarEstado)
		usuariosGroup.POST("/:id/rol", jwtMiddleware.RequireRole("Administrador", "admin"), usuarioHandler.AsignarRol)
		usuariosGroup.POST("/:id/password", jwtMiddleware.RequireRole("Administrador", "admin"), usuarioHandler.CambiarPassword)
		usuariosGroup.DELETE("/:id", jwtMiddleware.RequireRole("Administrador", "admin"), usuarioHandler.Eliminar)
	}

	rolesGroup := v1.Group("/roles", jwtMiddleware.Authenticate(), jwtMiddleware.RequireRole("Administrador", "admin", "Gestión Humana"))
	{
		rolesGroup.GET("", rolHandler.Listar)
		rolesGroup.GET("/:id", rolHandler.Obtener)
		rolesGroup.POST("", jwtMiddleware.RequireRole("Administrador", "admin"), rolHandler.Crear)
		rolesGroup.PUT("/:id", jwtMiddleware.RequireRole("Administrador", "admin"), rolHandler.Actualizar)
		rolesGroup.DELETE("/:id", jwtMiddleware.RequireRole("Administrador", "admin"), rolHandler.Eliminar)
	}
}
