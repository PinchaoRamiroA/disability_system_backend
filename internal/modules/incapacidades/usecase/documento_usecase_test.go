package usecase_test

import (
	"context"
	"testing"

	"disability_system_backend/internal/modules/incapacidades/domain"
	"disability_system_backend/internal/modules/incapacidades/ports"
	"disability_system_backend/internal/modules/incapacidades/usecase"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockDocumentoRepo struct {
	mock.Mock
}

func (m *mockDocumentoRepo) Create(ctx context.Context, documento *domain.Documento) error {
	args := m.Called(ctx, documento)
	if args.Error(0) == nil {
		documento.IDDocumento = 1
	}
	return args.Error(0)
}

func (m *mockDocumentoRepo) FindByID(ctx context.Context, id uint64) (*domain.Documento, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Documento), args.Error(1)
}

func (m *mockDocumentoRepo) List(ctx context.Context, incapacidadID uint64, estado, tipo string, page, limit int) ([]domain.Documento, int64, error) {
	args := m.Called(ctx, incapacidadID, estado, tipo, page, limit)
	return args.Get(0).([]domain.Documento), args.Get(1).(int64), args.Error(2)
}

func (m *mockDocumentoRepo) Update(ctx context.Context, documento *domain.Documento) error {
	args := m.Called(ctx, documento)
	return args.Error(0)
}

func (m *mockDocumentoRepo) Delete(ctx context.Context, id uint64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *mockDocumentoRepo) ExistsIncapacidad(ctx context.Context, id uint64) (bool, error) {
	args := m.Called(ctx, id)
	return args.Bool(0), args.Error(1)
}

type mockHistorialSvc struct {
	mock.Mock
}

func (m *mockHistorialSvc) CreateEntry(ctx context.Context, incapacidadID, tipoID uint64, descripcion string, gestorID *uint64) error {
	args := m.Called(ctx, incapacidadID, tipoID, descripcion, gestorID)
	return args.Error(0)
}

