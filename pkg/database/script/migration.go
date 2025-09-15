package main

import (
	"approval-system/config"
	migration "approval-system/pkg/database/migration"
	"context"
	"log"
	"os"
)

func main() {
	ctx := context.Background()
	_ = ctx

	cfg := config.LoadConfig(func() string {
		if len(os.Args) < 2 {
			log.Fatal("Error: .env path is required")
		}
		return os.Args[1]
	}())

	switch cfg.App.Name {
	case "doc":
		migration.DocumentMigration(&cfg)
	}
}
