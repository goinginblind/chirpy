package server

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/goinginblind/chirpy/internal/database"
)

func respondWithJSON(w http.ResponseWriter, code int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		log.Printf("Fail to encode json: %s\n", err)
	}
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	respondWithJSON(w, code, map[string]string{"error": msg})
}

func convertDBChirpToResponse(chirp database.Chirp) Chirp {
	return Chirp{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt.String(),
		UpdatedAt: chirp.UpdatedAt.String(),
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	}
}

func convertManyDBChirps(chirps []database.Chirp) []Chirp {
	out := make([]Chirp, len(chirps))
	for i, c := range chirps {
		out[i] = convertDBChirpToResponse(c)
	}
	return out
}

func convertDBUserRowToResponse(user database.CreateUserRow) User {
	return User{
		ID:        user.ID,
		Email:     user.Email,
		CreatedAt: user.CreatedAt.String(),
		UpdatedAt: user.UpdatedAt.String(),
	}
}

func convertDBUserToResponse(user database.User) User {
	return User{
		ID:        user.ID,
		Email:     user.Email,
		CreatedAt: user.CreatedAt.String(),
		UpdatedAt: user.UpdatedAt.String(),
	}
}
