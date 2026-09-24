package dto

type TranscribirIncapacidadRequest struct {
	ObservacionesTranscripcion *string `json:"observaciones_transcripcion,omitempty"`
	NumeroRadicado             *string `json:"numero_radicado,omitempty"`
	FechaTranscripcion         *string `json:"fecha_transcripcion,omitempty"`
	Observaciones              *string `json:"observaciones,omitempty"`
}

type MarcarTranscripcionRequest struct {
	Estado              string `json:"estado"`
	EstadoTranscripcion string `json:"estado_transcripcion,omitempty"`
}

type ListarTranscripcionesQuery struct {
	Estado string `form:"estado"`
	Page   int    `form:"page"`
	Limit  int    `form:"limit"`
}
