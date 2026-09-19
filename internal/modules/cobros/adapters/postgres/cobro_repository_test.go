package postgres_test

import (
	"context"
	"testing"

	"disability_system_backend/internal/modules/cobros/adapters/postgres"

	"github.com/stretchr/testify/assert"
)

func TestCobroRepository_GetIncapacidadesDetailed_EmptySlice_BUG14(t *testing.T) {
	ctx := context.Background()

	// CobroRepository with nil db should safely return empty map when ids slice is empty or nil
	repo := postgres.NewCobroRepository(nil)

	t.Run("should return empty map without error when ids is nil", func(t *testing.T) {
		res, err := repo.GetIncapacidadesDetailed(ctx, nil)
		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Empty(t, res)
	})

	t.Run("should return empty map without error when ids is empty slice", func(t *testing.T) {
		res, err := repo.GetIncapacidadesDetailed(ctx, []uint64{})
		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Empty(t, res)
	})

	t.Run("should return empty map without error when ids contains only zeros", func(t *testing.T) {
		res, err := repo.GetIncapacidadesDetailed(ctx, []uint64{0, 0})
		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Empty(t, res)
	})
}
