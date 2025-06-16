-- name: CreateUser :one
INSERT INTO users (email, hashed_password)
VALUES (
    $1,
    $2
)
RETURNING id, created_at, updated_at, email;

-- name: DeleteUsers :exec
DELETE FROM users;

-- name: GetUserByEmal :one
SELECT * FROM users
WHERE email = $1;