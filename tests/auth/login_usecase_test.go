package auth_test

import (
	"context"
	"testing"
	"time"

	"disability_system_backend/internal/modules/auth/domain"
	"disability_system_backend/internal/modules/auth/ports"
	"disability_system_backend/internal/modules/auth/usecase"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) FindByID(ctx context.Context, id uint64) (*domain.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepository) FindByDocumentNumber(ctx context.Context, docNumber string) (*domain.User, error) {
	args := m.Called(ctx, docNumber)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepository) Create(ctx context.Context, user *domain.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) Update(ctx context.Context, user *domain.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) Delete(ctx context.Context, id uint64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockUserRepository) EmailExists(ctx context.Context, email string) (bool, error) {
	args := m.Called(ctx, email)
	return args.Bool(0), args.Error(1)
}

func (m *MockUserRepository) DocumentExists(ctx context.Context, docNumber string) (bool, error) {
	args := m.Called(ctx, docNumber)
	return args.Bool(0), args.Error(1)
}

type MockRoleRepository struct {
	mock.Mock
}

func (m *MockRoleRepository) FindByID(ctx context.Context, id uint64) (*domain.Role, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Role), args.Error(1)
}

func (m *MockRoleRepository) FindByName(ctx context.Context, name string) (*domain.Role, error) {
	args := m.Called(ctx, name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Role), args.Error(1)
}

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

type MockPasswordHasher struct {
	mock.Mock
}

func (m *MockPasswordHasher) Hash(password string) (string, error) {
	args := m.Called(password)
	return args.String(0), args.Error(1)
}

func (m *MockPasswordHasher) Check(password, hash string) bool {
	args := m.Called(password, hash)
	return args.Bool(0)
}

func TestLoginUseCase_Execute_Success(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockRoleRepo := new(MockRoleRepository)
	mockTokenSvc := new(MockTokenService)
	mockHasher := new(MockPasswordHasher)

	uc := usecase.NewLoginUseCase(mockUserRepo, mockRoleRepo, mockTokenSvc, mockHasher, time.Hour)

	testUser := &domain.User{
		ID:           1,
		IDRol:        1,
		Nombre:       "Test User",
		Correo:       "test@example.com",
		PasswordHash: "hashedpassword",
		Estado:       true,
		IsDeleted:    false,
	}

	testRole := &domain.Role{
		ID:       1,
		Nombre:   "admin",
		Permisos: []string{"admin:read", "admin:write"},
	}

	testTokens := &ports.TokenPair{
		AccessToken:  "access_token_123",
		RefreshToken: "refresh_token_456",
	}

	mockUserRepo.On("FindByEmail", mock.Anything, "test@example.com").Return(testUser, nil)
	mockRoleRepo.On("FindByID", mock.Anything, uint64(1)).Return(testRole, nil)
	mockHasher.On("Check", "password123", "hashedpassword").Return(true)
	mockTokenSvc.On("GenerateTokenPair", uint64(1), "test@example.com", "admin").Return(testTokens, nil)

	tokens, user, role, err := uc.Execute(context.Background(), "test@example.com", "password123")

	assert.NoError(t, err)
	assert.NotNil(t, tokens)
	assert.NotNil(t, user)
	assert.NotNil(t, role)
	assert.Equal(t, "access_token_123", tokens.AccessToken)
	assert.Equal(t, "test@example.com", user.Correo)
	assert.Equal(t, "admin", role.Nombre)
	assert.Equal(t, []string{"admin:read", "admin:write"}, role.Permisos)

	mockUserRepo.AssertExpectations(t)
	mockRoleRepo.AssertExpectations(t)
	mockTokenSvc.AssertExpectations(t)
	mockHasher.AssertExpectations(t)
}

func TestLoginUseCase_Execute_PasswordMismatch(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockRoleRepo := new(MockRoleRepository)
	mockTokenSvc := new(MockTokenService)
	mockHasher := new(MockPasswordHasher)

	uc := usecase.NewLoginUseCase(mockUserRepo, mockRoleRepo, mockTokenSvc, mockHasher, time.Hour)

	testUser := &domain.User{
		ID:           1,
		IDRol:        1,
		Correo:       "test@example.com",
		PasswordHash: "hashedpassword",
		Estado:       true,
		IsDeleted:    false,
	}

	mockUserRepo.On("FindByEmail", mock.Anything, "test@example.com").Return(testUser, nil)
	mockHasher.On("Check", "wrongpassword", "hashedpassword").Return(false)

	tokens, user, role, err := uc.Execute(context.Background(), "test@example.com", "wrongpassword")

	assert.Error(t, err)
	assert.Nil(t, tokens)
	assert.Nil(t, user)
	assert.Nil(t, role)
	assert.Contains(t, err.Error(), "contraseña")
}

func TestLoginUseCase_GetExpirationSeconds(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockRoleRepo := new(MockRoleRepository)
	mockTokenSvc := new(MockTokenService)
	mockHasher := new(MockPasswordHasher)

	uc := usecase.NewLoginUseCase(mockUserRepo, mockRoleRepo, mockTokenSvc, mockHasher, 2*time.Hour)

	expiration := uc.GetExpirationSeconds()

	assert.Equal(t, int64(7200), expiration)
}