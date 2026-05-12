package storage

import (
	"errors"
	"fmt"
	"net/http"
)

var (
	ErrFileTooLarge       = &StorageError{Code: "FILE_TOO_LARGE", Message: "El archivo excede el tamaño máximo permitido"}
	ErrInvalidFileType   = &StorageError{Code: "INVALID_FILE_TYPE", Message: "Tipo de archivo no permitido"}
	ErrFileNotFound       = &StorageError{Code: "FILE_NOT_FOUND", Message: "Archivo no encontrado"}
	ErrUploadFailed       = &StorageError{Code: "UPLOAD_FAILED", Message: "Error al subir el archivo"}
	ErrDeleteFailed       = &StorageError{Code: "DELETE_FAILED", Message: "Error al eliminar el archivo"}
	ErrNotConfigured      = &StorageError{Code: "STORAGE_NOT_CONFIGURED", Message: "Almacenamiento no configurado"}
	ErrInvalidContentType = &StorageError{Code: "INVALID_CONTENT_TYPE", Message: "Tipo de contenido inválido"}
	ErrPresignFailed      = &StorageError{Code: "PRESIGN_FAILED", Message: "Error al generar URL prefirmada"}
)

type StorageError struct {
	Code       string
	Message    string
	Err        error
	HTTPStatus int
}

func (e *StorageError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s (%v)", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

func (e *StorageError) WithError(err error) *StorageError {
	return &StorageError{
		Code:       e.Code,
		Message:    e.Message,
		Err:        err,
		HTTPStatus: e.HTTPStatus,
	}
}

func (e *StorageError) WithMessage(msg string) *StorageError {
	return &StorageError{
		Code:       e.Code,
		Message:    msg,
		Err:        e.Err,
		HTTPStatus: e.HTTPStatus,
	}
}

func (e *StorageError) Status(status int) *StorageError {
	return &StorageError{
		Code:       e.Code,
		Message:    e.Message,
		Err:        e.Err,
		HTTPStatus: status,
	}
}

func (e *StorageError) Is(target error) bool {
	var se *StorageError
	if errors.As(target, &se) {
		return se.Code == e.Code
	}
	return false
}

func init() {
	ErrFileTooLarge.HTTPStatus = http.StatusRequestEntityTooLarge
	ErrInvalidFileType.HTTPStatus = http.StatusUnsupportedMediaType
	ErrFileNotFound.HTTPStatus = http.StatusNotFound
	ErrUploadFailed.HTTPStatus = http.StatusInternalServerError
	ErrDeleteFailed.HTTPStatus = http.StatusInternalServerError
	ErrNotConfigured.HTTPStatus = http.StatusServiceUnavailable
	ErrInvalidContentType.HTTPStatus = http.StatusUnsupportedMediaType
	ErrPresignFailed.HTTPStatus = http.StatusInternalServerError
}
