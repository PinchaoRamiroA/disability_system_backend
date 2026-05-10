package models

import (
	"encoding/json"
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
	Usuarios []UsuarioModel `gorm:"foreignKey:IDRol"`
}

func (RolModel) TableName() string {
	return "rol"
}

func (r *RolModel) GetPermisos() []string {
	if r.Permisos == nil {
		return nil
	}
	var permisos []string
	if err := json.Unmarshal(r.Permisos, &permisos); err != nil {
		return nil
	}
	return permisos
}

func (r *RolModel) SetPermisos(permisos []string) error {
	if permisos == nil {
		r.Permisos = nil
		return nil
	}
	data, err := json.Marshal(permisos)
	if err != nil {
		return err
	}
	r.Permisos = data
	return nil
}