package domain

import "time"

type TipoDocumento struct {
	IDTipoDocumento uint64
	Nombre          string
	Codigo          string
	Descripcion     *string
	Requerido       bool
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

type EstadoDocumento struct {
	IDEstadoDocumento uint64
	Nombre            string
	Descripcion       *string
	Color             *string
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

type TipoEntidad struct {
	IDTipoEntidad uint64
	Nombre        string
	Descripcion   *string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type CanalRecepcion struct {
	IDCanalRecepcion uint64
	Nombre           string
	Descripcion      *string
	Activo           bool
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

type CanalAtencionEntidad struct {
	IDCanalAtencion uint64
	Nombre          string
	Descripcion     *string
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

type TipoPago struct {
	IDTipoPago  uint64
	Nombre      string
	Descripcion *string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type EstadoPago struct {
	IDEstadoPago uint64
	Nombre       string
	Descripcion  *string
	Color        *string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type PeriodicidadReporte struct {
	IDPeriodicidad uint64
	Nombre         string
	Dias           *int
	Descripcion    *string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
