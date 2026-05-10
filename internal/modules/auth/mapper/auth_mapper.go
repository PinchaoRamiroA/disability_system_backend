package mapper

import (
	"disability_system_backend/internal/modules/auth/domain"
	"disability_system_backend/internal/modules/auth/dto"
)

func ToRoleResponse(role *domain.Role) *dto.RoleResponse {
	if role == nil {
		return nil
	}
	return &dto.RoleResponse{
		ID:       role.ID,
		Nombre:   role.Nombre,
		Permisos: role.Permisos,
	}
}

func ToUserResponse(user *domain.User, role *domain.Role) *dto.UserResponse {
	if user == nil {
		return nil
	}
	resp := &dto.UserResponse{
		ID:              user.ID,
		Nombre:          user.Nombre,
		Correo:          user.Correo,
		NumeroDocumento: user.NumeroDocumento,
		Estado:          user.Estado,
		CreatedAt:       user.CreatedAt,
	}
	if role != nil {
		resp.Rol = *ToRoleResponse(role)
	}
	return resp
}

func ToLoginResponse(user *domain.User, role *domain.Role, accessToken, refreshToken string, expiresIn int64) *dto.LoginResponse {
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