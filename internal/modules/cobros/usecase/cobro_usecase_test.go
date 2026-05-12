package usecase

import (
	"context"
	"testing"
	"time"

	"disability_system_backend/internal/modules/cobros/domain"
	"disability_system_backend/internal/modules/cobros/ports"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
)

func TestCrearPagoRequiresRegistrarPago(t *testing.T) {
	repo := newFakeCobroRepository()
	uc := NewCobroUseCase(repo)
	actor := ports.Actor{UserID: 10, Permisos: []string{"consultar_incapacidad"}}

	_, err := uc.CrearPago(context.Background(), actor, CrearPagoInput{
		IDIncapacidad: 1,
		IDEntidad:     1,
		TipoPago:      "Pago total",
		Valor:         "1000",
		FechaPago:     "2026-05-12",
	})

	require.Error(t, err)
	require.Contains(t, err.Error(), "no tienes permiso para registrar pagos")
}

func TestConciliarPagoRequiresPermission(t *testing.T) {
	repo := newFakeCobroRepository()
	uc := NewCobroUseCase(repo)
	actor := ports.Actor{UserID: 10, Permisos: []string{"registrar_pago"}}

	_, err := uc.ConciliarPago(context.Background(), actor, 1, true, nil, nil)

	require.Error(t, err)
	require.Contains(t, err.Error(), "no tienes permiso para conciliar pagos")
}

func TestConciliarPagoSetsEstadoConciliado(t *testing.T) {
	repo := newFakeCobroRepository()
	uc := NewCobroUseCase(repo)
	actor := ports.Actor{UserID: 10, Permisos: []string{"realizar_conciliacion"}}

	pago, err := uc.ConciliarPago(context.Background(), actor, 1, true, nil, nil)

	require.NoError(t, err)
	require.True(t, pago.Conciliado)
	require.Equal(t, "Conciliado", pago.EstadoPago)
}

func TestCrearSeguimientoJuridicoRequiresJuridicoPermission(t *testing.T) {
	repo := newFakeCobroRepository()
	uc := NewCobroUseCase(repo)
	actor := ports.Actor{UserID: 10, Permisos: []string{"gestionar_cobro_persuasivo"}}

	_, err := uc.CrearSeguimiento(context.Background(), actor, CrearSeguimientoInput{
		IDIncapacidad:   1,
		TipoSeguimiento: "Jurídico",
	})

	require.Error(t, err)
	require.Contains(t, err.Error(), "cobro jurídico")
}

type fakeCobroRepository struct {
	pagos        map[uint64]*domain.Pago
	seguimientos map[uint64]*domain.SeguimientoCobro
}

func newFakeCobroRepository() *fakeCobroRepository {
	now := time.Now()
	return &fakeCobroRepository{
		pagos: map[uint64]*domain.Pago{
			1: {
				IDPago:        1,
				IDIncapacidad: 1,
				IDEntidad:     1,
				TipoPago:      "Pago total",
				EstadoPago:    "Pendiente",
				Valor:         decimal.NewFromInt(1000),
				FechaPago:     now,
			},
		},
		seguimientos: map[uint64]*domain.SeguimientoCobro{},
	}
}

func (r *fakeCobroRepository) CreatePago(ctx context.Context, pago *domain.Pago) error {
	pago.IDPago = 2
	r.pagos[2] = pago
	return nil
}

func (r *fakeCobroRepository) FindPagoByID(ctx context.Context, id uint64) (*domain.Pago, error) {
	return r.pagos[id], nil
}

func (r *fakeCobroRepository) ListPagos(ctx context.Context, filters ports.PagoFilters) ([]domain.Pago, int64, error) {
	items := make([]domain.Pago, 0, len(r.pagos))
	for _, item := range r.pagos {
		items = append(items, *item)
	}
	return items, int64(len(items)), nil
}

func (r *fakeCobroRepository) UpdatePago(ctx context.Context, pago *domain.Pago) error {
	r.pagos[pago.IDPago] = pago
	return nil
}

func (r *fakeCobroRepository) SoftDeletePago(ctx context.Context, id uint64) error {
	delete(r.pagos, id)
	return nil
}

func (r *fakeCobroRepository) IncapacidadExists(ctx context.Context, id uint64) (bool, error) {
	return true, nil
}

func (r *fakeCobroRepository) EntidadExists(ctx context.Context, id uint64) (bool, error) {
	return true, nil
}

func (r *fakeCobroRepository) CreateSeguimiento(ctx context.Context, seguimiento *domain.SeguimientoCobro) error {
	seguimiento.IDSeguimiento = 1
	r.seguimientos[1] = seguimiento
	return nil
}

func (r *fakeCobroRepository) FindSeguimientoByID(ctx context.Context, id uint64) (*domain.SeguimientoCobro, error) {
	return r.seguimientos[id], nil
}

func (r *fakeCobroRepository) ListSeguimientos(ctx context.Context, filters ports.SeguimientoFilters) ([]domain.SeguimientoCobro, int64, error) {
	items := make([]domain.SeguimientoCobro, 0, len(r.seguimientos))
	for _, item := range r.seguimientos {
		items = append(items, *item)
	}
	return items, int64(len(items)), nil
}

func (r *fakeCobroRepository) UpdateSeguimiento(ctx context.Context, seguimiento *domain.SeguimientoCobro) error {
	r.seguimientos[seguimiento.IDSeguimiento] = seguimiento
	return nil
}
