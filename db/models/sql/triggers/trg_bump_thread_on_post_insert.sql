--Trigger to bump threads on new post inserts
CREATE OR REPLACE FUNCTION bump_thread_on_post_insert()
RETURNS trigger AS $tg_bump_thread_on_post_insert$
BEGIN
    IF NEW.deleted IS FALSE THEN
        UPDATE threads
        SET bumped_at = GREATEST(bumped_at, NEW.post_timestamp)
        WHERE id = NEW.thread_id
    END IF
    RETURN NEW
END
$tg_bump_thread_on_post_insert$ LANGUAGE plpgsql;