package main

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"os"
)

// serveMediaColumn looks up a file path stored in the given media column
// (e.g. "video_path" or "audio_path") and streams it with full HTTP range
// support via http.ServeContent (handles 206 Partial Content, ETags, etc.).
func (s *server) serveMediaColumn(w http.ResponseWriter, r *http.Request, column string) {
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET is accepted here", http.StatusMethodNotAllowed)
		return
	}

	id, ok := parseMediaID(w, r)
	if !ok {
		return
	}

	var path string
	query := fmt.Sprintf("SELECT %s FROM media WHERE id = ?", column)
	err := s.db.QueryRow(query, id).Scan(&path)
	if errors.Is(err, sql.ErrNoRows) {
		http.Error(w, "Media not found in database", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "Error querying the database", http.StatusInternalServerError)
		return
	}

	f, err := os.Open(path)
	if err != nil {
		http.Error(w, "Media file not found on disk", http.StatusNotFound)
		return
	}
	defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		http.Error(w, "Could not read media file", http.StatusInternalServerError)
		return
	}
	http.ServeContent(w, r, fi.Name(), fi.ModTime(), f)
}

func (s *server) handleStream(w http.ResponseWriter, r *http.Request) {
	s.serveMediaColumn(w, r, "video_path")
}

func (s *server) handleAudio(w http.ResponseWriter, r *http.Request) {
	s.serveMediaColumn(w, r, "audio_path")
}
