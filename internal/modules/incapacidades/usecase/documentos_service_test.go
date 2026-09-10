package usecase_test

import (
	"context"
	"testing"

	"disability_system_backend/internal/modules/incapacidades/domain"
	"disability_system_backend/internal/modules/incapacidades/ports"
	"disability_system_backend/internal/modules/incapacidades/usecase"
	apperrors "disability_system_backend/internal/shared/errors"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockIncapacidadRepo struct {
	mock.Mock
}

func (m *mockIncapacidadRepo) Create(ctx context.Context, incapacidad *domain.Incapacidad) error {
	return m.Called(ctx, incapacidad).Error(0)
}

func (m *mockIncapacidadRepo) FindByID(ctx context.Context, id uint64) (*domain.Incapacidad, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Incapacidad), args.Error(1)
}

func (m *mockIncapacidadRepo) List(ctx context.Context, filters ports.IncapacidadFilters) ([]domain.Incapacidad, int64, error) {
	args := m.Called(ctx, filters)
	return args.Get(0).([]domain.Incapacidad), args.Get(1).(int64), args.Error(2)
}

func (m *mockIncapacidadRepo) Update(ctx context.Context, incapacidad *domain.Incapacidad) error {
	return m.Called(ctx, incapacidad).Error(0)
}

func (m *mockIncapacidadRepo) SoftDelete(ctx context.Context, id uint64) error {
	return m.Called(ctx, id).Error(0)
}

func (m *mockIncapacidadRepo) ExistsUsuario(ctx context.Context, id uint64) (bool, error) {
	args := m.Called(ctx, id)
	return args.Bool(0), args.Error(1)
}

func (m *mockIncapacidadRepo) FindEstadoByID(ctx context.Context, id uint64) (*domain.EstadoIncapacidad, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.EstadoIncapacidad), args.Error(1)
}

func (m *mockIncapacidadRepo) FindEstadoByName(ctx context.Context, name string) (*domain.EstadoIncapacidad, error) {
	args := m.Called(ctx, name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.EstadoIncapacidad), args.Error(1)
}

func (m *mockIncapacidadRepo) FindTipoByID(ctx context.Context, id uint64) (*domain.TipoIncapacidad, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.TipoIncapacidad), args.Error(1)
}

func (m *mockIncapacidadRepo) FindEntidadByID(ctx context.Context, id uint64) (*domain.Entidad, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Entidad), args.Error(1)
}

func (m *mockIncapacidadRepo) ListEstados(ctx context.Context) ([]domain.EstadoIncapacidad, error) {
	args := m.Called(ctx)
	return args.Get(0).([]domain.EstadoIncapacidad), args.Error(1)
}

func (m *mockIncapacidadRepo) ListTipos(ctx context.Context) ([]domain.TipoIncapacidad, error) {
	args := m.Called(ctx)
	return args.Get(0).([]domain.TipoIncapacidad), args.Error(1)
}

func (m *mockIncapacidadRepo) ListEntidades(ctx context.Context) ([]domain.Entidad, error) {
	args := m.Called(ctx)
	return args.Get(0).([]domain.Entidad), args.Error(1)
}

func (m *mockIncapacidadRepo) ListEstadosDocumento(ctx context.Context) ([]domain.EstadoDocumento, error) {
	args := m.Called(ctx)
	return args.Get(0).([]domain.EstadoDocumento), args.Error(1)
}

func (m *mockIncapacidadRepo) ListTiposDocumento(ctx context.Context) ([]domain.TipoDocumento, error) {
	args := m.Called(ctx)
	return args.Get(0).([]domain.TipoDocumento), args.Error(1)
}

func (m *mockIncapacidadRepo) ListTiposPago(ctx context.Context) ([]domain.TipoPago, error) {
	args := m.Called(ctx)
	return args.Get(0).([]domain.TipoPago), args.Error(1)
}

func (m *mockIncapacidadRepo) FindTiposDocumentoByCodigo(ctx context.Context, codigos []string) ([]domain.TipoDocumento, error) {
	args := m.Called(ctx, codigos)
	return args.Get(0).([]domain.TipoDocumento), args.Error(1)
}

