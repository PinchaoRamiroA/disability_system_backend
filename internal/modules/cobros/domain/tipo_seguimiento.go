package domain

import "time"

type TipoSeguimiento struct {
	IDTipoSeguimiento uint64
	Nombre            string
	Descripcion       *string
	CreatedAt         time.Time
	UpdatedAt         time.Time
}
