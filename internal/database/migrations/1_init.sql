CREATE TABLE IF NOT EXISTS media (
    id TEXT PRIMARY KEY,
    day INTEGER NOT NULL,
    title TEXT NOT NULL,
    video_path TEXT NOT NULL,
    audio_path TEXT NOT NULL,
    progress_seconds REAL DEFAULT 0.0,
    total_seconds REAL DEFAULT 7200.0,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);