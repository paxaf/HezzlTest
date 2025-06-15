-- +goose Up
CREATE TABLE IF NOT EXISTS projects(
	id SERIAL PRIMARY KEY,
	name TEXT NOT NULL,
	created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- +goose Down
DROP TABLE IF EXISTS projects;