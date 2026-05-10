package domain

type TipoNotificacion string

const (
	TipoNotifAlerta             TipoNotificacion = "Alerta"
	TipoNotifRecordatorio       TipoNotificacion = "Recordatorio"
	TipoNotifDocFaltante        TipoNotificacion = "Documento faltante"
	TipoNotifVencimientoProximo TipoNotificacion = "Vencimiento próximo"
	TipoNotifIncapacidadRechazada TipoNotificacion = "Incapacidad rechazada"
	TipoNotifPagoRecibido       TipoNotificacion = "Pago recibido"
	TipoNotifPagoPendiente      TipoNotificacion = "Pago pendiente"
	TipoNotifConciliacionPendiente TipoNotificacion = "Conciliación pendiente"
	TipoNotifCobroJuridico      TipoNotificacion = "Cobro jurídico"
	TipoNotifSeguimientoRequerido TipoNotificacion = "Seguimiento requerido"
)

func (t TipoNotificacion) IsValid() bool {
	switch t {
	case TipoNotifAlerta, TipoNotifRecordatorio, TipoNotifDocFaltante,
		TipoNotifVencimientoProximo, TipoNotifIncapacidadRechazada,
		TipoNotifPagoRecibido, TipoNotifPagoPendiente,
		TipoNotifConciliacionPendiente, TipoNotifCobroJuridico,
		TipoNotifSeguimientoRequerido:
		return true
	}
	return false
}