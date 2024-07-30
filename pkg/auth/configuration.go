package auth

import (
	"github.com/jmoiron/sqlx"
	"quizzly/pkg/structs"
	"quizzly/pkg/transactional"
	"quizzly/pkg/variables"
)

type Configuration struct {
	SimpleAuth structs.Singleton[SimpleAuth]
}

func NewConfiguration(
	db *sqlx.DB,
	template transactional.Template,
	variablesRepo variables.Repository,
) *Configuration {
	return &Configuration{
		SimpleAuth: structs.NewSingleton(func() (auth SimpleAuth, err error) {
			return NewSimpleAuth(
				db,
				template,
				&EncryptorConfig{
					SecretKey: variablesRepo.GetString(variables.AuthSecretKey),
				},
				&SenderConfig{
					FromEmail: Email(variablesRepo.GetString(variables.AuthSenderFromEmail)),
					Host:      variablesRepo.GetString(variables.AuthSenderHost),
					Port:      variablesRepo.GetInt64(variables.AuthSenderPort),
					User:      variablesRepo.GetString(variables.AuthSenderUser),
					Password:  variablesRepo.GetString(variables.AuthSenderPassword),
					Debug:     variablesRepo.GetString(variables.AppEnvironmentVariable) == string(variables.EnvironmentLocal),
				},
			), nil
		}),
	}
}
