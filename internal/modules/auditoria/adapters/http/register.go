package http

import (
	authhttp "disability_system_backend/internal/modules/auth/adapters/http"
	"disability_system_backend/internal/modules/auditoria/adapters/postgres"
	"disability_system_backend/internal/modules/auditoria/usecase"
	"disability_system_backend/internal/shared/auth"
	"disability_system_backend/internal/shared/router"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Register(v1 *router.APIVersion, db *gorm.DB, jwtService *auth.JWTService) {
	repo := postgres.NewAuditoriaRepository(db)
	uc := usecase.NewAuditoriaUseCase(repo)
	handler := NewAuditoriaHandler(uc)

	jwtMiddleware := authhttp.NewJWTMiddleware(jwtService)

	group := v1.Group("/auditoria", jwtMiddleware.Authenticate(), extractActor())
	{
		group.GET("", handler.Listar)
		group.GET("/usuario/:id", handler.ListarPorUsuario)
	}
}

// extractActor is a simple middleware to extract permissions 
// similar to the one used in other modules
func extractActor() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetUint64("user_id")
		if userID == 0 {
			c.AbortWithStatusJSON(401, gin.H{"error": "Unauthorized"})
			return
		}

		rawPermisos, _ := c.Get("permisos")
		var permisos []string
		if p, ok := rawPermisos.([]string); ok {
			permisos = p
		}

		actor := &actorWrapper{
			userID:   userID,
			permisos: permisos,
		}

		// Inject the actor back just in case, but actorFromGin will extract it from the context directly
		// using c.GetUint64 and c.Get.
		_ = actor
		c.Next()
	}
}
