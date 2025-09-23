--Users table.

CREATE TABLE IF NOT EXISTS users (
id SERIAL PRIMARY KEY,
username TEXT UNIQUE NOT NULL,
password_hash TEXT NOT NULL,
email TEXT UNIQUE NOT NULL,
date_created TIMESTAMP NOT NULL DEFAULT NOW(),
date_updated TIMESTAMP,
user_type INT NOT NULL
);