-- name: Authorize :one
INSERT INTO users (username, full_name, phone_number) VALUES (sqlc.arg('username'), sqlc.narg('full_name')::varchar, sqlc.arg('phone_number'))
ON CONFLICT (username) DO UPDATE SET last_login=now() WHERE username=excluded.username AND phone_number=excluded.phone_number
RETURNING *;

-- name: Verify :exec
UPDATE users SET is_verified=true, updated_at=now() where username=$1;

-- name: ChangePhoneNumber :one
UPDATE users SET phone_number=$1, updated_at=now() where username=$2 RETURNING *;