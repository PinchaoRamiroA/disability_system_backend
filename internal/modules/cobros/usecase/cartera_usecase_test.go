package usecase_test

import (
	"context"
	"testing"
	"time"

	"disability_system_backend/internal/modules/cobros/domain"
	"disability_system_backend/internal/modules/cobros/ports"
	"disability_system_backend/internal/modules/cobros/usecase"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockPagoRepo struct {
	mock.Mock
}

func (m *mockPagoRepo) CreatePago(ctx context.Context, pago *domain.Pago) error {
	return m.Called(ctx, pago).Error(0)
}

func (m *mockPagoRepo) FindPagoByID(ctx context.Context, id uint64) (*domain.Pago, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Pago), args.Error(1)
}

func (m *mockPagoRepo) ListPagos(ctx context.Context, filters ports.PagoFilters) ([]domain.Pago, int64, error) {
	args := m.Called(ctx, filters)
	return args.Get(0).([]domain.Pago), args.Get(1).(int64), args.Error(2)
}

func (m *mockPagoRepo) UpdatePago(ctx context.Context, pago *domain.Pago) error {
	return m.Called(ctx, pago).Error(0)
}

func (m *mockPagoRepo) SoftDeletePago(ctx context.Context, id uint64) error {
	return m.Called(ctx, id).Error(0)
}

func (m *mockPagoRepo) IncapacidadExists(ctx context.Context, id uint64) (bool, error) {
	args := m.Called(ctx, id)
	return args.Bool(0), args.Error(1)
}

func (m *mockPagoRepo) EntidadExists(ctx context.Context, id uint64) (bool, error) {
	args := m.Called(ctx, id)
	return args.Bool(0), args.Error(1)
}

func (m *mockPagoRepo) GetEntidadInfo(ctx context.Context) (map[uint64]struct{ Nombre, Tipo string }, error) {
	args := m.Called(ctx)
	return args.Get(0).(map[uint64]struct{ Nombre, Tipo string }), args.Error(1)
}

func (m *mockPagoRepo) GetIncapacidadesDetailed(ctx context.Context, ids []uint64) (map[uint64]ports.IncapacidadInfo, error) {
	args := m.Called(ctx, ids)
	return args.Get(0).(map[uint64]ports.IncapacidadInfo), args.Error(1)
}

type mockSeguimientoRepo struct {
	mock.Mock
}

func (m *mockSeguimientoRepo) CreateSeguimiento(ctx context.Context, seguimiento *domain.SeguimientoCobro) error {
	return m.Called(ctx, seguimiento).Error(0)
}

func (m *mockSeguimientoRepo) FindSeguimientoByID(ctx context.Context, id uint64) (*domain.SeguimientoCobro, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.SeguimientoCobro), args.Error(1)
}

func (m *mockSeguimientoRepo) ListSeguimientos(ctx context.Context, filters ports.SeguimientoFilters) ([]domain.SeguimientoCobro, int64, error) {
	args := m.Called(ctx, filters)
	return args.Get(0).([]domain.SeguimientoCobro), args.Get(1).(int64), args.Error(2)
}

func (m *mockSeguimientoRepo) UpdateSeguimiento(ctx context.Context, seguimiento *domain.SeguimientoCobro) error {
	return m.Called(ctx, seguimiento).Error(0)
}

type combinedCobroRepo struct {
	*mockPagoRepo
	*mockSeguimientoRepo
}

