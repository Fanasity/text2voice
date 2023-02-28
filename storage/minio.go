package storage

import (
	"context"
	"io"
	"log"

	"github.com/Fanasity/text2voice/model"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func NewClient(config model.MinIOConfig) (*minio.Client, error) {
	// Initialize minio client object.
	return minio.New(config.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(config.AccessKey, config.Secret, ""),
		Secure: config.UseSecure,
	})
}

func InitBucket(ctx context.Context, minioClient *minio.Client, bucketName string) error {
	exists, err := minioClient.BucketExists(ctx, bucketName)
	if err == nil && exists {
		log.Printf("We already own %s\n", bucketName)
	} else {
		err = minioClient.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{Region: ""})
	}
	return err
}

func UploadFile(ctx context.Context, minioClient *minio.Client, bucketName, objectName, filePath, contentType string) (string, error) {

	// Upload the zip file with FPutObject
	info, err := minioClient.FPutObject(ctx, bucketName, objectName, filePath, minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		return "", err
	}
	return info.Key, nil
}

func UploadObject(ctx context.Context, minioClient *minio.Client, bucketName, objectName, contentType string, content io.Reader, objectSize int64) (string, error) {
	// Upload the file with io
	info, err := minioClient.PutObject(ctx, bucketName, objectName, content, objectSize, minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		return "", err
	}
	return info.Key, nil
}

func DownFile(ctx context.Context, minioClient *minio.Client, bucketName, objectName, filePath string) (string, error) {
	err := minioClient.FGetObject(ctx, bucketName, objectName, filePath, minio.GetObjectOptions{})
	if err != nil {
		return "", err
	}
	return objectName, nil
}

func DownObject(ctx context.Context, minioClient *minio.Client, bucketName, objectName string) (*minio.Object, error) {
	return minioClient.GetObject(ctx, bucketName, objectName, minio.GetObjectOptions{})
}
