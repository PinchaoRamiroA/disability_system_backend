package models

import "time"

type GestionHumanaModel struct {
	IDUsuario uint64 `gorm:"primaryKey;column:id_usuario"`

	CreatedAt time.Time `gorm:"autoCreateTime;column:created_at"`

	UpdatedAt time.Time `gorm:"autoUpdateTime;column:updated_at"`

	// Relaciones

	Usuario UsuarioModel `gorm:"foreignKey:IDUsuario"`
}

func (GestionHumanaModel) TableName() string {
	return "gestion_humana"
}
