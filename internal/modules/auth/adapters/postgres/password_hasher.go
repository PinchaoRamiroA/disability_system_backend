package postgres

import (
	"disability_system_backend/internal/modules/auth/ports"
	"disability_system_backend/internal/shared/auth"
)

type PasswordHasher struct{}

func NewPasswordHasher() *PasswordHasher {
	return &PasswordHasher{}
}

func (p *PasswordHasher) Hash(password string) (string, error) {
	return auth.HashPassword(password)
}

func (p *PasswordHasher) Check(password, hash string) bool {
	return auth.CheckPassword(password, hash)
}

var _ ports.PasswordHasher = (*PasswordHasher)(nil)