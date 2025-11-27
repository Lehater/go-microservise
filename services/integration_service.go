// services/integration_service.go
package services

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// IntegrationService — обёртка над MinIO/S3.
// Для ДЗ достаточно одного метода UploadTestObject.
type IntegrationService struct {
	client *minio.Client
	bucket string
}

func NewIntegrationService(endpoint, accessKey, secretKey, bucket string, useSSL bool) (*IntegrationService, error) {
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return nil, err
	}

	svc := &IntegrationService{
		client: client,
		bucket: bucket,
	}
	return svc, nil
}

// UploadTestObject — простая операция, чтобы продемонстрировать работоспособность.
// В реальном проекте сюда можно передавать audit-логи и т.п.
func (s *IntegrationService) UploadTestObject(ctx context.Context, objectName string, content string) error {
	reader := bytes.NewReader([]byte(content))
	_, err := s.client.PutObject(ctx, s.bucket, objectName, reader, int64(reader.Len()), minio.PutObjectOptions{
		ContentType: "text/plain",
	})
	if err != nil {
		return err
	}
	log.Printf("[MINIO] uploaded object %s to bucket %s\n", objectName, s.bucket)
	return nil
}

// Helper для чтения объекта (если нужно):
func (s *IntegrationService) GetObjectContent(ctx context.Context, objectName string) (string, error) {
	obj, err := s.client.GetObject(ctx, s.bucket, objectName, minio.GetObjectOptions{})
	if err != nil {
		return "", err
	}
	defer obj.Close()

	data, err := io.ReadAll(obj)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func (s *IntegrationService) EnsureBucket(ctx context.Context) error {
	exists, err := s.client.BucketExists(ctx, s.bucket)
	if err != nil {
		return err
	}
	if !exists {
		if err := s.client.MakeBucket(ctx, s.bucket, minio.MakeBucketOptions{}); err != nil {
			return fmt.Errorf("make bucket: %w", err)
		}
	}
	return nil
}
