package dto

type CrearPagoRequest struct {
	IDIncapacidad   uint64  `json:"id_incapacidad" binding:"required"`
	IDEntidad       uint64  `json:"id_entidad" binding:"required"`
	TipoPago        string  `json:"tipo_pago" binding:"required"`
	EstadoPago      string  `json:"estado_pago,omitempty"`
	Descripcion     *string `json:"descripcion,omitempty"`
	Valor           string  `json:"valor" binding:"required"`
	FechaPago       string  `json:"fecha_pago" binding:"required"`
	PeriodoContable *string `json:"periodo_contable,omitempty"`
}

type ActualizarPagoRequest struct {
	IDEntidad       *uint64 `json:"id_entidad,omitempty"`
	TipoPago        *string `json:"tipo_pago,omitempty"`
	EstadoPago      *string `json:"estado_pago,omitempty"`
	Descripcion     *string `json:"descripcion,omitempty"`
	Valor           *string `json:"valor,omitempty"`
	FechaPago       *string `json:"fecha_pago,omitempty"`
	PeriodoContable *string `json:"periodo_contable,omitempty"`
}

type ConciliarPagoRequest struct {
	Conciliado  bool    `json:"conciliado"`
	EstadoPago  *string `json:"estado_pago,omitempty"`
	Descripcion *string `json:"descripcion,omitempty"`
}

type ListarPagosQuery struct {
	IDIncapacidad *uint64 `form:"id_incapacidad"`
	IDEntidad     *uint64 `form:"id_entidad"`
	TipoPago      string  `form:"tipo_pago"`
	EstadoPago    string  `form:"estado_pago"`
	Conciliado    *bool   `form:"conciliado"`
	Page          int     `form:"page"`
	Limit         int     `form:"limit"`
}
