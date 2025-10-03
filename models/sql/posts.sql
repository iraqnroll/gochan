CREATE TABLE posts (
    id SERIAL PRIMARY KEY,
    thread_id INT REFERENCES threads(id) ON DELETE CASCADE,
    identifier TEXT,
    content TEXT,
    post_timestamp TIMESTAMP NOT NULL DEFAULT NOW(),
    is_op BOOLEAN NOT NULL DEFAULT(FALSE)
)
