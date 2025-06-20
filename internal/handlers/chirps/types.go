package chirps

import "github.com/google/uuid"

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
