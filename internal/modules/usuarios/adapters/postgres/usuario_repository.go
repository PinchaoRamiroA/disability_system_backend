package postgres

import (
	"context"
	"errors"

	"disability_system_backend/internal/modules/usuarios/adapters/postgres/models"
	"disability_system_backend/internal/modules/usuarios/domain"
	apperrors "disability_system_backend/internal/shared/errors"

	"gorm.io/gorm"
)

type UsuarioRepository struct {
	db *gorm.DB
}

func NewUsuarioRepository(db *gorm.DB) *UsuarioRepository {
	return &UsuarioRepository{db: db}
}

func (r *UsuarioRepository) FindByID(ctx context.Context, id uint64) (*domain.Usuario, error) {
	var model models.UsuarioModel
	err := r.db.WithContext(ctx).
		Where("id_usuario = ? AND is_deleted = false", id).
		First(&model).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.ErrUserNotFound.WithError(err)
		}
		return nil, err
	}
	return toDomainUsuario(&model), nil
}

func (r *UsuarioRepository) FindByEmail(ctx context.Context, email string) (*domain.Usuario, error) {
	var model models.UsuarioModel
	err := r.db.WithContext(ctx).
		Where("correo = ? AND is_deleted = false", email).
		First(&model).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.ErrUserNotFound.WithError(err)
		}
		return nil, err
	}
	return toDomainUsuario(&model), nil
}

func (r *UsuarioRepository) FindByDocumentNumber(ctx context.Context, docNumber string) (*domain.Usuario, error) {
	var model models.UsuarioModel
	err := r.db.WithContext(ctx).
		Where("numero_documento = ? AND is_deleted = false", docNumber).
		First(&model).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.ErrUserNotFound.WithError(err)
		}
		return nil, err
	}
	return toDomainUsuario(&model), nil
}

func (r *UsuarioRepository) FindAll(ctx context.Context, page, limit int, estado *bool, idRol *uint64, search string) ([]domain.Usuario, int64, error) {
	var modelList []models.UsuarioModel
	var total int64

	query := r.db.WithContext(ctx).Model(&models.UsuarioModel{}).Where("is_deleted = false")

	if estado != nil {
		query = query.Where("estado = ?", *estado)
	}
	if idRol != nil {
		query = query.Where("id_rol = ?", *idRol)
	}
	if search != "" {
		searchTerm := "%" + search + "%"
		query = query.Where("(nombre ILIKE ? OR correo ILIKE ? OR numero_documento ILIKE ?)", searchTerm, searchTerm, searchTerm)
	}

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	err = query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&modelList).Error
	if err != nil {
		return nil, 0, err
	}

	usuarios := make([]domain.Usuario, len(modelList))
	for i, m := range modelList {
		usuarios[i] = *toDomainUsuario(&m)
	}

	return usuarios, total, nil
}

func (r *UsuarioRepository) Create(ctx context.Context, usuario *domain.Usuario) error {
	model := toModelUsuario(usuario)
	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		return err
	}
	*usuario = *toDomainUsuario(model)
	return nil
}

func (r *UsuarioRepository) Update(ctx context.Context, usuario *domain.Usuario) error {
	model := toModelUsuario(usuario)
	if err := r.db.WithContext(ctx).Save(model).Error; err != nil {
		return err
	}
	*usuario = *toDomainUsuario(model)
	return nil
}

func (r *UsuarioRepository) SoftDelete(ctx context.Context, id uint64) error {
	return r.db.WithContext(ctx).
		Model(&models.UsuarioModel{}).
		Where("id_usuario = ?", id).
		Update("is_deleted", true).Error
}

func (r *UsuarioRepository) EmailExists(ctx context.Context, email string, excludeID *uint64) (bool, error) {
	var count int64
	query := r.db.WithContext(ctx).Model(&models.UsuarioModel{}).Where("correo = ? AND is_deleted = false", email)
	if excludeID != nil {
		query = query.Where("id_usuario != ?", *excludeID)
	}
	err := query.Count(&count).Error
	return count > 0, err
}

func (r *UsuarioRepository) DocumentExists(ctx context.Context, docNumber string, excludeID *uint64) (bool, error) {
	var count int64
	query := r.db.WithContext(ctx).Model(&models.UsuarioModel{}).Where("numero_documento = ? AND is_deleted = false", docNumber)
	if excludeID != nil {
		query = query.Where("id_usuario != ?", *excludeID)
	}
	err := query.Count(&count).Error
	return count > 0, err
}

func toDomainUsuario(m *models.UsuarioModel) *domain.Usuario {
	return &domain.Usuario{
		ID:              m.IDUsuario,
		IDRol:           m.IDRol,
		Nombre:          m.Nombre,
		Correo:          m.Correo,
		NumeroDocumento: m.NumeroDocumento,
		NumeroCelular:   m.NumeroCelular,
		Direccion:       m.Direccion,
		NumeroAcudiente: m.NumeroAcudiente,
		Estado:          m.Estado,
		PasswordHash:    m.PasswordHash,
		IsDeleted:       m.IsDeleted,
		CreatedAt:       m.CreatedAt,
		UpdatedAt:       m.UpdatedAt,
	}
}

func toModelUsuario(u *domain.Usuario) *models.UsuarioModel {
	return &models.UsuarioModel{
		IDUsuario:       u.ID,
		IDRol:           u.IDRol,
		Nombre:          u.Nombre,
		Correo:          u.Correo,
		NumeroDocumento: u.NumeroDocumento,
		NumeroCelular:   u.NumeroCelular,
		Direccion:       u.Direccion,
		NumeroAcudiente: u.NumeroAcudiente,
		Estado:          u.Estado,
		PasswordHash:    u.PasswordHash,
		IsDeleted:       u.IsDeleted,
		CreatedAt:       u.CreatedAt,
		UpdatedAt:       u.UpdatedAt,
	}
}
