package dto

type CrearSeguimientoRequest struct {
	IDIncapacidad   uint64  `json:"id_incapacidad" binding:"required"`
	TipoSeguimiento string  `json:"tipo_seguimiento" binding:"required"`
	Descripcion     *string `json:"descripcion,omitempty"`
	Resultado       *string `json:"resultado,omitempty"`
}

type ActualizarSeguimientoRequest struct {
	TipoSeguimiento *string `json:"tipo_seguimiento,omitempty"`
	Descripcion     *string `json:"descripcion,omitempty"`
	Resultado       *string `json:"resultado,omitempty"`
}

type ListarSeguimientosQuery struct {
	IDIncapacidad   *uint64 `form:"id_incapacidad"`
	TipoSeguimiento string  `form:"tipo_seguimiento"`
	Page            int     `form:"page"`
	Limit           int     `form:"limit"`
}
