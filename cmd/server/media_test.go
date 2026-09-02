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

	// Newest publish date first. The seed rows are inserted in one statement, so
	// they share an updated_at — the previous ordering only put clase-01-crypto
	// first by accident of tie-breaking, and asserted that accident.
	assert.Equal(t, []string{"podcast-01-tech", "clase-02-complexity", "clase-01-crypto"},
		[]string{resp[0].Title, resp[1].Title, resp[2].Title})

	assert.Equal(t, "2026-05-22T00:00:00Z", resp[0].Date)
	assert.Equal(t, int64(3600), resp[0].ProgressSeconds)

	// The dates themselves must come back in descending order.
	for i := 1; i < len(resp); i++ {
		assert.Greater(t, resp[i-1].Date, resp[i].Date,
			"catalog must be sorted newest-first by publish date")
	}
}
