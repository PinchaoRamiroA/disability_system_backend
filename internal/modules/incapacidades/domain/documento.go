package domain

import "time"

type Documento struct {
	IDDocumento     uint64
	IDIncapacidad   uint64
	Nombre          string
	Tipo            string
	URL             string
	Formato         string
	Estado          string
	Comentario      *string
	FechaCarga      time.Time
	ValidadoPor     *uint64
	FechaValidacion *time.Time
	IsDeleted       bool
	CreatedAt       time.Time
	UpdatedAt       time.Time
}
