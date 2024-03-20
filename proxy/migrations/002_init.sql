-- +goose Up
CREATE TABLE IF NOT EXISTS address (
    id SERIAL PRIMARY KEY,
    lat VARCHAR(20),
    lon VARCHAR(20)
    );
-- +goose Down
DROP TABLE IF EXISTS address