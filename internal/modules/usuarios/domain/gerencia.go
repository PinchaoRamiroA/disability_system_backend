package domain

import "time"

type Gerencia struct {
	IDUsuario      uint64
	PuestoTrabajo  string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}