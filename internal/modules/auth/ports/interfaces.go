package ports

import (
	"context"

	"disability_system_backend/internal/modules/auth/domain"
)

type UserRepository interface {
	FindByID(ctx context.Context, id uint64) (*domain.User, error)
	FindByEmail(ctx context.Context, email string) (*domain.User, error)
	FindByDocumentNumber(ctx context.Context, docNumber string) (*domain.User, error)
	Create(ctx context.Context, user *domain.User) error
	Update(ctx context.Context, user *domain.User) error
	Delete(ctx context.Context, id uint64) error
	EmailExists(ctx context.Context, email string) (bool, error)
	DocumentExists(ctx context.Context, docNumber string) (bool, error)
}

type RoleRepository interface {
	FindByID(ctx context.Context, id uint64) (*domain.Role, error)
	FindByName(ctx context.Context, name string) (*domain.Role, error)
}

type TokenService interface {
	GenerateTokenPair(userID uint64, email, role string) (*TokenPair, error)
	ValidateToken(token string) (*TokenClaims, error)
	RefreshToken(token string) (*TokenPair, error)
}

type PasswordHasher interface {
	Hash(password string) (string, error)
	Check(password, hash string) bool
}

type TokenPair struct {
	AccessToken  string
	RefreshToken string
}

type TokenClaims struct {
	UserID uint64
	Email  string
	Role   string
}