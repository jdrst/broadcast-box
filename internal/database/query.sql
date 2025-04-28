-- name: GetUser :one
SELECT * FROM users
WHERE name = ? LIMIT 1;

-- name: ListUsers :many
SELECT * FROM users
ORDER BY id;

-- name: GetUsernameForStreamKey :one
SELECT name FROM users
WHERE streamKey = ?;

-- name: CreateUser :one
INSERT INTO users (
  id, name, password, streamKey) VALUES (
  ?, ?, ?, ?
)
RETURNING *;

-- name: UpdateUser :exec
UPDATE users
SET name = ?,
password = ?,
streamKey = ?
WHERE id = ?;

-- name: DeleteUser :exec
DELETE FROM users
WHERE name = ?;
