package dto

type HistorialRequest struct {
	IDIncapacidad uint64 `json:"id_incapacidad" binding:"required"`
	IDTipoHistorial uint64 `json:"id_tipo_historial" binding:"required"`
	Descripcion   string `json:"descripcion" binding:"required"`
}

type HistorialResponse struct {
	IDHistorial    uint64  `json:"id_historial"`
	IDIncapacidad  uint64  `json:"id_incapacidad"`
	IDTipoHistorial uint64 `json:"id_tipo_historial"`
	NombreTipo     string  `json:"nombre_tipo,omitempty"`
	Descripcion    string  `json:"descripcion"`
	Fecha          string  `json:"fecha"`
	GestorID       *uint64 `json:"gestor_id,omitempty"`
}

type ListarHistorialQuery struct {
	IDIncapacidad uint64 `form:"id_incapacidad"`
	IDTipoHistorial uint64 `form:"id_tipo_historial"`
	Page          int    `form:"page"`
	Limit         int    `form:"limit"`
}