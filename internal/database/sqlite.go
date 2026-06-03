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

//go:embed seed.sql
var seedScript string

// InitDB opens the database and runs schema migrations. Development seed data is
// applied only when seed is true (server -seed, or tests) — production DBs stay
// empty so the catalog isn't polluted with placeholder rows.
func InitDB(dbPath string, reset, seed bool) *sql.DB {
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
	if seed {
		if _, err = db.Exec(seedScript); err != nil {
			log.Fatalf("Failed to apply seed data: %v", err)
		}
	}
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
