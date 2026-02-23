CREATE OR REPLACE VIEW recent_posts AS
SELECT DISTINCT ON (t.id)
    b.id AS board_id,
    b.uri AS board_uri,
    b.name AS board_name,
    t.id AS thread_id,
    t.topic AS thread_topic,
    p.id AS post_id,
    p.identifier AS post_ident,
    p.content AS post_content,
    p.post_timestamp AS post_timestamp
FROM posts AS p
INNER JOIN threads AS t ON t.id = p.thread_id
INNER JOIN boards AS b ON b.id = t.board_id
WHERE t.locked IS FALSE AND p.deleted IS FALSE
ORDER BY t.id, p.post_timestamp DESC
