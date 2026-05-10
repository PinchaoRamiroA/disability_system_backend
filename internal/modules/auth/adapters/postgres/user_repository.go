package postgres

import (
	"context"
	"errors"

	usuariosmodels "disability_system_backend/internal/modules/usuarios/adapters/postgres/models"
	usuariosdomain "disability_system_backend/internal/modules/usuarios/domain"
	apperrors "disability_system_backend/internal/shared/errors"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) FindByID(ctx context.Context, id uint64) (*usuariosdomain.Usuario, error) {
	var model usuariosmodels.UsuarioModel
	err := r.db.WithContext(ctx).Where("id_usuario = ?", id).First(&model).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.ErrUserNotFound.WithError(err)
		}
		return nil, err
	}
	return toDomainUser(&model), nil
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*usuariosdomain.Usuario, error) {
	var model usuariosmodels.UsuarioModel
	err := r.db.WithContext(ctx).Where("correo = ?", email).First(&model).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.ErrUserNotFound.WithError(err)
		}
		return nil, err
	}
	return toDomainUser(&model), nil
}

func (r *UserRepository) FindByDocumentNumber(ctx context.Context, docNumber string) (*usuariosdomain.Usuario, error) {
	var model usuariosmodels.UsuarioModel
	err := r.db.WithContext(ctx).Where("numero_documento = ?", docNumber).First(&model).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.ErrUserNotFound.WithError(err)
		}
		return nil, err
	}
	return toDomainUser(&model), nil
}

func (r *UserRepository) Create(ctx context.Context, user *usuariosdomain.Usuario) error {
	model := toModelUser(user)
	return r.db.WithContext(ctx).Create(model).Error
}

func (r *UserRepository) Update(ctx context.Context, user *usuariosdomain.Usuario) error {
	model := toModelUser(user)
	return r.db.WithContext(ctx).Save(model).Error
}

func (r *UserRepository) Delete(ctx context.Context, id uint64) error {
	return r.db.WithContext(ctx).Delete(&usuariosmodels.UsuarioModel{}, "id_usuario = ?", id).Error
}

func (r *UserRepository) EmailExists(ctx context.Context, email string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&usuariosmodels.UsuarioModel{}).Where("correo = ?", email).Count(&count).Error
	return count > 0, err
}

func (r *UserRepository) DocumentExists(ctx context.Context, docNumber string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&usuariosmodels.UsuarioModel{}).Where("numero_documento = ?", docNumber).Count(&count).Error
	return count > 0, err
}

func toDomainUser(m *usuariosmodels.UsuarioModel) *usuariosdomain.Usuario {
	return &usuariosdomain.Usuario{
		ID:              m.IDUsuario,
		IDRol:           m.IDRol,
		Nombre:          m.Nombre,
		Correo:          m.Correo,
		NumeroCelular:   m.NumeroCelular,
		Direccion:       m.Direccion,
		PasswordHash:    m.PasswordHash,
		NumeroDocumento: m.NumeroDocumento,
		NumeroAcudiente: m.NumeroAcudiente,
		Estado:          m.Estado,
		IsDeleted:       m.IsDeleted,
		CreatedAt:       m.CreatedAt,
		UpdatedAt:       m.UpdatedAt,
	}
}

func toModelUser(u *usuariosdomain.Usuario) *usuariosmodels.UsuarioModel {
	return &usuariosmodels.UsuarioModel{
		IDUsuario:       u.ID,
		IDRol:           u.IDRol,
		Nombre:          u.Nombre,
		Correo:          u.Correo,
		NumeroCelular:   u.NumeroCelular,
		Direccion:       u.Direccion,
		PasswordHash:    u.PasswordHash,
		NumeroDocumento: u.NumeroDocumento,
		NumeroAcudiente: u.NumeroAcudiente,
		Estado:          u.Estado,
		IsDeleted:       u.IsDeleted,
		CreatedAt:       u.CreatedAt,
		UpdatedAt:       u.UpdatedAt,
	}
}

var _ UserRepositoryI = (*UserRepository)(nil)

type UserRepositoryI interface {
	FindByID(ctx context.Context, id uint64) (*usuariosdomain.Usuario, error)
	FindByEmail(ctx context.Context, email string) (*usuariosdomain.Usuario, error)
	FindByDocumentNumber(ctx context.Context, docNumber string) (*usuariosdomain.Usuario, error)
	Create(ctx context.Context, user *usuariosdomain.Usuario) error
	Update(ctx context.Context, user *usuariosdomain.Usuario) error
	Delete(ctx context.Context, id uint64) error
	EmailExists(ctx context.Context, email string) (bool, error)
	DocumentExists(ctx context.Context, docNumber string) (bool, error)
}