package dto

import "time"

type RoleResponse struct {
	ID       uint64   `json:"id"`
	Nombre   string   `json:"nombre"`
	Permisos []string `json:"permisos"`
}

type UserResponse struct {
	ID              uint64        `json:"id"`
	Nombre          string        `json:"nombre"`
	Correo          string        `json:"correo"`
	NumeroCelular   *string       `json:"numero_celular,omitempty"`
	Direccion       *string       `json:"direccion,omitempty"`
	NumeroDocumento string        `json:"numero_documento"`
	Estado          bool          `json:"estado"`
	Rol             RoleResponse  `json:"rol"`
	CreatedAt       time.Time     `json:"created_at"`
}

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int64  `json:"expires_in"`
}

type LoginResponse struct {
	User         UserResponse `json:"user"`
	AccessToken  string       `json:"access_token"`
	RefreshToken string       `json:"refresh_token"`
	TokenType    string       `json:"token_type"`
	ExpiresIn    int64        `json:"expires_in"`
}