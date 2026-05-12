package models

import (
	"encoding/json"
	"time"

	"gorm.io/datatypes"
)

type EntidadModel struct {
	IDEntidad              uint64         `gorm:"primaryKey;autoIncrement;column:id_entidad"`
	Nombre                 string         `gorm:"type:varchar(150);uniqueIndex;not null;column:nombre"`
	Tipo                   string         `gorm:"type:varchar(100);not null;column:tipo"`
	PlazoTranscripcionDias *int           `gorm:"column:plazo_transcripcion_dias"`
	TiempoMaximoPagoDias   *int           `gorm:"column:tiempo_maximo_pago_dias"`
	CanalAtencion          *string        `gorm:"type:varchar(150);column:canal_atencion"`
	CanalesAtencion        datatypes.JSON `gorm:"type:jsonb;column:canales_atencion"`
	RequiereTranscripcion  bool           `gorm:"default:false;column:requiere_transcripcion"`
	CreatedAt              time.Time      `gorm:"autoCreateTime;column:created_at"`
	UpdatedAt              time.Time      `gorm:"autoUpdateTime;column:updated_at"`
}

func (EntidadModel) TableName() string {
	return "entidad"
}

func (e *EntidadModel) GetCanalesAtencion() []string {
	if e.CanalesAtencion == nil {
		return nil
	}
	var canales []string
	if err := json.Unmarshal(e.CanalesAtencion, &canales); err != nil {
		return nil
	}
	return canales
}

func (e *EntidadModel) SetCanalesAtencion(canales []string) error {
	if canales == nil {
		e.CanalesAtencion = nil
		return nil
	}
	data, err := json.Marshal(canales)
	if err != nil {
		return err
	}
	e.CanalesAtencion = data
	return nil
}
