package main

import (
	"database/sql"
	"flag"
	"log"
	"net/http"
	"streamlite/internal/database"
)

type server struct {
	db *sql.DB
}

func main() {
	dbPath := flag.String("db", "streamlite.db", "Path to the SQLite file")
	resetDB := flag.Bool("reset", false, "Drop the existing database and re-run migrations from scratch")
	flag.Parse()

	s := &server{db: database.InitDB(*dbPath, *resetDB)}

	fs := http.FileServer(http.Dir("./web"))
	http.Handle("/", fs)

	http.HandleFunc("/api/media", s.handleGetCatalog)
	http.HandleFunc("/api/stream/{id}", s.handleStream)
	http.HandleFunc("/api/progress/{id}", s.handleProgress)

	log.Println("streamlite running at http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
