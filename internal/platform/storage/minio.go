package storage

import (
	"context"
	"encoding/base64"
	"fmt"
	"go-service-boilerplate/configs"
	"net/url"
	"strings"
	"time"

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

// UploadImageToMinIO upload base64 image to minio and return object path
func UploadImageToMinIO(minioClient *minio.Client, bucketName string, base64Image string, folderPrefix string, filePrefix string) (string, error) {
	parts := strings.SplitN(base64Image, ",", 2)
	if len(parts) != 2 {
		return "", fmt.Errorf("invalid base64 image format")
	}

	metaParts := strings.Split(parts[0], ";")
	if len(metaParts) < 1 {
		return "", fmt.Errorf("invalid base64 image metadata")
	}
	contentType := strings.TrimPrefix(metaParts[0], "data:")

	ext := "png" // default
	switch {
	case strings.Contains(contentType, "jpeg"):
		ext = "jpg"
	case strings.Contains(contentType, "gif"):
		ext = "gif"
	case strings.Contains(contentType, "webp"):
		ext = "webp"
	}

	base64Data := parts[1]
	imgBytes, err := base64.StdEncoding.DecodeString(base64Data)
	if err != nil {
		return "", fmt.Errorf("failed to decode base64: %w", err)
	}

	ctx := context.Background()
	objectName := fmt.Sprintf("%s/%s-%d.%s", folderPrefix, filePrefix, time.Now().UnixNano(), ext)
	reader := strings.NewReader(string(imgBytes)) // Convert byte slice to reader

	_, err = minioClient.PutObject(
		ctx,
		bucketName,
		objectName,
		reader,
		int64(len(imgBytes)),
		minio.PutObjectOptions{ContentType: contentType},
	)
	if err != nil {
		return "", fmt.Errorf("failed to upload to MinIO: %w", err)
	}

	return objectName, nil
}

// Genereate Presigned URL for object minio.
// return example: https://cdn.example.com/bucket/file.jpg?X-Amz-...
func GetPresignedObject(minioClient *minio.Client, bucketName, objectName, minioEndpoint, publicEndpoint string, expiry time.Duration) (string, error) {
	if expiry <= 0 {
		expiry = 15 * time.Minute
	}

	reqParams := make(url.Values)
	presignedUrl, err := minioClient.PresignedGetObject(
		context.Background(),
		bucketName,
		objectName,
		expiry,
		reqParams,
	)
	if err != nil {
		return "", fmt.Errorf("error when get presigned url: %w", err)
	}

	// Replace minio endpoint with public endpoint
	urlStr := strings.Replace(presignedUrl.String(), minioEndpoint, publicEndpoint, 1)
	return urlStr, nil
}
