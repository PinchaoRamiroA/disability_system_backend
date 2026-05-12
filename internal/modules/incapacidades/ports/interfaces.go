package ports

import (
	"context"

	"disability_system_backend/internal/modules/incapacidades/domain"
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

func (a Actor) CanManageIncapacidades() bool {
	return a.HasPermission("editar_incapacidad") ||
		a.HasPermission("validar_documentos") ||
		a.HasPermission("registrar_pago") ||
		a.HasPermission("consultar_reportes") ||
		a.HasPermission("consultar_historial") ||
		a.HasPermission("gestionar_cobro_persuasivo") ||
		a.HasPermission("gestionar_cobro_juridico") ||
		a.HasPermission("archivar_incapacidad")
}

type IncapacidadFilters struct {
	UserID         *uint64
	EstadoID       *uint64
	TipoID         *uint64
	EntidadID      *uint64
	Origen         string
	CanalRecepcion string
	IncludeDeleted bool
	Page           int
	Limit          int
}

type IncapacidadRepository interface {
	Create(ctx context.Context, incapacidad *domain.Incapacidad) error
	FindByID(ctx context.Context, id uint64) (*domain.Incapacidad, error)
	List(ctx context.Context, filters IncapacidadFilters) ([]domain.Incapacidad, int64, error)
	Update(ctx context.Context, incapacidad *domain.Incapacidad) error
	SoftDelete(ctx context.Context, id uint64) error
	ExistsUsuario(ctx context.Context, id uint64) (bool, error)
	FindEstadoByID(ctx context.Context, id uint64) (*domain.EstadoIncapacidad, error)
	FindEstadoByName(ctx context.Context, name string) (*domain.EstadoIncapacidad, error)
	FindTipoByID(ctx context.Context, id uint64) (*domain.TipoIncapacidad, error)
	FindEntidadByID(ctx context.Context, id uint64) (*domain.Entidad, error)
	ListEstados(ctx context.Context) ([]domain.EstadoIncapacidad, error)
	ListTipos(ctx context.Context) ([]domain.TipoIncapacidad, error)
	ListEntidades(ctx context.Context) ([]domain.Entidad, error)
	ListEstadosDocumento(ctx context.Context) ([]domain.EstadoDocumento, error)
	ListTiposDocumento(ctx context.Context) ([]domain.TipoDocumento, error)
	ListTiposPago(ctx context.Context) ([]domain.TipoPago, error)
	FindTiposDocumentoByNombre(ctx context.Context, nombres []string) ([]domain.TipoDocumento, error)
}

type PermissionRepository interface {
	FindPermissionsByRoleName(ctx context.Context, role string) ([]string, error)
}
