CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    public_address VARCHAR(42) UNIQUE NOT NULL,
    nonce VARCHAR(100) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    last_login TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE sessions (
    id uuid PRIMARY KEY,
    public_address text NOT NULL,
    refresh_token text NOT NULL,
    user_agent text NOT NULL,
    client_ip text NOT NULL,
    is_blocked boolean NOT NULL DEFAULT false,
    expires_at timestamptz NOT NULL,
    created_at timestamptz DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_users_address ON users(public_address);
CREATE INDEX idx_users_nonce ON users(nonce);
CREATE INDEX idx_sessions_public_address ON sessions(public_address);
CREATE INDEX idx_sessions_token ON sessions(refresh_token);