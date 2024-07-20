package files

import (
	"bytes"
	"context"
	"io"
)

type (
	S3Config struct {
		Endpoint        string
		AccessKeyID     string
		SecretAccessKey string
		UseSSL          bool
		BucketName      string
	}

	UploadFile struct {
		Data io.Reader
		Name string
		Size int64
	}

	File struct {
		Name string
		Url  string
	}

	Manager interface {
		Upload(ctx context.Context, in *UploadFile) error
		Get(ctx context.Context, name string) (*bytes.Buffer, error)
	}
)
