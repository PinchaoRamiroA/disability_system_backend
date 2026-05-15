package domain

import "time"

type Auditoria struct {
	ID             uint64
	IDUsuario      *uint64
	UsuarioNombre  string // Opcional, para lectura
	IDIncapacidad  *uint64
	TipoAccion     string
	Modulo         string
	Descripcion    string
	CambioAnterior *string
	CambioNuevo    *string
	CreatedAt      time.Time
}
