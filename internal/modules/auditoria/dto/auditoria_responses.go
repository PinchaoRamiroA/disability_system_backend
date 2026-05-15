package dto

import "time"

type AuditoriaResponse struct {
	IDAuditoria    uint64    `json:"id_auditoria"`
	IDUsuario      *uint64   `json:"id_usuario"`
	UsuarioNombre  string    `json:"usuario_nombre"`
	IDIncapacidad  *uint64   `json:"id_incapacidad"`
	TipoAccion     string    `json:"tipo_accion"`
	Modulo         string    `json:"modulo"`
	Descripcion    string    `json:"descripcion"`
	CambioAnterior *string   `json:"cambio_anterior"`
	CambioNuevo    *string   `json:"cambio_nuevo"`
	CreatedAt      time.Time `json:"created_at"`
}
