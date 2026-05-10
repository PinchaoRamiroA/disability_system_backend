package domain

type OrigenIncapacidad string

const (
	OrigenComun      OrigenIncapacidad = "Común"
	OrigenLaboral    OrigenIncapacidad = "Laboral"
	OrigenTransito   OrigenIncapacidad = "Tránsito"
	OrigenMaternidad OrigenIncapacidad = "Maternidad"
	OrigenPaternidad OrigenIncapacidad = "Paternidad"
)

func (o OrigenIncapacidad) IsValid() bool {
	switch o {
	case OrigenComun, OrigenLaboral, OrigenTransito, OrigenMaternidad, OrigenPaternidad:
		return true
	}
	return false
}

func (o OrigenIncapacidad) String() string {
	return string(o)
}