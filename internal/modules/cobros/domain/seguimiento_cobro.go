package domain

import "time"

type SeguimientoCobro struct {
	IDSeguimiento     uint64
	IDIncapacidad     uint64
	IDTipoSeguimiento uint64
	Descripcion       *string
	Fecha             time.Time
	Resultado         *string
	GestionadoPor     *uint64
	CreatedAt         time.Time
	UpdatedAt         time.Time
}
