CREATE TABLE IF NOT EXISTS logs.events (
    timestamp    DateTime DEFAULT now(),
    action       Enum8('create' = 1, 'update' = 2, 'delete' = 3),
    entity       String,
    entity_id    Int32,
    project_id   Int32,
    name         String,
    description  String,
    priority     Int32,
    removed      UInt8,
    created_at   DateTime
) ENGINE = MergeTree()
ORDER BY (entity, entity_id, timestamp)
PARTITION BY toYYYYMM(timestamp);