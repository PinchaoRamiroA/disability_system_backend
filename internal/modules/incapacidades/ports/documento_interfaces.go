package ports

import (
	"context"

	"disability_system_backend/internal/modules/incapacidades/domain"
)

type DocumentoRepository interface {
	Create(ctx context.Context, documento *domain.Documento) error
	FindByID(ctx context.Context, id uint64) (*domain.Documento, error)
	List(ctx context.Context, incapacidadID uint64, estado, tipo string, page, limit int) ([]domain.Documento, int64, error)
	Update(ctx context.Context, documento *domain.Documento) error
	Delete(ctx context.Context, id uint64) error
	ExistsIncapacidad(ctx context.Context, id uint64) (bool, error)
}

type HistorialService interface {
	CreateEntry(ctx context.Context, incapacidadID, tipoID uint64, descripcion string, gestorID *uint64) error
}

type HistorialEntry struct {
	IDIncapacidad   uint64
	IDTipoHistorial uint64
	Descripcion     string
	GestorID        *uint64
}