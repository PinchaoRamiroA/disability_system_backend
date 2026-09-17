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

func TestObtenerAlertasVencimiento_BUG06(t *testing.T) {
	ctx := context.Background()
	actorAdmin := ports.Actor{
		UserID:   1,
		Permisos: []string{"registrar_pago", "generar_alertas"},
	}

	t.Run("should reach future expiration branch and return correct alert message and days", func(t *testing.T) {
		pagoRepo := new(mockPagoRepo)
		segRepo := new(mockSeguimientoRepo)
		service := usecase.NewCobroWorkflowService(pagoRepo, segRepo)

		now := time.Now()
		pagos := []domain.Pago{
			{
				IDPago:        1,
				IDIncapacidad: 101,
				NombreEntidad: "EPS Sanitas",
				EstadoPago:    "Pendiente",
				FechaPago:     now.AddDate(0, 0, 5), // Future: due in 5 days
			},
			{
				IDPago:        2,
				IDIncapacidad: 102,
				NombreEntidad: "Sura EPS",
				EstadoPago:    "Pendiente",
				FechaPago:     now.AddDate(0, 0, -5), // Past: overdue by 5 days
			},
			{
				IDPago:        3,
				IDIncapacidad: 103,
				NombreEntidad: "Compensar",
				EstadoPago:    "Pendiente",
				FechaPago:     now, // Today
			},
			{
				IDPago:        4,
				IDIncapacidad: 104,
				NombreEntidad: "Nueva EPS",
				EstadoPago:    "Pendiente",
				FechaPago:     now.AddDate(0, 0, 20), // Beyond 10-day window
			},
			{
				IDPago:        5,
				IDIncapacidad: 105,
				NombreEntidad: "Salud Total",
				EstadoPago:    "Pagado", // Already paid -> should be ignored
				FechaPago:     now.AddDate(0, 0, 2),
			},
			{
				IDPago:        6,
				IDIncapacidad: 106,
				NombreEntidad: "Famisanar",
				EstadoPago:    "Pendiente",
				FechaPago:     now.AddDate(0, 0, -20), // Past: overdue by 20 days -> Medio
			},
		}

		pagoRepo.On("ListPagos", ctx, ports.PagoFilters{Limit: 1000}).Return(pagos, int64(len(pagos)), nil)
		pagoRepo.On("GetIncapacidadesDetailed", ctx, mock.Anything).Return(map[uint64]ports.IncapacidadInfo{
			101: {ID: 101, Titulo: "Incapacidad 101"},
			102: {ID: 102, Titulo: "Incapacidad 102"},
			103: {ID: 103, Titulo: "Incapacidad 103"},
			106: {ID: 106, Titulo: "Incapacidad 106"},
		}, nil)

		// Ask with diasMinimos = 10 (window of next 10 days)
		alertas, err := service.ObtenerAlertasVencimiento(ctx, 10, actorAdmin)

		assert.NoError(t, err)
		assert.Len(t, alertas, 4, "Debe incluir las 4 alertas dentro de la ventana de 10 días pendientes")

		// Alerta 1: Vence en 5 días (before BUG-06, this branch was unreachable!)
		assert.Equal(t, uint64(101), alertas[0].IDIncapacidad)
		assert.Equal(t, "Incapacidad 101", alertas[0].Incapacidad.Titulo)
		assert.Equal(t, 5, alertas[0].DiasRestantes)
		assert.Equal(t, "Vence en 5 días", alertas[0].Mensaje)
		assert.Equal(t, "Bajo", alertas[0].TipoAlerta)

		// Alerta 2: Vencida hace 5 días
		assert.Equal(t, uint64(102), alertas[1].IDIncapacidad)
		assert.Equal(t, -5, alertas[1].DiasRestantes)
		assert.Equal(t, "Vencida hace 5 días", alertas[1].Mensaje)
		assert.Equal(t, "Bajo", alertas[1].TipoAlerta)

		// Alerta 3: Vence hoy
		assert.Equal(t, uint64(103), alertas[2].IDIncapacidad)
		assert.Equal(t, 0, alertas[2].DiasRestantes)
		assert.Equal(t, "Vence hoy", alertas[2].Mensaje)
		assert.Equal(t, "Bajo", alertas[2].TipoAlerta)

		// Alerta 4: Vencida hace 20 días
		assert.Equal(t, uint64(106), alertas[3].IDIncapacidad)
		assert.Equal(t, -20, alertas[3].DiasRestantes)
		assert.Equal(t, "Vencida hace 20 días", alertas[3].Mensaje)
		assert.Equal(t, "Medio", alertas[3].TipoAlerta)

		pagoRepo.AssertExpectations(t)
	})

	t.Run("should handle negative diasMinimos gracefully as absolute window", func(t *testing.T) {
		pagoRepo := new(mockPagoRepo)
		segRepo := new(mockSeguimientoRepo)
		service := usecase.NewCobroWorkflowService(pagoRepo, segRepo)

		now := time.Now()
		pagos := []domain.Pago{
			{
				IDPago:        1,
				IDIncapacidad: 101,
				EstadoPago:    "Pendiente",
				FechaPago:     now.AddDate(0, 0, 5),
			},
		}

		pagoRepo.On("ListPagos", ctx, ports.PagoFilters{Limit: 1000}).Return(pagos, int64(1), nil)
		pagoRepo.On("GetIncapacidadesDetailed", ctx, mock.Anything).Return(map[uint64]ports.IncapacidadInfo{}, nil)

		alertas, err := service.ObtenerAlertasVencimiento(ctx, -10, actorAdmin)

		assert.NoError(t, err)
		assert.Len(t, alertas, 1)
		assert.Equal(t, 5, alertas[0].DiasRestantes)
		assert.Equal(t, "Vence en 5 días", alertas[0].Mensaje)

		pagoRepo.AssertExpectations(t)
	})

	t.Run("should enforce permission check in CobroUseCase.ObtenerAlertasVencimiento", func(t *testing.T) {
		pagoRepo := new(mockPagoRepo)
		segRepo := new(mockSeguimientoRepo)
		combined := &combinedCobroRepo{
			mockPagoRepo:        pagoRepo,
			mockSeguimientoRepo: segRepo,
		}
		uc := usecase.NewCobroUseCase(combined)

		actorSinPermiso := ports.Actor{
			UserID:   2,
			Permisos: []string{"otro_permiso"},
		}
		_, err := uc.ObtenerAlertasVencimiento(ctx, actorSinPermiso, 0)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "no tienes permiso")

		actorConPermiso := ports.Actor{
			UserID:   2,
			Permisos: []string{"generar_alertas"},
		}
		pagoRepo.On("ListPagos", ctx, mock.Anything).Return([]domain.Pago{}, int64(0), nil)
		pagoRepo.On("GetIncapacidadesDetailed", ctx, mock.Anything).Return(map[uint64]ports.IncapacidadInfo{}, nil)

		alertas, err := uc.ObtenerAlertasVencimiento(ctx, actorConPermiso, 0)
		assert.NoError(t, err)
		assert.Empty(t, alertas)
	})
}

