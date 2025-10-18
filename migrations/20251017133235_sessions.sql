-- +goose Up
CREATE TABLE IF NOT EXISTS sessions (
                                        id TEXT PRIMARY KEY,
                                        created_at TIMESTAMPTZ NOT NULL,
                                        updated_at TIMESTAMPTZ NOT NULL,
                                        status TEXT NOT NULL,
                                        tally INT NOT NULL,
                                        final_message TEXT
);
-- +goose Down
DROP TABLE IF EXISTS sessions;
