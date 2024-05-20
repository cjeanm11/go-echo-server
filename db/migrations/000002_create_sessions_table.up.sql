CREATE TABLE sessions (
    id UUID PRIMARY KEY,
    user_id INTEGER REFERENCES users(id) NOT NULL,
    expires_at TIMESTAMP NOT NULL
);

CREATE INDEX idx_sessions_user_id ON sessions (user_id);


