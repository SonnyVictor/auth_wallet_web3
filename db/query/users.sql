-- name: CreateUser :one
INSERT INTO users (public_address, nonce)
VALUES ($1, $2)
RETURNING *;

-- name: GetUserBypublic_address :one
SELECT * FROM users
WHERE public_address = $1;

-- name: UpdateUserNonce :one
UPDATE users
SET nonce = $2
WHERE public_address = $1
RETURNING *;

-- name: UpdateLastLogin :one
UPDATE users 
SET last_login = CURRENT_TIMESTAMP 
WHERE id = $1 
RETURNING *;

-- name: UserExists :one
SELECT EXISTS(SELECT 1 FROM users WHERE public_address = $1);

-- name: CreateSessions :one
INSERT INTO sessions(
    id,
    public_address,
    refresh_token,
    user_agent,
    client_ip,
    is_blocked,
    expires_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7
) RETURNING *;

-- name: GetSession :one
SELECT * FROM sessions
WHERE id = $1 LIMIT 1;
