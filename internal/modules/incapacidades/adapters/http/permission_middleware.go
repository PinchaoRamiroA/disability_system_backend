package http

import (
	"disability_system_backend/internal/modules/incapacidades/ports"
	"disability_system_backend/internal/shared/response"

	"github.com/gin-gonic/gin"
)

type PermissionMiddleware struct {
	permissionRepo ports.PermissionRepository
}

func NewPermissionMiddleware(permissionRepo ports.PermissionRepository) *PermissionMiddleware {
	return &PermissionMiddleware{permissionRepo: permissionRepo}
}

func (m *PermissionMiddleware) LoadActor() gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDValue, ok := c.Get("user_id")
		if !ok {
			response.Unauthorized(c, "usuario no autenticado", "UNAUTHORIZED")
			c.Abort()
			return
		}
		roleValue, ok := c.Get("user_role")
		if !ok {
			response.Unauthorized(c, "rol no encontrado", "UNAUTHORIZED")
			c.Abort()
			return
		}
		userID, ok := userIDValue.(uint64)
		if !ok {
			response.InternalError(c, "usuario inválido en contexto", "INVALID_USER_CONTEXT")
			c.Abort()
			return
		}
		role, ok := roleValue.(string)
		if !ok {
			response.InternalError(c, "rol inválido en contexto", "INVALID_ROLE_CONTEXT")
			c.Abort()
			return
		}

		permisos, err := m.permissionRepo.FindPermissionsByRoleName(c.Request.Context(), role)
		if err != nil {
			handleError(c, err)
			c.Abort()
			return
		}

		actor := ports.Actor{
			UserID:   userID,
			Role:     role,
			Permisos: permisos,
		}
		c.Set("actor", actor)
		c.Request = c.Request.WithContext(contextWithActor(c.Request.Context(), actor))
		c.Next()
	}
}
