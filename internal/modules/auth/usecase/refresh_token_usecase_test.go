package usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"disability_system_backend/internal/modules/auth/ports"
	"disability_system_backend/internal/modules/auth/usecase"
	apperrors "disability_system_backend/internal/shared/errors"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockTokenService struct {
	mock.Mock
}

func (m *MockTokenService) GenerateTokenPair(userID uint64, email, role string) (*ports.TokenPair, error) {
	args := m.Called(userID, email, role)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ports.TokenPair), args.Error(1)
}

func (m *MockTokenService) ValidateToken(token string) (*ports.TokenClaims, error) {
	args := m.Called(token)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ports.TokenClaims), args.Error(1)
}

func (m *MockTokenService) RefreshToken(token string) (*ports.TokenPair, error) {
	args := m.Called(token)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ports.TokenPair), args.Error(1)
}

func TestRefreshTokenUseCase_Execute(t *testing.T) {
	ctx := context.Background()

	t.Run("should refresh tokens successfully", func(t *testing.T) {
		mockTokenSvc := new(MockTokenService)
		uc := usecase.NewRefreshTokenUseCase(mockTokenSvc, time.Hour)

		expectedTokens := &ports.TokenPair{
			AccessToken:  "new_access_token",
			RefreshToken: "new_refresh_token",
		}
		mockTokenSvc.On("RefreshToken", "valid_refresh_token").Return(expectedTokens, nil)

		tokens, err := uc.Execute(ctx, "valid_refresh_token")

		assert.NoError(t, err)
		assert.NotNil(t, tokens)
		assert.Equal(t, "new_access_token", tokens.AccessToken)
		assert.Equal(t, "new_refresh_token", tokens.RefreshToken)
		assert.Equal(t, int64(3600), uc.GetExpirationSeconds())

		mockTokenSvc.AssertExpectations(t)
	})

	t.Run("should return ErrTokenExpired when token service fails", func(t *testing.T) {
		mockTokenSvc := new(MockTokenService)
		uc := usecase.NewRefreshTokenUseCase(mockTokenSvc, time.Hour)

		mockTokenSvc.On("RefreshToken", "expired_token").Return(nil, errors.New("token expired"))

		tokens, err := uc.Execute(ctx, "expired_token")

		assert.Error(t, err)
		assert.Nil(t, tokens)

		var appErr *apperrors.AppError
		assert.ErrorAs(t, err, &appErr)
		assert.Equal(t, "TOKEN_EXPIRED", appErr.Code)

		mockTokenSvc.AssertExpectations(t)
	})

	t.Run("should return ErrTokenInvalid when tokens returned is nil without error", func(t *testing.T) {
		mockTokenSvc := new(MockTokenService)
		uc := usecase.NewRefreshTokenUseCase(mockTokenSvc)

		mockTokenSvc.On("RefreshToken", "bad_token").Return(nil, nil)

		tokens, err := uc.Execute(ctx, "bad_token")

		assert.Error(t, err)
		assert.Nil(t, tokens)

		var appErr *apperrors.AppError
		assert.ErrorAs(t, err, &appErr)
		assert.Equal(t, "TOKEN_INVALID", appErr.Code)

		mockTokenSvc.AssertExpectations(t)
	})

	t.Run("should return correct expiration seconds", func(t *testing.T) {
		mockTokenSvc := new(MockTokenService)
		uc := usecase.NewRefreshTokenUseCase(mockTokenSvc, 2*time.Hour)

		assert.Equal(t, int64(7200), uc.GetExpirationSeconds())
	})
}
