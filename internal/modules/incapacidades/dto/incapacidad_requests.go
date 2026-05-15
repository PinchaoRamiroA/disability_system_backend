package dto

type CrearIncapacidadRequest struct {
	IDUsuario       *uint64 `json:"id_usuario,omitempty"`
	IDEstado        *uint64 `json:"id_estado,omitempty"`
	IDTipo          uint64  `json:"id_tipo" binding:"required"`
	IDEntidad       uint64  `json:"id_entidad" binding:"required"`
	CanalRecepcion  string  `json:"canal_recepcion,omitempty"`
	Titulo          string  `json:"titulo" binding:"required"`
	FechaInicio     string  `json:"fecha_inicio" binding:"required"`
	FechaFin        *string `json:"fecha_fin,omitempty"`
	Origen          string  `json:"origen,omitempty"`
	FechaRadicacion *string `json:"fecha_radicacion,omitempty"`
	FechaPago       *string `json:"fecha_pago,omitempty"`
	Observaciones   *string `json:"observaciones,omitempty"`
	CreatedBy       *uint64 `json:"created_by,omitempty"`
}

type ActualizarIncapacidadRequest struct {
	IDUsuario       *uint64 `json:"id_usuario,omitempty"`
	IDTipo          *uint64 `json:"id_tipo,omitempty"`
	IDEntidad       *uint64 `json:"id_entidad,omitempty"`
	CanalRecepcion  *string `json:"canal_recepcion,omitempty"`
	Titulo          *string `json:"titulo,omitempty"`
	FechaInicio     *string `json:"fecha_inicio,omitempty"`
	FechaFin        *string `json:"fecha_fin,omitempty"`
	Origen          *string `json:"origen,omitempty"`
	FechaRadicacion *string `json:"fecha_radicacion,omitempty"`
	FechaPago       *string `json:"fecha_pago,omitempty"`
	Observaciones   *string `json:"observaciones,omitempty"`
}

type CambiarEstadoRequest struct {
	IDEstado      uint64  `json:"id_estado" binding:"required"`
	Observaciones *string `json:"observaciones,omitempty"`
}

type ListarIncapacidadesQuery struct {
	IDUsuario      *uint64 `form:"id_usuario"`
	IDEstado       *uint64 `form:"id_estado"`
	IDTipo         *uint64 `form:"id_tipo"`
	IDEntidad      *uint64 `form:"id_entidad"`
	Origen         string  `form:"origen"`
	CanalRecepcion string  `form:"canal_recepcion"`
	Search         string  `form:"search"`
	FechaDesde     string  `form:"fecha_desde"`
	FechaHasta     string  `form:"fecha_hasta"`
	Page           int     `form:"page"`
	Limit          int     `form:"limit"`
}
