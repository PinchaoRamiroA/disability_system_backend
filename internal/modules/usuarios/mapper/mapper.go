package mapper

import (
	"disability_system_backend/internal/modules/usuarios/adapters/postgres/models"
	"disability_system_backend/internal/modules/usuarios/domain"
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