package domain

import (
	"time"

	"github.com/shopspring/decimal"
)

type Pago struct {
	IDPago          uint64
	IDIncapacidad   uint64
	IDEntidad       uint64
	TipoPago        string
	EstadoPago      string
	Descripcion     *string
	Valor           decimal.Decimal
	FechaPago       time.Time
	PeriodoContable *string
	Conciliado      bool
	RegistradoPor   *uint64
	CreatedAt       time.Time
	UpdatedAt       time.Time
	IsDeleted       bool
	NombreEntidad   string
}
