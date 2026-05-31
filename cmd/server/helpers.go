package main

import (
	"net/http"
	"strconv"
)

func parseMediaID(w http.ResponseWriter, r *http.Request) (int, bool) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "The `id` parameter must be a valid integer", http.StatusBadRequest)
		return 0, false
	}
	return id, true
}
