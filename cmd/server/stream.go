package main

import (
	"database/sql"
	"errors"
	"net/http"
	"os"
)

func (s *server) handleStream(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET is accepted in /api/stream", http.StatusMethodNotAllowed)
		return
	}

	id, ok := parseMediaID(w, r)
	if !ok {
		return
	}

	var videoPath string
	err := s.db.QueryRow("SELECT video_path FROM media WHERE id = ?", id).Scan(&videoPath)
	if errors.Is(err, sql.ErrNoRows) {
		http.Error(w, "Media not found in database", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "Error querying the database", http.StatusInternalServerError)
		return
	}

	f, err := os.Open(videoPath)
	if err != nil {
		http.Error(w, "Media file not found on disk", http.StatusNotFound)
		return
	}
	defer f.Close()

	// http.ServeContent handles range requests, Content-Type, and ETags
	fi, err := f.Stat()
	if err != nil {
		http.Error(w, "Could not read media file", http.StatusInternalServerError)
		return
	}
	http.ServeContent(w, r, fi.Name(), fi.ModTime(), f)
}
