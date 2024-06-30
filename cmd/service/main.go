package main

import (
	"context"
	"github.com/joho/godotenv"
	"os"
	"quizzly/cmd"
	"quizzly/internal/quizzly"
	"quizzly/pkg/auth"
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

	quizzlyConfig := quizzly.NewConfiguration(
		db,
		template,
	)
	simpleAuthConfig := auth.NewConfiguration(db, template, variables.Repository.MustGet())

	web.ServerRun(quizzlyConfig, simpleAuthConfig.SimpleAuth.MustGet())
}
