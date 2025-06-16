CREATE TABLE IF NOT EXISTS logs.events
(
    timestamp   DateTime DEFAULT now(),
    action      String,
    entity      String,
    entity_id   Int32,
    project_id  Int32,
    name        String,
    description Nullable(String),
    priority    Nullable(Int32),
    removed     Nullable(UInt8),
    created_at  Nullable(DateTime)
)
ENGINE = MergeTree()
ORDER BY (entity, entity_id, timestamp);