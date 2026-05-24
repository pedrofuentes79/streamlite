package database

import (
	"database/sql"
	"embed"
	"log"
	"os"

	_ "github.com/ncruces/go-sqlite3/driver"
)

// Embebe la carpeta de migraciones dentro del binario
// 
//go:embed migrations/*.sql

var migrationFiles embed.FS

var DB *sql.DB

func InitDB(dbPath string, reset bool) {
	if reset {
		log.Println("Resetting full db")
		os.Remove(dbPath)
	}
	var err error
	DB, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatalf("Could not find sqlite3 file at %s with error %v", dbPath, err)
	}

	err = DB.Ping()
	if err != nil {
		log.Fatalf("Could not reach DB: %v", err)
	}

	runMigrations()
}

func runMigrations() {
	script, err := migrationFiles.ReadFile("migrations/1_init.sql")
	if err != nil {
		log.Fatalf("Could not find migrations file.")
	}

	_, err = DB.Exec(string(script))
	if err != nil {
		log.Fatalf("Failed to execute migrations with: %v", err)
	}
}
