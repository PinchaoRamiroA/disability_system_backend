package storage

import (
	"context"
)

func (c *Client) Delete(ctx context.Context, key string) error {
	if err := c.DeleteObject(ctx, key); err != nil {
		return ErrDeleteFailed.WithError(err)
	}
	return nil
}

func (c *Client) DeleteMultiple(ctx context.Context, keys []string) error {
	for _, key := range keys {
		if err := c.Delete(ctx, key); err != nil {
			return err
		}
	}
	return nil
}
