package server

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/goinginblind/chirpy/internal/auth"
	"github.com/goinginblind/chirpy/internal/database"
)

func (s *Server) HandlerCreateChirp(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}
	userID, err := auth.ValidateJWT(token, s.Cfg.TokenSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	var params createChirpRequest
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		respondWithError(w, http.StatusBadRequest, "Something went wrong")
		log.Printf("Fail to decode request body: %s\n", err)
		return
	}

	if len(params.Body) > s.Cfg.MaxChirpLen {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long")
		return
	}
	params.Body = replaceProfanity(params.Body)

	chirp, err := s.Cfg.DB.CreateChirp(r.Context(), database.CreateChirpParams{
		Body:   params.Body,
		UserID: userID,
	})
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Fail to create chirp")
		log.Printf("Fail to create a chirp in db: %s\n", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, convertDBChirpToResponse(chirp))
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
