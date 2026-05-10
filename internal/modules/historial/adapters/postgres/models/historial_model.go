package models

import "time"

type HistorialModel struct {
	IDHistorial uint64 `gorm:"primaryKey;autoIncrement;column:id_historial"`

	IDIncapacidad uint64 `gorm:"not null;column:id_incapacidad;index"`

	IDTipoHistorial uint64 `gorm:"not null;column:id_tipo_historial;index"`

	Descripcion string `gorm:"type:text;not null;column:descripcion"`

	Fecha time.Time `gorm:"autoCreateTime;column:fecha"`

	GestorID *uint64 `gorm:"column:gestor_id;index"`

	CreatedAt time.Time `gorm:"autoCreateTime;column:created_at"`

	UpdatedAt time.Time `gorm:"autoUpdateTime;column:updated_at"`
}

func (HistorialModel) TableName() string {
	return "historial"
}