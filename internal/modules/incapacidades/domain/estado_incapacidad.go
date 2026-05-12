package domain

import "time"

type EstadoIncapacidad struct {
	IDEstado          uint64
	Nombre            string
	Descripcion       *string
	PermiteTransicion bool
	CreatedAt         time.Time
	UpdatedAt         time.Time
}
