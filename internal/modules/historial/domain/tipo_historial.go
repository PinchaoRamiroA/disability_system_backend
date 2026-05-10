package domain

import "time"

type TipoHistorial struct {
	IDTipoHistorial   uint64
	Nombre             string
	Descripcion        *string
	CreatedAt          time.Time
	UpdatedAt          time.Time
}