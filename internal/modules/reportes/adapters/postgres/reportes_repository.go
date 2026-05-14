package postgres

import (
	"context"
	"fmt"
	"strings"
	"time"

	"disability_system_backend/internal/modules/reportes/domain"

	"gorm.io/gorm"
)

type ReportesRepository struct {
	db *gorm.DB
}

func NewReportesRepository(db *gorm.DB) *ReportesRepository {
	return &ReportesRepository{db: db}
}

func (r *ReportesRepository) buildBaseQuery(filtros domain.FiltrosReporte) *gorm.DB {
	query := r.db.Table("incapacidad i").
		Select("i.id_incapacidad, u.nombre as nombre_empleado, ti.nombre as tipo_incapacidad, e.nombre as estado, en.nombre as entidad, i.fecha_inicio, i.fecha_fin, i.observaciones").
		Joins("JOIN usuario u ON i.id_usuario = u.id_usuario").
		Joins("JOIN tipo_incapacidad ti ON i.id_tipo = ti.id_tipo").
		Joins("JOIN estado_incapacidad e ON i.id_estado = e.id_estado").
		Joins("JOIN entidad en ON i.id_entidad = en.id_entidad").
		Where("i.is_deleted = false")

	if filtros.FechaInicio != nil {
		query = query.Where("i.fecha_inicio >= ?", *filtros.FechaInicio)
	}
	if filtros.FechaFin != nil {
		query = query.Where("i.fecha_inicio <= ?", *filtros.FechaFin)
	}
	if filtros.IDEntidad != nil {
		query = query.Where("i.id_entidad = ?", *filtros.IDEntidad)
	}
	if filtros.IDTipo != nil {
		query = query.Where("i.id_tipo = ?", *filtros.IDTipo)
	}
	if filtros.IDEstado != nil {
		query = query.Where("i.id_estado = ?", *filtros.IDEstado)
	}
	if filtros.IDEmpleado != nil {
		query = query.Where("i.id_usuario = ?", *filtros.IDEmpleado)
	}
	if filtros.Origen != "" {
		query = query.Where("LOWER(i.origen) = LOWER(?)", filtros.Origen)
	}

	return query
}

func parseDate(dateStr string) (time.Time, error) {
	return time.Parse("2006-01-02", dateStr)
}

func (r *ReportesRepository) GetIncapacidadesReport(ctx context.Context, filtros domain.FiltrosReporte) ([]domain.ReporteIncapacidadDetalle, error) {
	query := r.buildBaseQuery(filtros)

	var results []struct {
		IDIncapacidad   uint64
		NombreEmpleado  string
		TipoIncapacidad string
		Estado          string
		Entidad         string
		FechaInicio     string
		FechaFin        *string
		Observaciones   *string
	}

	if err := query.Scan(&results).Error; err != nil {
		return nil, err
	}

	datos := make([]domain.ReporteIncapacidadDetalle, 0, len(results))
	for _, row := range results {
		datos = append(datos, domain.ReporteIncapacidadDetalle{
			IDIncapacidad:   row.IDIncapacidad,
			NombreEmpleado:  row.NombreEmpleado,
			TipoIncapacidad: row.TipoIncapacidad,
			Estado:          row.Estado,
			Entidad:         row.Entidad,
		})
	}

	return datos, nil
}

func (r *ReportesRepository) GetAusentismoReport(ctx context.Context, filtros domain.FiltrosReporte) ([]domain.ReporteAusentismoDetalle, []domain.ReporteTopIncapacidad, error) {
	query := r.db.WithContext(ctx).Table("incapacidad i").
		Select("i.id_usuario, u.nombre as nombre_empleado, i.fecha_inicio, i.fecha_fin, ti.nombre as tipo_incapacidad, en.nombre as entidad").
		Joins("JOIN usuario u ON i.id_usuario = u.id_usuario").
		Joins("JOIN tipo_incapacidad ti ON i.id_tipo = ti.id_tipo").
		Joins("JOIN entidad en ON i.id_entidad = en.id_entidad").
		Where("i.is_deleted = false")

	if filtros.FechaInicio != nil {
		query = query.Where("i.fecha_inicio >= ?", *filtros.FechaInicio)
	}
	if filtros.FechaFin != nil {
		query = query.Where("i.fecha_inicio <= ?", *filtros.FechaFin)
	}

	var results []struct {
		IDEmpleado      uint64
		NombreEmpleado  string
		FechaInicio     string
		FechaFin        *string
		TipoIncapacidad string
		Entidad         string
	}

	if err := query.Scan(&results).Error; err != nil {
		return nil, nil, err
	}

	datos := make([]domain.ReporteAusentismoDetalle, 0, len(results))
	for _, row := range results {
		datos = append(datos, domain.ReporteAusentismoDetalle{
			IDEmpleado:      row.IDEmpleado,
			NombreEmpleado:  row.NombreEmpleado,
			TipoIncapacidad: row.TipoIncapacidad,
			Entidad:         row.Entidad,
		})
	}

	topEmpleados, err := r.GetTopEmpleadosIncapacidades(ctx, filtros, 10)
	if err != nil {
		return nil, nil, err
	}

	return datos, topEmpleados, nil
}

