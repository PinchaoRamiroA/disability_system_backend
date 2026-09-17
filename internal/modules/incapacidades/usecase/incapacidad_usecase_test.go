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

func TestIncapacidadUseCase_CambiarEstado_BUG07(t *testing.T) {
	ctx := context.Background()

	actorEditor := ports.Actor{
		UserID:   1,
		Permisos: []string{"editar_incapacidad"},
	}
	actorArchivador := ports.Actor{
		UserID:   2,
		Permisos: []string{"editar_incapacidad", "archivar_incapacidad"},
	}
	actorSinPermiso := ports.Actor{
		UserID:   3,
		Permisos: []string{"consultar_incapacidad"},
	}

	estadoRecibida := &domain.EstadoIncapacidad{
		IDEstado:          1,
		Nombre:            "Recibida",
		PermiteTransicion: true,
	}
	estadoEnValidacion := &domain.EstadoIncapacidad{
		IDEstado:          2,
		Nombre:            "En validación documental",
		PermiteTransicion: true,
	}
	estadoAprobada := &domain.EstadoIncapacidad{
		IDEstado:          7,
		Nombre:            "Aprobada",
		PermiteTransicion: true,
	}
	estadoPagada := &domain.EstadoIncapacidad{
		IDEstado:          10,
		Nombre:            "Pagada",
		PermiteTransicion: true,
	}
	estadoRechazada := &domain.EstadoIncapacidad{
		IDEstado:          99,
		Nombre:            "Rechazada",
		PermiteTransicion: false,
	}
	estadoArchivada := &domain.EstadoIncapacidad{
		IDEstado:          100,
		Nombre:            "Archivada",
		PermiteTransicion: false,
	}

	t.Run("should allow valid state transition (Recibida -> En validación documental)", func(t *testing.T) {
		repo := new(mockIncapacidadRepo)
		uc := usecase.NewIncapacidadUseCase(repo)

		incapacidad := &domain.Incapacidad{
			IDIncapacidad: 1,
			IDEstado:      1,
			Estado:        estadoRecibida,
		}

		repo.On("FindByID", ctx, uint64(1)).Return(incapacidad, nil)
		repo.On("FindEstadoByID", ctx, uint64(2)).Return(estadoEnValidacion, nil)
		repo.On("Update", ctx, mock.MatchedBy(func(inc *domain.Incapacidad) bool {
			return inc.IDEstado == 2 && inc.Estado.Nombre == "En validación documental"
		})).Return(nil)

		obs := "Documentos recibidos completos"
		res, err := uc.CambiarEstado(ctx, actorEditor, 1, 2, &obs)

		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, uint64(2), res.IDEstado)
		assert.Equal(t, "En validación documental", res.Estado.Nombre)
		assert.Equal(t, &obs, res.Observaciones)

		repo.AssertExpectations(t)
	})

	t.Run("should reject illegal state transition (Pagada -> Recibida) - BUG-07", func(t *testing.T) {
		repo := new(mockIncapacidadRepo)
		uc := usecase.NewIncapacidadUseCase(repo)

		incapacidad := &domain.Incapacidad{
			IDIncapacidad: 1,
			IDEstado:      10,
			Estado:        estadoPagada,
		}

		repo.On("FindByID", ctx, uint64(1)).Return(incapacidad, nil)
		repo.On("FindEstadoByID", ctx, uint64(1)).Return(estadoRecibida, nil)

		res, err := uc.CambiarEstado(ctx, actorEditor, 1, 1, nil)

		assert.Error(t, err)
		assert.Nil(t, res)

		var appErr *apperrors.AppError
		assert.ErrorAs(t, err, &appErr)
		assert.Equal(t, "CONFLICT", appErr.Code)
		assert.Contains(t, err.Error(), "transición de estado no válida: de Pagada a Recibida")

		repo.AssertNotCalled(t, "Update", mock.Anything, mock.Anything)
		repo.AssertExpectations(t)
	})

	t.Run("should reject skipping required states (Recibida -> Aprobada)", func(t *testing.T) {
		repo := new(mockIncapacidadRepo)
		uc := usecase.NewIncapacidadUseCase(repo)

		incapacidad := &domain.Incapacidad{
			IDIncapacidad: 1,
			IDEstado:      1,
			Estado:        estadoRecibida,
		}

		repo.On("FindByID", ctx, uint64(1)).Return(incapacidad, nil)
		repo.On("FindEstadoByID", ctx, uint64(7)).Return(estadoAprobada, nil)

		res, err := uc.CambiarEstado(ctx, actorEditor, 1, 7, nil)

		assert.Error(t, err)
		assert.Nil(t, res)
		assert.Contains(t, err.Error(), "transición de estado no válida: de Recibida a Aprobada")

		repo.AssertNotCalled(t, "Update", mock.Anything, mock.Anything)
		repo.AssertExpectations(t)
	})

	t.Run("should allow transition to universal states (Rechazada)", func(t *testing.T) {
		repo := new(mockIncapacidadRepo)
		uc := usecase.NewIncapacidadUseCase(repo)

		incapacidad := &domain.Incapacidad{
			IDIncapacidad: 1,
			IDEstado:      1,
			Estado:        estadoRecibida,
		}

		repo.On("FindByID", ctx, uint64(1)).Return(incapacidad, nil)
		repo.On("FindEstadoByID", ctx, uint64(99)).Return(estadoRechazada, nil)
		repo.On("Update", ctx, mock.Anything).Return(nil)

		res, err := uc.CambiarEstado(ctx, actorEditor, 1, 99, nil)

		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, uint64(99), res.IDEstado)

		repo.AssertExpectations(t)
	})

	t.Run("should allow transition to Archivada when actor has archivar_incapacidad permission", func(t *testing.T) {
		repo := new(mockIncapacidadRepo)
		uc := usecase.NewIncapacidadUseCase(repo)

		incapacidad := &domain.Incapacidad{
			IDIncapacidad: 1,
			IDEstado:      1,
			Estado:        estadoRecibida,
		}

		repo.On("FindByID", ctx, uint64(1)).Return(incapacidad, nil)
		repo.On("FindEstadoByID", ctx, uint64(100)).Return(estadoArchivada, nil)
		repo.On("Update", ctx, mock.Anything).Return(nil)

		res, err := uc.CambiarEstado(ctx, actorArchivador, 1, 100, nil)

		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, uint64(100), res.IDEstado)

		repo.AssertExpectations(t)
	})

	t.Run("should reject transition to Archivada when actor lacks archivar_incapacidad permission", func(t *testing.T) {
		repo := new(mockIncapacidadRepo)
		uc := usecase.NewIncapacidadUseCase(repo)

		incapacidad := &domain.Incapacidad{
			IDIncapacidad: 1,
			IDEstado:      1,
			Estado:        estadoRecibida,
		}

		repo.On("FindByID", ctx, uint64(1)).Return(incapacidad, nil)
		repo.On("FindEstadoByID", ctx, uint64(100)).Return(estadoArchivada, nil)

		res, err := uc.CambiarEstado(ctx, actorEditor, 1, 100, nil)

		assert.Error(t, err)
		assert.Nil(t, res)
		assert.Contains(t, err.Error(), "no tienes permiso para archivar incapacidades")

		repo.AssertNotCalled(t, "Update", mock.Anything, mock.Anything)
	})

	t.Run("should reject transition when current state does not permit transitions", func(t *testing.T) {
		repo := new(mockIncapacidadRepo)
		uc := usecase.NewIncapacidadUseCase(repo)

		estadoCerrada := &domain.EstadoIncapacidad{
			IDEstado:          50,
			Nombre:            "Cerrada",
			PermiteTransicion: false,
		}
		incapacidad := &domain.Incapacidad{
			IDIncapacidad: 1,
			IDEstado:      50,
			Estado:        estadoCerrada,
		}

		repo.On("FindByID", ctx, uint64(1)).Return(incapacidad, nil)
		repo.On("FindEstadoByID", ctx, uint64(1)).Return(estadoRecibida, nil)

		res, err := uc.CambiarEstado(ctx, actorEditor, 1, 1, nil)

		assert.Error(t, err)
		assert.Nil(t, res)
		assert.Contains(t, err.Error(), "el estado actual no permite transiciones")

		repo.AssertNotCalled(t, "Update", mock.Anything, mock.Anything)
	})

	t.Run("should reject when actor has no permissions", func(t *testing.T) {
		repo := new(mockIncapacidadRepo)
		uc := usecase.NewIncapacidadUseCase(repo)

		res, err := uc.CambiarEstado(ctx, actorSinPermiso, 1, 2, nil)

		assert.Error(t, err)
		assert.Nil(t, res)
		assert.Contains(t, err.Error(), "no tienes permiso")
	})

	t.Run("should resolve current state from repository if incapacidad.Estado is nil", func(t *testing.T) {
		repo := new(mockIncapacidadRepo)
		uc := usecase.NewIncapacidadUseCase(repo)

		// Incapacidad without preloaded Estado relation
		incapacidad := &domain.Incapacidad{
			IDIncapacidad: 1,
			IDEstado:      1,
			Estado:        nil,
		}

		repo.On("FindByID", ctx, uint64(1)).Return(incapacidad, nil)
		repo.On("FindEstadoByID", ctx, uint64(2)).Return(estadoEnValidacion, nil)
		repo.On("FindEstadoByID", ctx, uint64(1)).Return(estadoRecibida, nil)
		repo.On("Update", ctx, mock.Anything).Return(nil)

		res, err := uc.CambiarEstado(ctx, actorEditor, 1, 2, nil)

		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, uint64(2), res.IDEstado)
		assert.Equal(t, "En validación documental", res.Estado.Nombre)

		repo.AssertExpectations(t)
	})
}
