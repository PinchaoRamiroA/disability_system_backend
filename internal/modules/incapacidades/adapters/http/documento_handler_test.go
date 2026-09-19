package http_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	httpadapter "disability_system_backend/internal/modules/incapacidades/adapters/http"
	"disability_system_backend/internal/modules/incapacidades/domain"
	"disability_system_backend/internal/modules/incapacidades/ports"
	"disability_system_backend/internal/modules/incapacidades/usecase"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockDocRepo struct {
	mock.Mock
}

func (m *mockDocRepo) Create(ctx context.Context, documento *domain.Documento) error {
	return m.Called(ctx, documento).Error(0)
}

func (m *mockDocRepo) FindByID(ctx context.Context, id uint64) (*domain.Documento, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Documento), args.Error(1)
}

func (m *mockDocRepo) List(ctx context.Context, incapacidadID uint64, estado, tipo string, page, limit int) ([]domain.Documento, int64, error) {
	args := m.Called(ctx, incapacidadID, estado, tipo, page, limit)
	return args.Get(0).([]domain.Documento), args.Get(1).(int64), args.Error(2)
}

func (m *mockDocRepo) Update(ctx context.Context, documento *domain.Documento) error {
	return m.Called(ctx, documento).Error(0)
}

func (m *mockDocRepo) Delete(ctx context.Context, id uint64) error {
	return m.Called(ctx, id).Error(0)
}

func (m *mockDocRepo) ExistsIncapacidad(ctx context.Context, id uint64) (bool, error) {
	args := m.Called(ctx, id)
	return args.Bool(0), args.Error(1)
}

func setupTestRouter(handler *httpadapter.DocumentoHandler, actor ports.Actor) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("actor", actor)
		c.Next()
	})

	router.GET("/incapacidades/:id/documentos", handler.Listar)
	router.GET("/documentos", handler.Listar)
	return router
}

func TestDocumentoHandler_Listar_BUG08(t *testing.T) {
	actor := ports.Actor{
		UserID:   1,
		Role:     "admin",
		Permisos: []string{"consultar_incapacidad"},
	}

	dummyDocs := []domain.Documento{
		{
			IDDocumento:   1,
			IDIncapacidad: 123,
			Nombre:        "epicrisis.pdf",
			Tipo:          "Epicrisis",
			Estado:        "Validado",
			URL:           "https://storage/123/epicrisis.pdf",
			Formato:       ".pdf",
			CreatedAt:     time.Now(),
		},
	}

	t.Run("should read route parameter :id when no query param is sent (BUG-08 fix)", func(t *testing.T) {
		repo := new(mockDocRepo)
		uc := usecase.NewDocumentoUseCase(repo, nil)
		handler := httpadapter.NewDocumentoHandler(uc, nil, nil)
		router := setupTestRouter(handler, actor)

		repo.On("List", mock.Anything, uint64(123), "", "", 1, 20).
			Return(dummyDocs, int64(1), nil)

		req := httptest.NewRequest(http.MethodGet, "/incapacidades/123/documentos", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var resp map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.True(t, resp["success"].(bool))

		repo.AssertExpectations(t)
	})

	t.Run("should read route parameter :id along with query filters", func(t *testing.T) {
		repo := new(mockDocRepo)
		uc := usecase.NewDocumentoUseCase(repo, nil)
		handler := httpadapter.NewDocumentoHandler(uc, nil, nil)
		router := setupTestRouter(handler, actor)

		repo.On("List", mock.Anything, uint64(123), "Validado", "Epicrisis", 2, 10).
			Return(dummyDocs, int64(1), nil)

		req := httptest.NewRequest(http.MethodGet, "/incapacidades/123/documentos?estado=Validado&tipo=Epicrisis&page=2&limit=10", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		repo.AssertExpectations(t)
	})

	t.Run("should fallback to query param id_incapacidad if path param is not present", func(t *testing.T) {
		repo := new(mockDocRepo)
		uc := usecase.NewDocumentoUseCase(repo, nil)
		handler := httpadapter.NewDocumentoHandler(uc, nil, nil)
		router := setupTestRouter(handler, actor)

		repo.On("List", mock.Anything, uint64(456), "", "", 1, 20).
			Return(dummyDocs, int64(1), nil)

		req := httptest.NewRequest(http.MethodGet, "/documentos?id_incapacidad=456", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		repo.AssertExpectations(t)
	})

	t.Run("should return 400 Bad Request when path param is invalid string", func(t *testing.T) {
		repo := new(mockDocRepo)
		uc := usecase.NewDocumentoUseCase(repo, nil)
		handler := httpadapter.NewDocumentoHandler(uc, nil, nil)
		router := setupTestRouter(handler, actor)

		req := httptest.NewRequest(http.MethodGet, "/incapacidades/abc/documentos", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "id de incapacidad inválido")

		repo.AssertNotCalled(t, "List", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything)
	})

	t.Run("should return 400 Bad Request when neither path param nor query param are provided", func(t *testing.T) {
		repo := new(mockDocRepo)
		uc := usecase.NewDocumentoUseCase(repo, nil)
		handler := httpadapter.NewDocumentoHandler(uc, nil, nil)
		router := setupTestRouter(handler, actor)

		req := httptest.NewRequest(http.MethodGet, "/documentos", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "id_incapacidad es requerido")

		repo.AssertNotCalled(t, "List", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything)
	})
}

func TestGetExtensionFromFilename_BUG15(t *testing.T) {
	testCases := []struct {
		filename string
		expected string
	}{
		{"documento.pdf", ".pdf"},
		{"foto.jpg", ".jpg"},
		{"foto.jpeg", ".jpeg"},
		{"imagen.png", ".png"},
		{"SCAN.PDF", ".pdf"},
		{"IMAGE.JPEG", ".jpeg"},
		{"FOTO.PNG", ".png"},
		{"a.go", ".go"},
		{"1.pdf", ".pdf"},
		{"reporte.2024.final.docx", ".docx"},
		{"archivo_sin_extension", ""},
		{"1234", ""},
		{"12345", ""},
		{".gitignore", ".gitignore"},
		{".env.local", ".local"},
	}

	for _, tc := range testCases {
		t.Run(tc.filename, func(t *testing.T) {
			actual := httpadapter.GetExtensionFromFilename(tc.filename)
			assert.Equal(t, tc.expected, actual, "filename %s should yield extension %s", tc.filename, tc.expected)
		})
	}
}
