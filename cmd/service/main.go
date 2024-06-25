package main

import (
	"context"
	"github.com/joho/godotenv"
	"os"
	"quizzly/cmd"
	"quizzly/internal/quizzly"
	"quizzly/pkg/transactional"
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

	log := cmd.MustInitLogger()
	defer log.Flush()

	quizzlyConfig := quizzly.NewConfiguration(
		db,
		template,
	)

	web.ServerRun(quizzlyConfig)
}
