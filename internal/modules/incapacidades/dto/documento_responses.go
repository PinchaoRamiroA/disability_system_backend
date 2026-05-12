package dto

type DocumentoResponse struct {
	IDDocumento     uint64  `json:"id_documento"`
	IDIncapacidad   uint64  `json:"id_incapacidad"`
	Nombre          string  `json:"nombre"`
	Tipo            string  `json:"tipo"`
	URL             string  `json:"url"`
	Formato         string  `json:"formato"`
	Estado          string  `json:"estado"`
	Comentario      *string `json:"comentario,omitempty"`
	FechaCarga      string  `json:"fecha_carga"`
	ValidadoPor     *uint64 `json:"validado_por,omitempty"`
	FechaValidacion *string `json:"fecha_validacion,omitempty"`
	CreatedAt       string  `json:"created_at"`
}