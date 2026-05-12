package dto

import "time"

type UsuarioResponse struct {
	IDUsuario        uint64     `json:"id_usuario"`
	IDRol            uint64     `json:"id_rol"`
	NombreRol        string     `json:"nombre_rol"`
	Nombre           string     `json:"nombre"`
	Correo           string     `json:"correo"`
	NumeroCelular    *string    `json:"numero_celular"`
	Direccion        *string    `json:"direccion"`
	NumeroDocumento  string     `json:"numero_documento"`
	NumeroAcudiente  *string    `json:"numero_acudiente"`
	Estado           bool       `json:"estado"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
}

type UsuarioConPasswordResponse struct {
	UsuarioResponse
	PasswordHash string `json:"-"`
}

type RolResponse struct {
	IDRol     uint64    `json:"id_rol"`
	Nombre   string    `json:"nombre"`
	Permisos []string  `json:"permisos"`
}

type PaginatedUsuariosResponse struct {
	Success bool              `json:"success"`
	Message string            `json:"message"`
	Data    PaginatedData     `json:"data"`
}

type PaginatedData struct {
	Items     interface{} `json:"items"`
	Total     int64       `json:"total"`
	Page      int         `json:"page"`
	Limit     int         `json:"limit"`
	TotalPages int        `json:"total_pages"`
}

type MensajeSimpleResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}
