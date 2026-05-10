package models

import "time"

type SeguimientoCobroModel struct {
	IDSeguimiento uint64 `gorm:"primaryKey;autoIncrement;column:id_seguimiento"`

	IDIncapacidad uint64 `gorm:"not null;column:id_incapacidad;index"`

	IDTipoSeguimiento uint64 `gorm:"not null;column:id_tipo_seguimiento;index"`

	Descripcion *string `gorm:"type:text;column:descripcion"`

	Fecha time.Time `gorm:"autoCreateTime;column:fecha"`

	Resultado *string `gorm:"type:text;column:resultado"`

	GestionadoPor *uint64 `gorm:"column:gestionado_por;index"`

	CreatedAt time.Time `gorm:"autoCreateTime;column:created_at"`

	UpdatedAt time.Time `gorm:"autoUpdateTime;column:updated_at"`
}

func (SeguimientoCobroModel) TableName() string {
	return "seguimiento_cobro"
}