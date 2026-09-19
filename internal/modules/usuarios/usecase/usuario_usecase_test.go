package usecase_test

import (
	"context"
	"errors"
	"testing"

	"disability_system_backend/internal/modules/usuarios/domain"
	"disability_system_backend/internal/modules/usuarios/usecase"
	"disability_system_backend/internal/shared/auth"
	apperrors "disability_system_backend/internal/shared/errors"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockUsuarioRepository struct {
	mock.Mock
}

func (m *MockUsuarioRepository) FindByID(ctx context.Context, id uint64) (*domain.Usuario, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Usuario), args.Error(1)
}

func (m *MockUsuarioRepository) FindByEmail(ctx context.Context, email string) (*domain.Usuario, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Usuario), args.Error(1)
}

func (m *MockUsuarioRepository) FindByDocumentNumber(ctx context.Context, docNumber string) (*domain.Usuario, error) {
	args := m.Called(ctx, docNumber)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Usuario), args.Error(1)
}

func (m *MockUsuarioRepository) FindAll(ctx context.Context, page, limit int, estado *bool, idRol *uint64, search string) ([]domain.Usuario, int64, error) {
	args := m.Called(ctx, page, limit, estado, idRol, search)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]domain.Usuario), args.Get(1).(int64), args.Error(2)
}

func (m *MockUsuarioRepository) Create(ctx context.Context, usuario *domain.Usuario) error {
	args := m.Called(ctx, usuario)
	return args.Error(0)
}

func (m *MockUsuarioRepository) Update(ctx context.Context, usuario *domain.Usuario) error {
	args := m.Called(ctx, usuario)
	return args.Error(0)
}

func (m *MockUsuarioRepository) SoftDelete(ctx context.Context, id uint64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockUsuarioRepository) EmailExists(ctx context.Context, email string, excludeID *uint64) (bool, error) {
	args := m.Called(ctx, email, excludeID)
	return args.Bool(0), args.Error(1)
}

func (m *MockUsuarioRepository) DocumentExists(ctx context.Context, docNumber string, excludeID *uint64) (bool, error) {
	args := m.Called(ctx, docNumber, excludeID)
	return args.Bool(0), args.Error(1)
}

type MockRolRepository struct {
	mock.Mock
}

func (m *MockRolRepository) FindByID(ctx context.Context, id uint64) (*domain.Rol, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Rol), args.Error(1)
}

func (m *MockRolRepository) FindByName(ctx context.Context, name string) (*domain.Rol, error) {
	args := m.Called(ctx, name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Rol), args.Error(1)
}

func (m *MockRolRepository) FindAll(ctx context.Context, page, limit int) ([]domain.Rol, int64, error) {
	args := m.Called(ctx, page, limit)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]domain.Rol), args.Get(1).(int64), args.Error(2)
}

func (m *MockRolRepository) Create(ctx context.Context, rol *domain.Rol) error {
	args := m.Called(ctx, rol)
	return args.Error(0)
}

func (m *MockRolRepository) Update(ctx context.Context, rol *domain.Rol) error {
	args := m.Called(ctx, rol)
	return args.Error(0)
}

