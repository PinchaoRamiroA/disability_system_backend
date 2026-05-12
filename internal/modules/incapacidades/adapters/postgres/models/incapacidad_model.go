package models

import (
	"time"
)

type IncapacidadModel struct {
	IDIncapacidad              uint64     `gorm:"primaryKey;autoIncrement;column:id_incapacidad"`
	IDUsuario                  uint64     `gorm:"not null;column:id_usuario;index"`
	IDEstado                   uint64     `gorm:"not null;column:id_estado;index"`
	IDTipo                     uint64     `gorm:"not null;column:id_tipo;index"`
	IDEntidad                  uint64     `gorm:"not null;column:id_entidad;index"`
	CanalRecepcion            string     `gorm:"type:varchar(100);column:canal_recepcion"`
	Titulo                    string     `gorm:"type:varchar(200);not null;column:titulo"`
	FechaInicio                time.Time  `gorm:"type:date;not null;column:fecha_inicio"`
	FechaFin                  *time.Time `gorm:"type:date;column:fecha_fin"`
	Origen                    string     `gorm:"type:varchar(100);not null;column:origen"`
	FechaRadicacion          *time.Time `gorm:"type:date;column:fecha_radicacion"`
	FechaPago                 *time.Time `gorm:"type:date;column:fecha_pago"`
	Observaciones             *string    `gorm:"type:text;column:observaciones"`
	CreatedBy                *uint64    `gorm:"column:created_by;index"`
	CreatedAt                time.Time  `gorm:"autoCreateTime;column:created_at"`
	UpdatedAt                time.Time  `gorm:"autoUpdateTime;column:updated_at"`
	IsDeleted                bool       `gorm:"default:false;column:is_deleted"`
	FechaTranscripcion        *time.Time `gorm:"type:timestamp;column:fecha_transcripcion"`
	TranscritoPor             *uint64    `gorm:"column:transcrito_por"`
	ObservacionesTranscripcion *string    `gorm:"type:text;column:observaciones_transcripcion"`
	FechaLimiteTranscripcion *time.Time `gorm:"type:date;column:fecha_limite_transcripcion"`
	EstadoTranscripcion      string     `gorm:"type:varchar(50);column:estado_transcripcion"`

	Estado  EstadoIncapacidadModel `gorm:"foreignKey:IDEstado;references:IDEstado"`
	Tipo    TipoIncapacidadModel   `gorm:"foreignKey:IDTipo;references:IDTipo"`
	Entidad EntidadModel           `gorm:"foreignKey:IDEntidad;references:IDEntidad"`
}

func (IncapacidadModel) TableName() string {
	return "incapacidad"
}