func TestDocumentoUseCase_BUG03_NilHistorialSvc(t *testing.T) {
	ctx := context.Background()
	actor := ports.Actor{
		UserID:   10,
		Permisos: []string{"crear_incapacidad", "editar_incapacidad", "validar_documentos"},
	}

	t.Run("Subir: should not panic when historialSvc is untyped nil", func(t *testing.T) {
		repo := new(mockDocumentoRepo)
		uc := usecase.NewDocumentoUseCase(repo, nil)

		repo.On("ExistsIncapacidad", ctx, uint64(1)).Return(true, nil)
		repo.On("Create", ctx, mock.MatchedBy(func(d *domain.Documento) bool {
			return d.IDIncapacidad == 1 && d.Nombre == "orden_medica.pdf"
		})).Return(nil)

		input := struct {
			IDIncapacidad uint64
			Nombre        string
			Tipo          string
			URL           string
			Formato       string
		}{
			IDIncapacidad: 1,
			Nombre:        "orden_medica.pdf",
			Tipo:          "ORDEN_MEDICA",
			URL:           "http://storage/orden_medica.pdf",
			Formato:       "application/pdf",
		}

		assert.NotPanics(t, func() {
			doc, err := uc.Subir(ctx, actor, input)
			assert.NoError(t, err)
			assert.NotNil(t, doc)
			assert.Equal(t, "orden_medica.pdf", doc.Nombre)
		})

		repo.AssertExpectations(t)
	})

	t.Run("Subir: should not panic when historialSvc is typed nil (*mockHistorialSvc)(nil)", func(t *testing.T) {
		repo := new(mockDocumentoRepo)
		var nilConcreteSvc *mockHistorialSvc = nil
		var svc ports.HistorialService = nilConcreteSvc

		uc := usecase.NewDocumentoUseCase(repo, svc)

		repo.On("ExistsIncapacidad", ctx, uint64(1)).Return(true, nil)
		repo.On("Create", ctx, mock.Anything).Return(nil)

		input := struct {
			IDIncapacidad uint64
			Nombre        string
			Tipo          string
			URL           string
			Formato       string
		}{
			IDIncapacidad: 1,
			Nombre:        "historia_clinica.pdf",
			Tipo:          "HISTORIA_CLINICA",
			URL:           "http://storage/hc.pdf",
			Formato:       "application/pdf",
		}

		assert.NotPanics(t, func() {
			doc, err := uc.Subir(ctx, actor, input)
			assert.NoError(t, err)
			assert.NotNil(t, doc)
		})

		repo.AssertExpectations(t)
	})

	t.Run("Subir: should call CreateEntry when historialSvc is valid", func(t *testing.T) {
		repo := new(mockDocumentoRepo)
		historial := new(mockHistorialSvc)
		uc := usecase.NewDocumentoUseCase(repo, historial)

		repo.On("ExistsIncapacidad", ctx, uint64(1)).Return(true, nil)
		repo.On("Create", ctx, mock.Anything).Return(nil)
		historial.On("CreateEntry", ctx, uint64(1), uint64(1), "Documento 'epicrisis.pdf' subido al sistema", mock.MatchedBy(func(id *uint64) bool {
			return id != nil && *id == 10
		})).Return(nil)

		input := struct {
			IDIncapacidad uint64
			Nombre        string
			Tipo          string
			URL           string
			Formato       string
		}{
			IDIncapacidad: 1,
			Nombre:        "epicrisis.pdf",
			Tipo:          "EPICRISIS",
			URL:           "http://storage/epicrisis.pdf",
			Formato:       "application/pdf",
		}

		doc, err := uc.Subir(ctx, actor, input)
		assert.NoError(t, err)
		assert.NotNil(t, doc)

		repo.AssertExpectations(t)
		historial.AssertExpectations(t)
	})

	t.Run("Validar: should not panic when historialSvc is nil", func(t *testing.T) {
		repo := new(mockDocumentoRepo)
		uc := usecase.NewDocumentoUseCase(repo, nil)

		existingDoc := &domain.Documento{
			IDDocumento:   1,
			IDIncapacidad: 100,
			Nombre:        "certificado.pdf",
			Estado:        "Pendiente",
		}

		repo.On("FindByID", ctx, uint64(1)).Return(existingDoc, nil)
		repo.On("Update", ctx, mock.MatchedBy(func(d *domain.Documento) bool {
			return d.Estado == "Validado"
		})).Return(nil)

		assert.NotPanics(t, func() {
			doc, err := uc.Validar(ctx, actor, 1, "Validado", "Documento correcto")
			assert.NoError(t, err)
			assert.NotNil(t, doc)
			assert.Equal(t, "Validado", doc.Estado)
		})

		repo.AssertExpectations(t)
	})

	t.Run("Validar: should not panic when historialSvc is typed nil", func(t *testing.T) {
		repo := new(mockDocumentoRepo)
		var nilConcreteSvc *mockHistorialSvc = nil
		var svc ports.HistorialService = nilConcreteSvc

		uc := usecase.NewDocumentoUseCase(repo, svc)

		existingDoc := &domain.Documento{
			IDDocumento:   1,
			IDIncapacidad: 100,
			Nombre:        "certificado.pdf",
			Estado:        "Pendiente",
		}

		repo.On("FindByID", ctx, uint64(1)).Return(existingDoc, nil)
		repo.On("Update", ctx, mock.Anything).Return(nil)

		assert.NotPanics(t, func() {
			doc, err := uc.Validar(ctx, actor, 1, "Validado", "OK")
			assert.NoError(t, err)
			assert.NotNil(t, doc)
		})

		repo.AssertExpectations(t)
	})

	t.Run("Validar: should call CreateEntry when historialSvc is valid", func(t *testing.T) {
		repo := new(mockDocumentoRepo)
		historial := new(mockHistorialSvc)
		uc := usecase.NewDocumentoUseCase(repo, nil)
		uc.SetHistorialService(historial)

		existingDoc := &domain.Documento{
			IDDocumento:   1,
			IDIncapacidad: 100,
			Nombre:        "certificado.pdf",
			Estado:        "Pendiente",
		}

		repo.On("FindByID", ctx, uint64(1)).Return(existingDoc, nil)
		repo.On("Update", ctx, mock.Anything).Return(nil)
		historial.On("CreateEntry", ctx, uint64(100), uint64(1), "Documento 'certificado.pdf' Validado: Aprobado por médico", mock.MatchedBy(func(id *uint64) bool {
			return id != nil && *id == 10
		})).Return(nil)

		doc, err := uc.Validar(ctx, actor, 1, "Validado", "Aprobado por médico")
		assert.NoError(t, err)
		assert.NotNil(t, doc)

		repo.AssertExpectations(t)
		historial.AssertExpectations(t)
	})
}
