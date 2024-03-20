-- +goose Up
CREATE TABLE IF NOT EXISTS search_history (
    id SERIAL PRIMARY KEY,
    query_text TEXT NOT NULL,
);
-- +goose Down
DROP TABLE IF EXISTS search_history