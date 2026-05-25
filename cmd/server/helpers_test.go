package main

import (
	"testing"

	"streamlite/internal/database"
)

func newTestServer(t *testing.T) *server {
	t.Helper()
	db := database.InitDB(":memory:", false)
	return &server{db: db}
}
