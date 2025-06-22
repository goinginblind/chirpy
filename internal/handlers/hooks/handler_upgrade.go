package hooks

import (
	"crypto/subtle"
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/goinginblind/chirpy/internal/auth"
	"github.com/goinginblind/chirpy/internal/config"
	"github.com/google/uuid"
)

func UpgradeToChirpyRed(cfg *config.APIConfig, w http.ResponseWriter, r *http.Request) {
	apiKey, err := auth.GetAPIKey(r.Header)
	if err != nil {
		log.Printf("Failed to get an api key from the request: %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if subtle.ConstantTimeCompare([]byte(apiKey), []byte(cfg.PolkaKey)) != 1 {
		log.Printf("Incorrect polka api key: %v\n", apiKey)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	defer r.Body.Close()
	var req upgradeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Failed to decode request: %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if req.Event != "user.upgraded" {
		log.Printf("Ignoring irrelevant event: %s", req.Event)
		w.WriteHeader(http.StatusNoContent)
		return
	}

	userID, err := uuid.Parse(req.Data.UserID)
	if err != nil {
		log.Printf("Invalid UUID in request: %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	_, err = cfg.DB.UpgradeUser(r.Context(), userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Printf("Failed to find user in the DB: %v\n", err)
			w.WriteHeader(http.StatusNotFound)
			return
		}
		log.Printf("Failed to upgrade user in the DB: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	log.Printf("Successfully upgraded user %s to Chirpy Red", userID)
	w.WriteHeader(http.StatusNoContent)
}
