package ports

import (
	"context"

	"disability_system_backend/internal/modules/historial/domain"
)

type HistorialRepository interface {
	Create(ctx context.Context, historial *domain.Historial) error
	List(ctx context.Context, incapacidadID uint64, tipoID *uint64, page, limit int) ([]domain.Historial, int64, error)
	FindByID(ctx context.Context, id uint64) (*domain.Historial, error)
	FindTipoByID(ctx context.Context, id uint64) (*domain.TipoHistorial, error)
}