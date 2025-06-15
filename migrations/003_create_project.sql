-- +goose Up
INSERT INTO PROJECTS (name) VALUES ('Первая запись');

-- +goose Down
TRUNCATE PROJECTS;