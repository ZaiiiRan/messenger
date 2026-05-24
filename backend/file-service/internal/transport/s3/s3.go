package s3

import (
	"bytes"
	"context"
	"fmt"
	"time"

	"github.com/ZaiiiRan/messenger/backend/file-service/internal/config/settings"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type S3Client struct {
	client     *minio.Client
	bucket     string
	presignTTL time.Duration
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
		client:     client,
		bucket:     cfg.Bucket,
		presignTTL: time.Duration(cfg.PresignTTL) * time.Second,
	}, nil
}

func (c *S3Client) PresignUpload(ctx context.Context, key string, contentType string, maxSize int64) (string, map[string]string, error) {
	policy := minio.NewPostPolicy()
	if err := policy.SetBucket(c.bucket); err != nil {
		return "", nil, fmt.Errorf("set bucket: %w", err)
	}
	if err := policy.SetKey(key); err != nil {
		return "", nil, fmt.Errorf("set key: %w", err)
	}
	if err := policy.SetExpires(time.Now().Add(c.presignTTL)); err != nil {
		return "", nil, fmt.Errorf("set expires: %w", err)
	}
	if err := policy.SetContentType(contentType); err != nil {
		return "", nil, fmt.Errorf("set content type: %w", err)
	}
	if err := policy.SetContentLengthRange(1, maxSize); err != nil {
		return "", nil, fmt.Errorf("set content length range: %w", err)
	}

	url, formData, err := c.client.PresignedPostPolicy(ctx, policy)
	if err != nil {
		return "", nil, fmt.Errorf("presign upload: %w", err)
	}

	return url.String(), formData, nil
}

func (c *S3Client) PresignDownload(ctx context.Context, key string) (string, error) {
	url, err := c.client.PresignedGetObject(ctx, c.bucket, key, c.presignTTL, nil)
	if err != nil {
		return "", fmt.Errorf("presign download: %w", err)
	}
	return url.String(), nil
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

	info, err := obj.Stat()
	if err != nil {
		return nil, fmt.Errorf("stat object: %w", err)
	}

	buf := make([]byte, info.Size)
	if _, err := obj.Read(buf); err != nil {
		return nil, fmt.Errorf("read object: %w", err)
	}

	return buf, nil
}

func (c *S3Client) Delete(ctx context.Context, key string) error {
	if err := c.client.RemoveObject(ctx, c.bucket, key, minio.RemoveObjectOptions{}); err != nil {
		return fmt.Errorf("delete object: %w", err)
	}
	return nil
}
