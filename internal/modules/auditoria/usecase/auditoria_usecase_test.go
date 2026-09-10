package usecase_test

import (
	"context"
	"testing"

	"disability_system_backend/internal/modules/auditoria/domain"
	"disability_system_backend/internal/modules/auditoria/dto"
	"disability_system_backend/internal/modules/auditoria/ports"
	"disability_system_backend/internal/modules/auditoria/usecase"
	apperrors "disability_system_backend/internal/shared/errors"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockAuditoriaRepo struct {
	mock.Mock
}

func (m *mockAuditoriaRepo) Create(ctx context.Context, auditoria *domain.Auditoria) error {
	args := m.Called(ctx, auditoria)
	if args.Error(0) == nil {
		auditoria.ID = 1
	}
	return args.Error(0)
}

func (m *mockAuditoriaRepo) List(ctx context.Context, filters dto.ListarAuditoriaQuery) ([]domain.Auditoria, int64, error) {
	args := m.Called(ctx, filters)
	return args.Get(0).([]domain.Auditoria), args.Get(1).(int64), args.Error(2)
}

// customActorImpl is a concrete type implementing ports.Actor to test typed nil pointers
type customActorImpl struct {
	userID   uint64
	permisos []string
}

func (a *customActorImpl) GetUserID() uint64 {
	return a.userID
}

func (a *customActorImpl) HasPermission(p string) bool {
	for _, perm := range a.permisos {
		if perm == p {
			return true
		}
	}
	return false
}

func TestAuditoriaUseCase_BUG02_NilActor(t *testing.T) {
	ctx := context.Background()

	t.Run("should not panic when actor is typed nil (*customActorImpl)(nil) in Crear", func(t *testing.T) {
		repo := new(mockAuditoriaRepo)
		uc := usecase.NewAuditoriaUseCase(repo)

		var nilConcreteActor *customActorImpl = nil
		var actor ports.Actor = nilConcreteActor

		req := dto.CrearAuditoriaRequest{
			TipoAccion:  "CREAR",
			Modulo:      "INCAPACIDADES",
			Descripcion: "Creación de prueba con actor typed nil",
		}

		repo.On("Create", ctx, mock.MatchedBy(func(a *domain.Auditoria) bool {
			return a.IDUsuario == nil && a.TipoAccion == "CREAR"
		})).Return(nil)

		assert.NotPanics(t, func() {
			auditoria, err := uc.Crear(ctx, actor, req)
			assert.NoError(t, err)
			assert.NotNil(t, auditoria)
			assert.Nil(t, auditoria.IDUsuario)
		})

		repo.AssertExpectations(t)
	})

	t.Run("should not panic when actor is untyped nil in Crear", func(t *testing.T) {
		repo := new(mockAuditoriaRepo)
		uc := usecase.NewAuditoriaUseCase(repo)

		req := dto.CrearAuditoriaRequest{
			TipoAccion:  "CREAR",
			Modulo:      "INCAPACIDADES",
			Descripcion: "Creación sin actor",
		}

		repo.On("Create", ctx, mock.MatchedBy(func(a *domain.Auditoria) bool {
			return a.IDUsuario == nil
		})).Return(nil)

		assert.NotPanics(t, func() {
			auditoria, err := uc.Crear(ctx, nil, req)
			assert.NoError(t, err)
			assert.NotNil(t, auditoria)
			assert.Nil(t, auditoria.IDUsuario)
		})

		repo.AssertExpectations(t)
	})

	t.Run("should extract userID from valid actor in Crear", func(t *testing.T) {
		repo := new(mockAuditoriaRepo)
		uc := usecase.NewAuditoriaUseCase(repo)

		validActor := &customActorImpl{
			userID: 42,
		}

		req := dto.CrearAuditoriaRequest{
			TipoAccion:  "CREAR",
			Modulo:      "INCAPACIDADES",
			Descripcion: "Creación con actor válido",
		}

		repo.On("Create", ctx, mock.MatchedBy(func(a *domain.Auditoria) bool {
			return a.IDUsuario != nil && *a.IDUsuario == 42
		})).Return(nil)

		auditoria, err := uc.Crear(ctx, validActor, req)
		assert.NoError(t, err)
		assert.NotNil(t, auditoria)
		assert.NotNil(t, auditoria.IDUsuario)
		assert.Equal(t, uint64(42), *auditoria.IDUsuario)

		repo.AssertExpectations(t)
	})

	t.Run("should respect req.IDUsuario when explicitly provided even with actor", func(t *testing.T) {
		repo := new(mockAuditoriaRepo)
		uc := usecase.NewAuditoriaUseCase(repo)

		validActor := &customActorImpl{
			userID: 42,
		}

		customUserID := uint64(99)
		req := dto.CrearAuditoriaRequest{
			IDUsuario:   &customUserID,
			TipoAccion:  "CREAR",
			Modulo:      "INCAPACIDADES",
			Descripcion: "Creación con id_usuario explícito",
		}

		repo.On("Create", ctx, mock.MatchedBy(func(a *domain.Auditoria) bool {
			return a.IDUsuario != nil && *a.IDUsuario == 99
		})).Return(nil)

		auditoria, err := uc.Crear(ctx, validActor, req)
		assert.NoError(t, err)
		assert.NotNil(t, auditoria)
		assert.Equal(t, uint64(99), *auditoria.IDUsuario)

		repo.AssertExpectations(t)
	})

	t.Run("should not panic when actor is typed nil in Listar and reject with ErrForbidden", func(t *testing.T) {
		repo := new(mockAuditoriaRepo)
		uc := usecase.NewAuditoriaUseCase(repo)

		var nilConcreteActor *customActorImpl = nil
		var actor ports.Actor = nilConcreteActor

		query := dto.ListarAuditoriaQuery{
			Page:  1,
			Limit: 20,
		}

		assert.NotPanics(t, func() {
			items, total, err := uc.Listar(ctx, actor, query)
			assert.Error(t, err)
			assert.Equal(t, apperrors.ErrForbidden.Code, err.(*apperrors.AppError).Code)
			assert.Nil(t, items)
			assert.Equal(t, int64(0), total)
		})
	})

	t.Run("should allow Listar when valid actor has permission", func(t *testing.T) {
		repo := new(mockAuditoriaRepo)
		uc := usecase.NewAuditoriaUseCase(repo)

		actor := &customActorImpl{
			userID:   1,
			permisos: []string{"consultar_auditoria"},
		}

		query := dto.ListarAuditoriaQuery{
			Page:  1,
			Limit: 20,
		}

		repo.On("List", ctx, query).Return([]domain.Auditoria{
			{ID: 1, TipoAccion: "CREAR"},
		}, int64(1), nil)

		items, total, err := uc.Listar(ctx, actor, query)
		assert.NoError(t, err)
		assert.Equal(t, int64(1), total)
		assert.Len(t, items, 1)

		repo.AssertExpectations(t)
	})

	t.Run("should reject Listar when actor lacks permission", func(t *testing.T) {
		repo := new(mockAuditoriaRepo)
		uc := usecase.NewAuditoriaUseCase(repo)

		actor := &customActorImpl{
			userID:   1,
			permisos: []string{"otra_cosa"},
		}

		query := dto.ListarAuditoriaQuery{
			Page:  1,
			Limit: 20,
		}

		items, total, err := uc.Listar(ctx, actor, query)
		assert.Error(t, err)
		assert.Equal(t, apperrors.ErrForbidden.Code, err.(*apperrors.AppError).Code)
		assert.Nil(t, items)
		assert.Equal(t, int64(0), total)
	})
}
