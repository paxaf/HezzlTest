-- +goose Up
CREATE TABLE IF NOT EXISTS GOODS (
	id SERIAL,
	project_id INT REFERENCES PROJECTS(id) ON DELETE CASCADE,
	name TEXT NOT NULL,
	description TEXT,
	priority INT NOT NULL,
	removed BOOLEAN NOT NULL DEFAULT false,
	created_at TIMESTAMP NOT NULL DEFAULT NOW(),
	PRIMARY KEY(id, project_id)
);

CREATE INDEX IF NOT EXISTS idx_goods_id ON GOODS(id);
CREATE INDEX IF NOT EXISTS idx_goods_project_id ON GOODS(project_id);
CREATE INDEX IF NOT EXISTS idx_goods_name ON GOODS(name);

-- +goose Down
DROP INDEX IF EXISTS idx_goods_id;
DROP INDEX IF EXISTS idx_goods_project_id;
DROP INDEX IF EXISTS idx_goods_name;
DROP TABLE IF EXISTS GOODS; 
