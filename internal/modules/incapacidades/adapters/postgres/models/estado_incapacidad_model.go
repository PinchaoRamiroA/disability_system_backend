package models

import (
	"time"

	"gorm.io/datatypes"
)

type TipoIncapacidadModel struct {
	IDTipo               uint64         `gorm:"primaryKey;autoIncrement;column:id_tipo"`
	Nombre               string         `gorm:"type:varchar(100);uniqueIndex;not null;column:nombre"`
	DocumentosRequeridos datatypes.JSON `gorm:"type:jsonb;column:documentos_requeridos"`
	CreatedAt            time.Time      `gorm:"autoCreateTime;column:created_at"`
	UpdatedAt            time.Time      `gorm:"autoUpdateTime;column:updated_at"`
}

func (TipoIncapacidadModel) TableName() string {
	return "tipo_incapacidad"
}

type EstadoIncapacidadModel struct {
	IDEstado          uint64    `gorm:"primaryKey;autoIncrement;column:id_estado"`
	Nombre            string    `gorm:"type:varchar(100);uniqueIndex;not null;column:nombre"`
	Descripcion       *string   `gorm:"type:text;column:descripcion"`
	PermiteTransicion bool      `gorm:"default:true;column:permite_transicion"`
	CreatedAt         time.Time `gorm:"autoCreateTime;column:created_at"`
	UpdatedAt         time.Time `gorm:"autoUpdateTime;column:updated_at"`
}

func (EstadoIncapacidadModel) TableName() string {
	return "estado_incapacidad"
}
