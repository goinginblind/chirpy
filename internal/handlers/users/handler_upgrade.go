package users

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/goinginblind/chirpy/internal/config"
	"github.com/goinginblind/chirpy/internal/handlers"
	"github.com/google/uuid"
)

func UpgradeToChirpyRed(cfg *config.APIConfig, w http.ResponseWriter, r *http.Request) {
	var req upgradeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Failed to decode request: %v", err)
		handlers.RespondWithError(w, http.StatusBadRequest, "")
		return
	}
	defer r.Body.Close()

	if req.Event != "user.upgraded" {
		log.Printf("Ignoring irrelevant event: %s", req.Event)
		w.WriteHeader(http.StatusNoContent)
		return
	}

	userID, err := uuid.Parse(req.Data.UserID)
	if err != nil {
		log.Printf("Invalid UUID in request: %v", err)
		handlers.RespondWithError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}
	_, err = cfg.DB.UpgradeUser(r.Context(), userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Printf("Failed to find user in the DB: %v", err)
			handlers.RespondWithError(w, http.StatusNotFound, "")
			return
		}
		log.Printf("Failed to upgrade user in the DB: %v", err)
		handlers.RespondWithError(w, http.StatusInternalServerError, "Internal error")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
