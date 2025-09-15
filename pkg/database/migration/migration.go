package migration

import (
	"log"
	"approval-system/config"
	"approval-system/internal/domain"
	"approval-system/pkg/database"
)

func DocumentMigration(cfg *config.Config) {
	db := database.NewPostgresDatabase(cfg.Postgres).Connect()

	if err := db.AutoMigrate(&domain.Document{}); err != nil {
		panic(err)
	}

	db.Exec(`CREATE INDEX IF NOT EXISTS idx_document_status ON documents (status)`)
	db.Exec(`CREATE INDEX IF NOT EXISTS idx_document_created ON documents (created_at DESC)`)
	db.Exec(`CREATE INDEX IF NOT EXISTS idx_document_status_created ON documents (status, created_at DESC)`)

	log.Println("✅ Document migration completed successfully")
}