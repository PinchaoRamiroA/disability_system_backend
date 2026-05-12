package dto

type CrearNotificacionRequest struct {
	IDUsuario        uint64  `json:"id_usuario" binding:"required"`
	IDIncapacidad    *uint64 `json:"id_incapacidad,omitempty"`
	TipoNotificacion string  `json:"tipo_notificacion" binding:"required"`
	Mensaje          string  `json:"mensaje" binding:"required"`
}

type ListarNotificacionesQuery struct {
	IDUsuario        *uint64 `form:"id_usuario"`
	IDIncapacidad    *uint64 `form:"id_incapacidad"`
	TipoNotificacion string  `form:"tipo_notificacion"`
	Leida            *bool   `form:"leida"`
	Page             int     `form:"page"`
	Limit            int     `form:"limit"`
}
