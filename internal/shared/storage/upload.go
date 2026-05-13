package storage

import (
	"context"
	"fmt"
	"log"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type UploadResult struct {
	Key      string
	URL      string
	Filename string
}

func (c *Client) Upload(ctx context.Context, data []byte, filename string, contentType string, incapacidadID uint64) (*UploadResult, error) {
	key := GenerateKey(incapacidadID, filename)

	if err := c.PutObject(ctx, key, data, contentType); err != nil {
		log.Printf("R2 Upload Error: %v", err)
		return nil, ErrUploadFailed.WithError(err)
	}

	return &UploadResult{
		Key:      key,
		URL:      c.PublicURL(key),
		Filename: filename,
	}, nil
}

func (c *Client) UploadWithValidation(ctx context.Context, data []byte, filename string, contentType string, size int64, incapacidadID uint64, validator *FileValidator) (*UploadResult, error) {
	if err := validator.Validate(contentType, size); err != nil {
		return nil, err
	}

	return c.Upload(ctx, data, filename, contentType, incapacidadID)
}

func GenerateKey(incapacidadID uint64, filename string) string {
	ext := filepath.Ext(filename)
	originalName := strings.TrimSuffix(filename, ext)
	originalName = sanitizeFilename(originalName)
	timestamp := time.Now().Unix()
	
	return fmt.Sprintf("incapacidad/%d/%d_%s%s", incapacidadID, timestamp, originalName, ext)
}

func sanitizeFilename(filename string) string {
	filename = strings.ReplaceAll(filename, " ", "_")
	filename = strings.ReplaceAll(filename, "/", "_")
	filename = strings.ReplaceAll(filename, "\\", "_")
	filename = strings.ReplaceAll(filename, ":", "_")
	filename = strings.ReplaceAll(filename, "*", "_")
	filename = strings.ReplaceAll(filename, "?", "_")
	filename = strings.ReplaceAll(filename, "\"", "_")
	filename = strings.ReplaceAll(filename, "<", "_")
	filename = strings.ReplaceAll(filename, ">", "_")
	filename = strings.ReplaceAll(filename, "|", "_")
	
	filename = strings.ToLower(filename)
	
	if len(filename) > 50 {
		filename = filename[:50]
	}
	
	return filename
}

func GenerateUniqueFilename(originalName, extension string) string {
	timestamp := strconv.FormatInt(time.Now().UnixNano(), 10)
	return fmt.Sprintf("%s_%s%s", originalName, timestamp[len(timestamp)-8:], extension)
}
