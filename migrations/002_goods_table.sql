-- +goose Up
CREATE TABLE IF NOT EXISTS goods(
	id SERIAL,
	project_id INT REFERENCES projects(id) ON DELETE CASCADE,
	name TEXT NOT NULL,
	description TEXT,
	priority INT NOT NULL,
	removed BOOLEAN NOT NULL DEFAULT false,
	created_at TIMESTAMP NOT NULL DEFAULT NOW(),
	PRIMARY KEY(id, project_id)
);

CREATE INDEX IF NOT EXISTS idx_goods_id ON goods(id);
CREATE INDEX IF NOT EXISTS idx_goods_project_id ON goods(project_id);
CREATE INDEX IF NOT EXISTS idx_goods_name ON goods(name);
CREATE INDEX IF NOT EXISTS idx_goods_priority ON goods(priority);

-- +goose Down
DROP INDEX IF EXISTS idx_goods_id;
DROP INDEX IF EXISTS idx_goods_project_id;
DROP INDEX IF EXISTS idx_goods_name;
DROP INDEX IF EXISTS idx_goods_priority;
DROP TABLE IF EXISTS goods; 
