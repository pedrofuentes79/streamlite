-- Insertamos datos de prueba para testing local en streamlite
INSERT OR IGNORE INTO media (id, day, title, video_path, audio_path, progress_seconds, total_seconds)
VALUES
    (
        'clase-01-crypto',
        20591, -- 2026-05-18 (Lunes)
        'Introducción a la Criptografía Moderna',
        'media_store/crypto_clase_1.mp4',
        'media_store/crypto_clase_1.mp3',
        450.5, -- Ya vio 7.5 minutos
        5400.0 -- 1 hora y media en total
    ),
    (
        'clase-02-complexity',
        20593, -- 2026-05-20 (Miércoles)
        'Complejidad: Clases P vs NP y coNP',
        'media_store/complexity_clase_2.mp4',
        'media_store/complexity_clase_2.mp3',
        0.0, -- No empezado
        7200.0 -- 2 horas en total
    ),
    (
        'podcast-01-tech',
        20595, -- 2026-05-22 (Viernes)
        'Diseño de Arquitecturas Asíncronas',
        'media_store/podcast_1.mp4',
        'media_store/podcast_1.mp3',
        3600.0, -- Mitad de camino
        7200.0  -- 2 horas en total
    );