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
		w.WriteHeader(500) //does this actually do anything at this point tho..
		return
	}
}

func convertDBChripToResponse(chirp database.Chirp) Chirp {
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
		out[i] = convertDBChripToResponse(c)
	}
	return out
}
