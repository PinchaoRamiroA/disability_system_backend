package authhttp

import (
	"strings"

	"disability_system_backend/internal/shared/auth"
	apperrors "disability_system_backend/internal/shared/errors"
	"disability_system_backend/internal/shared/response"

	"github.com/gin-gonic/gin"
)

type JWTMiddleware struct {
	jwtService *auth.JWTService
}

func NewJWTMiddleware(jwtService *auth.JWTService) *JWTMiddleware {
	return &JWTMiddleware{jwtService: jwtService}
}

func (m *JWTMiddleware) Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Unauthorized(c, "autorización requerida", apperrors.ErrUnauthorized.Code)
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			response.Error(c, apperrors.ErrTokenInvalid.HTTPStatus, "formato de token inválido", apperrors.ErrTokenInvalid.Code, nil)
			c.Abort()
			return
		}

		tokenString := parts[1]
		claims, err := m.jwtService.ValidateToken(tokenString)
		if err != nil {
			response.Error(c, apperrors.ErrTokenInvalid.HTTPStatus, "token inválido", apperrors.ErrTokenInvalid.Code, nil)
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("user_email", claims.Email)
		c.Set("user_role", claims.Role)

		c.Next()
	}
}

func (m *JWTMiddleware) RequireRole(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("user_role")
		if !exists {
			response.Unauthorized(c, "rol no encontrado en contexto", apperrors.ErrUnauthorized.Code)
			c.Abort()
			return
		}

		role, ok := userRole.(string)
		if !ok {
			response.InternalError(c, "error al procesar rol", apperrors.ErrInternal.Code)
			c.Abort()
			return
		}

		for _, r := range roles {
			if role == r {
				c.Next()
				return
			}
		}

		response.Forbidden(c, "no tienes permiso para acceder a este recurso", apperrors.ErrForbidden.Code)
		c.Abort()
	}
}
