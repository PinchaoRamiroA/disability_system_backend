package ports

import (
	"context"

	"disability_system_backend/internal/modules/auditoria/domain"
	"disability_system_backend/internal/modules/auditoria/dto"
)

type Actor interface {
	GetUserID() uint64
	HasPermission(string) bool
}

type AuditoriaRepository interface {
	Create(ctx context.Context, auditoria *domain.Auditoria) error
	List(ctx context.Context, filters dto.ListarAuditoriaQuery) ([]domain.Auditoria, int64, error)
}
