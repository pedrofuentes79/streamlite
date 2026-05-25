package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type MediaResponse struct {
	Title           string `json:"title"`
	Date            string `json:"date"`
	ProgressSeconds int64  `json:"progress_seconds"`
}

func (s *server) handleGetCatalog(w http.ResponseWriter, r *http.Request) {
	rows, err := s.db.Query("SELECT title, date, progress_seconds FROM media ORDER BY updated_at DESC LIMIT 10")
	if err != nil {
		http.Error(w, "Error querying the database", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	catalog := []MediaResponse{}
	for rows.Next() {
		var m MediaResponse
		if err := rows.Scan(&m.Title, &m.Date, &m.ProgressSeconds); err != nil {
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
