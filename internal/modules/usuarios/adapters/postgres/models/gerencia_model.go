package models

import "time"

type GerenciaModel struct {
	IDUsuario uint64 `gorm:"primaryKey;column:id_usuario"`

	PuestoTrabajo string `gorm:"type:varchar(150);not null;column:puesto_trabajo"`

	Usuario   UsuarioModel `gorm:"foreignKey:IDUsuario"`
	CreatedAt time.Time    `gorm:"autoCreateTime;column:created_at"`

	UpdatedAt time.Time `gorm:"autoUpdateTime;column:updated_at"`
}

func (GerenciaModel) TableName() string {
	return "gerencia"
}
