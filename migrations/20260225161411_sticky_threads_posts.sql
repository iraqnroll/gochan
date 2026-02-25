-- +goose Up

-- +goose StatementBegin
ALTER TABLE threads ADD COLUMN bumped_at TIMESTAMP NOT NULL DEFAULT NOW();
ALTER TABLE threads ADD COLUMN sticky BOOLEAN NOT NULL DEFAULT(FALSE);
ALTER TABLE posts ADD COLUMN sticky BOOLEAN NOT NULL DEFAULT(FALSE);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION bump_thread_on_post_insert()
RETURNS trigger AS $trg_bump_thread_on_post_insert$
BEGIN
    IF NEW.deleted IS FALSE THEN
        UPDATE threads 
        SET bumped_at = GREATEST(bumped_at, NEW.post_timestamp)
        WHERE id = NEW.thread_id;
    END IF;
    RETURN NEW;
END;
$trg_bump_thread_on_post_insert$ LANGUAGE plpgsql;

CREATE TRIGGER trg_posts_bump_thread
AFTER INSERT ON posts
FOR EACH ROW
EXECUTE FUNCTION bump_thread_on_post_insert();
-- +goose StatementEnd

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION recompute_thread_bump_on_post_deleted_change()
RETURNS trigger AS $trg_recompute_thread_bump_on_post_deleted_change$
DECLARE
    new_bump TIMESTAMP;
BEGIN
    IF NEW.deleted <> OLD.deleted THEN
        SELECT 
            COALESCE(MAX(p.post_timestamp), t.date_created) --if no posts are present, fallback to date_created.
        INTO new_bump
        FROM threads t
        LEFT JOIN posts p
            ON p.thread_id = t.id AND p.deleted IS FALSE
        WHERE t.id = NEW.thread_id
        GROUP BY t.date_created;
        
        UPDATE threads SET bumped_at = new_bump WHERE id = NEW.thread_id;
    END IF;
    RETURN NEW;
END;
$trg_recompute_thread_bump_on_post_deleted_change$ LANGUAGE plpgsql;

CREATE TRIGGER trg_posts_recompute_thread_bump
AFTER UPDATE OF deleted ON posts
FOR EACH ROW
EXECUTE FUNCTION recompute_thread_bump_on_post_deleted_change();
-- +goose StatementEnd

-- +goose StatementBegin
CREATE INDEX IF NOT EXISTS idx_posts_thread_id_post_timestamp_desc
    ON posts (thread_id, post_timestamp DESC);
-- +goose StatementEnd

-- +goose Down

-- +goose StatementBegin
DROP INDEX IF EXISTS idx_posts_thread_id_post_timestamp_desc;
-- +goose StatementEnd

-- +goose StatementBegin
DROP TRIGGER IF EXISTS trg_posts_recompute_thread_bump ON posts;
DROP TRIGGER IF EXISTS trg_posts_bump_thread ON posts;
-- +goose StatementEnd

-- +goose StatementBegin
DROP FUNCTION IF EXISTS recompute_thread_bump_on_post_deleted_change();
DROP FUNCTION IF EXISTS bump_thread_on_post_insert();
-- +goose StatementEnd

-- +goose StatementBegin
ALTER TABLE threads DROP COLUMN IF EXISTS bumped_at;
ALTER TABLE threads DROP COLUMN IF EXISTS sticky;
ALTER TABLE posts DROP COLUMN IF EXISTS sticky;
-- +goose StatementEnd
