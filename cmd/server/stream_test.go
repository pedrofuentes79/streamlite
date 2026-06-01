// Some tests
// 1) Not found media
// 2) Media is found on the DB but not on disk
// 3) Happy path: what does this return? How do I inspect the payload to verify its correctness?
package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"os"
	"testing"
	"path/filepath"
	"strconv"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_MediaNotFoundInDB(t *testing.T) {  
	s := newTestServer(t)

	testIdStr := "120"
	body := strings.NewReader(`{"id": ` + testIdStr + `}`) // this id is not on the test DB
	req := httptest.NewRequest(http.MethodGet, "/api/stream/" + testIdStr, body)
	req.SetPathValue("id", testIdStr)
	rec := httptest.NewRecorder()

	s.handleStream(rec, req)
	assert.Equal(t, http.StatusNotFound, rec.Code)
	assert.Contains(t, rec.Body.String(), "Media not found in database")
}

func Test_MediaFoundInDB_ButNotOnDisk(t *testing.T) {  
	s := newTestServer(t)

	testIdStr := "1"
	body := strings.NewReader(`{"id": ` + testIdStr + `}`) // this id is not on the test DB
	req := httptest.NewRequest(http.MethodGet, "/api/stream/" + testIdStr, body)
	req.SetPathValue("id", testIdStr)
	rec := httptest.NewRecorder()

	s.handleStream(rec, req)
	assert.Equal(t, http.StatusNotFound, rec.Code)
	assert.Contains(t, rec.Body.String(), "Media file not found on disk")
}

func Test_FullFileRequestHappyPath(t *testing.T) {  
	s := newTestServer(t)

	// Add a fictitious row
	d := t.TempDir()
	videoPath := filepath.Join(d, "test.mp4")
	audioPath := filepath.Join(d, "test.mp3")

	// Write some bytes in the file
	os.WriteFile(videoPath, []byte("test"), 0644)
	
	result, err := s.db.Exec(
		`
		INSERT INTO media 
		(date, title, video_path, audio_path, total_seconds) 
		VALUES (?, ?, ?, ?, ?)
		`, "2026-06-01", "test", videoPath, audioPath, 120,
	)
	require.NoError(t, err) // query didn't fail
	id, err := result.LastInsertId()
	require.NoError(t, err) // query returns an id
	
	testIdStr := strconv.Itoa(int(id))
	req := httptest.NewRequest(http.MethodGet, "/api/stream/" + testIdStr, nil)
	req.SetPathValue("id", testIdStr)
	rec := httptest.NewRecorder()	
	
	s.handleStream(rec, req)
	
	assert.Equal(t, http.StatusOK, rec.Code)

	// Assert that the bytes are correct
	expected := []byte("test")
	got := rec.Body.Bytes()
	assert.Equal(t, expected, got)
}

func Test_AudioServesAudioColumn(t *testing.T) {
	s := newTestServer(t)

	// Distinct contents per column so we can prove handleAudio reads audio_path.
	d := t.TempDir()
	videoPath := filepath.Join(d, "test.mp4")
	audioPath := filepath.Join(d, "test.m4a")
	os.WriteFile(videoPath, []byte("video-bytes"), 0644)
	os.WriteFile(audioPath, []byte("audio-bytes"), 0644)

	result, err := s.db.Exec(
		`
		INSERT INTO media
		(date, title, video_path, audio_path, total_seconds)
		VALUES (?, ?, ?, ?, ?)
		`, "2026-06-01", "test", videoPath, audioPath, 120,
	)
	require.NoError(t, err)
	id, err := result.LastInsertId()
	require.NoError(t, err)

	testIdStr := strconv.Itoa(int(id))
	req := httptest.NewRequest(http.MethodGet, "/api/audio/"+testIdStr, nil)
	req.SetPathValue("id", testIdStr)
	rec := httptest.NewRecorder()

	s.handleAudio(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, []byte("audio-bytes"), rec.Body.Bytes())
}

func Test_RangeRequestHappyPath(t *testing.T) {
	s := newTestServer(t)

	// Add a fictitious row
	d := t.TempDir()
	videoPath := filepath.Join(d, "test.mp4")
	audioPath := filepath.Join(d, "test.mp3")

	// Write some bytes in the file
	os.WriteFile(videoPath, []byte("test-file"), 0644)
	
	result, err := s.db.Exec(
		`
		INSERT INTO media 
		(date, title, video_path, audio_path, total_seconds) 
		VALUES (?, ?, ?, ?, ?)
		`, "2026-06-01", "test", videoPath, audioPath, 120,
	)
	require.NoError(t, err) // query didn't fail
	id, err := result.LastInsertId()
	require.NoError(t, err) // query returns an id
	
	testIdStr := strconv.Itoa(int(id))
	req := httptest.NewRequest(http.MethodGet, "/api/stream/" + testIdStr, nil)
	req.SetPathValue("id", testIdStr)
	req.Header.Set("Range", "bytes=0-5")
	rec := httptest.NewRecorder()	
	
	s.handleStream(rec, req)
	
	assert.Equal(t, http.StatusPartialContent, rec.Code)

	// Assert that the bytes are correct
	expected := []byte("test-f")
	got := rec.Body.Bytes()
	assert.Equal(t, expected, got)
}