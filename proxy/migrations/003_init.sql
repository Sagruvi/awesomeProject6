-- +goose Up
CREATE TABLE IF NOT EXISTS history_search_address (
    id SERIAL PRIMARY KEY,
    search_history_id INTEGER REFERENCES search_history(id),
    address_id INTEGER REFERENCES address(id),
    );
-- +goose Down
DROP TABLE IF EXISTS history_search_address