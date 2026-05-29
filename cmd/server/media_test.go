package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"encoding/json"

	"github.com/stretchr/testify/assert"
)

func TestGetCatalog_OK(t *testing.T) {
	s := newTestServer(t)

	body := strings.NewReader(`{"new_progress": 120}`)
	req := httptest.NewRequest(http.MethodPatch, "/api/progress/1", body)
	req.SetPathValue("id", "1")
	rec := httptest.NewRecorder()

	s.handleGetCatalog(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)

	var resp []MediaResponse
	
	// json parses fine
	assert.NoError(t, json.NewDecoder(rec.Body).Decode(&resp))

	assert.Equal(t, resp[0].Title, "clase-01-crypto")
	assert.Equal(t, resp[0].Date,  "2026-05-18T00:00:00Z")
	assert.Equal(t, resp[0].ProgressSeconds, int64(450))
	
}

