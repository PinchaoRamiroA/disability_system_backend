package usecase

import (
	"context"

	"disability_system_backend/internal/modules/usuarios/adapters/postgres"
	"disability_system_backend/internal/modules/usuarios/domain"
	"disability_system_backend/internal/shared/auth"
	apperrors "disability_system_backend/internal/shared/errors"

	"gorm.io/gorm"
)

type UsuarioUseCase struct {
	usuarioRepo *postgres.UsuarioRepository
	rolRepo     *postgres.RolRepository
	db          *gorm.DB
}

func NewUsuarioUseCase(db *gorm.DB) *UsuarioUseCase {
	return &UsuarioUseCase{
		usuarioRepo: postgres.NewUsuarioRepository(db),
		rolRepo:     postgres.NewRolRepository(db),
		db:          db,
	}
}

func (uc *UsuarioUseCase) Crear(ctx context.Context, req struct {
	IDRol           uint64
	Nombre          string
	Correo          string
	NumeroCelular   *string
	Direccion       *string
	Password        string
	NumeroDocumento string
	NumeroAcudiente *string
}) (*domain.Usuario, string, error) {
	_, err := uc.rolRepo.FindByID(ctx, req.IDRol)
	if err != nil {
		return nil, "", apperrors.ErrRolNotFound.WithError(err)
	}

	exists, err := uc.usuarioRepo.EmailExists(ctx, req.Correo, nil)
	if err != nil {
		return nil, "", err
	}
	if exists {
		return nil, "", apperrors.ErrEmailAlreadyExists
	}

	exists, err = uc.usuarioRepo.DocumentExists(ctx, req.NumeroDocumento, nil)
	if err != nil {
		return nil, "", err
	}
	if exists {
		return nil, "", apperrors.ErrBadRequest.WithMessage("El número de documento ya está registrado")
	}

	passwordHash, err := auth.HashPassword(req.Password)
	if err != nil {
		return nil, "", err
	}

	usuario := &domain.Usuario{
		IDRol:           req.IDRol,
		Nombre:          req.Nombre,
		Correo:          req.Correo,
		NumeroCelular:   req.NumeroCelular,
		Direccion:       req.Direccion,
		PasswordHash:    passwordHash,
		NumeroDocumento: req.NumeroDocumento,
		NumeroAcudiente: req.NumeroAcudiente,
		Estado:          true,
	}

	if err := uc.usuarioRepo.Create(ctx, usuario); err != nil {
		return nil, "", err
	}

	return usuario, passwordHash, nil
}

func (uc *UsuarioUseCase) Obtener(ctx context.Context, id uint64) (*domain.Usuario, *domain.Rol, error) {
	usuario, err := uc.usuarioRepo.FindByID(ctx, id)
	if err != nil {
		return nil, nil, err
	}

	rol, err := uc.rolRepo.FindByID(ctx, usuario.IDRol)
	if err != nil {
		return nil, nil, err
	}

	return usuario, rol, nil
}

func (uc *UsuarioUseCase) Listar(ctx context.Context, page, limit int, estado *bool, idRol *uint64) ([]domain.Usuario, []domain.Rol, int64, error) {
	usuarios, total, err := uc.usuarioRepo.FindAll(ctx, page, limit, estado, idRol)
	if err != nil {
		return nil, nil, 0, err
	}

	rolesMap := make(map[uint64]*domain.Rol)
	for _, u := range usuarios {
		if _, ok := rolesMap[u.IDRol]; !ok {
			rol, err := uc.rolRepo.FindByID(ctx, u.IDRol)
			if err == nil {
				rolesMap[u.IDRol] = rol
			}
		}
	}

	roles := make([]domain.Rol, 0, len(rolesMap))
	for _, r := range rolesMap {
		roles = append(roles, *r)
	}

	return usuarios, roles, total, nil
}

func (uc *UsuarioUseCase) Actualizar(ctx context.Context, id uint64, req struct {
	IDRol         *uint64
	Nombre        string
	Correo        string
	NumeroCelular *string
	Direccion     *string
}) (*domain.Usuario, error) {
	usuario, err := uc.usuarioRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.IDRol != nil {
		_, err = uc.rolRepo.FindByID(ctx, *req.IDRol)
		if err != nil {
			return nil, apperrors.ErrRolNotFound.WithError(err)
		}
		usuario.IDRol = *req.IDRol
	}

	if req.Correo != "" && req.Correo != usuario.Correo {
		exists, err := uc.usuarioRepo.EmailExists(ctx, req.Correo, &id)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, apperrors.ErrEmailAlreadyExists
		}
		usuario.Correo = req.Correo
	}

	if req.Nombre != "" {
		usuario.Nombre = req.Nombre
	}

	if req.NumeroCelular != nil {
		usuario.NumeroCelular = req.NumeroCelular
	}

	if req.Direccion != nil {
		usuario.Direccion = req.Direccion
	}

	if err := uc.usuarioRepo.Update(ctx, usuario); err != nil {
		return nil, err
	}

	return usuario, nil
}

func (uc *UsuarioUseCase) CambiarEstado(ctx context.Context, id uint64, estado bool) error {
	_, err := uc.usuarioRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	return uc.db.WithContext(ctx).Exec(
		"UPDATE usuario SET estado = ? WHERE id_usuario = ?",
		estado, id,
	).Error
}

func (uc *UsuarioUseCase) CambiarPassword(ctx context.Context, id uint64, passwordActual, passwordNuevo string) error {
	usuario, err := uc.usuarioRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	if !auth.CheckPassword(passwordActual, usuario.PasswordHash) {
		return apperrors.ErrUnauthorized.WithMessage("La contraseña actual es incorrecta")
	}

	passwordHash, err := auth.HashPassword(passwordNuevo)
	if err != nil {
		return err
	}

	return uc.db.WithContext(ctx).Exec(
		"UPDATE usuario SET password_hash = ? WHERE id_usuario = ?",
		passwordHash, id,
	).Error
}

func (uc *UsuarioUseCase) AsignarRol(ctx context.Context, id uint64, idRol uint64) error {
	usuario, err := uc.usuarioRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	_, err = uc.rolRepo.FindByID(ctx, idRol)
	if err != nil {
		return apperrors.ErrRolNotFound.WithError(err)
	}

	usuario.IDRol = idRol
	return uc.usuarioRepo.Update(ctx, usuario)
}

func (uc *UsuarioUseCase) Eliminar(ctx context.Context, id uint64) error {
	_, err := uc.usuarioRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	return uc.usuarioRepo.SoftDelete(ctx, id)
}
