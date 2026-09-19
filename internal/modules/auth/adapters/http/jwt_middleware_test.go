package authhttp_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	authhttp "disability_system_backend/internal/modules/auth/adapters/http"
	"disability_system_backend/internal/shared/auth"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupTestJWTMiddleware() (*authhttp.JWTMiddleware, *auth.JWTService) {
	jwtSvc := auth.NewJWTService("test_secret_key_1234567890_32_bytes!", time.Hour, 24*time.Hour)
	middleware := authhttp.NewJWTMiddleware(jwtSvc)
	return middleware, jwtSvc
}

func TestJWTMiddleware_Authenticate_BUG09(t *testing.T) {
	gin.SetMode(gin.TestMode)
	middleware, jwtSvc := setupTestJWTMiddleware()

	setupRouter := func() *gin.Engine {
		r := gin.New()
		r.GET("/protected", middleware.Authenticate(), func(c *gin.Context) {
			userID, _ := c.Get("user_id")
			userEmail, _ := c.Get("user_email")
			userRole, _ := c.Get("user_role")
			c.JSON(http.StatusOK, gin.H{
				"user_id":    userID,
				"user_email": userEmail,
				"user_role":  userRole,
			})
		})
		return r
	}

	t.Run("should return standardized error structure when Authorization header is missing (BUG-09 fix)", func(t *testing.T) {
		router := setupRouter()

		req := httptest.NewRequest(http.MethodGet, "/protected", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)

		var rawResp map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &rawResp)
		assert.NoError(t, err)

		// Verification of standard response structure:
		// success should be false
		assert.Equal(t, false, rawResp["success"])
		assert.Equal(t, "autorización requerida", rawResp["message"])

		// code should NOT be at the root level (which was the bug)
		assert.Nil(t, rawResp["code"], "code should NOT be at the root level")

		// code should be nested in error: { code: ... }
		errorObj, ok := rawResp["error"].(map[string]interface{})
		assert.True(t, ok, "error must be an object")
		assert.Equal(t, "UNAUTHORIZED", errorObj["code"])
	})

	t.Run("should return standardized error structure when Bearer format is invalid", func(t *testing.T) {
		router := setupRouter()

		req := httptest.NewRequest(http.MethodGet, "/protected", nil)
		req.Header.Set("Authorization", "Basic invalidtoken123")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)

		var rawResp map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &rawResp)
		assert.NoError(t, err)

		assert.Equal(t, false, rawResp["success"])
		assert.Equal(t, "formato de token inválido", rawResp["message"])
		assert.Nil(t, rawResp["code"])

		errorObj, ok := rawResp["error"].(map[string]interface{})
		assert.True(t, ok)
		assert.Equal(t, "TOKEN_INVALID", errorObj["code"])
	})

	t.Run("should return standardized error structure when token is invalid", func(t *testing.T) {
		router := setupRouter()

		req := httptest.NewRequest(http.MethodGet, "/protected", nil)
		req.Header.Set("Authorization", "Bearer invalid.jwt.token")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)

		var rawResp map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &rawResp)
		assert.NoError(t, err)

		assert.Equal(t, false, rawResp["success"])
		assert.Equal(t, "token inválido", rawResp["message"])
		assert.Nil(t, rawResp["code"])

		errorObj, ok := rawResp["error"].(map[string]interface{})
		assert.True(t, ok)
		assert.Equal(t, "TOKEN_INVALID", errorObj["code"])
	})

	t.Run("should succeed when valid Bearer token is provided", func(t *testing.T) {
		router := setupRouter()

		token, err := jwtSvc.GenerateToken(42, "user@example.com", "admin")
		assert.NoError(t, err)

		req := httptest.NewRequest(http.MethodGet, "/protected", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var rawResp map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &rawResp)
		assert.NoError(t, err)

		assert.Equal(t, float64(42), rawResp["user_id"])
		assert.Equal(t, "user@example.com", rawResp["user_email"])
		assert.Equal(t, "admin", rawResp["user_role"])
	})
}

func TestJWTMiddleware_RequireRole_BUG09(t *testing.T) {
	gin.SetMode(gin.TestMode)
	middleware, _ := setupTestJWTMiddleware()

	t.Run("should return standardized error structure when role does not match", func(t *testing.T) {
		r := gin.New()
		r.GET("/admin-only", func(c *gin.Context) {
			c.Set("user_role", "colaborador")
			c.Next()
		}, middleware.RequireRole("admin"), func(c *gin.Context) {
			c.Status(http.StatusOK)
		})

		req := httptest.NewRequest(http.MethodGet, "/admin-only", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusForbidden, w.Code)

		var rawResp map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &rawResp)
		assert.NoError(t, err)

		assert.Equal(t, false, rawResp["success"])
		assert.Equal(t, "no tienes permiso para acceder a este recurso", rawResp["message"])
		assert.Nil(t, rawResp["code"])

		errorObj, ok := rawResp["error"].(map[string]interface{})
		assert.True(t, ok)
		assert.Equal(t, "FORBIDDEN", errorObj["code"])
	})

	t.Run("should return standardized error structure when role is missing in context", func(t *testing.T) {
		r := gin.New()
		r.GET("/admin-only", middleware.RequireRole("admin"), func(c *gin.Context) {
			c.Status(http.StatusOK)
		})

		req := httptest.NewRequest(http.MethodGet, "/admin-only", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)

		var rawResp map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &rawResp)
		assert.NoError(t, err)

		assert.Equal(t, false, rawResp["success"])
		assert.Equal(t, "rol no encontrado en contexto", rawResp["message"])
		assert.Nil(t, rawResp["code"])

		errorObj, ok := rawResp["error"].(map[string]interface{})
		assert.True(t, ok)
		assert.Equal(t, "UNAUTHORIZED", errorObj["code"])
	})

	t.Run("should pass when role matches", func(t *testing.T) {
		r := gin.New()
		r.GET("/admin-only", func(c *gin.Context) {
			c.Set("user_role", "admin")
			c.Next()
		}, middleware.RequireRole("admin", "superadmin"), func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"ok": true})
		})

		req := httptest.NewRequest(http.MethodGet, "/admin-only", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}
