package domain

type CanalRecepcion string

const (
	CanalRecepFisica    CanalRecepcion = "Recepción física"
	CanalRecepOficina   CanalRecepcion = "Oficina principal"
	CanalRecepLider     CanalRecepcion = "Líder comercial"
	CanalRecepCorreo    CanalRecepcion = "Correo electrónico"
	CanalRecepPortalEPS CanalRecepcion = "Portal EPS"
	CanalRecepPortalARL CanalRecepcion = "Portal ARL"
)

func (c CanalRecepcion) IsValid() bool {
	switch c {
	case CanalRecepFisica, CanalRecepOficina, CanalRecepLider,
		CanalRecepCorreo, CanalRecepPortalEPS, CanalRecepPortalARL:
		return true
	}
	return false
}

type CanalAtencionEntidad string

const (
	CanalEntWeb        CanalAtencionEntidad = "Portal web"
	CanalEntEmail      CanalAtencionEntidad = "Correo electrónico"
	CanalEntPresencial CanalAtencionEntidad = "Atención presencial"
	CanalEntTelefono   CanalAtencionEntidad = "Línea telefónica"
	CanalEntMesaAyuda  CanalAtencionEntidad = "Mesa de ayuda"
)

func (c CanalAtencionEntidad) IsValid() bool {
	switch c {
	case CanalEntWeb, CanalEntEmail, CanalEntPresencial,
		CanalEntTelefono, CanalEntMesaAyuda:
		return true
	}
	return false
}

type PeriodicidadReporte string

const (
	PeriodoDiario     PeriodicidadReporte = "Diario"
	PeriodoSemanal    PeriodicidadReporte = "Semanal"
	PeriodoMensual    PeriodicidadReporte = "Mensual"
	PeriodoTrimestral PeriodicidadReporte = "Trimestral"
	PeriodoAnual      PeriodicidadReporte = "Anual"
)

func (p PeriodicidadReporte) IsValid() bool {
	switch p {
	case PeriodoDiario, PeriodoSemanal, PeriodoMensual, PeriodoTrimestral, PeriodoAnual:
		return true
	}
	return false
}

func (p PeriodicidadReporte) Dias() int {
	switch p {
	case PeriodoDiario:
		return 1
	case PeriodoSemanal:
		return 7
	case PeriodoMensual:
		return 30
	case PeriodoTrimestral:
		return 90
	case PeriodoAnual:
		return 365
	}
	return 0
}