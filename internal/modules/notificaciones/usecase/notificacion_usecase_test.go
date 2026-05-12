package usecase

import (
	"context"
	"testing"
	"time"

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

type fakeNotificacionRepository struct {
	items       map[uint64]*domain.Notificacion
	lastFilters ports.NotificacionFilters
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
	}
}

func (r *fakeNotificacionRepository) Create(ctx context.Context, notificacion *domain.Notificacion) error {
	notificacion.IDNotificacion = 2
	r.items[2] = notificacion
	return nil
}

func (r *fakeNotificacionRepository) FindByID(ctx context.Context, id uint64) (*domain.Notificacion, error) {
	return r.items[id], nil
}

func (r *fakeNotificacionRepository) List(ctx context.Context, filters ports.NotificacionFilters) ([]domain.Notificacion, int64, error) {
	r.lastFilters = filters
	items := make([]domain.Notificacion, 0, len(r.items))
	for _, item := range r.items {
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
