package main

import (
	"testing"

	"streamlite/internal/database"
)

func newTestServer(t *testing.T) *server {
	t.Helper()
	db := database.InitDB(":memory:", false, true)
	return &server{db: db}
}
