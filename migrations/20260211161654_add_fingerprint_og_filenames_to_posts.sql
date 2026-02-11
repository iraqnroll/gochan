-- +goose Up

-- +goose StatementBegin
ALTER TABLE posts ADD COLUMN IF NOT EXISTS og_media TEXT;
ALTER TABLE posts ADD COLUMN IF NOT EXISTS fingerprint TEXT;
-- +goose StatementEnd

-- +goose Down

-- +goose StatementBegin
ALTER TABLE posts DROP COLUMN IF EXISTS og_media CASCADE;
ALTER TABLE posts DROP COLUMN IF EXISTS fingerprint CASCADE;
-- +goose StatementEnd
