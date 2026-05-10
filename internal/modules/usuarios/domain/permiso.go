package domain

import "time"

type Permiso struct {
	IDPermiso    uint64
	Nombre       string
	Descripcion  *string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}