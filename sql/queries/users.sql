-- name: CreateUser :one
INSERT INTO users (email, hashed_password)
VALUES (
    $1,
    $2
)
RETURNING id, created_at, updated_at, email;

-- name: DeleteUsers :exec
DELETE FROM users;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1;

-- name: ChangeUserLoginInfo :one
UPDATE users
SET email = $3, 
    hashed_password = $2,
    updated_at = now()
WHERE id = $1
RETURNING id, created_at, updated_at, email;