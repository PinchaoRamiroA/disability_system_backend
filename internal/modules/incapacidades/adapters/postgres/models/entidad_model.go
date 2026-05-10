package models

import "time"

type EntidadModel struct {
	IDEntidad uint64 `gorm:"primaryKey;autoIncrement;column:id_entidad"`

	Nombre string `gorm:"type:varchar(150);uniqueIndex;not null;column:nombre"`

	Tipo string `gorm:"type:varchar(100);not null;column:tipo"`

	PlazoTranscripcionDias *int `gorm:"column:plazo_transcripcion_dias"`

	TiempoMaximoPagoDias *int `gorm:"column:tiempo_maximo_pago_dias"`

	CanalAtencion *string `gorm:"type:varchar(150);column:canal_atencion"`

	RequiereTranscripcion bool `gorm:"default:false;column:requiere_transcripcion"`

	CreatedAt time.Time `gorm:"autoCreateTime;column:created_at"`

	UpdatedAt time.Time `gorm:"autoUpdateTime;column:updated_at"`
}

func (EntidadModel) TableName() string {
	return "entidad"
}
