package incapacidades_test

import (
	"context"
	"testing"
	"time"

	"disability_system_backend/internal/modules/incapacidades/domain"
	"disability_system_backend/internal/modules/incapacidades/ports"
	"disability_system_backend/internal/modules/incapacidades/usecase"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockIncapacidadRepository struct {
	mock.Mock
}

func (m *MockIncapacidadRepository) Create(ctx context.Context, incapacidad *domain.Incapacidad) error {
	args := m.Called(ctx, incapacidad)
	if args.Error(0) == nil {
		incapacidad.IDIncapacidad = 1
	}
	return args.Error(0)
}

func (m *MockIncapacidadRepository) FindByID(ctx context.Context, id uint64) (*domain.Incapacidad, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Incapacidad), args.Error(1)
}

func (m *MockIncapacidadRepository) List(ctx context.Context, filters ports.IncapacidadFilters) ([]domain.Incapacidad, int64, error) {
	args := m.Called(ctx, filters)
	return args.Get(0).([]domain.Incapacidad), args.Get(1).(int64), args.Error(2)
}

func (m *MockIncapacidadRepository) Update(ctx context.Context, incapacidad *domain.Incapacidad) error {
	args := m.Called(ctx, incapacidad)
	return args.Error(0)
}

func (m *MockIncapacidadRepository) SoftDelete(ctx context.Context, id uint64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockIncapacidadRepository) ExistsUsuario(ctx context.Context, id uint64) (bool, error) {
	args := m.Called(ctx, id)
	return args.Bool(0), args.Error(1)
}

func (m *MockIncapacidadRepository) FindEstadoByID(ctx context.Context, id uint64) (*domain.EstadoIncapacidad, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.EstadoIncapacidad), args.Error(1)
}

func (m *MockIncapacidadRepository) FindEstadoByName(ctx context.Context, name string) (*domain.EstadoIncapacidad, error) {
	args := m.Called(ctx, name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.EstadoIncapacidad), args.Error(1)
}

func (m *MockIncapacidadRepository) FindTipoByID(ctx context.Context, id uint64) (*domain.TipoIncapacidad, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.TipoIncapacidad), args.Error(1)
}

func (m *MockIncapacidadRepository) FindEntidadByID(ctx context.Context, id uint64) (*domain.Entidad, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Entidad), args.Error(1)
}

func (m *MockIncapacidadRepository) ListEstados(ctx context.Context) ([]domain.EstadoIncapacidad, error) {
	args := m.Called(ctx)
	return args.Get(0).([]domain.EstadoIncapacidad), args.Error(1)
}

func (m *MockIncapacidadRepository) ListTipos(ctx context.Context) ([]domain.TipoIncapacidad, error) {
	args := m.Called(ctx)
	return args.Get(0).([]domain.TipoIncapacidad), args.Error(1)
}

func (m *MockIncapacidadRepository) ListEntidades(ctx context.Context) ([]domain.Entidad, error) {
	args := m.Called(ctx)
	return args.Get(0).([]domain.Entidad), args.Error(1)
}

func (m *MockIncapacidadRepository) ListEstadosDocumento(ctx context.Context) ([]domain.EstadoDocumento, error) {
	args := m.Called(ctx)
	return args.Get(0).([]domain.EstadoDocumento), args.Error(1)
}

func (m *MockIncapacidadRepository) ListTiposDocumento(ctx context.Context) ([]domain.TipoDocumento, error) {
	args := m.Called(ctx)
	return args.Get(0).([]domain.TipoDocumento), args.Error(1)
}

func (m *MockIncapacidadRepository) ListTiposPago(ctx context.Context) ([]domain.TipoPago, error) {
	args := m.Called(ctx)
	return args.Get(0).([]domain.TipoPago), args.Error(1)
}

func (m *MockIncapacidadRepository) FindTiposDocumentoByNombre(ctx context.Context, nombres []string) ([]domain.TipoDocumento, error) {
	args := m.Called(ctx, nombres)
	return args.Get(0).([]domain.TipoDocumento), args.Error(1)
}

type MockDocumentoFaltanteNotifier struct {
	mock.Mock
}

func (m *MockDocumentoFaltanteNotifier) NotificarDocumentosFaltantes(ctx context.Context, userID, incapacidadID uint64, documentos []domain.TipoDocumento) error {
	args := m.Called(ctx, userID, incapacidadID, documentos)
	return args.Error(0)
}

