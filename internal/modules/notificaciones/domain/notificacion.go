package domain

import "time"

type Notificacion struct {
	IDNotificacion   uint64
	IDUsuario        uint64
	IDIncapacidad    *uint64
	TipoNotificacion string
	Mensaje          string
	Fecha            time.Time
	Leida            bool
	CreatedAt        time.Time
	UpdatedAt        time.Time
}
