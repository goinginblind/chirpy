package chirps

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/goinginblind/chirpy/internal/auth"
	"github.com/goinginblind/chirpy/internal/config"
	"github.com/goinginblind/chirpy/internal/database"
	"github.com/goinginblind/chirpy/internal/handlers"
)

func Create(cfg *config.APIConfig, w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		handlers.RespondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}
	userID, err := auth.ValidateJWT(token, cfg.TokenSecret)
	if err != nil {
		handlers.RespondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	var params createChirpRequest
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		handlers.RespondWithError(w, http.StatusBadRequest, "Something went wrong")
		log.Printf("Fail to decode request body: %s\n", err)
		return
	}

	if len(params.Body) > cfg.MaxChirpLen {
		handlers.RespondWithError(w, http.StatusBadRequest, "Chirp is too long")
		return
	}
	params.Body = replaceProfanity(params.Body)

	chirp, err := cfg.DB.CreateChirp(r.Context(), database.CreateChirpParams{
		Body:   params.Body,
		UserID: userID,
	})
	if err != nil {
		handlers.RespondWithError(w, http.StatusBadRequest, "Fail to create chirp")
		log.Printf("Fail to create a chirp in db: %s\n", err)
		return
	}

	handlers.RespondWithJSON(w, http.StatusCreated, convertDBChirpToResponse(chirp))
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
