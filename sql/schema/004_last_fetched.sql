-- +goose Up
ALTER TABLE feeds
ADD COLUMN last_fetched TIMESTAMP;

-- +goose Down
