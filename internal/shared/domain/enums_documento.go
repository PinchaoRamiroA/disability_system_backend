package domain

type TipoDocumento string

const (
	TipoDocumentoCertificadoIncapacidad TipoDocumento = "Certificado de incapacidad"
	TipoDocumentoEpicrisis               TipoDocumento = "Epicrisis"
	TipoDocumentoFURIPS                 TipoDocumento = "FURIPS"
	TipoDocumentoHistoriaClinica       TipoDocumento = "Historia clínica"
	TipoDocumentoCertificadoNacidoVivo  TipoDocumento = "Certificado de nacido vivo"
	TipoDocumentoRegistroCivil           TipoDocumento = "Registro civil"
	TipoDocumentoDocumentoIdentidad      TipoDocumento = "Documento de identidad"
	TipoDocumentoSoporteAtencion         TipoDocumento = "Soporte de atención médica"
	TipoDocumentoConceptoRehabilitacion  TipoDocumento = "Concepto de rehabilitación"
	TipoDocumentoEvidenciaRadicacion      TipoDocumento = "Evidencia de radicación"
	TipoDocumentoSoportePago             TipoDocumento = "Soporte de pago"
	TipoDocumentoFormatoSeguimiento      TipoDocumento = "Formato de seguimiento"
)

func (t TipoDocumento) IsValid() bool {
	switch t {
	case TipoDocumentoCertificadoIncapacidad, TipoDocumentoEpicrisis,
		TipoDocumentoFURIPS, TipoDocumentoHistoriaClinica,
		TipoDocumentoCertificadoNacidoVivo, TipoDocumentoRegistroCivil,
		TipoDocumentoDocumentoIdentidad, TipoDocumentoSoporteAtencion,
		TipoDocumentoConceptoRehabilitacion, TipoDocumentoEvidenciaRadicacion,
		TipoDocumentoSoportePago, TipoDocumentoFormatoSeguimiento:
		return true
	}
	return false
}

type EstadoDocumento string

const (
	EstadoDocPendiente  EstadoDocumento = "Pendiente"
	EstadoDocValidado   EstadoDocumento = "Validado"
	EstadoDocRechazado   EstadoDocumento = "Rechazado"
	EstadoDocIncompleto  EstadoDocumento = "Incompleto"
	EstadoDocVencido     EstadoDocumento = "Vencido"
	EstadoDocArchivado   EstadoDocumento = "Archivado"
)

func (e EstadoDocumento) IsValid() bool {
	switch e {
	case EstadoDocPendiente, EstadoDocValidado, EstadoDocRechazado,
		EstadoDocIncompleto, EstadoDocVencido, EstadoDocArchivado:
		return true
	}
	return false
}