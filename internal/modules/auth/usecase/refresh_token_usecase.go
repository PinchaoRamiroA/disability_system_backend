package usecase

import (
	"context"
	"time"

	"disability_system_backend/internal/modules/auth/ports"
	apperrors "disability_system_backend/internal/shared/errors"
)

type RefreshTokenUseCase struct {
	tokenService ports.TokenService
	tokenExpiry  time.Duration
}

func NewRefreshTokenUseCase(tokenService ports.TokenService, tokenExpiry ...time.Duration) *RefreshTokenUseCase {
	var expiry time.Duration
	if len(tokenExpiry) > 0 {
		expiry = tokenExpiry[0]
	}
	return &RefreshTokenUseCase{
		tokenService: tokenService,
		tokenExpiry:  expiry,
	}
}

func (uc *RefreshTokenUseCase) GetExpirationSeconds() int64 {
	return int64(uc.tokenExpiry.Seconds())
}

func (uc *RefreshTokenUseCase) Execute(ctx context.Context, refreshToken string) (*ports.TokenPair, error) {
	tokens, err := uc.tokenService.RefreshToken(refreshToken)
	if err != nil {
		return nil, apperrors.ErrTokenExpired.WithError(err)
	}
	if tokens == nil {
		return nil, apperrors.ErrTokenInvalid.WithMessage("token inválido")
	}
	return tokens, nil
}