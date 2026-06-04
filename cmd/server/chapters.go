package main

import (
	"database/sql"
	"errors"
	"net/http"
)

// handleChapters returns a media item's chapter markers as a JSON array of
// {start, title} objects (seconds + label). It's fetched lazily by the player
// when an item is opened, so the catalog list stays lean. The column already
// holds a JSON string (written by the scraper's ingest), so we stream it
// straight through rather than decode/re-encode it.
func (s *server) handleChapters(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET is accepted here", http.StatusMethodNotAllowed)
		return
	}

	id, ok := parseMediaID(w, r)
	if !ok {
		return
	}

	var chapters string
	err := s.db.QueryRow("SELECT chapters FROM media WHERE id = ?", id).Scan(&chapters)
	if errors.Is(err, sql.ErrNoRows) {
		http.Error(w, "Media not found in database", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "Error querying the database", http.StatusInternalServerError)
		return
	}
	if chapters == "" {
		chapters = "[]"
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(chapters))
}
