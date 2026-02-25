--Recompute thread last bump after a soft-delete of a post.
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