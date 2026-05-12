package models

import (
	"time"
)

type TipoDocumentoModel struct {
	IDTipoDocumento uint64    `gorm:"primaryKey;autoIncrement;column:id_tipo_documento"`
	Nombre          string    `gorm:"type:varchar(100);uniqueIndex;not null;column:nombre"`
	Descripcion     *string   `gorm:"type:text;column:descripcion"`
	Requerido       bool      `gorm:"default:true;column:requerido"`
	CreatedAt       time.Time `gorm:"autoCreateTime;column:created_at"`
	UpdatedAt       time.Time `gorm:"autoUpdateTime;column:updated_at"`
}

func (TipoDocumentoModel) TableName() string {
	return "tipo_documento"
}

type EstadoDocumentoModel struct {
	IDEstadoDocumento uint64    `gorm:"primaryKey;autoIncrement;column:id_estado_documento"`
	Nombre            string    `gorm:"type:varchar(50);uniqueIndex;not null;column:nombre"`
	Descripcion       *string   `gorm:"type:text;column:descripcion"`
	Color             *string   `gorm:"type:varchar(20);column:color"`
	CreatedAt         time.Time `gorm:"autoCreateTime;column:created_at"`
	UpdatedAt         time.Time `gorm:"autoUpdateTime;column:updated_at"`
}

func (EstadoDocumentoModel) TableName() string {
	return "estado_documento"
}

type TipoEntidadModel struct {
	IDTipoEntidad uint64    `gorm:"primaryKey;autoIncrement;column:id_tipo_entidad"`
	Nombre        string    `gorm:"type:varchar(50);uniqueIndex;not null;column:nombre"`
	Descripcion   *string   `gorm:"type:text;column:descripcion"`
	CreatedAt     time.Time `gorm:"autoCreateTime;column:created_at"`
	UpdatedAt     time.Time `gorm:"autoUpdateTime;column:updated_at"`
}

func (TipoEntidadModel) TableName() string {
	return "tipo_entidad"
}

type CanalRecepcionModel struct {
	IDCanalRecepcion uint64    `gorm:"primaryKey;autoIncrement;column:id_canal_recepcion"`
	Nombre           string    `gorm:"type:varchar(100);uniqueIndex;not null;column:nombre"`
	Descripcion      *string   `gorm:"type:text;column:descripcion"`
	Activo           bool      `gorm:"default:true;column:activo"`
	CreatedAt        time.Time `gorm:"autoCreateTime;column:created_at"`
	UpdatedAt        time.Time `gorm:"autoUpdateTime;column:updated_at"`
}

func (CanalRecepcionModel) TableName() string {
	return "canal_recepcion"
}

type CanalAtencionEntidadModel struct {
	IDCanalAtencion uint64    `gorm:"primaryKey;autoIncrement;column:id_canal_atencion"`
	Nombre          string    `gorm:"type:varchar(100);uniqueIndex;not null;column:nombre"`
	Descripcion     *string   `gorm:"type:text;column:descripcion"`
	CreatedAt       time.Time `gorm:"autoCreateTime;column:created_at"`
	UpdatedAt       time.Time `gorm:"autoUpdateTime;column:updated_at"`
}

func (CanalAtencionEntidadModel) TableName() string {
	return "canal_atencion_entidad"
}

type TipoPagoModel struct {
	IDTipoPago  uint64    `gorm:"primaryKey;autoIncrement;column:id_tipo_pago"`
	Nombre      string    `gorm:"type:varchar(100);uniqueIndex;not null;column:nombre"`
	Descripcion *string   `gorm:"type:text;column:descripcion"`
	CreatedAt   time.Time `gorm:"autoCreateTime;column:created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime;column:updated_at"`
}

func (TipoPagoModel) TableName() string {
	return "tipo_pago"
}

type EstadoPagoModel struct {
	IDEstadoPago uint64    `gorm:"primaryKey;autoIncrement;column:id_estado_pago"`
	Nombre       string    `gorm:"type:varchar(50);uniqueIndex;not null;column:nombre"`
	Descripcion  *string   `gorm:"type:text;column:descripcion"`
	Color        *string   `gorm:"type:varchar(20);column:color"`
	CreatedAt    time.Time `gorm:"autoCreateTime;column:created_at"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime;column:updated_at"`
}

func (EstadoPagoModel) TableName() string {
	return "estado_pago"
}

type PeriodicidadReporteModel struct {
	IDPeriodicidad uint64    `gorm:"primaryKey;autoIncrement;column:id_periodicidad"`
	Nombre         string    `gorm:"type:varchar(50);uniqueIndex;not null;column:nombre"`
	Dias           *int      `gorm:"column:dias"`
	Descripcion    *string   `gorm:"type:text;column:descripcion"`
	CreatedAt      time.Time `gorm:"autoCreateTime;column:created_at"`
	UpdatedAt      time.Time `gorm:"autoUpdateTime;column:updated_at"`
}

func (PeriodicidadReporteModel) TableName() string {
	return "periodicidad_reporte"
}
