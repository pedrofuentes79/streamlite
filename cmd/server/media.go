package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type MediaResponse struct {
	ID              int64   `json:"id"`
	Title           string  `json:"title"`
	Date            string  `json:"date"`
	ProgressSeconds int64   `json:"progress_seconds"`
	TotalSeconds    float64 `json:"total_seconds"`
}

func (s *server) handleGetCatalog(w http.ResponseWriter, r *http.Request) {
	// Newest broadcast first. `date` is the video's publish date; ordering by
	// updated_at instead surfaced whatever row the ingest last touched, so the list
	// reshuffled every time a download landed and a back-filled old episode could
	// sit at the top. id DESC only breaks ties inside a single day.
	//
	// The limit sits comfortably above the downloader's KEEP_LATEST so the catalog
	// shows the whole library; at LIMIT 10 a 14-episode library hid four of them.
	rows, err := s.db.Query("SELECT id, title, date, progress_seconds, total_seconds FROM media ORDER BY date DESC, id DESC LIMIT 50")
	if err != nil {
		http.Error(w, "Error querying the database", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	catalog := []MediaResponse{}
	for rows.Next() {
		var m MediaResponse
		if err := rows.Scan(&m.ID, &m.Title, &m.Date, &m.ProgressSeconds, &m.TotalSeconds); err != nil {
			log.Printf("Error scanning row: %v", err)
			http.Error(w, "Error processing results", http.StatusInternalServerError)
			return
		}
		catalog = append(catalog, m)
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(catalog); err != nil {
		log.Printf("Error encoding JSON response: %v", err)
	}
}
