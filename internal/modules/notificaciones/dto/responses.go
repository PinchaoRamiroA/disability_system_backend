package dto

type NotificacionResponse struct {
	IDNotificacion   uint64  `json:"id_notificacion"`
	IDUsuario        uint64  `json:"id_usuario"`
	IDIncapacidad    *uint64 `json:"id_incapacidad,omitempty"`
	TipoNotificacion string  `json:"tipo_notificacion"`
	Mensaje          string  `json:"mensaje"`
	Fecha            string  `json:"fecha"`
	Leida            bool    `json:"leida"`
	CreatedAt        string  `json:"created_at"`
	UpdatedAt        string  `json:"updated_at"`
}

type ConteoNoLeidasResponse struct {
	Total int64 `json:"total"`
}
