package domain

import "time"

type TipoIncapacidad struct {
	IDTipo                  uint64
	Nombre                  string
	DocumentosRequeridos     map[string]interface{}
	CreatedAt               time.Time
	UpdatedAt               time.Time
}