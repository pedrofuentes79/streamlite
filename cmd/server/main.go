package main

import (
	"flag"
	"json"
	"log"
	"net/http"
	"streamlite/internal/database"
)

func main() {
	// CLI flags
	dbPath := flag.String("db", "streamlite.db", "Ruta al archivo de SQLite")
	resetDB := flag.Bool("reset", false, "Borra la base de datos existente y aplica migraciones de cero")
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
	Day             int    `json:"day"`
	ProgressSeconds int64  `json:"progress_seconds"`
}


// Los stubs de los handlers quedan iguales por ahora...
func handleGetCatalog(w http.ResponseWriter, r *http.Request) {
	rows, err := database.DB.Query("SELECT title, day, progress_seconds FROM media ORDER BY updated_at DESC limit 10")
	if err != nil {
		http.Error(w, "Querying the db raised an error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var catalog []MediaResponse
	
	for rows.Next() {
		var m MediaResponse
		err := rows.Scan(&m.Title, &m.Day, &m.ProgressSeconds)
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
		log.Printf("Error al codificar JSON: %v", err)
	}
	

}
func handleStream(w http.ResponseWriter, r *http.Request)   {}
func handleProgress(w http.ResponseWriter, r *http.Request) {}
