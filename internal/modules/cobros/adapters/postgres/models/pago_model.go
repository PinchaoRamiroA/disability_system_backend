package models

import (
	"time"

	"github.com/shopspring/decimal"
)

type PagoModel struct {
	IDPago uint64 `gorm:"primaryKey;autoIncrement;column:id_pago"`

	IDIncapacidad uint64 `gorm:"not null;column:id_incapacidad;index"`

	IDEntidad uint64 `gorm:"not null;column:id_entidad;index"`

	Descripcion *string `gorm:"type:text;column:descripcion"`

	Valor decimal.Decimal `gorm:"type:numeric(14,2);not null;column:valor"`

	FechaPago time.Time `gorm:"type:date;not null;column:fecha_pago"`

	PeriodoContable *string `gorm:"type:varchar(20);column:periodo_contable"`

	Conciliado bool `gorm:"default:false;column:conciliado"`

	RegistradoPor *uint64 `gorm:"column:registrado_por;index"`

	CreatedAt time.Time `gorm:"autoCreateTime;column:created_at"`

	UpdatedAt time.Time `gorm:"autoUpdateTime;column:updated_at"`

	IsDeleted bool `gorm:"default:false;column:is_deleted"`
}

func (PagoModel) TableName() string {
	return "pago"
}