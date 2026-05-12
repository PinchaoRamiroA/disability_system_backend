package dto

type IncapacidadResponse struct {
	IDIncapacidad   uint64                     `json:"id_incapacidad"`
	IDUsuario       uint64                     `json:"id_usuario"`
	CanalRecepcion  string                     `json:"canal_recepcion,omitempty"`
	Titulo          string                     `json:"titulo"`
	FechaInicio     string                     `json:"fecha_inicio"`
	FechaFin        *string                    `json:"fecha_fin,omitempty"`
	Origen          string                     `json:"origen"`
	FechaRadicacion *string                    `json:"fecha_radicacion,omitempty"`
	FechaPago       *string                    `json:"fecha_pago,omitempty"`
	Observaciones   *string                    `json:"observaciones,omitempty"`
	Estado          *EstadoIncapacidadResponse `json:"estado,omitempty"`
	Tipo            *TipoIncapacidadResponse   `json:"tipo,omitempty"`
	Entidad         *EntidadResponse           `json:"entidad,omitempty"`
	CreatedBy       *uint64                    `json:"created_by,omitempty"`
	CreatedAt       string                     `json:"created_at"`
	UpdatedAt       string                     `json:"updated_at"`
}

type EstadoIncapacidadResponse struct {
	IDEstado          uint64  `json:"id_estado"`
	Nombre            string  `json:"nombre"`
	Descripcion       *string `json:"descripcion,omitempty"`
	PermiteTransicion bool    `json:"permite_transicion"`
}

type TipoIncapacidadResponse struct {
	IDTipo               uint64   `json:"id_tipo"`
	Nombre               string   `json:"nombre"`
	DocumentosRequeridos []string `json:"documentos_requeridos,omitempty"`
}

type EntidadResponse struct {
	IDEntidad              uint64   `json:"id_entidad"`
	Nombre                 string   `json:"nombre"`
	Tipo                   string   `json:"tipo"`
	PlazoTranscripcionDias *int     `json:"plazo_transcripcion_dias,omitempty"`
	TiempoMaximoPagoDias   *int     `json:"tiempo_maximo_pago_dias,omitempty"`
	CanalAtencion          *string  `json:"canal_atencion,omitempty"`
	CanalesAtencion        []string `json:"canales_atencion,omitempty"`
	RequiereTranscripcion  bool     `json:"requiere_transcripcion"`
}

type EstadoDocumentoResponse struct {
	IDEstadoDocumento uint64  `json:"id_estado_documento"`
	Nombre            string  `json:"nombre"`
	Descripcion       *string `json:"descripcion,omitempty"`
	Color             *string `json:"color,omitempty"`
}

type TipoDocumentoResponse struct {
	IDTipoDocumento uint64  `json:"id_tipo_documento"`
	Nombre          string  `json:"nombre"`
	Descripcion     *string `json:"descripcion,omitempty"`
	Requerido       bool    `json:"requerido"`
}

type TipoPagoResponse struct {
	IDTipoPago  uint64  `json:"id_tipo_pago"`
	Nombre      string  `json:"nombre"`
	Descripcion *string `json:"descripcion,omitempty"`
}
