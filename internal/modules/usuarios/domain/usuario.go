package domain

import "time"

type Usuario struct {
	ID                uint64
	IDRol             uint64
	Nombre            string
	Correo            string
	NumeroCelular     *string
	Direccion         *string
	PasswordHash      string
	NumeroDocumento   string
	NumeroAcudiente   *string
	Estado            bool
	IsDeleted         bool
	CreatedAt         time.Time
	UpdatedAt         time.Time
}