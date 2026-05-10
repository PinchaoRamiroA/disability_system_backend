package ports

import (
	"context"

	usuariosdomain "disability_system_backend/internal/modules/usuarios/domain"
)

type UserRepository interface {
	FindByID(ctx context.Context, id uint64) (*usuariosdomain.Usuario, error)
	FindByEmail(ctx context.Context, email string) (*usuariosdomain.Usuario, error)
	FindByDocumentNumber(ctx context.Context, docNumber string) (*usuariosdomain.Usuario, error)
	Create(ctx context.Context, user *usuariosdomain.Usuario) error
	Update(ctx context.Context, user *usuariosdomain.Usuario) error
	Delete(ctx context.Context, id uint64) error
	EmailExists(ctx context.Context, email string) (bool, error)
	DocumentExists(ctx context.Context, docNumber string) (bool, error)
}

type RoleRepository interface {
	FindByID(ctx context.Context, id uint64) (*usuariosdomain.Rol, error)
	FindByName(ctx context.Context, name string) (*usuariosdomain.Rol, error)
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