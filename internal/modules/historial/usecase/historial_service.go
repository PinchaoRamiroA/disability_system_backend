package usecase

import (
	"context"
	"time"

	"disability_system_backend/internal/modules/historial/domain"
	"disability_system_backend/internal/modules/historial/ports"
)

type HistorialService struct {
	repo ports.HistorialRepository
}

func NewHistorialService(repo ports.HistorialRepository) *HistorialService {
	return &HistorialService{repo: repo}
}

func (s *HistorialService) CreateEntry(ctx context.Context, incapacidadID, tipoID uint64, descripcion string, gestorID *uint64) error {
	historial := &domain.Historial{
		IDIncapacidad:   incapacidadID,
		IDTipoHistorial: tipoID,
		Descripcion:     descripcion,
		Fecha:           time.Now(),
		GestorID:        gestorID,
	}
	return s.repo.Create(ctx, historial)
}

func (s *HistorialService) List(ctx context.Context, incapacidadID uint64, tipoID *uint64, page, limit int) ([]domain.Historial, int64, error) {
	return s.repo.List(ctx, incapacidadID, tipoID, page, limit)
}