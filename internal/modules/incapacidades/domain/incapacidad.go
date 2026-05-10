package domain

import "time"

type Incapacidad struct {
	IDIncapacidad      uint64
	IDUsuario          uint64
	IDEstado           uint64
	IDTipo             uint64
	IDEntidad          uint64
	CanalRecepcion     string
	Titulo             string
	FechaInicio        time.Time
	FechaFin           *time.Time
	Origen             string
	FechaRadicacion    *time.Time
	FechaPago          *time.Time
	Observaciones      *string
	CreatedBy          *uint64
	CreatedAt          time.Time
	UpdatedAt          time.Time
	IsDeleted          bool
}