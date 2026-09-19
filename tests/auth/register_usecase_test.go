package auth_test

import (
	"context"
	"testing"

	"disability_system_backend/internal/modules/auth/domain"
	"disability_system_backend/internal/modules/auth/usecase"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRegisterUseCase_Execute_Success(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockHasher := new(MockPasswordHasher)

	uc := usecase.NewRegisterUseCase(mockUserRepo, mockHasher)

	mockUserRepo.On("EmailExists", mock.Anything, "newuser@example.com").Return(false, nil)
	mockUserRepo.On("DocumentExists", mock.Anything, "12345678").Return(false, nil)
	mockHasher.On("Hash", "password123").Return("hashedpassword123", nil)
	mockUserRepo.On("Create", mock.Anything, mock.MatchedBy(func(u *domain.User) bool {
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

	uc := usecase.NewRegisterUseCase(mockUserRepo, mockHasher)

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

	uc := usecase.NewRegisterUseCase(mockUserRepo, mockHasher)

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

	uc := usecase.NewRegisterUseCase(mockUserRepo, mockHasher)

	mockUserRepo.On("EmailExists", mock.Anything, "test@example.com").Return(false, nil)
	mockUserRepo.On("DocumentExists", mock.Anything, "99999999").Return(false, nil)
	mockHasher.On("Hash", "password123").Return("hashed", nil)
	mockUserRepo.On("Create", mock.Anything, mock.MatchedBy(func(u *domain.User) bool {
		return u.IDRol == 4
	})).Run(func(args mock.Arguments) {
		user := args.Get(1).(*domain.User)
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

func TestRegisterUseCase_Execute_AssignsRoleDynamicallyByName_BUG12(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockHasher := new(MockPasswordHasher)
	mockRoleRepo := new(MockRoleRepository)

	uc := usecase.NewRegisterUseCase(mockUserRepo, mockHasher, mockRoleRepo)

	mockUserRepo.On("EmailExists", mock.Anything, "dynamic@example.com").Return(false, nil)
	mockUserRepo.On("DocumentExists", mock.Anything, "88888888").Return(false, nil)
	mockRoleRepo.On("FindByName", mock.Anything, "Empleado").Return(&domain.Role{
		ID:     77,
		Nombre: "Empleado",
	}, nil)
	mockHasher.On("Hash", "password123").Return("hashed", nil)
	mockUserRepo.On("Create", mock.Anything, mock.MatchedBy(func(u *domain.User) bool {
		return u.IDRol == 77
	})).Return(nil)

	user, err := uc.Execute(
		context.Background(),
		"Dynamic User",
		"dynamic@example.com",
		"password123",
		"88888888",
	)

	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, uint64(77), user.IDRol, "IDRol must be dynamically resolved from role repository, not hardcoded to 4")

	mockUserRepo.AssertExpectations(t)
	mockHasher.AssertExpectations(t)
	mockRoleRepo.AssertExpectations(t)
}

func TestRegisterUseCase_Execute_DefaultRoleNotFound_BUG12(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockHasher := new(MockPasswordHasher)
	mockRoleRepo := new(MockRoleRepository)

	uc := usecase.NewRegisterUseCase(mockUserRepo, mockHasher, mockRoleRepo)

	mockUserRepo.On("EmailExists", mock.Anything, "dynamic@example.com").Return(false, nil)
	mockUserRepo.On("DocumentExists", mock.Anything, "88888888").Return(false, nil)
	mockRoleRepo.On("FindByName", mock.Anything, "Empleado").Return(nil, assert.AnError)

	user, err := uc.Execute(
		context.Background(),
		"Dynamic User",
		"dynamic@example.com",
		"password123",
		"88888888",
	)

	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Contains(t, err.Error(), "rol por defecto no encontrado")

	mockUserRepo.AssertExpectations(t)
	mockRoleRepo.AssertExpectations(t)
}