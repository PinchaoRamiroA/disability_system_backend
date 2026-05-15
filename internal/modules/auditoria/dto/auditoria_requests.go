package dto

type CrearAuditoriaRequest struct {
	IDUsuario      *uint64 `json:"id_usuario"`
	IDIncapacidad  *uint64 `json:"id_incapacidad"`
	TipoAccion     string  `json:"tipo_accion" binding:"required"`
	Modulo         string  `json:"modulo" binding:"required"`
	Descripcion    string  `json:"descripcion" binding:"required"`
	CambioAnterior *string `json:"cambio_anterior"`
	CambioNuevo    *string `json:"cambio_nuevo"`
}

type ListarAuditoriaQuery struct {
	IDUsuario     *uint64 `form:"id_usuario"`
	IDIncapacidad *uint64 `form:"id_incapacidad"`
	TipoAccion    string  `form:"tipo_accion"`
	Modulo        string  `form:"modulo"`
	FechaInicio   string  `form:"fecha_inicio"`
	FechaFin      string  `form:"fecha_fin"`
	Page          int     `form:"page,default=1" binding:"gte=1"`
	Limit         int     `form:"limit,default=20" binding:"gte=1,lte=100"`
}
