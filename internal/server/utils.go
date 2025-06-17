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

func convertDBChirpToResponse(chirp database.Chirp) chirpParams {
	return chirpParams{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt.String(),
		UpdatedAt: chirp.UpdatedAt.String(),
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	}
}

func convertManyDBChirps(chirps []database.Chirp) []chirpParams {
	out := make([]chirpParams, len(chirps))
	for i, c := range chirps {
		out[i] = convertDBChirpToResponse(c)
	}
	return out
}

func dbUserRowToCreateParams(user database.CreateUserRow) createUserParams {
	return createUserParams{
		ID:        user.ID,
		Email:     user.Email,
		CreatedAt: user.CreatedAt.String(),
		UpdatedAt: user.UpdatedAt.String(),
	}
}

func dbUserToLoginParams(user database.User, token, refreshToken string) loginUserParams {
	return loginUserParams{
		ID:           user.ID,
		Email:        user.Email,
		CreatedAt:    user.CreatedAt.String(),
		UpdatedAt:    user.UpdatedAt.String(),
		Token:        token,
		RefreshToken: refreshToken,
	}
}
