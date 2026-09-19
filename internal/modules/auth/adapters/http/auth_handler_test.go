package authhttp_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	authhttp "disability_system_backend/internal/modules/auth/adapters/http"
	"disability_system_backend/internal/modules/auth/ports"
	"disability_system_backend/internal/modules/auth/usecase"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockTokenService struct {
	mock.Mock
}

func (m *mockTokenService) GenerateTokenPair(userID uint64, email, role string) (*ports.TokenPair, error) {
	args := m.Called(userID, email, role)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ports.TokenPair), args.Error(1)
}

func (m *mockTokenService) ValidateToken(token string) (*ports.TokenClaims, error) {
	args := m.Called(token)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ports.TokenClaims), args.Error(1)
}

func (m *mockTokenService) RefreshToken(token string) (*ports.TokenPair, error) {
	args := m.Called(token)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ports.TokenPair), args.Error(1)
}

func setupAuthRouter(handler *authhttp.AuthHandler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.POST("/auth/refresh", handler.RefreshToken)
	return r
}

func TestAuthHandler_RefreshToken_BUG10(t *testing.T) {
	tokenExpiry := 2 * time.Hour // 7200 seconds

	t.Run("should return positive expires_in instead of 0 on token refresh (BUG-10 fix)", func(t *testing.T) {
		tokenSvc := new(mockTokenService)
		refreshUC := usecase.NewRefreshTokenUseCase(tokenSvc, tokenExpiry)
		handler := authhttp.NewAuthHandler(nil, nil, refreshUC)
		router := setupAuthRouter(handler)

		expectedTokens := &ports.TokenPair{
			AccessToken:  "new_access_token_123",
			RefreshToken: "new_refresh_token_456",
		}
		tokenSvc.On("RefreshToken", "valid_refresh_token").Return(expectedTokens, nil)

		body := bytes.NewBufferString(`{"refresh_token": "valid_refresh_token"}`)
		req := httptest.NewRequest(http.MethodPost, "/auth/refresh", body)
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var resp map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.True(t, resp["success"].(bool))

		data, ok := resp["data"].(map[string]interface{})
		assert.True(t, ok)
		assert.Equal(t, "new_access_token_123", data["access_token"])
		assert.Equal(t, "new_refresh_token_456", data["refresh_token"])
		assert.Equal(t, "Bearer", data["token_type"])

		// Critical assertion for BUG-10: expires_in must be 7200 and NOT 0!
		expiresIn, ok := data["expires_in"].(float64)
		assert.True(t, ok)
		assert.Equal(t, float64(7200), expiresIn, "expires_in must equal token expiration seconds and not 0")

		tokenSvc.AssertExpectations(t)
	})

	t.Run("should fallback to loginUseCase expiration if refreshUseCase has 0", func(t *testing.T) {
		tokenSvc := new(mockTokenService)
		// refreshUC without explicit expiry
		refreshUC := usecase.NewRefreshTokenUseCase(tokenSvc)
		// loginUC with 3600s expiry
		loginUC := usecase.NewLoginUseCase(nil, nil, tokenSvc, nil, time.Hour)
		handler := authhttp.NewAuthHandler(loginUC, nil, refreshUC)
		router := setupAuthRouter(handler)

		expectedTokens := &ports.TokenPair{
			AccessToken:  "access_token",
			RefreshToken: "refresh_token",
		}
		tokenSvc.On("RefreshToken", "valid_token").Return(expectedTokens, nil)

		body := bytes.NewBufferString(`{"refresh_token": "valid_token"}`)
		req := httptest.NewRequest(http.MethodPost, "/auth/refresh", body)
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var resp map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)

		data := resp["data"].(map[string]interface{})
		assert.Equal(t, float64(3600), data["expires_in"])
	})

	t.Run("should return error when refresh token is invalid or expired", func(t *testing.T) {
		tokenSvc := new(mockTokenService)
		refreshUC := usecase.NewRefreshTokenUseCase(tokenSvc, tokenExpiry)
		handler := authhttp.NewAuthHandler(nil, nil, refreshUC)
		router := setupAuthRouter(handler)

		tokenSvc.On("RefreshToken", "expired_token").Return(nil, errors.New("token expired"))

		body := bytes.NewBufferString(`{"refresh_token": "expired_token"}`)
		req := httptest.NewRequest(http.MethodPost, "/auth/refresh", body)
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Contains(t, w.Body.String(), "TOKEN_EXPIRED")
	})

	t.Run("should return error when body is invalid", func(t *testing.T) {
		tokenSvc := new(mockTokenService)
		refreshUC := usecase.NewRefreshTokenUseCase(tokenSvc, tokenExpiry)
		handler := authhttp.NewAuthHandler(nil, nil, refreshUC)
		router := setupAuthRouter(handler)

		body := bytes.NewBufferString(`{}`)
		req := httptest.NewRequest(http.MethodPost, "/auth/refresh", body)
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
	})
}
