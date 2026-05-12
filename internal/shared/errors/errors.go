package errors

import (
	"errors"
	"fmt"
	"net/http"
)

type AppError struct {
	Code       string
	Message    string
	HTTPStatus int
	Details    any
	Err        error
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s", e.Message, e.Err.Error())
	}
	return e.Message
}

func (e *AppError) Unwrap() error {
	return e.Err
}

func (e *AppError) WithDetails(details any) *AppError {
	return &AppError{
		Code:       e.Code,
		Message:    e.Message,
		HTTPStatus: e.HTTPStatus,
		Details:    details,
		Err:        e.Err,
	}
}

func (e *AppError) WithMessage(message string) *AppError {
	return &AppError{
		Code:       e.Code,
		Message:    message,
		HTTPStatus: e.HTTPStatus,
		Details:    e.Details,
		Err:        e.Err,
	}
}

func (e *AppError) WithError(err error) *AppError {
	return &AppError{
		Code:       e.Code,
		Message:    e.Message,
		HTTPStatus: e.HTTPStatus,
		Details:    e.Details,
		Err:        err,
	}
}

var (
	ErrUnauthorized = &AppError{
		Code:       "UNAUTHORIZED",
		Message:    "No autorizado",
		HTTPStatus: http.StatusUnauthorized,
	}

	ErrInvalidCredentials = &AppError{
		Code:       "INVALID_CREDENTIALS",
		Message:    "Credenciales inválidas",
		HTTPStatus: http.StatusUnauthorized,
	}

	ErrTokenExpired = &AppError{
		Code:       "TOKEN_EXPIRED",
		Message:    "Token expirado",
		HTTPStatus: http.StatusUnauthorized,
	}

	ErrTokenInvalid = &AppError{
		Code:       "TOKEN_INVALID",
		Message:    "Token inválido",
		HTTPStatus: http.StatusUnauthorized,
	}

	ErrForbidden = &AppError{
		Code:       "FORBIDDEN",
		Message:    "Acceso prohibido",
		HTTPStatus: http.StatusForbidden,
	}

	ErrNotFound = &AppError{
		Code:       "NOT_FOUND",
		Message:    "Recurso no encontrado",
		HTTPStatus: http.StatusNotFound,
	}

	ErrUserNotFound = &AppError{
		Code:       "USER_NOT_FOUND",
		Message:    "Usuario no encontrado",
		HTTPStatus: http.StatusNotFound,
	}

	ErrIncapacidadNotFound = &AppError{
		Code:       "INCAPACIDAD_NOT_FOUND",
		Message:    "Incapacidad no encontrada",
		HTTPStatus: http.StatusNotFound,
	}

	ErrRolNotFound = &AppError{
		Code:       "ROL_NOT_FOUND",
		Message:    "Rol no encontrado",
		HTTPStatus: http.StatusNotFound,
	}

	ErrConflict = &AppError{
		Code:       "CONFLICT",
		Message:    "Conflicto de datos",
		HTTPStatus: http.StatusConflict,
	}

	ErrEmailAlreadyExists = &AppError{
		Code:       "EMAIL_ALREADY_EXISTS",
		Message:    "El email ya está registrado",
		HTTPStatus: http.StatusConflict,
	}

	ErrBadRequest = &AppError{
		Code:       "BAD_REQUEST",
		Message:    "Solicitud inválida",
		HTTPStatus: http.StatusBadRequest,
	}

	ErrValidation = &AppError{
		Code:       "VALIDATION_ERROR",
		Message:    "Error de validación",
		HTTPStatus: http.StatusUnprocessableEntity,
	}

	ErrInternal = &AppError{
		Code:       "INTERNAL_ERROR",
		Message:    "Error interno del servidor",
		HTTPStatus: http.StatusInternalServerError,
	}

	ErrDatabase = &AppError{
		Code:       "DATABASE_ERROR",
		Message:    "Error de base de datos",
		HTTPStatus: http.StatusInternalServerError,
	}

	ErrHashPassword = &AppError{
		Code:       "HASH_PASSWORD_ERROR",
		Message:    "Error al hashear contraseña",
		HTTPStatus: http.StatusInternalServerError,
	}
)

func New(code, message string, status int) *AppError {
	return &AppError{
		Code:       code,
		Message:    message,
		HTTPStatus: status,
	}
}

func Is(err, target error) bool {
	return errors.Is(err, target)
}