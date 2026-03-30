package migrations

import (
	"log"
	"os"

	"otas/config"
)

func Run() {
	query, err := os.ReadFile("migrations/001_init.sql")
	if err != nil {
		log.Fatal("Could not read migration file:", err)
	}

	if _, err := config.DB.Exec(string(query)); err != nil {
		log.Fatal("Migration failed:", err)
	}

	log.Println("Migrations ran successfully")
}
