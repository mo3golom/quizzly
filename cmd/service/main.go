package main

import (
	"context"
	trmsqlx "github.com/avito-tech/go-transaction-manager/drivers/sqlx/v2"
	txmanager "github.com/avito-tech/go-transaction-manager/trm/v2/manager"
	"github.com/joho/godotenv"
	"os"
	"quizzly/cmd"
	"quizzly/internal/quizzly"
	"quizzly/pkg/auth"
	"quizzly/pkg/files"
	jobs2 "quizzly/pkg/jobs"

	variables2 "quizzly/pkg/variables"
	"quizzly/web"
)

func main() {
	ctx := context.Background()
	if _, err := os.Stat(".env"); err == nil {
		// path/to/whatever exists
		err := godotenv.Load(".env")
		if err != nil {
			panic(err)
		}
	}

	db := cmd.MustInitDB(ctx)
	trm := txmanager.Must(trmsqlx.NewDefaultFactory(db))
	variables, err := variables2.NewConfiguration()
	if err != nil {
		panic(err)
	}

	variablesRepo := variables.Repository.MustGet()

	log := cmd.MustInitLogger()
	defer log.Flush()

	jobs := jobs2.NewDefaultRunner(log)

	filesConfig := files.NewConfiguration(variables.Repository.MustGet())
	simpleAuth := auth.NewSimpleAuth(
		db,
		trm,
		&auth.Config{
			SecretKey:      variablesRepo.GetString(variables2.AuthSecretKey),
			CookieBlockKey: variablesRepo.GetString(variables2.AuthCookieBlockKey),
			FromEmail:      variablesRepo.GetString(variables2.AuthSenderFromEmail),
			Host:           variablesRepo.GetString(variables2.AuthSenderHost),
			Port:           variablesRepo.GetInt64(variables2.AuthSenderPort),
			User:           variablesRepo.GetString(variables2.AuthSenderUser),
			Password:       variablesRepo.GetString(variables2.AuthSenderPassword),
			Debug:          variablesRepo.GetString(variables2.AppEnvironmentVariable) == string(variables2.EnvironmentLocal),
		},
	)
	err = jobs.Register(simpleAuth.Cleaner())
	if err != nil {
		panic(err)
	}

	quizzlyConfig := quizzly.NewConfiguration(
		db,
		trm,
	)

	server := web.NewServer(
		log,
		variables.Repository.MustGet(),
		quizzlyConfig,
		simpleAuth,
		filesConfig.S3.MustGet(),
		web.ServerTypeHttp,
	)

	runner := cmd.NewRunner(log)
	runner.Start(
		server,
		jobs,
	)
}
