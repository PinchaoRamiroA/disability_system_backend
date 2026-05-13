package ports

import (
	"context"

	"disability_system_backend/internal/modules/reportes/domain"
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

type ReportesRepository interface {
	GetIncapacidadesReport(ctx context.Context, filtros domain.FiltrosReporte) ([]domain.ReporteIncapacidadDetalle, error)
	GetAusentismoReport(ctx context.Context, filtros domain.FiltrosReporte) ([]domain.ReporteAusentismoDetalle, []domain.ReporteTopIncapacidad, error)
	GetCarteraReport(ctx context.Context, filtros domain.FiltrosReporte) (*domain.ReporteCartera, error)
	GetEntidadReport(ctx context.Context, entidadID uint64, filtros domain.FiltrosReporte) (*domain.ReporteEntidad, error)
	GetVencimientosReport(ctx context.Context, diasMinimos int) (*domain.ReporteVencimientos, error)

	GetTotalIncapacidades(ctx context.Context, filtros domain.FiltrosReporte) (total, activas, cerradas int64, err error)
	GetTotalDiasPerdidos(ctx context.Context, filtros domain.FiltrosReporte) (int64, error)
	GetDiasPorTipo(ctx context.Context, filtros domain.FiltrosReporte) (map[string]int64, error)
	GetDiasPorEntidad(ctx context.Context, filtros domain.FiltrosReporte) (map[string]int64, error)
	GetTopEmpleadosIncapacidades(ctx context.Context, filtros domain.FiltrosReporte, limit int) ([]domain.ReporteTopIncapacidad, error)

	CountIncapacidadesActivas(ctx context.Context) (int64, error)
	CountPagosPendientes(ctx context.Context) (int64, error)
	CountPagosVencidos(ctx context.Context, diasMinimos int) (int64, error)
	SumValorCartera(ctx context.Context) (pagado, pendiente float64, err error)
	GetCarteraPorEntidad(ctx context.Context) ([]domain.ReporteCarteraEntidad, error)
}

type PermissionRepository interface {
	FindPermissionsByRoleName(ctx context.Context, role string) ([]string, error)
}
