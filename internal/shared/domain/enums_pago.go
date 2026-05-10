package domain

type TipoPago string

const (
	TipoPagoTransferencia TipoPago = "Transferencia bancaria"
	TipoPagoConsignacion  TipoPago = "Consignación"
	TipoPagoParcial       TipoPago = "Pago parcial"
	TipoPagoTotal         TipoPago = "Pago total"
	TipoPagoReintegro     TipoPago = "Reintegro"
)

func (t TipoPago) IsValid() bool {
	switch t {
	case TipoPagoTransferencia, TipoPagoConsignacion, TipoPagoParcial,
		TipoPagoTotal, TipoPagoReintegro:
		return true
	}
	return false
}

type EstadoPago string

const (
	EstadoPagoPendiente  EstadoPago = "Pendiente"
	EstadoPagoEnProceso   EstadoPago = "En proceso"
	EstadoPagoPagado      EstadoPago = "Pagado"
	EstadoPagoConciliado  EstadoPago = "Conciliado"
	EstadoPagoRechazado   EstadoPago = "Rechazado"
	EstadoPagoParcial     EstadoPago = "Parcial"
	EstadoPagoAnulado     EstadoPago = "Anulado"
)

func (e EstadoPago) IsValid() bool {
	switch e {
	case EstadoPagoPendiente, EstadoPagoEnProceso, EstadoPagoPagado,
		EstadoPagoConciliado, EstadoPagoRechazado, EstadoPagoParcial, EstadoPagoAnulado:
		return true
	}
	return false
}