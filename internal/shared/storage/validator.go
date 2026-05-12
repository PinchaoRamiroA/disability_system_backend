package storage

import (
	"io"
	"path/filepath"
	"strconv"
	"strings"
)

type FileValidator struct {
	allowedTypes map[string]bool
	maxSize      int64
}

func NewFileValidator(maxSize int64) *FileValidator {
	return &FileValidator{
		allowedTypes: map[string]bool{
			"application/pdf": true,
			"image/jpeg":      true,
			"image/jpg":       true,
			"image/png":       true,
		},
		maxSize: maxSize,
	}
}

func (v *FileValidator) Validate(contentType string, size int64) error {
	if size > v.maxSize {
		return ErrFileTooLarge.WithMessage(
			"El archivo excede el tamaño máximo de " + formatBytes(v.maxSize),
		)
	}

	if !v.IsAllowedType(contentType) {
		return ErrInvalidFileType.WithMessage(
			"Tipo de archivo no permitido. Solo se aceptan: PDF, JPG, PNG",
		)
	}

	return nil
}

func (v *FileValidator) IsAllowedType(contentType string) bool {
	contentType = strings.ToLower(contentType)
	return v.allowedTypes[contentType]
}

func (v *FileValidator) GetAllowedExtensions() []string {
	return []string{".pdf", ".jpg", ".jpeg", ".png"}
}

func (v *FileValidator) GetAllowedMimeTypes() []string {
	types := make([]string, 0, len(v.allowedTypes))
	for t := range v.allowedTypes {
		types = append(types, t)
	}
	return types
}

func (v *FileValidator) ValidateExtension(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	for _, allowed := range v.GetAllowedExtensions() {
		if ext == allowed {
			return true
		}
	}
	return false
}

func (v *FileValidator) GetMimeTypeFromExtension(filename string) string {
	ext := strings.ToLower(filepath.Ext(filename))
	switch ext {
	case ".pdf":
		return "application/pdf"
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".png":
		return "image/png"
	default:
		return ""
	}
}

func (v *FileValidator) MaxSize() int64 {
	return v.maxSize
}

type LimitedReader struct {
	reader io.Reader
	limit  int64
	read   int64
}

func NewLimitedReader(r io.Reader, maxSize int64) *LimitedReader {
	return &LimitedReader{
		reader: r,
		limit:  maxSize,
	}
}

func (lr *LimitedReader) Read(p []byte) (int, error) {
	if lr.read >= lr.limit {
		return 0, ErrFileTooLarge.WithMessage("Límite de tamaño excedido durante la lectura")
	}

	remaining := lr.limit - lr.read
	if int64(len(p)) > remaining {
		p = p[:remaining]
	}

	n, err := lr.reader.Read(p)
	lr.read += int64(n)

	return n, err
}

func formatBytes(bytes int64) string {
	const unit = int64(1024)
	if bytes < unit {
		return strconv.FormatInt(bytes, 10) + " B"
	}
	if bytes%unit == 0 {
		return strconv.FormatInt(bytes/unit, 10) + " KB"
	}
	f := float64(bytes) / float64(unit)
	return strconv.FormatFloat(f, 'f', 1, 64) + " KB"
}
