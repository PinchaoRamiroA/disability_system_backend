package dto

type PagoResponse struct {
	IDPago          uint64  `json:"id_pago"`
	IDIncapacidad   uint64  `json:"id_incapacidad"`
	IDEntidad       uint64  `json:"id_entidad"`
	TipoPago        string  `json:"tipo_pago"`
	EstadoPago      string  `json:"estado_pago"`
	Descripcion     *string `json:"descripcion,omitempty"`
	Valor           string  `json:"valor"`
	FechaPago       string  `json:"fecha_pago"`
	PeriodoContable *string `json:"periodo_contable,omitempty"`
	Conciliado      bool    `json:"conciliado"`
	RegistradoPor   *uint64 `json:"registrado_por,omitempty"`
	CreatedAt       string  `json:"created_at"`
	UpdatedAt       string  `json:"updated_at"`
}

type SeguimientoCobroResponse struct {
	IDSeguimiento   uint64  `json:"id_seguimiento"`
	IDIncapacidad   uint64  `json:"id_incapacidad"`
	TipoSeguimiento string  `json:"tipo_seguimiento"`
	Descripcion     *string `json:"descripcion,omitempty"`
	Fecha           string  `json:"fecha"`
	Resultado       *string `json:"resultado,omitempty"`
	GestionadoPor   *uint64 `json:"gestionado_por,omitempty"`
	CreatedAt       string  `json:"created_at"`
	UpdatedAt       string  `json:"updated_at"`
}
