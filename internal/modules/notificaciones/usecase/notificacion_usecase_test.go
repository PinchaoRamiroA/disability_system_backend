package usecase

import (
	"context"
	"testing"
	"time"

	incdomain "disability_system_backend/internal/modules/incapacidades/domain"
	"disability_system_backend/internal/modules/notificaciones/domain"
	"disability_system_backend/internal/modules/notificaciones/ports"

	"github.com/stretchr/testify/require"
)

func TestCrearRequiresGenerarAlertas(t *testing.T) {
	repo := newFakeNotificacionRepository()
	uc := NewNotificacionUseCase(repo)
	actor := ports.Actor{UserID: 1, Permisos: []string{"consultar_incapacidad"}}

	_, err := uc.Crear(context.Background(), actor, CrearNotificacionInput{
		IDUsuario:        2,
		TipoNotificacion: "Alerta",
		Mensaje:          "Mensaje",
	})

	require.Error(t, err)
	require.Contains(t, err.Error(), "no tienes permiso")
}

func TestCrearWithPermission(t *testing.T) {
	repo := newFakeNotificacionRepository()
	uc := NewNotificacionUseCase(repo)
	actor := ports.Actor{UserID: 1, Permisos: []string{"generar_alertas"}}

	notificacion, err := uc.Crear(context.Background(), actor, CrearNotificacionInput{
		IDUsuario:        2,
		TipoNotificacion: "Alerta",
		Mensaje:          "Mensaje",
	})

	require.NoError(t, err)
	require.Equal(t, uint64(2), notificacion.IDUsuario)
	require.Equal(t, "Alerta", notificacion.TipoNotificacion)
}

func TestListarScopesRegularUserToOwnNotifications(t *testing.T) {
	repo := newFakeNotificacionRepository()
	uc := NewNotificacionUseCase(repo)
	actor := ports.Actor{UserID: 7, Permisos: []string{"consultar_incapacidad"}}

	_, _, err := uc.Listar(context.Background(), actor, ports.NotificacionFilters{})

	require.NoError(t, err)
	require.NotNil(t, repo.lastFilters.IDUsuario)
	require.Equal(t, uint64(7), *repo.lastFilters.IDUsuario)
}

func TestObtenerRejectsOtherUserNotification(t *testing.T) {
	repo := newFakeNotificacionRepository()
	uc := NewNotificacionUseCase(repo)
	actor := ports.Actor{UserID: 7, Permisos: []string{"consultar_incapacidad"}}

	_, err := uc.Obtener(context.Background(), actor, 1)

	require.Error(t, err)
	require.Contains(t, err.Error(), "propias notificaciones")
}

func TestDocumentoFaltanteNotifierCreatesNotification(t *testing.T) {
	repo := newFakeNotificacionRepository()
	notifier := NewDocumentoFaltanteNotifier(repo)

	err := notifier.NotificarDocumentosFaltantes(context.Background(), 7, 10, incapacidadTipoDocumentos{
		{Nombre: "Certificado de incapacidad", Requerido: true},
	}.toDomain())

	require.NoError(t, err)
	require.Len(t, repo.items, 2)
	require.Equal(t, "Documento faltante", repo.items[2].TipoNotificacion)
	require.Contains(t, repo.items[2].Mensaje, "Certificado de incapacidad")
}

func TestDocumentoFaltanteNotifierAvoidsDuplicateUnreadNotification(t *testing.T) {
	repo := newFakeNotificacionRepository()
	notifier := NewDocumentoFaltanteNotifier(repo)
	documentos := incapacidadTipoDocumentos{
		{Nombre: "Certificado de incapacidad", Requerido: true},
	}.toDomain()

	require.NoError(t, notifier.NotificarDocumentosFaltantes(context.Background(), 7, 10, documentos))
	require.NoError(t, notifier.NotificarDocumentosFaltantes(context.Background(), 7, 10, documentos))

	require.Len(t, repo.items, 2)
}

type fakeNotificacionRepository struct {
	items       map[uint64]*domain.Notificacion
	lastFilters ports.NotificacionFilters
	nextID      uint64
}

func newFakeNotificacionRepository() *fakeNotificacionRepository {
	now := time.Now()
	return &fakeNotificacionRepository{
		items: map[uint64]*domain.Notificacion{
			1: {
				IDNotificacion:   1,
				IDUsuario:        99,
				TipoNotificacion: "Alerta",
				Mensaje:          "Mensaje",
				Fecha:            now,
			},
		},
		nextID: 2,
	}
}

func (r *fakeNotificacionRepository) Create(ctx context.Context, notificacion *domain.Notificacion) error {
	notificacion.IDNotificacion = r.nextID
	r.items[r.nextID] = notificacion
	r.nextID++
	return nil
}

func (r *fakeNotificacionRepository) FindByID(ctx context.Context, id uint64) (*domain.Notificacion, error) {
	return r.items[id], nil
}

func (r *fakeNotificacionRepository) List(ctx context.Context, filters ports.NotificacionFilters) ([]domain.Notificacion, int64, error) {
	r.lastFilters = filters
	items := make([]domain.Notificacion, 0, len(r.items))
	for _, item := range r.items {
		if filters.IDUsuario != nil && item.IDUsuario != *filters.IDUsuario {
			continue
		}
		if filters.IDIncapacidad != nil && (item.IDIncapacidad == nil || *item.IDIncapacidad != *filters.IDIncapacidad) {
			continue
		}
		if filters.TipoNotificacion != "" && item.TipoNotificacion != filters.TipoNotificacion {
			continue
		}
		if filters.Leida != nil && item.Leida != *filters.Leida {
			continue
		}
		items = append(items, *item)
	}
	return items, int64(len(items)), nil
}

func (r *fakeNotificacionRepository) CountUnread(ctx context.Context, userID uint64) (int64, error) {
	return 1, nil
}

func (r *fakeNotificacionRepository) MarkAsRead(ctx context.Context, id uint64) error {
	r.items[id].Leida = true
	return nil
}

func (r *fakeNotificacionRepository) MarkAllAsRead(ctx context.Context, userID uint64) error {
	return nil
}

func (r *fakeNotificacionRepository) Delete(ctx context.Context, id uint64) error {
	delete(r.items, id)
	return nil
}

func (r *fakeNotificacionRepository) UserExists(ctx context.Context, id uint64) (bool, error) {
	return true, nil
}

func (r *fakeNotificacionRepository) IncapacidadExists(ctx context.Context, id uint64) (bool, error) {
	return true, nil
}

type incapacidadTipoDocumento struct {
	Nombre    string
	Requerido bool
}

type incapacidadTipoDocumentos []incapacidadTipoDocumento

func (items incapacidadTipoDocumentos) toDomain() []incdomain.TipoDocumento {
	result := make([]incdomain.TipoDocumento, 0, len(items))
	for _, item := range items {
		result = append(result, incdomain.TipoDocumento{
			Nombre:    item.Nombre,
			Requerido: item.Requerido,
		})
	}
	return result
}