func TestObtenerEstadisticasGenerales_BUG04(t *testing.T) {
	ctx := context.Background()

	t.Run("should calculate IncapacidadesActivas correctly when total == len(pagos)", func(t *testing.T) {
		pagoRepo := new(mockPagoRepo)
		segRepo := new(mockSeguimientoRepo)
		service := usecase.NewCobroWorkflowService(pagoRepo, segRepo)

		// 3 incapacidades:
		// Incapacidad 1: Pendiente ($1000)
		// Incapacidad 2: Pagado ($500) y Pendiente ($500)
		// Incapacidad 3: Pagado ($2000)
		pagos := []domain.Pago{
			{
				IDPago:        1,
				IDIncapacidad: 1,
				IDEntidad:     10,
				EstadoPago:    "Pendiente",
				Valor:         decimal.NewFromInt(1000),
				FechaPago:     time.Now().AddDate(0, 0, 5),
			},
			{
				IDPago:        2,
				IDIncapacidad: 2,
				IDEntidad:     10,
				EstadoPago:    "Pagado",
				Valor:         decimal.NewFromInt(500),
				FechaPago:     time.Now().AddDate(0, 0, -5),
			},
			{
				IDPago:        3,
				IDIncapacidad: 2,
				IDEntidad:     10,
				EstadoPago:    "Pendiente",
				Valor:         decimal.NewFromInt(500),
				FechaPago:     time.Now().AddDate(0, 0, 5),
			},
			{
				IDPago:        4,
				IDIncapacidad: 3,
				IDEntidad:     20,
				EstadoPago:    "Pagado",
				Valor:         decimal.NewFromInt(2000),
				FechaPago:     time.Now().AddDate(0, 0, -10),
			},
		}

		pagoRepo.On("ListPagos", ctx, ports.PagoFilters{Limit: 10000}).Return(pagos, int64(len(pagos)), nil)
		segRepo.On("ListSeguimientos", ctx, ports.SeguimientoFilters{Limit: 1000}).Return([]domain.SeguimientoCobro{}, int64(0), nil)

		stats, err := service.ObtenerEstadisticasGenerales(ctx)

		assert.NoError(t, err)
		assert.NotNil(t, stats)

		// With BUG-04, IncapacidadesActivas evaluated to total - len(pagos) = 4 - 4 = 0!
		// Now, Incapacidad 1 and Incapacidad 2 are active (have pending payments). Incapacidad 3 is fully paid.
		assert.Equal(t, int64(3), stats.TotalIncapacidades, "Debe haber 3 incapacidades únicas en total")
		assert.Equal(t, int64(2), stats.IncapacidadesActivas, "Debe haber 2 incapacidades activas con pagos pendientes")
		assert.Equal(t, int64(2), stats.PagosPendientes)
		assert.Equal(t, "4000", stats.TotalValorCartera)
		assert.Equal(t, "2500", stats.TotalValorCobrado)
		assert.Equal(t, "1500", stats.TotalValorPendiente)

		pagoRepo.AssertExpectations(t)
		segRepo.AssertExpectations(t)
	})

	t.Run("should return 0 IncapacidadesActivas when all payments are Pagado or Conciliado", func(t *testing.T) {
		pagoRepo := new(mockPagoRepo)
		segRepo := new(mockSeguimientoRepo)
		service := usecase.NewCobroWorkflowService(pagoRepo, segRepo)

		pagos := []domain.Pago{
			{
				IDPago:        1,
				IDIncapacidad: 1,
				EstadoPago:    "Pagado",
				Valor:         decimal.NewFromInt(1000),
				FechaPago:     time.Now(),
			},
			{
				IDPago:        2,
				IDIncapacidad: 2,
				EstadoPago:    "Conciliado",
				Valor:         decimal.NewFromInt(1500),
				FechaPago:     time.Now(),
			},
		}

		pagoRepo.On("ListPagos", ctx, ports.PagoFilters{Limit: 10000}).Return(pagos, int64(2), nil)
		segRepo.On("ListSeguimientos", ctx, ports.SeguimientoFilters{Limit: 1000}).Return([]domain.SeguimientoCobro{}, int64(0), nil)

		stats, err := service.ObtenerEstadisticasGenerales(ctx)

		assert.NoError(t, err)
		assert.NotNil(t, stats)
		assert.Equal(t, int64(2), stats.TotalIncapacidades)
		assert.Equal(t, int64(0), stats.IncapacidadesActivas)
		assert.Equal(t, int64(0), stats.PagosPendientes)
	})

	t.Run("should check permissions in CobroUseCase.ObtenerEstadisticasGenerales", func(t *testing.T) {
		pagoRepo := new(mockPagoRepo)
		segRepo := new(mockSeguimientoRepo)
		combined := &combinedCobroRepo{
			mockPagoRepo:        pagoRepo,
			mockSeguimientoRepo: segRepo,
		}

		uc := usecase.NewCobroUseCase(combined)

		// Actor without permission
		actorSinPermiso := ports.Actor{
			UserID:   1,
			Permisos: []string{"otra_cosa"},
		}
		_, err := uc.ObtenerEstadisticasGenerales(ctx, actorSinPermiso)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "no tienes permiso")

		// Actor with permission
		actorConPermiso := ports.Actor{
			UserID:   1,
			Permisos: []string{"consultar_incapacidad"},
		}
		pagoRepo.On("ListPagos", ctx, ports.PagoFilters{Limit: 10000}).Return([]domain.Pago{}, int64(0), nil)
		segRepo.On("ListSeguimientos", ctx, ports.SeguimientoFilters{Limit: 1000}).Return([]domain.SeguimientoCobro{}, int64(0), nil)

		stats, err := uc.ObtenerEstadisticasGenerales(ctx, actorConPermiso)
		assert.NoError(t, err)
		assert.NotNil(t, stats)
		assert.Equal(t, int64(0), stats.TotalIncapacidades)
		assert.Equal(t, int64(0), stats.IncapacidadesActivas)
	})
}

