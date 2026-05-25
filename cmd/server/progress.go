package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	sqlite3 "github.com/ncruces/go-sqlite3"
)

type UpdateFieldRequest struct {
	NewProgress int `json:"new_progress"`
}

func (s *server) handleProgress(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch {
		http.Error(w, "Only PATCH is allowed here", http.StatusMethodNotAllowed)
		return
	}

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "The `id` parameter must be a valid integer", http.StatusBadRequest)
		return
	}

	var req UpdateFieldRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	result, err := s.db.Exec("UPDATE media SET progress_seconds = ? WHERE id = ?", req.NewProgress, id)
	if err != nil {
		var sqliteErr *sqlite3.Error
		if errors.As(err, &sqliteErr) && sqliteErr.ExtendedCode() == sqlite3.CONSTRAINT_CHECK {
			http.Error(w, "Progress exceeds total duration", http.StatusBadRequest)
			return
		}
		http.Error(w, "Error updating the database", http.StatusInternalServerError)
		return
	}

	affected, err := result.RowsAffected()
	if err != nil {
		http.Error(w, "Error checking the result", http.StatusInternalServerError)
		return
	}
	if affected < 1 {
		http.Error(w, "No rows found for that id", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "Progress updated successfully"}`))
}
