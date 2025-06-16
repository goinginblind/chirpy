package server

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/goinginblind/chirpy/internal/database"
)

func (s *Server) HandlerCreateChirp(w http.ResponseWriter, r *http.Request) {
	var params createChirpParams
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		respondWithJSON(w, http.StatusBadRequest, map[string]string{"error": "Something went wrong"})
		log.Printf("Fail to decode request body: %s\n", err)
		return
	}
	if len(params.Body) > s.Cfg.MaxChirpLen {
		respondWithJSON(w, http.StatusBadRequest, map[string]string{"error": "Chirp is too long"})
		return
	}
	params.Body = replaceProfanity(params.Body)

	chirp, err := s.Cfg.DB.CreateChirp(r.Context(), database.CreateChirpParams{
		Body:   params.Body,
		UserID: params.UserID,
	})
	if err != nil {
		respondWithJSON(w, http.StatusBadRequest, map[string]string{"error": "Fail to create chirp"})
		log.Printf("Fail to create a chirp in db: %s\n", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, convertDBChripToResponse(chirp))
}

func replaceProfanity(chirp string) string {
	profane := map[string]bool{
		"kerfuffle": true,
		"sharbert":  true,
		"fornax":    true,
	}
	sep := strings.Split(chirp, " ")
	for i, word := range sep {
		if ok := profane[strings.ToLower(word)]; ok {
			sep[i] = "****"
		}
	}
	return strings.Join(sep, " ")
}
