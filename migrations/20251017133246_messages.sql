-- +goose Up
CREATE TABLE IF NOT EXISTS messages (
                                        id TEXT PRIMARY KEY,
                                        session_id TEXT NOT NULL REFERENCES sessions(id) ON DELETE CASCADE,
                                        role TEXT NOT NULL,
                                        content TEXT NOT NULL,
                                        created_at TIMESTAMPTZ NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_messages_session_created ON messages(session_id, created_at);

-- +goose Down
DROP INDEX IF EXISTS idx_messages_session_created;
DROP TABLE IF EXISTS messages;
