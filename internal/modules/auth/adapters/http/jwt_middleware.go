package authhttp

import (
	"strings"

	"disability_system_backend/internal/shared/auth"
	apperrors "disability_system_backend/internal/shared/errors"

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
			c.AbortWithStatusJSON(apperrors.ErrUnauthorized.HTTPStatus, gin.H{
				"success": false,
				"code":    apperrors.ErrUnauthorized.Code,
				"message": "autorización requerida",
			})
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			c.AbortWithStatusJSON(apperrors.ErrTokenInvalid.HTTPStatus, gin.H{
				"success": false,
				"code":    apperrors.ErrTokenInvalid.Code,
				"message": "formato de token inválido",
			})
			return
		}

		tokenString := parts[1]
		claims, err := m.jwtService.ValidateToken(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(apperrors.ErrTokenInvalid.HTTPStatus, gin.H{
				"success": false,
				"code":    apperrors.ErrTokenInvalid.Code,
				"message": "token inválido",
			})
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
			c.AbortWithStatusJSON(apperrors.ErrUnauthorized.HTTPStatus, gin.H{
				"success": false,
				"code":    apperrors.ErrUnauthorized.Code,
				"message": "rol no encontrado en contexto",
			})
			return
		}

		role, ok := userRole.(string)
		if !ok {
			c.AbortWithStatusJSON(apperrors.ErrInternal.HTTPStatus, gin.H{
				"success": false,
				"code":    apperrors.ErrInternal.Code,
				"message": "error al procesar rol",
			})
			return
		}

		for _, r := range roles {
			if role == r {
				c.Next()
				return
			}
		}

		c.AbortWithStatusJSON(apperrors.ErrForbidden.HTTPStatus, gin.H{
			"success": false,
			"code":    apperrors.ErrForbidden.Code,
			"message": "no tienes permiso para acceder a este recurso",
		})
	}
}
