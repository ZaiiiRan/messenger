package s3

import (
	"bytes"
	"context"
	"fmt"
	"io"

	"github.com/ZaiiiRan/messenger/backend/media-service/internal/config/settings"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type S3Client struct {
	client *minio.Client
	bucket string
}

func New(ctx context.Context, cfg settings.S3Settings) (*S3Client, error) {
	client, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKeyID, cfg.SecretAccessKey, cfg.SecretAccessToken),
		Secure: cfg.UseSSL,
		Region: cfg.Region,
	})
	if err != nil {
		return nil, fmt.Errorf("new s3 client: %w", err)
	}

	exists, err := client.BucketExists(ctx, cfg.Bucket)
	if err != nil {
		return nil, fmt.Errorf("check bucket exists: %w", err)
	}
	if !exists {
		return nil, fmt.Errorf("bucket %s does not exist", cfg.Bucket)
	}

	return &S3Client{
		client: client,
		bucket: cfg.Bucket,
	}, nil
}

func (c *S3Client) Stat(ctx context.Context, key string) (*ObjectInfo, error) {
	info, err := c.client.StatObject(ctx, c.bucket, key, minio.StatObjectOptions{})
	if err != nil {
		errResp := minio.ToErrorResponse(err)
		if errResp.Code == "NoSuchKey" {
			return nil, nil
		}
		return nil, fmt.Errorf("stat object: %w", err)
	}
	return &ObjectInfo{
		Size:        info.Size,
		ContentType: info.ContentType,
	}, nil
}

func (c *S3Client) Put(ctx context.Context, key string, data []byte, contentType string) error {
	_, err := c.client.PutObject(
		ctx,
		c.bucket,
		key,
		bytes.NewReader(data),
		int64(len(data)),
		minio.PutObjectOptions{ContentType: contentType},
	)
	if err != nil {
		return fmt.Errorf("put object: %w", err)
	}
	return nil
}

func (c *S3Client) Get(ctx context.Context, key string) ([]byte, error) {
	obj, err := c.client.GetObject(ctx, c.bucket, key, minio.GetObjectOptions{})
	if err != nil {
		return nil, fmt.Errorf("get object: %w", err)
	}
	defer obj.Close()

	data, err := io.ReadAll(obj)
	if err != nil {
		return nil, fmt.Errorf("read object: %w", err)
	}

	return data, nil
}

func (c *S3Client) Delete(ctx context.Context, key string) error {
	if err := c.client.RemoveObject(ctx, c.bucket, key, minio.RemoveObjectOptions{}); err != nil {
		return fmt.Errorf("delete object: %w", err)
	}
	return nil
}

func (c *S3Client) DeleteMany(ctx context.Context, keys []string) error {
	objectsCh := make(chan minio.ObjectInfo)
	go func() {
		defer close(objectsCh)
		for _, key := range keys {
			select {
			case objectsCh <- minio.ObjectInfo{Key: key}:
			case <-ctx.Done():
				return
			}
		}
	}()

	var firstErr error
	for removeErr := range c.client.RemoveObjects(ctx, c.bucket, objectsCh, minio.RemoveObjectsOptions{}) {
		if minio.ToErrorResponse(removeErr.Err).Code == "NoSuchKey" {
			continue
		}
		if firstErr == nil {
			firstErr = removeErr.Err
		}
	}
	if firstErr != nil {
		return fmt.Errorf("delete objects: %w", firstErr)
	}
	return nil
}
