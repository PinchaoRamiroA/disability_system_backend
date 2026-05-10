package models

import "time"

type EmpleadoModel struct {
	IDUsuario uint64 `gorm:"primaryKey;column:id_usuario"`

	PuestoTrabajo string `gorm:"type:varchar(150);not null;column:puesto_trabajo"`

	CreatedAt time.Time `gorm:"autoCreateTime;column:created_at"`

	UpdatedAt time.Time `gorm:"autoUpdateTime;column:updated_at"`

	// Relaciones

	Usuario UsuarioModel `gorm:"foreignKey:IDUsuario"`
}

func (EmpleadoModel) TableName() string {
	return "empleado"
}