func TestObtenerResumenPorEntidad_BUG05_Precision(t *testing.T) {
	ctx := context.Background()

	t.Run("should preserve precision when accumulating financial decimal values", func(t *testing.T) {
		pagoRepo := new(mockPagoRepo)
		segRepo := new(mockSeguimientoRepo)
		service := usecase.NewCobroWorkflowService(pagoRepo, segRepo)

		// 10 pagos de 0.10 cada uno. En float64 clásico, 0.1 + 0.1 + ... != 1.0 (acumula error 0.9999999999999999)
		pagos := make([]domain.Pago, 10)
		for i := 0; i < 10; i++ {
			pagos[i] = domain.Pago{
				IDPago:        uint64(i + 1),
				IDIncapacidad: uint64(i + 1),
				IDEntidad:     1,
				EstadoPago:    "Pagado",
				Valor:         decimal.NewFromFloat(0.10),
				FechaPago:     time.Now(),
			}
		}

		// Agregar pagos con decimales exactos
		pagoExtra := domain.Pago{
			IDPago:        11,
			IDIncapacidad: 11,
			IDEntidad:     1,
			EstadoPago:    "Pendiente",
			Valor:         decimal.RequireFromString("1234567.89"),
			FechaPago:     time.Now(),
		}
		pagos = append(pagos, pagoExtra)

		pagoRepo.On("ListPagos", ctx, ports.PagoFilters{Limit: 10000}).Return(pagos, int64(len(pagos)), nil)
		pagoRepo.On("GetEntidadInfo", ctx).Return(map[uint64]struct{ Nombre, Tipo string }{
			1: {Nombre: "Sura EPS", Tipo: "EPS"},
		}, nil)

		resumen, err := service.ObtenerResumenPorEntidad(ctx)

		assert.NoError(t, err)
		assert.Len(t, resumen, 1)

		entidadResumen := resumen[0]
		assert.Equal(t, uint64(1), entidadResumen.IDEntidad)
		assert.Equal(t, "Sura EPS", entidadResumen.Nombre)
		assert.Equal(t, int64(11), entidadResumen.CantidadINC)
		// 10 * 0.10 = 1.00 cobrado
		assert.Equal(t, "1", entidadResumen.ValorCobrado)
		// 1234567.89 pendiente
		assert.Equal(t, "1234567.89", entidadResumen.ValorPendiente)
		// 1.00 + 1234567.89 = 1234568.89 total exacto
		assert.Equal(t, "1234568.89", entidadResumen.ValorTotal)

		pagoRepo.AssertExpectations(t)
	})
}
