package mapper

import (
	"disability_system_backend/internal/modules/usuarios/adapters/postgres/models"
	"disability_system_backend/internal/modules/usuarios/domain"
	"disability_system_backend/internal/modules/usuarios/dto"
)

func ToDomain(model models.UsuarioModel) domain.Usuario {
	return domain.Usuario{
		ID:              model.IDUsuario,
		IDRol:           model.IDRol,
		Nombre:          model.Nombre,
		Correo:          model.Correo,
		NumeroDocumento: model.NumeroDocumento,
		NumeroCelular:   model.NumeroCelular,
		Direccion:       model.Direccion,
		NumeroAcudiente: model.NumeroAcudiente,
		Estado:          model.Estado,
		PasswordHash:    model.PasswordHash,
		IsDeleted:       model.IsDeleted,
		CreatedAt:       model.CreatedAt,
		UpdatedAt:       model.UpdatedAt,
	}
}

func ToModel(usuario domain.Usuario) models.UsuarioModel {
	return models.UsuarioModel{
		IDUsuario:       usuario.ID,
		IDRol:           usuario.IDRol,
		Nombre:          usuario.Nombre,
		Correo:          usuario.Correo,
		NumeroDocumento: usuario.NumeroDocumento,
		NumeroCelular:   usuario.NumeroCelular,
		Direccion:       usuario.Direccion,
		NumeroAcudiente: usuario.NumeroAcudiente,
		Estado:         usuario.Estado,
		PasswordHash:    usuario.PasswordHash,
		IsDeleted:       usuario.IsDeleted,
		CreatedAt:       usuario.CreatedAt,
		UpdatedAt:       usuario.UpdatedAt,
	}
}

func ToDomainRol(model models.RolModel) domain.Rol {
	return domain.Rol{
		ID:        model.IDRol,
		Nombre:    model.Nombre,
		Permisos:  model.GetPermisos(),
		IsDeleted: model.IsDeleted,
		CreatedAt: model.CreatedAt,
		UpdatedAt: model.UpdatedAt,
	}
}

func ToModelRol(rol domain.Rol) models.RolModel {
	model := models.RolModel{
		IDRol:     rol.ID,
		Nombre:    rol.Nombre,
		IsDeleted: rol.IsDeleted,
		CreatedAt: rol.CreatedAt,
		UpdatedAt: rol.UpdatedAt,
	}
	model.SetPermisos(rol.Permisos)
	return model
}

func ToUsuarioResponse(usuario domain.Usuario, nombreRol string) dto.UsuarioResponse {
	return dto.UsuarioResponse{
		IDUsuario:       usuario.ID,
		IDRol:           usuario.IDRol,
		NombreRol:       nombreRol,
		Nombre:          usuario.Nombre,
		Correo:          usuario.Correo,
		NumeroCelular:   usuario.NumeroCelular,
		Direccion:       usuario.Direccion,
		NumeroDocumento: usuario.NumeroDocumento,
		NumeroAcudiente: usuario.NumeroAcudiente,
		Estado:         usuario.Estado,
		CreatedAt:       usuario.CreatedAt,
		UpdatedAt:       usuario.UpdatedAt,
	}
}

func ToRolResponse(rol domain.Rol) dto.RolResponse {
	return dto.RolResponse{
		IDRol:    rol.ID,
		Nombre:   rol.Nombre,
		Permisos: rol.Permisos,
	}
}
