package server

import "github.com/google/uuid"

// Used to decode *into* from the request body when user tries to login or register
type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Sent back as json upon user registration
type createUserParams struct {
	ID        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	CreatedAt string    `json:"created_at"`
	UpdatedAt string    `json:"updated_at"`
}

// Sent back as json upon user login
type loginUserParams struct {
	ID           uuid.UUID `json:"id"`
	Email        string    `json:"email"`
	CreatedAt    string    `json:"created_at"`
	UpdatedAt    string    `json:"updated_at"`
	Token        string    `json:"token"`
	RefreshToken string    `json:"refresh_token"`
}

// Used to decode *into* from the request body when user tries to post a chirp
type createChirpRequest struct {
	Body string `json:"body"`
}

// Sent back as json upon posting a chirp
type chirpParams struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt string    `json:"created_at"`
	UpdatedAt string    `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}
