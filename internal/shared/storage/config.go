package storage

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
)

type R2Config struct {
	AccountID        string
	AccessKeyID      string
	SecretAccessKey  string
	Bucket           string
	Endpoint         string
	PublicURL        string
	MaxFileSize      int64
	PresignExpiry    time.Duration
}

func LoadR2Config() R2Config {
	return R2Config{
		AccountID:       getEnv("R2_ACCOUNT_ID", ""),
		AccessKeyID:     getEnv("R2_ACCESS_KEY_ID", ""),
		SecretAccessKey: getEnv("R2_SECRET_ACCESS_KEY", ""),
		Bucket:          getEnv("R2_BUCKET", ""),
		Endpoint:        getEnv("R2_ENDPOINT", ""),
		PublicURL:       getEnv("R2_PUBLIC_URL", ""),
		MaxFileSize:     parseFileSize(getEnv("R2_MAX_FILE_SIZE", "10MB")),
		PresignExpiry:   parseDuration(getEnv("R2_PRESIGN_EXPIRY", "15m")),
	}
}

func (c *R2Config) IsConfigured() bool {
	return c.AccountID != "" && c.AccessKeyID != "" && c.SecretAccessKey != "" && c.Bucket != ""
}

func (c *R2Config) AWSConfig(ctx context.Context) (aws.Config, error) {
	return config.LoadDefaultConfig(ctx,
		config.WithRegion("auto"),
		config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(
				c.AccessKeyID,
				c.SecretAccessKey,
				"",
			),
		),
	)
}

func (c *R2Config) EndpointURL() string {
	if c.Endpoint != "" {
		return c.Endpoint
	}
	return fmt.Sprintf("https://%s.r2.cloudflarestorage.com", c.AccountID)
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}

func parseFileSize(s string) int64 {
	if s == "" {
		return 10 * 1024 * 1024
	}

	if len(s) < 2 {
		return 10 * 1024 * 1024
	}

	unit := s[len(s)-2:]
	numStr := s[:len(s)-2]

	num, err := strconv.ParseInt(numStr, 10, 64)
	if err != nil {
		return 10 * 1024 * 1024
	}

	switch unit {
	case "KB":
		return num * 1024
	case "MB":
		return num * 1024 * 1024
	case "GB":
		return num * 1024 * 1024 * 1024
	default:
		return num
	}
}

func parseDuration(s string) time.Duration {
	d, err := time.ParseDuration(s)
	if err != nil {
		return 15 * time.Minute
	}
	return d
}
