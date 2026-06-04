-- Per-video chapter markers (YouTube chapters, captured by the scraper from
-- yt-dlp's info JSON). Stored as a JSON array of {start, title} — display-only
-- data we never query into, so a column beats a child table here.
--
-- Migrations are re-exec'd on every boot, so this ADD COLUMN runs again on an
-- already-migrated DB and fails with "duplicate column name"; runMigrations
-- tolerates exactly that error to keep the step idempotent.
ALTER TABLE media ADD COLUMN chapters TEXT NOT NULL DEFAULT '[]';