func (r *ReportesRepository) GetCarteraReport(ctx context.Context, filtros domain.FiltrosReporte) (*domain.ReporteCartera, error) {
	reporte := &domain.ReporteCartera{}

	pagado, pendiente, err := r.SumValorCartera(ctx)
	if err != nil {
		return nil, err
	}

	reporte.TotalValorCobrado = formatCurrency(pagado)
	reporte.TotalValorPendiente = formatCurrency(pendiente)
	reporte.TotalValorCartera = formatCurrency(pagado + pendiente)

	pagPendientes, err := r.CountPagosPendientes(ctx)
	if err == nil {
		reporte.PagosPendientes = pagPendientes
	}

	pagVencidos, err := r.CountPagosVencidos(ctx, 0)
	if err == nil {
		reporte.PagosVencidos = pagVencidos
	}

	porEntidad, err := r.GetCarteraPorEntidad(ctx)
	if err == nil {
		reporte.PorEntidad = porEntidad
	}

	return reporte, nil
}

func (r *ReportesRepository) GetEntidadReport(ctx context.Context, entidadID uint64, filtros domain.FiltrosReporte) (*domain.ReporteEntidad, error) {
	var entidad struct {
		ID     uint64
		Nombre string
		Tipo   string
	}

	if err := r.db.WithContext(ctx).Table("entidad").
		Select("id_entidad, nombre, tipo").
		Where("id_entidad = ?", entidadID).
		Scan(&entidad).Error; err != nil {
		return nil, err
	}

	var stats struct {
		Total      int64
		Activas    int64
		Pagadas    int64
		Rechazadas int64
	}

	r.db.WithContext(ctx).Table("incapacidad i").
		Select("COUNT(*) as total, "+
			"SUM(CASE WHEN e.nombre NOT IN ('Pagada', 'Archivada', 'Cerrada') THEN 1 ELSE 0 END) as activas, "+
			"SUM(CASE WHEN e.nombre = 'Pagada' THEN 1 ELSE 0 END) as pagadas, "+
			"SUM(CASE WHEN e.nombre = 'Rechazada' THEN 1 ELSE 0 END) as rechazadas").
		Joins("JOIN estado_incapacidad e ON i.id_estado = e.id_estado").
		Where("i.id_entidad = ? AND i.is_deleted = false", entidadID).
		Scan(&stats)

	return &domain.ReporteEntidad{
		EntidadID:     entidad.ID,
		NombreEntidad: entidad.Nombre,
		Tipo:          entidad.Tipo,
		Estadisticas: domain.EntidadEstadisticas{
			TotalIncapacidades:      stats.Total,
			IncapacidadesActivas:    stats.Activas,
			IncapacidadesPagadas:    stats.Pagadas,
			IncapacidadesRechazadas: stats.Rechazadas,
		},
	}, nil
}

