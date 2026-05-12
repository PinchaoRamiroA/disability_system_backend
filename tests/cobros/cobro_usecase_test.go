package cobros_test

import (
	"context"
	"testing"

	"disability_system_backend/internal/modules/cobros/ports"
	"disability_system_backend/internal/modules/cobros/usecase"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockPagoRepository struct {
	mock.Mock
}

func (m *MockPagoRepository) CreatePago(ctx context.Context, pago interface{}) error {
	args := m.Called(ctx, pago)
	return args.Error(0)
}

func (m *MockPagoRepository) FindPagoByID(ctx context.Context, id uint64) (interface{}, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0), args.Error(1)
}

func (m *MockPagoRepository) ListPagos(ctx context.Context, filters ports.PagoFilters) (interface{}, int64, error) {
	args := m.Called(ctx, filters)
	return args.Get(0), args.Get(1).(int64), args.Error(2)
}

func (m *MockPagoRepository) UpdatePago(ctx context.Context, pago interface{}) error {
	args := m.Called(ctx, pago)
	return args.Error(0)
}

func (m *MockPagoRepository) SoftDeletePago(ctx context.Context, id uint64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockPagoRepository) IncapacidadExists(ctx context.Context, id uint64) (bool, error) {
	args := m.Called(ctx, id)
	return args.Bool(0), args.Error(1)
}

func (m *MockPagoRepository) EntidadExists(ctx context.Context, id uint64) (bool, error) {
	args := m.Called(ctx, id)
	return args.Bool(0), args.Error(1)
}

type MockSeguimientoRepository struct {
	mock.Mock
}

func (m *MockSeguimientoRepository) CreateSeguimiento(ctx context.Context, seguimiento interface{}) error {
	args := m.Called(ctx, seguimiento)
	return args.Error(0)
}

func (m *MockSeguimientoRepository) FindSeguimientoByID(ctx context.Context, id uint64) (interface{}, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0), args.Error(1)
}

func (m *MockSeguimientoRepository) ListSeguimientos(ctx context.Context, filters ports.SeguimientoFilters) (interface{}, int64, error) {
	args := m.Called(ctx, filters)
	return args.Get(0), args.Get(1).(int64), args.Error(2)
}

func (m *MockSeguimientoRepository) UpdateSeguimiento(ctx context.Context, seguimiento interface{}) error {
	args := m.Called(ctx, seguimiento)
	return args.Error(0)
}

type CombinedRepo struct {
	MockPagoRepository
	MockSeguimientoRepository
}

func TestValidarTransicion_CobradaAPendientePago(t *testing.T) {
	workflowService := &usecase.CobroWorkflowService{}

	err := workflowService.ValidarTransicionEstado("Cobrada", "Pendiente pago")

	assert.NoError(t, err)
}

func TestValidarTransicion_Rechazada(t *testing.T) {
	workflowService := &usecase.CobroWorkflowService{}

	err := workflowService.ValidarTransicionEstado("Cobrada", "Rechazada")

	assert.NoError(t, err)
}

func TestValidarTransicion_Invalida(t *testing.T) {
	workflowService := &usecase.CobroWorkflowService{}

	err := workflowService.ValidarTransicionEstado("Recibida", "Pagada")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "transición de estado no válida")
}

func TestValidarTransicion_CobroJuridicoAPagada(t *testing.T) {
	workflowService := &usecase.CobroWorkflowService{}

	err := workflowService.ValidarTransicionEstado("Cobro jurídico", "Pagada")

	assert.NoError(t, err)
}

func TestValidarTransicion_PagadaAConciliada(t *testing.T) {
	workflowService := &usecase.CobroWorkflowService{}

	err := workflowService.ValidarTransicionEstado("Pagada", "En conciliación")

	assert.NoError(t, err)
}
