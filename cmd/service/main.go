package main

import (
	"context"
	"github.com/joho/godotenv"
	"os"
	"quizzly/cmd"
	"quizzly/internal/quizzly"
	"quizzly/pkg/auth"
	"quizzly/pkg/files"
	"quizzly/pkg/transactional"
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
	template := transactional.NewTemplate(db)
	variables, err := variables2.NewConfiguration()
	if err != nil {
		panic(err)
	}

	log := cmd.MustInitLogger()
	defer log.Flush()

	filesConfig := files.NewConfiguration(variables.Repository.MustGet())
	simpleAuthConfig := auth.NewConfiguration(db, template, variables.Repository.MustGet())

	quizzlyConfig := quizzly.NewConfiguration(
		db,
		template,
	)

	web.ServerRun(
		log,
		quizzlyConfig,
		simpleAuthConfig.SimpleAuth.MustGet(),
		filesConfig.S3.MustGet(),
	)
}
