CREATE TABLE IF NOT EXISTS comments (
    id          SERIAL PRIMARY KEY,
    name        VARCHAR(255) NOT NULL,
    comment     TEXT NOT NULL,
    created_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);