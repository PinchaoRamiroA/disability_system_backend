package ports

import (
	"context"

	"disability_system_backend/internal/modules/cobros/domain"
)

type Actor struct {
	UserID   uint64
	Role     string
	Permisos []string
}

func (a Actor) HasPermission(permission string) bool {
	for _, p := range a.Permisos {
		if p == permission {
			return true
		}
	}
	return false
}

type PagoFilters struct {
	IDIncapacidad *uint64
	IDEntidad     *uint64
	TipoPago      string
	EstadoPago    string
	Conciliado    *bool
	Page          int
	Limit         int
}

type SeguimientoFilters struct {
	IDIncapacidad   *uint64
	TipoSeguimiento string
	Page            int
	Limit           int
}

type PagoRepository interface {
	CreatePago(ctx context.Context, pago *domain.Pago) error
	FindPagoByID(ctx context.Context, id uint64) (*domain.Pago, error)
	ListPagos(ctx context.Context, filters PagoFilters) ([]domain.Pago, int64, error)
	UpdatePago(ctx context.Context, pago *domain.Pago) error
	SoftDeletePago(ctx context.Context, id uint64) error
	IncapacidadExists(ctx context.Context, id uint64) (bool, error)
	EntidadExists(ctx context.Context, id uint64) (bool, error)
	GetEntidadInfo(ctx context.Context) (map[uint64]struct{ Nombre, Tipo string }, error)
	GetIncapacidadesDetailed(ctx context.Context, ids []uint64) (map[uint64]IncapacidadInfo, error)
}

type IncapacidadInfo struct {
	ID     uint64
	Titulo string
}

type SeguimientoRepository interface {
	CreateSeguimiento(ctx context.Context, seguimiento *domain.SeguimientoCobro) error
	FindSeguimientoByID(ctx context.Context, id uint64) (*domain.SeguimientoCobro, error)
	ListSeguimientos(ctx context.Context, filters SeguimientoFilters) ([]domain.SeguimientoCobro, int64, error)
	UpdateSeguimiento(ctx context.Context, seguimiento *domain.SeguimientoCobro) error
}

type PermissionRepository interface {
	FindPermissionsByRoleName(ctx context.Context, role string) ([]string, error)
}
