package domain

import "time"

type Rol struct {
	ID        uint64
	Nombre    string
	Permisos  map[string]interface{}
	IsDeleted bool
	CreatedAt time.Time
	UpdatedAt time.Time
}