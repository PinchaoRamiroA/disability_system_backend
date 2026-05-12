package storage

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type PresignedResult struct {
	URL       string
	Key      string
	ExpiresAt time.Time
}

func (c *Client) GenerateUploadURL(ctx context.Context, filename string, contentType string, incapacidadID uint64, expiry time.Duration) (*PresignedResult, error) {
	key := GenerateKey(incapacidadID, filename)

	presignClient := s3.NewPresignClient(c.s3Client)

	request, err := presignClient.PresignPutObject(ctx, &s3.PutObjectInput{
		Bucket:      &c.bucket,
		Key:         &key,
		ContentType: &contentType,
	}, s3.WithPresignExpires(expiry))
	if err != nil {
		return nil, ErrPresignFailed.WithError(err)
	}

	return &PresignedResult{
		URL:       request.URL,
		Key:       key,
		ExpiresAt: time.Now().Add(expiry),
	}, nil
}

func (c *Client) GenerateDownloadURL(ctx context.Context, key string, expiry time.Duration) (string, error) {
	presignClient := s3.NewPresignClient(c.s3Client)

	request, err := presignClient.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: &c.bucket,
		Key:    &key,
	}, s3.WithPresignExpires(expiry))
	if err != nil {
		return "", ErrPresignFailed.WithError(err)
	}

	return request.URL, nil
}

func (c *Client) GenerateDownloadURLWithFilename(ctx context.Context, key string, filename string, expiry time.Duration) (string, error) {
	presignClient := s3.NewPresignClient(c.s3Client)

	downloadFilename := fmt.Sprintf("attachment; filename=\"%s\"", filename)

	request, err := presignClient.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: &c.bucket,
		Key:    &key,
	}, s3.WithPresignExpires(expiry))
	if err != nil {
		return "", ErrPresignFailed.WithError(err)
	}

	return request.URL + "&response-content-disposition=" + downloadFilename, nil
}