func (r *ReportesRepository) GetVencimientosReport(ctx context.Context, diasMinimos int) (*domain.ReporteVencimientos, error) {
	reporte := &domain.ReporteVencimientos{}

	docQuery := r.db.WithContext(ctx).Table("documento d").
		Select("d.id_incapacidad, u.nombre as nombre_empleado, td.nombre as tipo_documento, d.fecha_carga, ed.nombre as estado").
		Joins("JOIN incapacidad i ON d.id_incapacidad = i.id_incapacidad").
		Joins("JOIN usuario u ON i.id_usuario = u.id_usuario").
		Joins("JOIN tipo_documento td ON d.tipo_documento = td.nombre").
		Joins("JOIN estado_documento ed ON d.estado_documento = ed.nombre").
		Where("ed.nombre IN ('Pendiente', 'Incompleto')")

	var docResults []struct {
		IDIncapacidad  uint64
		NombreEmpleado string
		TipoDocumento  string
		FechaCarga     time.Time
		Estado         string
	}
	if err := docQuery.Scan(&docResults).Error; err == nil {
		for _, row := range docResults {
			reporte.AlertasDocumentos = append(reporte.AlertasDocumentos, domain.ReporteVencimientoDoc{
				IDIncapacidad:  row.IDIncapacidad,
				NombreEmpleado: row.NombreEmpleado,
				TipoDocumento:  row.TipoDocumento,
				Estado:         row.Estado,
			})
		}
	}

	pagoQuery := r.db.WithContext(ctx).Table("pago p").
		Select("p.id_incapacidad, p.id_pago, e.nombre as nombre_entidad, p.valor, p.fecha_pago as fecha_limite_pago, ep.nombre as estado").
		Joins("JOIN entidad e ON p.id_entidad = e.id_entidad").
		Joins("JOIN estado_pago ep ON p.estado_pago = ep.nombre").
		Where("ep.nombre NOT IN ('Pagado', 'Anulado', 'Conciliado')")

	var pagoResults []struct {
		IDIncapacidad   uint64
		IDPago          uint64
		NombreEntidad   string
		Valor           string
		FechaLimitePago time.Time
		Estado          string
	}
	if err := pagoQuery.Scan(&pagoResults).Error; err == nil {
		now := time.Now()
		for _, row := range pagoResults {
			diasVencido := 0
			if row.FechaLimitePago.Before(now) {
				diasVencido = int(now.Sub(row.FechaLimitePago).Hours() / 24)
			}

			reporte.AlertasPagos = append(reporte.AlertasPagos, domain.ReporteVencimientoPago{
				IDIncapacidad:   row.IDIncapacidad,
				IDPago:          row.IDPago,
				NombreEntidad:   row.NombreEntidad,
				Valor:           row.Valor,
				FechaLimitePago: row.FechaLimitePago,
				DiasVencido:     diasVencido,
				Estado:          row.Estado,
			})
		}
	}

	incQuery := r.db.WithContext(ctx).Table("incapacidad i").
		Select("i.id_incapacidad, u.nombre as nombre_empleado, i.fecha_inicio, e.nombre as estado").
		Joins("JOIN usuario u ON i.id_usuario = u.id_usuario").
		Joins("JOIN estado_incapacidad e ON i.id_estado = e.id_estado").
		Where("e.nombre NOT IN ('Pagada', 'Archivada', 'Cerrada') AND i.is_deleted = false AND i.fecha_fin IS NULL")

	var incResults []struct {
		IDIncapacidad  uint64
		NombreEmpleado string
		FechaInicio    time.Time
		Estado         string
	}
	if err := incQuery.Scan(&incResults).Error; err == nil {
		now := time.Now()
		for _, row := range incResults {
			reporte.AlertasIncapacidades = append(reporte.AlertasIncapacidades, domain.ReporteVencimientoINC{
				IDIncapacidad:     row.IDIncapacidad,
				NombreEmpleado:    row.NombreEmpleado,
				FechaInicio:       row.FechaInicio,
				DiasTranscurridos: int(now.Sub(row.FechaInicio).Hours() / 24),
				Estado:            row.Estado,
			})
		}
	}

	return reporte, nil
}

func (r *ReportesRepository) GetTotalIncapacidades(ctx context.Context, filtros domain.FiltrosReporte) (int64, int64, int64, error) {
	var results struct {
		Total    int64
		Activas  int64
		Cerradas int64
	}

	query := r.db.WithContext(ctx).Table("incapacidad i").
		Joins("JOIN estado_incapacidad e ON i.id_estado = e.id_estado").
		Where("i.is_deleted = false")

	if filtros.FechaInicio != nil {
		query = query.Where("i.fecha_inicio >= ?", *filtros.FechaInicio)
	}
	if filtros.FechaFin != nil {
		query = query.Where("i.fecha_inicio <= ?", *filtros.FechaFin)
	}
	if filtros.IDEntidad != nil {
		query = query.Where("i.id_entidad = ?", *filtros.IDEntidad)
	}

	query.Select("COUNT(*) as total, " +
		"SUM(CASE WHEN e.nombre NOT IN ('Pagada', 'Archivada', 'Cerrada') THEN 1 ELSE 0 END) as activas, " +
		"SUM(CASE WHEN e.nombre IN ('Pagada', 'Archivada', 'Cerrada') THEN 1 ELSE 0 END) as cerradas").
		Scan(&results)

	return results.Total, results.Activas, results.Cerradas, nil
}

