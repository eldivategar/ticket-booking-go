package storage

import (
	"context"
	"fmt"
	"go-service-boilerplate/configs"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func NewMinIOClient(cfg configs.Config) (*minio.Client, error) {

	minioClient, err := minio.New(cfg.MinioEndpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.MinioAccessKey, cfg.MinioSecretKey, ""),
		Secure: cfg.MinioUseSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to initialize minio client: %w", err)
	}

	// Test connection and create bucket if not exists
	ctx := context.Background()
	err = minioClient.MakeBucket(ctx, cfg.MinioBucket, minio.MakeBucketOptions{Region: cfg.MinioRegion})
	if err != nil {
		exists, errBucketExists := minioClient.BucketExists(ctx, cfg.MinioBucket)
		if errBucketExists == nil && exists {
			// bucket is already exists
		} else {
			return nil, fmt.Errorf("failed to make or find bucket '%s': %w", cfg.MinioBucket, err)
		}
	}

	fmt.Printf("MinIO client initialized and bucket '%s' is ready\n", cfg.MinioBucket)

	return minioClient, nil
}
