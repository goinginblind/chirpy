package server

type parameters struct {
	Body string `json:"body"`
}

type createUserParams struct {
	Email string `json:"email"`
}

type User struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}
