-- +goose Up

CREATE TABLE boards (
	id UUID PRIMARY KEY,
	title TEXT NOT NULL DEFAULT 'Untitled',
	data JSONB NOT NULL DEFAULT '{}'::jsonb,
	created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- +goose Down
DROP TABLE IF EXISTS boards;