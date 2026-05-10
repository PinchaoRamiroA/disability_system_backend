package models

import (
	"time"

	"gorm.io/datatypes"
)

type RolModel struct {
	IDRol uint64 `gorm:"primaryKey;autoIncrement;column:id_rol"`

	Nombre string `gorm:"type:varchar(100);uniqueIndex;not null;column:nombre"`

	Permisos datatypes.JSON `gorm:"type:jsonb;not null;column:permisos"`

	CreatedAt time.Time `gorm:"autoCreateTime;column:created_at"`

	UpdatedAt time.Time `gorm:"autoUpdateTime;column:updated_at"`

	IsDeleted bool `gorm:"default:false;column:is_deleted"`

	// Relaciones

	Usuarios []UsuarioModel `gorm:"foreignKey:IDRol"`
}

func (RolModel) TableName() string {
	return "rol"
}
