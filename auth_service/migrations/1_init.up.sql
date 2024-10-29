
BEGIN;
CREATE TABLE IF NOT EXISTS users
(
    id           SERIAL	 PRIMARY KEY,
    email        TEXT    NOT NULL UNIQUE,
    password    BYTEA    NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_email ON users (email);

CREATE TABLE IF NOT EXISTS sessions
(
    id     TEXT NOT NULL UNIQUE,
    user_id   BIGINT NOT NULL
);

COMMIT;
