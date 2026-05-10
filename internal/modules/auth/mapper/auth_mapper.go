package mapper

import (
	"disability_system_backend/internal/modules/auth/dto"
	usuariosdomain "disability_system_backend/internal/modules/usuarios/domain"
)

func ToRoleResponse(role *usuariosdomain.Rol) *dto.RoleResponse {
	if role == nil {
		return nil
	}
	return &dto.RoleResponse{
		ID:       role.ID,
		Nombre:   role.Nombre,
		Permisos: role.Permisos,
	}
}

func ToUserResponse(user *usuariosdomain.Usuario, role *usuariosdomain.Rol) *dto.UserResponse {
	if user == nil {
		return nil
	}
	resp := &dto.UserResponse{
		ID:              user.ID,
		Nombre:          user.Nombre,
		Correo:          user.Correo,
		NumeroCelular:   user.NumeroCelular,
		Direccion:       user.Direccion,
		NumeroDocumento: user.NumeroDocumento,
		Estado:          user.Estado,
		CreatedAt:       user.CreatedAt,
	}
	if role != nil {
		resp.Rol = *ToRoleResponse(role)
	}
	return resp
}

func ToLoginResponse(user *usuariosdomain.Usuario, role *usuariosdomain.Rol, accessToken, refreshToken string, expiresIn int64) *dto.LoginResponse {
	userResp := ToUserResponse(user, role)
	if userResp == nil {
		return nil
	}
	return &dto.LoginResponse{
		User:         *userResp,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    expiresIn,
	}
}

func ToTokenResponse(accessToken, refreshToken string, expiresIn int64) *dto.TokenResponse {
	return &dto.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    expiresIn,
	}
}