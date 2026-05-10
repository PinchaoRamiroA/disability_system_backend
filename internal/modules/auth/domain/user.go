package domain

import "time"

type User struct {
	ID              uint64
	IDRol           uint64
	Nombre          string
	Correo          string
	PasswordHash    string
	NumeroDocumento string
	Estado          bool
	IsDeleted       bool
	CreatedAt       time.Time
}
