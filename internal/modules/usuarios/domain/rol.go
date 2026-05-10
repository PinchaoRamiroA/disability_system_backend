package domain

import "time"

type Rol struct {
	ID        uint64
	Nombre    string
	Permisos  []string
	IsDeleted bool
	CreatedAt time.Time
	UpdatedAt time.Time
}