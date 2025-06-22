-- name: CreateUser :one
INSERT INTO users (email, hashed_password)
VALUES (
    $1,
    $2
)
RETURNING id, email, is_chirpy_red, created_at, updated_at;

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
RETURNING id, email, is_chirpy_red, created_at, updated_at;

-- name: UpgradeUser :one
UPDATE users
SET is_chirpy_red = true,
    updated_at = now()
WHERE id = $1
RETURNING id, email, is_chirpy_red, created_at, updated_at;