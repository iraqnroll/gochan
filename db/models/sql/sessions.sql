--Table to track active user sessions
--TODO: Add timestamp for a history of previous active sessions.
--TODO: Add session expiration logic, we dont want to keep sessions active indefinitely...

CREATE TABLE sessions (
id SERIAL PRIMARY KEY,
user_id INT UNIQUE REFERENCES users(id) ON DELETE CASCADE,
token_hash TEXT UNIQUE NOT NULL
);