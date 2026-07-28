-- +goose Up
ALTER TABLE conversations
ALTER COLUMN created_at TYPE timestamptz
USING created_at::timestamptz;

-- +goose Down
ALTER TABLE conversations
ALTER COLUMN created_at TYPE text
USING created_at::text;
