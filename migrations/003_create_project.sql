-- +goose Up
INSERT INTO projects (name) VALUES ('Первая запись');

-- +goose Down
TRUNCATE projects;