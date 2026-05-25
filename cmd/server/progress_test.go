package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHandleProgress_OK(t *testing.T) {
	s := newTestServer(t)

	body := strings.NewReader(`{"new_progress": 120}`)
	req := httptest.NewRequest(http.MethodPatch, "/api/progress/1", body)
	req.SetPathValue("id", "1")
	rec := httptest.NewRecorder()

	s.handleProgress(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)

	var progress int
	err := s.db.QueryRow("SELECT progress_seconds FROM media WHERE id = ?", 1).Scan(&progress)
	assert.NoError(t, err)
	assert.Equal(t, 120, progress)

	var resp map[string]string
	assert.NoError(t, json.NewDecoder(rec.Body).Decode(&resp))
	assert.Equal(t, "Progress updated successfully", resp["message"])
}

func TestHandleProgress_UnknownID(t *testing.T) {
	s := newTestServer(t)

	body := strings.NewReader(`{"new_progress": 120}`)
	req := httptest.NewRequest(http.MethodPatch, "/api/progress/100", body)
	req.SetPathValue("id", "100")
	rec := httptest.NewRecorder()

	s.handleProgress(rec, req)
	assert.Equal(t, http.StatusBadRequest, rec.Code)

	var progress int
	err := s.db.QueryRow("SELECT progress_seconds FROM media WHERE id = ?", 100).Scan(&progress)
	assert.Equal(t, sql.ErrNoRows, err)
}

func TestHandleProgress_NewProgressHigherThanTotalSeconds(t *testing.T) {
	s := newTestServer(t)

	body := strings.NewReader(`{"new_progress": 9999999}`)
	req := httptest.NewRequest(http.MethodPatch, "/api/progress/1", body)
	req.SetPathValue("id", "1")
	rec := httptest.NewRecorder()

	s.handleProgress(rec, req)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}
