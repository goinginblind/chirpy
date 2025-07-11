-- +goose Up
CREATE TABLE users(
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP NOT NULL DEFAULT now(),
    email TEXT UNIQUE NOT NULL
);

-- +goose Down
DROP TABLE IF EXISTS users;