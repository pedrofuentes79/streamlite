package main

import (
	"flag"
	"encoding/json"
	"log"
	"net/http"
	"streamlite/internal/database"
)

func main() {
	// CLI flags
	dbPath := flag.String("db", "streamlite.db", "Path to the SQLite file")
	resetDB := flag.Bool("reset", false, "Drop the existing database and re-run migrations from scratch")
	flag.Parse()

	// Initialize db with migrations.
	database.InitDB(*dbPath, *resetDB)

	// Base endpoints
	fs := http.FileServer(http.Dir("./web"))
	http.Handle("/", fs)

	http.HandleFunc("/api/media", handleGetCatalog)
	http.HandleFunc("/api/stream", handleStream)
	http.HandleFunc("/api/progress", handleProgress)

	log.Println("🚀 streamlite running at http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

type MediaResponse struct {
	Title           string `json:"title"`
	Date            string `json:"date"`
	ProgressSeconds int64  `json:"progress_seconds"`
}


func handleGetCatalog(w http.ResponseWriter, r *http.Request) {
	rows, err := database.DB.Query("SELECT title, date, progress_seconds FROM media ORDER BY updated_at DESC limit 10")
	if err != nil {
		http.Error(w, "Querying the db raised an error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var catalog []MediaResponse

	for rows.Next() {
		var m MediaResponse
		err := rows.Scan(&m.Title, &m.Date, &m.ProgressSeconds)
		if err != nil {
			log.Fatalf("Error while scanning %v", err)
			http.Error(w, "Processing rows data raised an error.", http.StatusInternalServerError)
		}
		catalog = append(catalog, m)
	}

	if catalog == nil {
		// empty db
		catalog = []MediaResponse{}
	}

	// Set headers and return the struct
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(catalog); err != nil {
		log.Printf("Error encoding JSON response: %v", err)
	}


}
func handleStream(w http.ResponseWriter, r *http.Request)   {}
func handleProgress(w http.ResponseWriter, r *http.Request) {}
