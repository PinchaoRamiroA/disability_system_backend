package models

import (
	"time"
)

type DocumentoModel struct {
	IDDocumento uint64 `gorm:"primaryKey;autoIncrement;column:id_documento"`
	IDIncapacidad uint64 `gorm:"not null;column:id_incapacidad;index"`
	Nombre string `gorm:"type:varchar(255);not null;column:nombre"`
	Tipo string `gorm:"type:varchar(100);not null;column:tipo"`
	URL string `gorm:"type:text;not null;column:url"`
	Formato string `gorm:"type:varchar(20);not null;column:formato"`
	Estado string `gorm:"type:varchar(50);not null;column:estado"`
	Comentario *string `gorm:"type:text;column:comentario"`
	FechaCarga time.Time `gorm:"autoCreateTime;column:fecha_carga"`
	ValidadoPor *uint64 `gorm:"column:validado_por;index"`
	FechaValidacion *time.Time `gorm:"column:fecha_validacion"`
	IsDeleted bool `gorm:"default:false;column:is_deleted"`
	CreatedAt time.Time `gorm:"autoCreateTime;column:created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime;column:updated_at"`
}

func (DocumentoModel) TableName() string {
	return "documento"
}