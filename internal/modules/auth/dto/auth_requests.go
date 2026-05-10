package dto

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type RegisterRequest struct {
	Nombre           string `json:"nombre" binding:"required,min=2,max=150"`
	Email            string `json:"email" binding:"required,email"`
	Password         string `json:"password" binding:"required,min=6"`
	NumeroDocumento  string `json:"numero_documento" binding:"required,min=5"`
	NumeroCelular    string `json:"numero_celular,omitempty"`
	Direccion        string `json:"direccion,omitempty"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}