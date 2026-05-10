package domain

import "time"

type Historial struct {
	IDHistorial      uint64
	IDIncapacidad     uint64
	IDTipoHistorial   uint64
	Descripcion       string
	Fecha             time.Time
	GestorID          *uint64
	CreatedAt         time.Time
	UpdatedAt         time.Time
}