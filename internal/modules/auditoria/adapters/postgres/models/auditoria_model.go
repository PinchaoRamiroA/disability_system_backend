package models

import "time"

type AuditoriaModel struct {
	IDAuditoria    uint64    `gorm:"primaryKey;column:id_auditoria"`
	IDUsuario      *uint64   `gorm:"column:id_usuario"`
	UsuarioNombre  string    `gorm:"->;column:usuario_nombre"` // Read-only from join
	IDIncapacidad  *uint64   `gorm:"column:id_incapacidad"`
	TipoAccion     string    `gorm:"column:tipo_accion"`
	Modulo         string    `gorm:"column:modulo"`
	Descripcion    string    `gorm:"column:descripcion"`
	CambioAnterior *string   `gorm:"column:cambio_anterior"`
	CambioNuevo    *string   `gorm:"column:cambio_nuevo"`
	CreatedAt      time.Time `gorm:"column:created_at"`
}

func (AuditoriaModel) TableName() string {
	return "auditoria"
}
