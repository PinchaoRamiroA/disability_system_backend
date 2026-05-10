package domain

type TipoHistorial string

const (
	TipoHistCreacion          TipoHistorial = "Creación"
	TipoHistActualizacion     TipoHistorial = "Actualización"
	TipoHistCambioEstado      TipoHistorial = "Cambio de estado"
	TipoHistValidacionDoc      TipoHistorial = "Validación documental"
	TipoHistRechazoDoc         TipoHistorial = "Rechazo documental"
	TipoHistTranscripcion      TipoHistorial = "Transcripción"
	TipoHistRadicacion         TipoHistorial = "Radicación"
	TipoHistSeguimiento        TipoHistorial = "Seguimiento"
	TipoHistPago               TipoHistorial = "Pago"
	TipoHistConciliacion       TipoHistorial = "Conciliación"
	TipoHistObservacion        TipoHistorial = "Observación"
	TipoHistArchivo            TipoHistorial = "Archivo"
	TipoHistNotificacion       TipoHistorial = "Notificación"
	TipoHistCobroPersuasivo    TipoHistorial = "Cobro persuasivo"
	TipoHistCobroJuridico      TipoHistorial = "Cobro jurídico"
)

func (t TipoHistorial) IsValid() bool {
	switch t {
	case TipoHistCreacion, TipoHistActualizacion, TipoHistCambioEstado,
		TipoHistValidacionDoc, TipoHistRechazoDoc, TipoHistTranscripcion,
		TipoHistRadicacion, TipoHistSeguimiento, TipoHistPago,
		TipoHistConciliacion, TipoHistObservacion, TipoHistArchivo,
		TipoHistNotificacion, TipoHistCobroPersuasivo, TipoHistCobroJuridico:
		return true
	}
	return false
}