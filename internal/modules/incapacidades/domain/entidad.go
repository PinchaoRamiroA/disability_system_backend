package domain

import "time"

type Entidad struct {
	IDEntidad                uint64
	Nombre                   string
	Tipo                     string
	PlazoTranscripcionDias   *int
	TiempoMaximoPagoDias     *int
	CanalAtencion            *string
	RequiereTranscripcion    bool
	CreatedAt                time.Time
	UpdatedAt                time.Time
}