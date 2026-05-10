package models

import (
	"time"
)

type PermisoModel struct {
	IDPermiso uint64 `gorm:"primaryKey;autoIncrement;column:id_permiso"`
	Nombre string `gorm:"type:varchar(100);uniqueIndex;not null;column:nombre"`
	Descripcion *string `gorm:"type:text;column:descripcion"`
	CreatedAt time.Time `gorm:"autoCreateTime;column:created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime;column:updated_at"`
}

func (PermisoModel) TableName() string {
	return "permisos"
}