
BEGIN;
CREATE TABLE IF NOT EXISTS users
(
    id           SERIAL	 PRIMARY KEY,
    email        TEXT    NOT NULL UNIQUE,
    pass_hash    BYTEA    NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_email ON users (email);

CREATE TABLE IF NOT EXISTS sessions
(
    id     TEXT NOT NULL UNIQUE,
    user_id   TEXT NOT NULL UNIQUE
);

COMMIT;
