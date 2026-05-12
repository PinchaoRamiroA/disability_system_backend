package storage

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type StorageService struct {
	client    *Client
	config    *R2Config
	validator *FileValidator
}

func NewStorageService(ctx context.Context, cfg R2Config) (*StorageService, error) {
	if !cfg.IsConfigured() {
		return nil, ErrNotConfigured
	}

	awsCfg, err := cfg.AWSConfig(ctx)
	if err != nil {
		return nil, ErrNotConfigured.WithError(err)
	}

	s3Client := s3.NewFromConfig(awsCfg, func(o *s3.Options) {
		o.UsePathStyle = true
	})

	client := NewClient(s3Client, cfg.Bucket, cfg.PublicURL)
	validator := NewFileValidator(cfg.MaxFileSize)

	return &StorageService{
		client:    client,
		config:    &cfg,
		validator: validator,
	}, nil
}

func (s *StorageService) Upload(ctx context.Context, data []byte, filename string, contentType string, incapacidadID uint64) (*UploadResult, error) {
	return s.client.Upload(ctx, data, filename, contentType, incapacidadID)
}

func (s *StorageService) UploadWithValidation(ctx context.Context, data []byte, filename string, contentType string, size int64, incapacidadID uint64) (*UploadResult, error) {
	return s.client.UploadWithValidation(ctx, data, filename, contentType, size, incapacidadID, s.validator)
}

func (s *StorageService) Delete(ctx context.Context, key string) error {
	return s.client.Delete(ctx, key)
}

func (s *StorageService) DeleteMultiple(ctx context.Context, keys []string) error {
	return s.client.DeleteMultiple(ctx, keys)
}

func (s *StorageService) GenerateUploadURL(ctx context.Context, filename string, contentType string, incapacidadID uint64) (*PresignedResult, error) {
	return s.client.GenerateUploadURL(ctx, filename, contentType, incapacidadID, s.config.PresignExpiry)
}

func (s *StorageService) GenerateDownloadURL(ctx context.Context, key string) (string, error) {
	return s.client.GenerateDownloadURL(ctx, key, s.config.PresignExpiry)
}

func (s *StorageService) GenerateDownloadURLWithFilename(ctx context.Context, key string, filename string) (string, error) {
	return s.client.GenerateDownloadURLWithFilename(ctx, key, filename, s.config.PresignExpiry)
}

func (s *StorageService) GetValidator() *FileValidator {
	return s.validator
}

func (s *StorageService) IsConfigured() bool {
	return s.config.IsConfigured()
}

func (s *StorageService) GetPublicURL(key string) string {
	return s.client.PublicURL(key)
}

func (s *StorageService) Validate(contentType string, size int64) error {
	return s.validator.Validate(contentType, size)
}

func (s *StorageService) GetAllowedMimeTypes() []string {
	return s.validator.GetAllowedMimeTypes()
}

func (s *StorageService) GetMaxFileSize() int64 {
	return s.config.MaxFileSize
}
