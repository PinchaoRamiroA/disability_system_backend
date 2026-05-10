package domain

import "time"

type GestionHumana struct {
	IDUsuario      uint64
	CreatedAt      time.Time
	UpdatedAt      time.Time
}