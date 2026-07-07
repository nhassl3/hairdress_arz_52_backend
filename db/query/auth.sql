-- name: CreateUser :one
INSERT INTO users(username, email, phone_number)
VALUES (sqlc.narg('username')::varchar, sqlc.arg('email')::text, sqlc.arg('phone_number')::varchar)
RETURNING *;

-- name: GetUserByUsername :one
SELECT * FROM users WHERE username=$1;

-- name: GetUserByPhone :one
SELECT * FROM users WHERE phone_number=$1;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email=$1;

-- name: ExistsByUsername :one
SELECT EXISTS(SELECT 1 FROM users WHERE username=$1) AS exists;

-- name: ExistsByPhoneNumber :one
SELECT EXISTS(SELECT 1 FROM users WHERE phone_number=$1) AS exists;

-- name: ExistsByEmail :one
SELECT EXISTS(SELECT 1 FROM users WHERE email=$1) AS exists;

-- name: VerifyUser :one
UPDATE users SET is_verified=true, last_login=now(), updated_at=now() WHERE
    (sqlc.narg('phone_number')::varchar is null or phone_number = sqlc.narg('phone_number')::varchar)
    AND (sqlc.narg('email')::text is null or email = sqlc.narg('email')::text) RETURNING *;

-- name: UpdateLastLogin :exec
UPDATE users SET last_login=now(), updated_at=now() WHERE username = $1;

-- name: ChangePhoneNumber :one
UPDATE users SET phone_number = $1, updated_at=now() WHERE username=$2 RETURNING *;
