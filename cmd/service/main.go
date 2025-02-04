package main

import (
	"context"
	"os"
	"quizzly/cmd"
	"quizzly/internal/quizzly"
	"quizzly/pkg/cookie"
	"quizzly/pkg/files"
	"quizzly/pkg/supabase"

	trmsqlx "github.com/avito-tech/go-transaction-manager/drivers/sqlx/v2"
	txmanager "github.com/avito-tech/go-transaction-manager/trm/v2/manager"
	"github.com/joho/godotenv"
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

	log := cmd.MustInitLogger()
	defer log.Flush()

	cookieService := cookie.NewService(variables.Repository.MustGet())

	filesConfig := files.NewConfiguration(variables.Repository.MustGet())
	authClient := supabase.NewAuth(cookieService, variables.Repository.MustGet())

	quizzlyConfig := quizzly.NewConfiguration(
		db,
		trm,
	)

	server := web.NewServer(
		log,
		variables.Repository.MustGet(),
		quizzlyConfig,
		authClient,
		cookieService,
		filesConfig.S3.MustGet(),
		web.ServerTypeHttp,
	)

	runner := cmd.NewRunner(log)
	runner.Start(
		server,
	)
}
