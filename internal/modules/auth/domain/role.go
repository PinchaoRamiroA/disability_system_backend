package domain

type Role struct {
	ID       uint64
	Nombre   string
	Permisos []string
}
