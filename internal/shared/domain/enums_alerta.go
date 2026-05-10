package domain

type TipoAlerta string

const (
	TipoAlertaVencimientoTranscripcion  TipoAlerta = "Vencimiento de transcripción"
	TipoAlertaDocIncompleta             TipoAlerta = "Documentación incompleta"
	TipoAlertaPagoRetrasado             TipoAlerta = "Pago retrasado"
	TipoAlertaSeguimientoPendiente      TipoAlerta = "Seguimiento pendiente"
	TipoAlertaIncapacidad90Dias         TipoAlerta = "Incapacidad superior a 90 días"
	TipoAlertaIncapacidad120Dias         TipoAlerta = "Incapacidad superior a 120 días"
	TipoAlertaIncapacidad150Dias         TipoAlerta = "Incapacidad superior a 150 días"
	TipoAlertaIncapacidad180Dias         TipoAlerta = "Incapacidad superior a 180 días"
	TipoAlertaConciliacionPendiente      TipoAlerta = "Conciliación pendiente"
	TipoAlertaCobroJuridicoRequerido     TipoAlerta = "Cobro jurídico requerido"
)

func (t TipoAlerta) IsValid() bool {
	switch t {
	case TipoAlertaVencimientoTranscripcion, TipoAlertaDocIncompleta,
		TipoAlertaPagoRetrasado, TipoAlertaSeguimientoPendiente,
		TipoAlertaIncapacidad90Dias, TipoAlertaIncapacidad120Dias,
		TipoAlertaIncapacidad150Dias, TipoAlertaIncapacidad180Dias,
		TipoAlertaConciliacionPendiente, TipoAlertaCobroJuridicoRequerido:
		return true
	}
	return false
}

func (t TipoAlerta) DiasAnelacion() int {
	switch t {
	case TipoAlertaVencimientoTranscripcion:
		return 2
	case TipoAlertaPagoRetrasado:
		return 5
	case TipoAlertaSeguimientoPendiente:
		return 3
	case TipoAlertaIncapacidad90Dias:
		return 90
	case TipoAlertaIncapacidad120Dias:
		return 120
	case TipoAlertaIncapacidad150Dias:
		return 150
	case TipoAlertaIncapacidad180Dias:
		return 180
	case TipoAlertaConciliacionPendiente:
		return 10
	}
	return 0
}