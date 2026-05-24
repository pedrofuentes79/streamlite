-- Test data for local development
INSERT OR IGNORE INTO media (id, date, title, video_path, audio_path, progress_seconds, total_seconds)
VALUES
    (
        'clase-01-crypto',
        '2026-05-18',
        'Modern Cryptography',
        'media_store/crypto_class_1.mp4',
        'media_store/crypto_class_1.m4a',
        450.5,  -- 7.5 minutes in
        5400.0  -- 1h30m
    ),
    (
        'clase-02-complexity',
        '2026-05-20',
        'Complexity: P vs NP and coNP',
        'media_store/crypto_class_1.mp4',
        'media_store/crypto_class_1.m4a',
        0.0,    -- not started
        7200.0  -- 2 hours
    ),
    (
        'podcast-01-tech',
        '2026-05-22',
        'Designing Asynchronous Architectures',
        'media_store/podcast_1.mp4',
        'media_store/podcast_1.mp3',
        3600.0, -- halfway
        7200.0  -- 2 hours
    );
