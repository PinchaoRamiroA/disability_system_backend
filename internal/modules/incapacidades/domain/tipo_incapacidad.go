package domain

import "time"

type TipoIncapacidad struct {
	IDTipo               uint64
	Nombre               string
	Origen               string
	DocumentosRequeridos []string
	CreatedAt            time.Time
	UpdatedAt            time.Time
}
