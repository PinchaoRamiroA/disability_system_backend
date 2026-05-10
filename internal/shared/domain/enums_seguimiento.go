package domain

type TipoSeguimiento string

const (
	TipoSegNormal          TipoSeguimiento = "Normal"
	TipoSegPersuasivo      TipoSeguimiento = "Persuasivo"
	TipoSegJuridico        TipoSeguimiento = "Jurídico"
	TipoSegPreventivo       TipoSeguimiento = "Preventivo"
	TipoSegCobroAdministrativo TipoSeguimiento = "Cobro administrativo"
)

func (t TipoSeguimiento) IsValid() bool {
	switch t {
	case TipoSegNormal, TipoSegPersuasivo, TipoSegJuridico,
		TipoSegPreventivo, TipoSegCobroAdministrativo:
		return true
	}
	return false
}

type ResultadoSeguimiento string

const (
	ResultadoPendienteRespuesta    ResultadoSeguimiento = "Pendiente respuesta"
	ResultadoEnRevision            ResultadoSeguimiento = "En revisión"
	ResultadoAprobado             ResultadoSeguimiento = "Aprobado"
	ResultadoRechazado            ResultadoSeguimiento = "Rechazado"
	ResultadoPagoProgramado        ResultadoSeguimiento = "Pago programado"
	ResultadoPagoRealizado        ResultadoSeguimiento = "Pago realizado"
	ResultadoEscaladoJuridica     ResultadoSeguimiento = "Escalado a jurídica"
	ResultadoRequiereSubsanacion  ResultadoSeguimiento = "Requiere subsanación"
	ResultadoSinRespuesta         ResultadoSeguimiento = "Sin respuesta"
)

func (r ResultadoSeguimiento) IsValid() bool {
	switch r {
	case ResultadoPendienteRespuesta, ResultadoEnRevision, ResultadoAprobado,
		ResultadoRechazado, ResultadoPagoProgramado, ResultadoPagoRealizado,
		ResultadoEscaladoJuridica, ResultadoRequiereSubsanacion, ResultadoSinRespuesta:
		return true
	}
	return false
}