func TestVerificarEstadoTransicion_BUG01_NilEstado(t *testing.T) {
	ctx := context.Background()

	t.Run("should not panic when incapacidad is nil", func(t *testing.T) {
		repo := new(mockIncapacidadRepo)
		service := usecase.NewIncapacidadDocumentosService(repo)

		nuevoEstado := &domain.EstadoIncapacidad{
			IDEstado:          2,
			Nombre:            "En validación documental",
			PermiteTransicion: true,
		}

		assert.NotPanics(t, func() {
			err := service.VerificarEstadoTransicion(ctx, nil, nuevoEstado)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), "la incapacidad es requerida")
		})
	})

	t.Run("should not panic when nuevoEstado is nil", func(t *testing.T) {
		repo := new(mockIncapacidadRepo)
		service := usecase.NewIncapacidadDocumentosService(repo)

		incapacidad := &domain.Incapacidad{
			IDIncapacidad: 1,
			IDEstado:      1,
			Estado: &domain.EstadoIncapacidad{
				IDEstado:          1,
				Nombre:            "Recibida",
				PermiteTransicion: true,
			},
		}

		assert.NotPanics(t, func() {
			err := service.VerificarEstadoTransicion(ctx, incapacidad, nil)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), "el nuevo estado es requerido")
		})
	})

	t.Run("should load estado from repo when incapacidad.Estado is nil and validate successfully", func(t *testing.T) {
		repo := new(mockIncapacidadRepo)
		service := usecase.NewIncapacidadDocumentosService(repo)

		incapacidad := &domain.Incapacidad{
			IDIncapacidad: 1,
			IDEstado:      1,
			Estado:        nil, // simulates non-preloaded state
		}

		estadoRecibida := &domain.EstadoIncapacidad{
			IDEstado:          1,
			Nombre:            "Recibida",
			PermiteTransicion: true,
		}

		nuevoEstado := &domain.EstadoIncapacidad{
			IDEstado:          2,
			Nombre:            "En validación documental",
			PermiteTransicion: true,
		}

		repo.On("FindEstadoByID", ctx, uint64(1)).Return(estadoRecibida, nil)

		assert.NotPanics(t, func() {
			err := service.VerificarEstadoTransicion(ctx, incapacidad, nuevoEstado)
			assert.NoError(t, err)
			assert.Equal(t, estadoRecibida, incapacidad.Estado)
		})

		repo.AssertExpectations(t)
	})

	t.Run("should not panic and return conflict error when incapacidad.Estado is nil and IDEstado is 0", func(t *testing.T) {
		repo := new(mockIncapacidadRepo)
		service := usecase.NewIncapacidadDocumentosService(repo)

		incapacidad := &domain.Incapacidad{
			IDIncapacidad: 1,
			IDEstado:      0,
			Estado:        nil,
		}

		nuevoEstado := &domain.EstadoIncapacidad{
			IDEstado:          2,
			Nombre:            "En validación documental",
			PermiteTransicion: true,
		}

		assert.NotPanics(t, func() {
			err := service.VerificarEstadoTransicion(ctx, incapacidad, nuevoEstado)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), "la incapacidad no tiene un estado actual asignado")
		})
	})

	t.Run("should reject transition if current state does not allow transitions", func(t *testing.T) {
		repo := new(mockIncapacidadRepo)
		service := usecase.NewIncapacidadDocumentosService(repo)

		incapacidad := &domain.Incapacidad{
			IDIncapacidad: 1,
			IDEstado:      1,
			Estado: &domain.EstadoIncapacidad{
				IDEstado:          1,
				Nombre:            "Cerrada",
				PermiteTransicion: false,
			},
		}

		nuevoEstado := &domain.EstadoIncapacidad{
			IDEstado:          2,
			Nombre:            "Recibida",
			PermiteTransicion: true,
		}

		err := service.VerificarEstadoTransicion(ctx, incapacidad, nuevoEstado)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "el estado actual no permite transiciones")
	})

	t.Run("should allow transition to universal states (Rechazada, Archivada, Cerrada)", func(t *testing.T) {
		repo := new(mockIncapacidadRepo)
		service := usecase.NewIncapacidadDocumentosService(repo)

		incapacidad := &domain.Incapacidad{
			IDIncapacidad: 1,
			IDEstado:      1,
			Estado: &domain.EstadoIncapacidad{
				IDEstado:          1,
				Nombre:            "Recibida",
				PermiteTransicion: true,
			},
		}

		for _, nombre := range []string{"Rechazada", "Archivada", "Cerrada"} {
			nuevoEstado := &domain.EstadoIncapacidad{
				Nombre: nombre,
			}
			err := service.VerificarEstadoTransicion(ctx, incapacidad, nuevoEstado)
			assert.NoError(t, err)
		}
	})

	t.Run("should reject invalid transition", func(t *testing.T) {
		repo := new(mockIncapacidadRepo)
		service := usecase.NewIncapacidadDocumentosService(repo)

		incapacidad := &domain.Incapacidad{
			IDIncapacidad: 1,
			IDEstado:      1,
			Estado: &domain.EstadoIncapacidad{
				IDEstado:          1,
				Nombre:            "Recibida",
				PermiteTransicion: true,
			},
		}

		nuevoEstado := &domain.EstadoIncapacidad{
			IDEstado: 10,
			Nombre:   "Pagada",
		}

		err := service.VerificarEstadoTransicion(ctx, incapacidad, nuevoEstado)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "transición de estado no válida: de Recibida a Pagada")
	})

	t.Run("should return repository error if FindEstadoByID fails", func(t *testing.T) {
		repo := new(mockIncapacidadRepo)
		service := usecase.NewIncapacidadDocumentosService(repo)

		incapacidad := &domain.Incapacidad{
			IDIncapacidad: 1,
			IDEstado:      999,
			Estado:        nil,
		}

		nuevoEstado := &domain.EstadoIncapacidad{
			IDEstado: 2,
			Nombre:   "En validación documental",
		}

		repo.On("FindEstadoByID", ctx, uint64(999)).Return(nil, apperrors.ErrNotFound)

		err := service.VerificarEstadoTransicion(ctx, incapacidad, nuevoEstado)
		assert.Error(t, err)
		assert.Equal(t, apperrors.ErrNotFound, err)
		repo.AssertExpectations(t)
	})
}
