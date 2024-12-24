package files

import (
	"bytes"
	"context"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/smithy-go/middleware"
	"quizzly/pkg/structs"
)

type S3Manager struct {
	client     *s3.Client
	bucketName string
}

func NewS3Manager(
	config *S3Config,
) (Manager, error) {
	s3Client := s3.New(s3.Options{
		BaseEndpoint:       structs.Pointer(config.Endpoint),
		EndpointResolverV2: &endpoint{},
		Credentials: credentials.NewStaticCredentialsProvider(
			config.AccessKeyID,
			config.SecretAccessKey,
			"",
		),
		APIOptions: []func(*middleware.Stack) error{
			func(stack *middleware.Stack) error {
				_, err := stack.Finalize.Remove("DisableAcceptEncodingGzip")
				return err
			},
		},
	})

	return &S3Manager{
		client:     s3Client,
		bucketName: config.BucketName,
	}, nil
}

func (m *S3Manager) Upload(ctx context.Context, in *UploadFile) error {
	_, err := m.client.PutObject(
		ctx,
		&s3.PutObjectInput{
			Bucket:        &m.bucketName,
			Body:          in.Data,
			Key:           &in.Name,
			ContentType:   structs.Pointer("application/octet-stream"),
			ContentLength: &in.Size,
		},
	)

	return err
}

func (m *S3Manager) Get(ctx context.Context, name string) (*bytes.Buffer, error) {
	file, err := m.client.GetObject(
		ctx,
		&s3.GetObjectInput{
			Bucket: &m.bucketName,
			Key:    &name,
		},
	)
	if err != nil {
		return nil, err
	}

	defer file.Body.Close()

	buffer := bytes.Buffer{}
	_, err = buffer.ReadFrom(file.Body)
	if err != nil {
		return nil, err
	}

	return &buffer, nil
}
