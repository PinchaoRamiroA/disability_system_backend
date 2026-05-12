package domain

import "time"

type TipoNotificacion struct {
	IDTipoNotificacion uint64
	Nombre             string
	Descripcion        *string
	CreatedAt          time.Time
	UpdatedAt          time.Time
}
