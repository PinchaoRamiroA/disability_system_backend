package dto

type TranscribirIncapacidadRequest struct {
	ObservacionesTranscripcion *string `json:"observaciones_transcripcion,omitempty"`
}

type MarcarTranscripcionRequest struct {
	Estado string  `json:"estado" binding:"required,oneof=en_proceso completado"`
}

type ListarTranscripcionesQuery struct {
	Estado string `form:"estado"`
	Page   int    `form:"page"`
	Limit  int    `form:"limit"`
}