func (r *ReportesRepository) GetTotalDiasPerdidos(ctx context.Context, filtros domain.FiltrosReporte) (int64, error) {
	var total int64
	query := r.db.WithContext(ctx).Table("incapacidad i").Where("i.is_deleted = false")

	if filtros.FechaInicio != nil {
		query = query.Where("i.fecha_inicio >= ?", *filtros.FechaInicio)
	}
	if filtros.FechaFin != nil {
		query = query.Where("i.fecha_inicio <= ?", *filtros.FechaFin)
	}

	query.Select("COALESCE(SUM(EXTRACT(DAY FROM (COALESCE(i.fecha_fin, CURRENT_DATE) - i.fecha_inicio))), 0)").Scan(&total)
	return total, nil
}

func (r *ReportesRepository) GetDiasPorTipo(ctx context.Context, filtros domain.FiltrosReporte) (map[string]int64, error) {
	type Result struct {
		Tipo string
		Dias int64
	}
	var results []Result

	query := r.db.WithContext(ctx).Table("incapacidad i").
		Select("ti.nombre, COALESCE(SUM(EXTRACT(DAY FROM (COALESCE(i.fecha_fin, CURRENT_DATE) - i.fecha_inicio))), 0) as dias").
		Joins("JOIN tipo_incapacidad ti ON i.id_tipo = ti.id_tipo").
		Where("i.is_deleted = false").
		Group("ti.nombre")

	if filtros.FechaInicio != nil {
		query = query.Where("i.fecha_inicio >= ?", *filtros.FechaInicio)
	}
	if filtros.FechaFin != nil {
		query = query.Where("i.fecha_inicio <= ?", *filtros.FechaFin)
	}

	if err := query.Scan(&results).Error; err != nil {
		return nil, err
	}

	diasPorTipo := make(map[string]int64)
	for _, row := range results {
		diasPorTipo[row.Tipo] = row.Dias
	}
	return diasPorTipo, nil
}

func (r *ReportesRepository) GetDiasPorEntidad(ctx context.Context, filtros domain.FiltrosReporte) (map[string]int64, error) {
	type Result struct {
		Entidad string
		Dias    int64
	}
	var results []Result

	query := r.db.WithContext(ctx).Table("incapacidad i").
		Select("en.nombre, COALESCE(SUM(EXTRACT(DAY FROM (COALESCE(i.fecha_fin, CURRENT_DATE) - i.fecha_inicio))), 0) as dias").
		Joins("JOIN entidad en ON i.id_entidad = en.id_entidad").
		Where("i.is_deleted = false").
		Group("en.nombre")

	if filtros.FechaInicio != nil {
		query = query.Where("i.fecha_inicio >= ?", *filtros.FechaInicio)
	}
	if filtros.FechaFin != nil {
		query = query.Where("i.fecha_inicio <= ?", *filtros.FechaFin)
	}

	if err := query.Scan(&results).Error; err != nil {
		return nil, err
	}

	diasPorEntidad := make(map[string]int64)
	for _, row := range results {
		diasPorEntidad[row.Entidad] = row.Dias
	}
	return diasPorEntidad, nil
}

func (r *ReportesRepository) GetTopEmpleadosIncapacidades(ctx context.Context, filtros domain.FiltrosReporte, limit int) ([]domain.ReporteTopIncapacidad, error) {
	type Result struct {
		IDEmpleado     uint64
		NombreEmpleado string
		TotalDias      int64
		CantidadINC    int64
	}
	var results []Result

	query := r.db.WithContext(ctx).Table("incapacidad i").
		Select("i.id_usuario, u.nombre, SUM(EXTRACT(DAY FROM (COALESCE(i.fecha_fin, CURRENT_DATE) - i.fecha_inicio))) as total_dias, COUNT(*) as cantidad_inc").
		Joins("JOIN usuario u ON i.id_usuario = u.id_usuario").
		Where("i.is_deleted = false").
		Group("i.id_usuario, u.nombre").
		Order("SUM(EXTRACT(DAY FROM (COALESCE(i.fecha_fin, CURRENT_DATE) - i.fecha_inicio))) DESC")

	if limit > 0 {
		query = query.Limit(limit)
	}

	if filtros.FechaInicio != nil {
		query = query.Where("i.fecha_inicio >= ?", *filtros.FechaInicio)
	}
	if filtros.FechaFin != nil {
		query = query.Where("i.fecha_inicio <= ?", *filtros.FechaFin)
	}

	if err := query.Scan(&results).Error; err != nil {
		return nil, err
	}

	top := make([]domain.ReporteTopIncapacidad, len(results))
	for i, row := range results {
		top[i] = domain.ReporteTopIncapacidad{
			IDEmpleado:     row.IDEmpleado,
			NombreEmpleado: row.NombreEmpleado,
			TotalDias:      row.TotalDias,
			CantidadINC:    row.CantidadINC,
		}
	}
	return top, nil
}

