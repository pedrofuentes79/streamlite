package database

import (
	"database/sql"
	"embed"
	"log"
	"os"

	_ "github.com/ncruces/go-sqlite3/driver"
)

//go:embed migrations/*.sql
var migrationFiles embed.FS

func InitDB(dbPath string, reset bool) *sql.DB {
	if reset {
		log.Println("Resetting full db")
		os.Remove(dbPath)
	}

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatalf("Could not open sqlite3 file at %s: %v", dbPath, err)
	}
	if err = db.Ping(); err != nil {
		log.Fatalf("Could not reach DB: %v", err)
	}

	runMigrations(db)
	return db
}

func runMigrations(db *sql.DB) {
	entries, err := migrationFiles.ReadDir("migrations")
	if err != nil {
		log.Fatalf("Could not read migrations directory: %v", err)
	}

	for _, entry := range entries {
		script, err := migrationFiles.ReadFile("migrations/" + entry.Name())
		if err != nil {
			log.Fatalf("Could not read migration file %s: %v", entry.Name(), err)
		}
		if _, err = db.Exec(string(script)); err != nil {
			log.Fatalf("Failed to execute migration %s: %v", entry.Name(), err)
		}
	}
}
