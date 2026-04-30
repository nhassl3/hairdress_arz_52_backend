-- name: CreateUser :one
INSERT INTO users(username, full_name, phone_number)
VALUES (sqlc.arg('username'), sqlc.narg('full_name')::varchar, sqlc.arg('phone_number'))
RETURNING *;

-- name: GetUserByUsername :one
SELECT * FROM users WHERE username=$1;

-- name: GetUserByPhone :one
SELECT * FROM users WHERE phone_number=$1;

-- name: ExistsByUsername :one
SELECT EXISTS(SELECT 1 FROM users WHERE username=$1) AS exists;

-- name: ExistsByPhoneNumber :one
SELECT EXISTS(SELECT 1 FROM users WHERE phone_number=$1) AS exists;

-- name: VerifyUser :exec
UPDATE users SET is_verified=true, last_login=now(), updated_at=now() WHERE username = $1;

-- name: UpdateLastLogin :exec
UPDATE users SET last_login=now(), updated_at=now() WHERE username = $1;

-- name: ChangePhoneNumber :one
UPDATE users SET phone_number = $1, updated_at=now() WHERE username=$2 RETURNING *;