func (r *ReportesRepository) CountIncapacidadesActivas(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Table("incapacidad i").
		Joins("JOIN estado_incapacidad e ON i.id_estado = e.id_estado").
		Where("e.nombre NOT IN ('Pagada', 'Archivada', 'Cerrada') AND i.is_deleted = false").
		Count(&count).Error
	return count, err
}

func (r *ReportesRepository) CountPagosPendientes(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Table("pago").
		Where("estado_pago NOT IN ('Pagado', 'Anulado', 'Conciliado') AND is_deleted = false").
		Count(&count).Error
	return count, err
}

func (r *ReportesRepository) CountPagosVencidos(ctx context.Context, diasMinimos int) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Table("pago").
		Where("fecha_pago < CURRENT_DATE AND estado_pago NOT IN ('Pagado', 'Anulado', 'Conciliado') AND is_deleted = false").
		Count(&count).Error
	return count, err
}

func (r *ReportesRepository) SumValorCartera(ctx context.Context) (float64, float64, error) {
	type Result struct {
		Pagado    float64
		Pendiente float64
	}
	var result Result

	err := r.db.WithContext(ctx).Table("pago").
		Select("COALESCE(SUM(CASE WHEN estado_pago IN ('Pagado', 'Conciliado') THEN valor ELSE 0 END), 0) as pagado, " +
			"COALESCE(SUM(CASE WHEN estado_pago NOT IN ('Pagado', 'Anulado', 'Conciliado') THEN valor ELSE 0 END), 0) as pendiente").
		Where("is_deleted = false").
		Scan(&result).Error

	return result.Pagado, result.Pendiente, err
}

func (r *ReportesRepository) GetCarteraPorEntidad(ctx context.Context) ([]domain.ReporteCarteraEntidad, error) {
	type Result struct {
		NombreEntidad   string
		Pagado          float64
		Pendiente       float64
		CantidadINC     int64
		PagosPendientes int64
	}

	var results []Result

	err := r.db.WithContext(ctx).Table("pago p").
		Select("e.nombre, " +
			"COALESCE(SUM(CASE WHEN p.estado_pago IN ('Pagado', 'Conciliado') THEN p.valor ELSE 0 END), 0) as pagado, " +
			"COALESCE(SUM(CASE WHEN p.estado_pago NOT IN ('Pagado', 'Anulado', 'Conciliado') THEN p.valor ELSE 0 END), 0) as pendiente, " +
			"COUNT(DISTINCT p.id_incapacidad) as cantidad_inc, " +
			"SUM(CASE WHEN p.estado_pago NOT IN ('Pagado', 'Anulado', 'Conciliado') THEN 1 ELSE 0 END) as pagos_pendientes").
		Joins("JOIN entidad e ON p.id_entidad = e.id_entidad").
		Where("p.is_deleted = false").
		Group("e.nombre").
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	porEntidad := make([]domain.ReporteCarteraEntidad, len(results))
	for i, row := range results {
		porEntidad[i] = domain.ReporteCarteraEntidad{
			NombreEntidad:   row.NombreEntidad,
			ValorCobrado:    formatCurrency(row.Pagado),
			ValorPendiente:  formatCurrency(row.Pendiente),
			ValorTotal:      formatCurrency(row.Pagado + row.Pendiente),
			CantidadINC:     row.CantidadINC,
			PagosPendientes: row.PagosPendientes,
		}
	}
	return porEntidad, nil
}

func formatCurrency(value float64) string {
	return strings.TrimRight(strings.TrimRight(formatFloat(value), "0"), ".")
}

func formatFloat(f float64) string {
	return strings.Replace(strings.Replace(formatMoney(f), "$", "", 1), ",", "", 1)
}

func formatMoney(f float64) string {
	return fmt.Sprintf("%.2f", f)
}
