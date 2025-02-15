package files

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"quizzly/pkg/structs"

	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/smithy-go/middleware"
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

	// Optionally verify bucket exists
	_, err := s3Client.HeadBucket(context.Background(), &s3.HeadBucketInput{
		Bucket: &config.BucketName,
	})
	if err != nil {
		return nil, fmt.Errorf("bucket %s not found or not accessible: %w", config.BucketName, err)
	}

	return &S3Manager{
		client:     s3Client,
		bucketName: config.BucketName,
	}, nil
}

func (m *S3Manager) Upload(ctx context.Context, in *UploadFile) error {
	// Read all data into a buffer to make it seekable
	buffer := bytes.NewBuffer(make([]byte, 0))
	_, err := io.Copy(buffer, in.Data)
	if err != nil {
		return err
	}

	_, err = m.client.PutObject(
		ctx,
		&s3.PutObjectInput{
			Bucket:        &m.bucketName,
			Body:          bytes.NewReader(buffer.Bytes()),
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
