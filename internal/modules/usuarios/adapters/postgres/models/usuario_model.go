package models

import "time"

type UsuarioModel struct {
	IDUsuario uint64 `gorm:"primaryKey;autoIncrement;column:id_usuario"`

	IDRol uint64 `gorm:"not null;column:id_rol;index"`

	Nombre string `gorm:"type:varchar(150);not null;column:nombre"`

	Correo string `gorm:"type:varchar(150);uniqueIndex;not null;column:correo"`

	NumeroCelular *string `gorm:"type:varchar(20);column:numero_celular"`

	Direccion *string `gorm:"type:varchar(255);column:direccion"`

	PasswordHash string `gorm:"type:text;not null;column:password_hash"`

	NumeroDocumento string `gorm:"type:varchar(30);uniqueIndex;not null;column:numero_documento"`

	NumeroAcudiente *string `gorm:"type:varchar(20);column:numero_acudiente"`

	Estado bool `gorm:"default:true;not null;column:estado"`

	CreatedAt time.Time `gorm:"autoCreateTime;column:created_at"`

	UpdatedAt time.Time `gorm:"autoUpdateTime;column:updated_at"`

	IsDeleted bool `gorm:"default:false;column:is_deleted"`
}

func (UsuarioModel) TableName() string {
	return "usuario"
}