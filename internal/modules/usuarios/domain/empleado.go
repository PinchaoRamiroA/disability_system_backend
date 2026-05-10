package domain

import "time"

type Empleado struct {
	IDUsuario      uint64
	PuestoTrabajo  string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}