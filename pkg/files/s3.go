package files

import (
	"bytes"
	"context"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type S3Manager struct {
	client     *minio.Client
	bucketName string
}

func NewS3Manager(
	config *S3Config,
) (Manager, error) {
	minioClient, err := minio.New(config.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(config.AccessKeyID, config.SecretAccessKey, ""),
		Secure: config.UseSSL,
	})
	if err != nil {
		return nil, err
	}

	return &S3Manager{
		client:     minioClient,
		bucketName: config.BucketName,
	}, nil
}

func (m *S3Manager) Upload(ctx context.Context, in *UploadFile) error {
	_, err := m.client.PutObject(
		ctx,
		m.bucketName,
		in.Name,
		in.Data,
		in.Size,
		minio.PutObjectOptions{
			ContentType: "application/octet-stream",
		},
	)

	return err
}

func (m *S3Manager) Get(ctx context.Context, name string) (*bytes.Buffer, error) {
	file, err := m.client.GetObject(ctx, m.bucketName, name, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}

	defer file.Close()

	buffer := bytes.Buffer{}
	_, err = buffer.ReadFrom(file)
	if err != nil {
		return nil, err
	}

	return &buffer, nil
}
