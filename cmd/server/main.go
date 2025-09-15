package main

import (
	"approval-system/config"
	"approval-system/pkg/database"
	"approval-system/server"
	"context"
	"log"
	"os"
)

func main() {
	ctx := context.Background()

	cfg := config.LoadConfig(func() string {
		if len(os.Args) < 2 {
			log.Fatal("Error: .env path is required")
		}
		return os.Args[1]
	}())

	var dbPost database.DatabasesPostgres
	if cfg.Postgres != nil {
		dbPost = database.NewPostgresDatabase(cfg.Postgres)
	}

	server.Start(ctx, &cfg, dbPost)
}
