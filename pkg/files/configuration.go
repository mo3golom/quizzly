package files

import (
	"quizzly/pkg/structs"
	"quizzly/pkg/variables"
)

type Configuration struct {
	S3 structs.Singleton[Manager]
}

func NewConfiguration(
	variablesRepo variables.Repository,
) *Configuration {
	return &Configuration{
		S3: structs.NewSingleton(func() (Manager, error) {
			return NewS3Manager(&S3Config{
				Endpoint:        variablesRepo.GetString(variables.S3Endpoint),
				AccessKeyID:     variablesRepo.GetString(variables.S3AccessKey),
				SecretAccessKey: variablesRepo.GetString(variables.S3SecretKey),
				BucketName:      variablesRepo.GetString(variables.S3Bucket),
				UseSSL:          variablesRepo.GetBool(variables.S3UseSSL),
			})
		}),
	}
}
