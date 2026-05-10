package postgres

import (
	"time"

	"disability_system_backend/internal/modules/auth/ports"
	"disability_system_backend/internal/shared/auth"
)

type TokenService struct {
	jwtService *auth.JWTService
}

func NewTokenService(jwtService *auth.JWTService) *TokenService {
	return &TokenService{jwtService: jwtService}
}

func (s *TokenService) GenerateTokenPair(userID uint64, email, role string) (*ports.TokenPair, error) {
	tokens, err := s.jwtService.GenerateTokenPair(userID, email, role)
	if err != nil {
		return nil, err
	}
	return &ports.TokenPair{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	}, nil
}

func (s *TokenService) ValidateToken(token string) (*ports.TokenClaims, error) {
	claims, err := s.jwtService.ValidateToken(token)
	if err != nil {
		return nil, err
	}
	return &ports.TokenClaims{
		UserID: claims.UserID,
		Email:  claims.Email,
		Role:   claims.Role,
	}, nil
}

func (s *TokenService) RefreshToken(token string) (*ports.TokenPair, error) {
	tokens, err := s.jwtService.RefreshToken(token)
	if err != nil {
		return nil, err
	}
	return &ports.TokenPair{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	}, nil
}

func (s *TokenService) GetExpiration() time.Duration {
	return s.jwtService.Expiration()
}

var _ ports.TokenService = (*TokenService)(nil)