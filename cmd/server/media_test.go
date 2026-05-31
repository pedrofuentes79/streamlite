package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetCatalog_OK(t *testing.T) {
	s := newTestServer(t)

	req := httptest.NewRequest(http.MethodGet, "/api/media", nil)
	rec := httptest.NewRecorder()

	s.handleGetCatalog(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)

	var resp []MediaResponse

	// json parses fine
	assert.NoError(t, json.NewDecoder(rec.Body).Decode(&resp))

	assert.Equal(t, "clase-01-crypto", resp[0].Title)
	assert.Equal(t, "2026-05-18T00:00:00Z", resp[0].Date)
	assert.Equal(t, int64(450), resp[0].ProgressSeconds)
}
