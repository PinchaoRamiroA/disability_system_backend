package usecase

import (
	"context"
	"testing"

	usuariosdomain "disability_system_backend/internal/modules/usuarios/domain"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRegisterUseCase_Execute_Success(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockHasher := new(MockPasswordHasher)

	uc := NewRegisterUseCase(mockUserRepo, mockHasher)

	mockUserRepo.On("EmailExists", mock.Anything, "newuser@example.com").Return(false, nil)
	mockUserRepo.On("DocumentExists", mock.Anything, "12345678").Return(false, nil)
	mockHasher.On("Hash", "password123").Return("hashedpassword123", nil)
	mockUserRepo.On("Create", mock.Anything, mock.MatchedBy(func(u *usuariosdomain.Usuario) bool {
		return u.Correo == "newuser@example.com" && u.Nombre == "New User" && u.IDRol == 4
	})).Return(nil)

	user, err := uc.Execute(
		context.Background(),
		"New User",
		"newuser@example.com",
		"password123",
		"12345678",
	)

	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, "newuser@example.com", user.Correo)
	assert.Equal(t, "New User", user.Nombre)
	assert.Equal(t, uint64(4), user.IDRol)

	mockUserRepo.AssertExpectations(t)
	mockHasher.AssertExpectations(t)
}

func TestRegisterUseCase_Execute_EmailExists(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockHasher := new(MockPasswordHasher)

	uc := NewRegisterUseCase(mockUserRepo, mockHasher)

	mockUserRepo.On("EmailExists", mock.Anything, "existing@example.com").Return(true, nil)

	user, err := uc.Execute(
		context.Background(),
		"New User",
		"existing@example.com",
		"password123",
		"12345678",
	)

	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Contains(t, err.Error(), "email")
}

func TestRegisterUseCase_Execute_DocumentExists(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockHasher := new(MockPasswordHasher)

	uc := NewRegisterUseCase(mockUserRepo, mockHasher)

	mockUserRepo.On("EmailExists", mock.Anything, "new@example.com").Return(false, nil)
	mockUserRepo.On("DocumentExists", mock.Anything, "existingdoc").Return(true, nil)

	user, err := uc.Execute(
		context.Background(),
		"New User",
		"new@example.com",
		"password123",
		"existingdoc",
	)

	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Contains(t, err.Error(), "documento")
}

func TestRegisterUseCase_Execute_AssignsDefaultRole(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockHasher := new(MockPasswordHasher)

	uc := NewRegisterUseCase(mockUserRepo, mockHasher)

	mockUserRepo.On("EmailExists", mock.Anything, "test@example.com").Return(false, nil)
	mockUserRepo.On("DocumentExists", mock.Anything, "99999999").Return(false, nil)
	mockHasher.On("Hash", "password123").Return("hashed", nil)
	mockUserRepo.On("Create", mock.Anything, mock.MatchedBy(func(u *usuariosdomain.Usuario) bool {
		return u.IDRol == 4
	})).Run(func(args mock.Arguments) {
		user := args.Get(1).(*usuariosdomain.Usuario)
		user.ID = 1
	}).Return(nil)

	user, err := uc.Execute(
		context.Background(),
		"Test User",
		"test@example.com",
		"password123",
		"99999999",
	)

	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, uint64(4), user.IDRol)
}