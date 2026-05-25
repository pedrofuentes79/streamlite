-- Test data for local development
INSERT OR IGNORE INTO media (date, title, video_path, audio_path, progress_seconds, total_seconds)
VALUES
    (
        '2026-05-18',
        'clase-01-crypto',
        'media_store/crypto_class_1.mp4',
        'media_store/crypto_class_1.m4a',
        450,    -- 7.5 minutes in
        5400.0  -- 1h30m
    ),
    (
        '2026-05-20',
        'clase-02-complexity',
        'media_store/crypto_class_1.mp4',
        'media_store/crypto_class_1.m4a',
        0,      -- not started
        7200.0  -- 2 hours
    ),
    (
        '2026-05-22',
        'podcast-01-tech',
        'media_store/podcast_1.mp4',
        'media_store/podcast_1.mp3',
        3600,   -- halfway
        7200.0  -- 2 hours
    );
