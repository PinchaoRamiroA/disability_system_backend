package storage_test

import (
	"bytes"
	"context"
	"testing"

	"disability_system_backend/internal/shared/logger"
	"disability_system_backend/internal/shared/storage"

	"github.com/stretchr/testify/assert"
)

func TestStorage_GenerateKey(t *testing.T) {
	key := storage.GenerateKey(123, "certificado_medico.pdf")
	assert.Contains(t, key, "incapacidad/123/")
	assert.Contains(t, key, "certificado_medico.pdf")
}

func TestStorage_ClientWithLogger_BUG13(t *testing.T) {
	var buf bytes.Buffer
	log := logger.NewWithWriter(&buf, "info")

	client := storage.NewClient(nil, "test-bucket", "https://cdn.example.com", log)
	assert.Equal(t, "test-bucket", client.GetBucket())
	assert.Equal(t, "https://cdn.example.com", client.GetPublicURL())

	// Upload with nil s3 client returns error without panicking, logging structured error instead of log.Printf
	ctx := context.Background()
	res, err := client.Upload(ctx, []byte("test data"), "test.pdf", "application/pdf", 1)
	assert.Error(t, err)
	assert.Nil(t, res)
}
