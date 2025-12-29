-- +goose Up

-- +goose statementbegin
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username TEXT UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    email TEXT UNIQUE NOT NULL,
    date_created TIMESTAMP NOT NULL DEFAULT NOW(),
    date_updated TIMESTAMP NULL DEFAULT NOW(),
    user_type INT NOT NULL
);
-- +goose statementend

-- +goose statementbegin
CREATE TABLE IF NOT EXISTS boards (
    id SERIAL PRIMARY KEY,
    uri TEXT NOT NULL,
    name TEXT NOT NULL,
    description TEXT,
    date_created TIMESTAMP NOT NULL DEFAULT NOW(),
    date_updated TIMESTAMP,
    ownerId INT REFERENCES users(id) ON DELETE CASCADE
);
-- +goose statementend

-- +goose statementbegin
CREATE TABLE IF NOT EXISTS threads (
    id SERIAL PRIMARY KEY,
    board_id INT REFERENCES boards(id) ON DELETE CASCADE,
    topic TEXT NOT NULL,
    date_created TIMESTAMP NOT NULL DEFAULT NOW(),
    locked BOOLEAN NOT NULL DEFAULT(FALSE)
);
-- +goose statementend

-- +goose statementbegin
CREATE TABLE IF NOT EXISTS posts (
    id SERIAL PRIMARY KEY,
    thread_id INT REFERENCES threads(id) ON DELETE CASCADE,
    identifier TEXT,
    content TEXT,
    post_timestamp TIMESTAMP NOT NULL DEFAULT NOW(),
    is_op BOOLEAN NOT NULL DEFAULT(FALSE),
    has_media TEXT
);
-- +goose statementend

-- +goose statementbegin
CREATE TABLE IF NOT EXISTS "sessions" (
    id SERIAL PRIMARY KEY,
    user_id INT UNIQUE REFERENCES users(id) ON DELETE CASCADE,
    token_hash TEXT UNIQUE NOT NULL
);
-- +goose statementend

-- +goose statementbegin
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
WHERE t.locked IS FALSE
ORDER BY t.id, p.post_timestamp DESC;
-- +goose statementend

-- +goose statementbegin
INSERT INTO users (username, password_hash, email, user_type)
VALUES ('admin', '$2a$10$4qz2nL9BK6LW2yz7wwRXHeMXMFWQCc4KHo0pP9UqylOSoP6vuALz.', 'test@email.com', 1);
-- +goose statementend

-- +goose Down
DROP VIEW IF EXISTS recent_posts;
DROP TABLE IF EXISTS "sessions";
DROP TABLE IF EXISTS posts;
DROP TABLE IF EXISTS threads;
DROP TABLE IF EXISTS boards;
DROP TABLE IF EXISTS users;