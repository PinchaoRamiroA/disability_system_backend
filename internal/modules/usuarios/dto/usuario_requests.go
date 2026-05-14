package dto

type CrearUsuarioRequest struct {
	IDRol           uint64  `json:"id_rol" binding:"required,gt=0"`
	Nombre          string  `json:"nombre" binding:"required,min=3,max=150"`
	Correo          string  `json:"correo" binding:"required,email,max=150"`
	NumeroCelular   *string `json:"numero_celular,omitempty"`
	Direccion       *string `json:"direccion,omitempty"`
	Password        string  `json:"password" binding:"required,min=8,max=100"`
	NumeroDocumento string  `json:"numero_documento" binding:"required,min=5,max=30"`
	NumeroAcudiente *string `json:"numero_acudiente,omitempty"`
}

type ActualizarUsuarioRequest struct {
	IDRol         *uint64 `json:"id_rol,omitempty"`
	Nombre        string  `json:"nombre"`
	Correo        string  `json:"correo"`
	NumeroCelular *string `json:"numero_celular,omitempty"`
	Direccion     *string `json:"direccion,omitempty"`
}

type CambiarPasswordRequest struct {
	PasswordActual string `json:"password_actual" binding:"required"`
	PasswordNuevo  string `json:"password_nuevo" binding:"required,min=8,max=100"`
}

type CambiarEstadoUsuarioRequest struct {
	Estado bool `json:"estado"`
}

type AsignarRolRequest struct {
	IDRol uint64 `json:"id_rol" binding:"required,gt=0"`
}

type ListarUsuariosQuery struct {
	Estado *bool   `form:"estado"`
	IDRol  *uint64  `form:"id_rol"`
	Search string   `form:"search"`
	Page   int      `form:"page,default=1" binding:"gte=1"`
	Limit  int      `form:"limit,default=20" binding:"gte=1,lte=100"`
}

type ListarRolesQuery struct {
	Page  int `form:"page,default=1" binding:"gte=1"`
	Limit int `form:"limit,default=20" binding:"gte=1,lte=100"`
}
