package chirps

import (
	"database/sql"
	"errors"
	"log"
	"net/http"

	"github.com/goinginblind/chirpy/internal/config"
	"github.com/goinginblind/chirpy/internal/handlers"
	"github.com/google/uuid"
)

func GetOneByID(cfg *config.APIConfig, w http.ResponseWriter, r *http.Request) {
	chirpIDString := r.PathValue("chirpID")
	chirpID, err := uuid.Parse(chirpIDString)
	if err != nil {
		log.Printf("Fail to convert path id into a uuid value: %v\n", err)
		handlers.RespondWithError(w, http.StatusBadRequest, "Error when processing chirpID")
		return
	}
	chirp, err := cfg.DB.GetChirp(r.Context(), chirpID)
	if errors.Is(err, sql.ErrNoRows) {
		log.Printf("Chirp not found: %v\n", chirpID)
		handlers.RespondWithError(w, http.StatusNotFound, "Chirp doesn't exist")
		return
	}
	if err != nil {
		log.Printf("Fail to get chirp by id from the db: %v\n", err)
		handlers.RespondWithError(w, http.StatusInternalServerError, "Failed to get chirp (internal error)")
		return
	}
	handlers.RespondWithJSON(w, http.StatusOK, convertDBChirpToResponse(chirp))
}
