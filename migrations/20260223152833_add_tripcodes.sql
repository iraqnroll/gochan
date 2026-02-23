-- +goose Up
-- +goose StatementBegin
ALTER TABLE posts ADD COLUMN IF NOT EXISTS tripcode VARCHAR(64);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE posts DROP COLUMN IF EXISTS tripcode CASCADE;
-- +goose StatementEnd
