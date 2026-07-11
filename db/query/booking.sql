-- name: CreateBooking :one
INSERT INTO bookings(username, hairdresser_id, service_id, salon_id, starts_at, ends_at, description, status)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;

-- name: GetBooking :many
SELECT * FROM bookings WHERE
                           (sqlc.narg('username')::varchar is null or username=sqlc.narg('username')::varchar)
                            AND (sqlc.narg('id')::bigint is null or id=sqlc.narg('id')::bigint)
                            AND (sqlc.narg('hairdresser_id')::uuid is null or hairdresser_id=sqlc.narg('hairdresser_id')::uuid)
                            AND (sqlc.narg('service_id')::integer is null or service_id=sqlc.narg('service_id')::integer)
                            AND (sqlc.narg('salon_id')::integer is null or salon_id=sqlc.narg('salon_id')::integer);