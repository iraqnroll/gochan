
CREATE TABLE boards (
    id SERIAL PRIMARY KEY,
    uri TEXT NOT NULL,
    name TEXT NOT NULL,
    description TEXT,
    date_created TIMESTAMP NOT NULL DEFAULT NOW(),
    date_updated TIMESTAMP,
    ownerId INT UNIQUE REFERENCES users(id) ON DELETE CASCADE
);