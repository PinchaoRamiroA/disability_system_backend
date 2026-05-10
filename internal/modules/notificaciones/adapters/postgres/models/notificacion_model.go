package models

import "time"

type NotificacionModel struct {
	IDNotificacion uint64 `gorm:"primaryKey;autoIncrement;column:id_notificacion"`

	IDUsuario uint64 `gorm:"not null;column:id_usuario;index"`

	IDIncapacidad *uint64 `gorm:"column:id_incapacidad;index"`

	IDTipoNotificacion uint64 `gorm:"not null;column:id_tipo_notificacion;index"`

	Mensaje string `gorm:"type:text;not null;column:mensaje"`

	Fecha time.Time `gorm:"autoCreateTime;column:fecha"`

	Leida bool `gorm:"default:false;column:leida"`

	CreatedAt time.Time `gorm:"autoCreateTime;column:created_at"`

	UpdatedAt time.Time `gorm:"autoUpdateTime;column:updated_at"`
}

func (NotificacionModel) TableName() string {
	return "notificacion"
}