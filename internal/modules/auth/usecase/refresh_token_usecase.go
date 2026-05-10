package usecase

import (
	"context"

	"disability_system_backend/internal/modules/auth/ports"
	apperrors "disability_system_backend/internal/shared/errors"
)

type RefreshTokenUseCase struct {
	tokenService ports.TokenService
}

func NewRefreshTokenUseCase(tokenService ports.TokenService) *RefreshTokenUseCase {
	return &RefreshTokenUseCase{tokenService: tokenService}
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