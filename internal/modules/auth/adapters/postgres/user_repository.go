package postgres

import (
	"context"
	"errors"

	authdomain "disability_system_backend/internal/modules/auth/domain"
	usuariosmodels "disability_system_backend/internal/modules/usuarios/adapters/postgres/models"
	apperrors "disability_system_backend/internal/shared/errors"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) FindByID(ctx context.Context, id uint64) (*authdomain.User, error) {
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

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*authdomain.User, error) {
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

func (r *UserRepository) FindByDocumentNumber(ctx context.Context, docNumber string) (*authdomain.User, error) {
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

func (r *UserRepository) Create(ctx context.Context, user *authdomain.User) error {
	model := toModelUser(user)
	return r.db.WithContext(ctx).Create(model).Error
}

func (r *UserRepository) Update(ctx context.Context, user *authdomain.User) error {
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

func toDomainUser(m *usuariosmodels.UsuarioModel) *authdomain.User {
	return &authdomain.User{
		ID:              m.IDUsuario,
		IDRol:           m.IDRol,
		Nombre:          m.Nombre,
		Correo:          m.Correo,
		PasswordHash:    m.PasswordHash,
		NumeroDocumento: m.NumeroDocumento,
		Estado:          m.Estado,
		IsDeleted:       m.IsDeleted,
		CreatedAt:       m.CreatedAt,
	}
}

func toModelUser(u *authdomain.User) *usuariosmodels.UsuarioModel {
	return &usuariosmodels.UsuarioModel{
		IDUsuario:       u.ID,
		IDRol:           u.IDRol,
		Nombre:          u.Nombre,
		Correo:          u.Correo,
		PasswordHash:    u.PasswordHash,
		NumeroDocumento: u.NumeroDocumento,
		Estado:          u.Estado,
		IsDeleted:       u.IsDeleted,
		CreatedAt:       u.CreatedAt,
	}
}