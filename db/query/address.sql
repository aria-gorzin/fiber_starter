-- name: CreateAddress :one
INSERT INTO addresses (
  client_id,
  title,
  city,
  street,
  phone,
  zip,
  lat,
  long
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8
) RETURNING *;

-- name: GetAddress :one
SELECT * FROM addresses
WHERE id = $1 LIMIT 1;

-- name: UpdateAddress :one
UPDATE addresses
SET
  title = $2,
  city = $3,
  street = $4,
  phone = $5,
  zip = $6,
  lat = $7,
  long = $8,
  updated_at = now()
WHERE
  id = $1
RETURNING *;

-- name: ListAddresses :many
SELECT * FROM addresses
WHERE client_id = $1;

-- name: DeleteAddress :exec
DELETE FROM addresses
WHERE id = $1;
