package domain

import "time"

type TipoIncapacidad struct {
	IDTipo               uint64
	Nombre               string
	DocumentosRequeridos []string
	CreatedAt            time.Time
	UpdatedAt            time.Time
}
