package storage

import (
	"bytes"
	"context"
	"io"

	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type Client struct {
	s3Client *s3.Client
	bucket   string
	publicURL string
}

func NewClient(s3Client *s3.Client, bucket, publicURL string) *Client {
	return &Client{
		s3Client: s3Client,
		bucket:   bucket,
		publicURL: publicURL,
	}
}

func (c *Client) GetBucket() string {
	return c.bucket
}

func (c *Client) GetPublicURL() string {
	return c.publicURL
}

func (c *Client) PublicURL(key string) string {
	if c.publicURL == "" {
		return key
	}
	return c.publicURL + "/" + key
}

func (c *Client) PutObject(ctx context.Context, key string, body []byte, contentType string) error {
	input := &s3.PutObjectInput{
		Bucket:      &c.bucket,
		Key:         &key,
		Body:        bytes.NewReader(body),
		ContentType: &contentType,
	}

	_, err := c.s3Client.PutObject(ctx, input)
	return err
}

func (c *Client) DeleteObject(ctx context.Context, key string) error {
	input := &s3.DeleteObjectInput{
		Bucket: &c.bucket,
		Key:    &key,
	}

	_, err := c.s3Client.DeleteObject(ctx, input)
	return err
}

func (c *Client) GetObject(ctx context.Context, key string) ([]byte, error) {
	input := &s3.GetObjectInput{
		Bucket: &c.bucket,
		Key:    &key,
	}

	result, err := c.s3Client.GetObject(ctx, input)
	if err != nil {
		return nil, err
	}
	defer result.Body.Close()

	body, err := io.ReadAll(result.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func (c *Client) HeadObject(ctx context.Context, key string) (*s3.HeadObjectOutput, error) {
	input := &s3.HeadObjectInput{
		Bucket: &c.bucket,
		Key:    &key,
	}

	return c.s3Client.HeadObject(ctx, input)
}

func (c *Client) ObjectExists(ctx context.Context, key string) (bool, error) {
	_, err := c.HeadObject(ctx, key)
	if err != nil {
		return false, nil
	}
	return true, nil
}
