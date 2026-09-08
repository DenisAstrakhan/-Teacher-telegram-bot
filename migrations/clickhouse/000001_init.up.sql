CREATE TABLE IF NOT EXISTS app_logs (
    timestamp DateTime DEFAULT now(),
    level String,
    service String,
    user_id String,
    message String
) ENGINE = MergeTree()
PARTITION BY toYYYYMM(timestamp)
ORDER BY (timestamp, level);