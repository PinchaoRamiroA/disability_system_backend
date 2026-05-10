package domain

type TipoEntidad string

const (
	TipoEntidadEPS TipoEntidad = "EPS"
	TipoEntidadARL TipoEntidad = "ARL"
	TipoEntidadAFP TipoEntidad = "AFP"
)

func (t TipoEntidad) IsValid() bool {
	switch t {
	case TipoEntidadEPS, TipoEntidadARL, TipoEntidadAFP:
		return true
	}
	return false
}

func (t TipoEntidad) String() string {
	return string(t)
}