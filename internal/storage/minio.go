package storage

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"github.com/minio/minio-go/v7"
	"net/url"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// MinIOClient (4.2.1) 封装了 MinIO 客户端
type MinIOClient struct {
	client   *minio.Client
	endpoint string
	useSSL   bool
}

func NewMinIOClient(endpoint, accessKey, secretKey string, useSSL bool) (*MinIOClient, error) {
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return nil, err
	}
	
	log.Println("INFO: MinIO client connected")
	return &MinIOClient{client: minioClient, endpoint: endpoint, useSSL: useSSL}, nil
}

// Upload (4.2.3) 上传文件到 MinIO
func (s *MinIOClient) Upload(ctx context.Context, bucketName, objectName string, data []byte, contentType string) (string, error) {
	
	// (确保 Bucket 存在)
	err := s.client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
	if err != nil {
		exists, errBucketExists := s.client.BucketExists(ctx, bucketName)
		if errBucketExists != nil || !exists {
			return "", fmt.Errorf("failed to check/create bucket: %w", err)
		}
	}

	reader := bytes.NewReader(data)
	
	info, err := s.client.PutObject(ctx, bucketName, objectName, reader, int64(len(data)), minio.PutObjectOptions{
		ContentType: contentType,
	})
	
	if err != nil {
		return "", err
	}

	log.Printf("INFO: Successfully uploaded %s to MinIO. Size: %d", objectName, info.Size)
	
	// 返回可公开访问的 URL (或签名 URL)
	url := &url.URL{
		Scheme: "http",
		Host:   s.endpoint,
		Path:   fmt.Sprintf("/%s/%s", bucketName, objectName),
	}
	if s.useSSL {
		url.Scheme = "https"
	}
	return url.String(), nil
}
