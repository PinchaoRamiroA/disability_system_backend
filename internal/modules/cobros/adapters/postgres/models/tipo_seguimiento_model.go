package models

import "time"

type TipoSeguimientoModel struct {
	IDTipoSeguimiento uint64 `gorm:"primaryKey;autoIncrement;column:id_tipo_seguimiento"`

	Nombre string `gorm:"type:varchar(50);uniqueIndex;not null;column:nombre"`

	Descripcion *string `gorm:"type:text;column:descripcion"`

	CreatedAt time.Time `gorm:"autoCreateTime;column:created_at"`

	UpdatedAt time.Time `gorm:"autoUpdateTime;column:updated_at"`
}

func (TipoSeguimientoModel) TableName() string {
	return "tipo_seguimiento"
}