func (m *MockRolRepository) Delete(ctx context.Context, id uint64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestUsuarioUseCase_CambiarEstado_BUG11(t *testing.T) {
	ctx := context.Background()

	t.Run("should change estado via repository Update without raw SQL (BUG-11 fix)", func(t *testing.T) {
		mockUserRepo := new(MockUsuarioRepository)
		mockRolRepo := new(MockRolRepository)
		uc := usecase.NewUsuarioUseCase(mockUserRepo, mockRolRepo)

		existingUser := &domain.Usuario{
			ID:     1,
			Nombre: "Juan Perez",
			Estado: true,
		}

		mockUserRepo.On("FindByID", ctx, uint64(1)).Return(existingUser, nil)
		mockUserRepo.On("Update", ctx, mock.MatchedBy(func(u *domain.Usuario) bool {
			return u.ID == 1 && u.Estado == false
		})).Return(nil)

		err := uc.CambiarEstado(ctx, 1, false)

		assert.NoError(t, err)
		assert.False(t, existingUser.Estado)
		mockUserRepo.AssertExpectations(t)
	})

	t.Run("should return error if user not found", func(t *testing.T) {
		mockUserRepo := new(MockUsuarioRepository)
		mockRolRepo := new(MockRolRepository)
		uc := usecase.NewUsuarioUseCase(mockUserRepo, mockRolRepo)

		mockUserRepo.On("FindByID", ctx, uint64(99)).Return(nil, apperrors.ErrUserNotFound)

		err := uc.CambiarEstado(ctx, 99, false)

		assert.Error(t, err)
		assert.Equal(t, apperrors.ErrUserNotFound, err)
		mockUserRepo.AssertExpectations(t)
	})

	t.Run("should return error if repository Update fails", func(t *testing.T) {
		mockUserRepo := new(MockUsuarioRepository)
		mockRolRepo := new(MockRolRepository)
		uc := usecase.NewUsuarioUseCase(mockUserRepo, mockRolRepo)

		existingUser := &domain.Usuario{
			ID:     1,
			Estado: true,
		}

		mockUserRepo.On("FindByID", ctx, uint64(1)).Return(existingUser, nil)
		mockUserRepo.On("Update", ctx, mock.Anything).Return(errors.New("db error"))

		err := uc.CambiarEstado(ctx, 1, false)

		assert.Error(t, err)
		assert.Equal(t, "db error", err.Error())
		mockUserRepo.AssertExpectations(t)
	})
}

func TestUsuarioUseCase_CambiarPassword_BUG11(t *testing.T) {
	ctx := context.Background()

	t.Run("should change password via repository Update without raw SQL (BUG-11 fix)", func(t *testing.T) {
		mockUserRepo := new(MockUsuarioRepository)
		mockRolRepo := new(MockRolRepository)
		uc := usecase.NewUsuarioUseCase(mockUserRepo, mockRolRepo)

		oldHash, err := auth.HashPassword("oldPassword123")
		assert.NoError(t, err)

		existingUser := &domain.Usuario{
			ID:           1,
			Nombre:       "Juan Perez",
			PasswordHash: oldHash,
		}

		mockUserRepo.On("FindByID", ctx, uint64(1)).Return(existingUser, nil)
		mockUserRepo.On("Update", ctx, mock.MatchedBy(func(u *domain.Usuario) bool {
			return u.ID == 1 && auth.CheckPassword("newPassword456", u.PasswordHash)
		})).Return(nil)

		err = uc.CambiarPassword(ctx, 1, "oldPassword123", "newPassword456")

		assert.NoError(t, err)
		assert.True(t, auth.CheckPassword("newPassword456", existingUser.PasswordHash))
		mockUserRepo.AssertExpectations(t)
	})

	t.Run("should return ErrUnauthorized when old password does not match", func(t *testing.T) {
		mockUserRepo := new(MockUsuarioRepository)
		mockRolRepo := new(MockRolRepository)
		uc := usecase.NewUsuarioUseCase(mockUserRepo, mockRolRepo)

		oldHash, err := auth.HashPassword("correctPassword")
		assert.NoError(t, err)

		existingUser := &domain.Usuario{
			ID:           1,
			PasswordHash: oldHash,
		}

		mockUserRepo.On("FindByID", ctx, uint64(1)).Return(existingUser, nil)

		err = uc.CambiarPassword(ctx, 1, "wrongPassword", "newPassword456")

		assert.Error(t, err)
		var appErr *apperrors.AppError
		assert.ErrorAs(t, err, &appErr)
		assert.Equal(t, "UNAUTHORIZED", appErr.Code)
		mockUserRepo.AssertExpectations(t)
	})

	t.Run("should return error if user not found", func(t *testing.T) {
		mockUserRepo := new(MockUsuarioRepository)
		mockRolRepo := new(MockRolRepository)
		uc := usecase.NewUsuarioUseCase(mockUserRepo, mockRolRepo)

		mockUserRepo.On("FindByID", ctx, uint64(99)).Return(nil, apperrors.ErrUserNotFound)

		err := uc.CambiarPassword(ctx, 99, "any", "newPassword")

		assert.Error(t, err)
		mockUserRepo.AssertExpectations(t)
	})
}

func TestUsuarioUseCase_AsignarRol(t *testing.T) {
	ctx := context.Background()

	mockUserRepo := new(MockUsuarioRepository)
	mockRolRepo := new(MockRolRepository)
	uc := usecase.NewUsuarioUseCase(mockUserRepo, mockRolRepo)

	existingUser := &domain.Usuario{
		ID:    1,
		IDRol: 2,
	}

	mockUserRepo.On("FindByID", ctx, uint64(1)).Return(existingUser, nil)
	mockRolRepo.On("FindByID", ctx, uint64(3)).Return(&domain.Rol{ID: 3, Nombre: "Gerencia"}, nil)
	mockUserRepo.On("Update", ctx, mock.MatchedBy(func(u *domain.Usuario) bool {
		return u.ID == 1 && u.IDRol == 3
	})).Return(nil)

	err := uc.AsignarRol(ctx, 1, 3)

	assert.NoError(t, err)
	assert.Equal(t, uint64(3), existingUser.IDRol)
	mockUserRepo.AssertExpectations(t)
	mockRolRepo.AssertExpectations(t)
}

func TestUsuarioUseCase_Eliminar(t *testing.T) {
	ctx := context.Background()

	mockUserRepo := new(MockUsuarioRepository)
	mockRolRepo := new(MockRolRepository)
	uc := usecase.NewUsuarioUseCase(mockUserRepo, mockRolRepo)

	mockUserRepo.On("FindByID", ctx, uint64(1)).Return(&domain.Usuario{ID: 1}, nil)
	mockUserRepo.On("SoftDelete", ctx, uint64(1)).Return(nil)

	err := uc.Eliminar(ctx, 1)

	assert.NoError(t, err)
	mockUserRepo.AssertExpectations(t)
}
