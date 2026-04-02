package migrations

import (
	"log"
	"os"

	"otas/config"
)

func Run() {
	files := []string{
		"migrations/001_init.sql",
		"migrations/002_transactions.sql",
	}

	for _, file := range files {
		query, err := os.ReadFile(file)
		if err != nil {
			log.Fatalf("Could not read migration file %s: %v", file, err)
		}

		if _, err := config.DB.Exec(string(query)); err != nil {
			log.Fatalf("Migration failed for %s: %v", file, err)
		}

		log.Printf("Migrations ran successfully: %s", file)

	}
}