func TestCrearIncapacidad_Success(t *testing.T) {
	mockRepo := new(MockIncapacidadRepository)
	mockNotifier := new(MockDocumentoFaltanteNotifier)

	uc := usecase.NewIncapacidadUseCase(mockRepo)
	uc.SetDocumentoFaltanteNotifier(mockNotifier)

	actor := ports.Actor{UserID: 1, Role: "admin", Permisos: []string{"crear_incapacidad"}}

	estadoRecibida := &domain.EstadoIncapacidad{IDEstado: 1, Nombre: "Recibida"}
	tipo := &domain.TipoIncapacidad{IDTipo: 1, Nombre: "Enfermedad General"}
	entidad := &domain.Entidad{IDEntidad: 1, Nombre: "EPS Test"}

	mockRepo.On("FindEstadoByName", mock.Anything, "Recibida").Return(estadoRecibida, nil)
	mockRepo.On("ExistsUsuario", mock.Anything, uint64(1)).Return(true, nil)
	mockRepo.On("FindTipoByID", mock.Anything, uint64(1)).Return(tipo, nil)
	mockRepo.On("FindEntidadByID", mock.Anything, uint64(1)).Return(entidad, nil)
	mockRepo.On("Create", mock.Anything, mock.Anything).Return(nil)

	input := usecase.CrearIncapacidadInput{
		IDUsuario:      1,
		IDTipo:         1,
		IDEntidad:      1,
		CanalRecepcion: "correo",
		Titulo:         "Test Incapacidad",
		FechaInicio:    "2024-01-15",
		Origen:         "Enfermedad General",
	}

	incapacidad, err := uc.Crear(context.Background(), actor, input)

	assert.NoError(t, err)
	assert.NotNil(t, incapacidad)
	assert.Equal(t, "Test Incapacidad", incapacidad.Titulo)
	assert.Equal(t, uint64(1), incapacidad.IDEstado)

	mockRepo.AssertExpectations(t)
}

func TestCrearIncapacidad_SinPermiso(t *testing.T) {
	mockRepo := new(MockIncapacidadRepository)

	uc := usecase.NewIncapacidadUseCase(mockRepo)

	actor := ports.Actor{UserID: 1, Role: "viewer", Permisos: []string{}}

	input := usecase.CrearIncapacidadInput{
		IDUsuario:   1,
		IDTipo:      1,
		IDEntidad:   1,
		Titulo:      "Test",
		FechaInicio: "2024-01-15",
		Origen:     "EG",
	}

	incapacidad, err := uc.Crear(context.Background(), actor, input)

	assert.Error(t, err)
	assert.Nil(t, incapacidad)
	assert.Contains(t, err.Error(), "permiso")
}

func TestTranscribirIncapacidad_Success(t *testing.T) {
	mockRepo := new(MockIncapacidadRepository)

	uc := usecase.NewTranscripcionUseCase(mockRepo)

	now := time.Now()
	fechaLimite := now.Add(24 * time.Hour)

	incapacidad := &domain.Incapacidad{
		IDIncapacidad:              1,
		Titulo:                    "Test Incapacidad",
		EstadoTranscripcion:       "pendiente",
		FechaLimiteTranscripcion:   &fechaLimite,
	}

	mockRepo.On("FindByID", mock.Anything, uint64(1)).Return(incapacidad, nil)
	mockRepo.On("Update", mock.Anything, mock.Anything).Return(nil)

	actor := ports.Actor{UserID: 1, Role: "admin", Permisos: []string{"transcribir_incapacidad"}}

	result, err := uc.Transcribir(context.Background(), 1, actor, nil)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "completado", result.EstadoTranscripcion)

	mockRepo.AssertExpectations(t)
}

func TestTranscribirIncapacidad_YaCompletada(t *testing.T) {
	mockRepo := new(MockIncapacidadRepository)

	uc := usecase.NewTranscripcionUseCase(mockRepo)

	incapacidad := &domain.Incapacidad{
		IDIncapacidad:        1,
		EstadoTranscripcion:  "completado",
	}

	mockRepo.On("FindByID", mock.Anything, uint64(1)).Return(incapacidad, nil)

	actor := ports.Actor{UserID: 1, Role: "admin", Permisos: []string{"transcribir_incapacidad"}}

	result, err := uc.Transcribir(context.Background(), 1, actor, nil)

	assert.Error(t, err)
	assert.Nil(t, result)

	mockRepo.AssertExpectations(t)
}

func TestMarcarEnProceso_Success(t *testing.T) {
	mockRepo := new(MockIncapacidadRepository)

	uc := usecase.NewTranscripcionUseCase(mockRepo)

	incapacidad := &domain.Incapacidad{
		IDIncapacidad:       1,
		EstadoTranscripcion: "pendiente",
	}

	mockRepo.On("FindByID", mock.Anything, uint64(1)).Return(incapacidad, nil)
	mockRepo.On("Update", mock.Anything, mock.Anything).Return(nil)

	actor := ports.Actor{UserID: 1, Role: "admin", Permisos: []string{"transcribir_incapacidad"}}

	result, err := uc.MarcarEnProceso(context.Background(), 1, actor)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "en_proceso", result.EstadoTranscripcion)

	mockRepo.AssertExpectations(t)
}
