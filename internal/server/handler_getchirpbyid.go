package server

import (
	"database/sql"
	"errors"
	"log"
	"net/http"

	"github.com/google/uuid"
)

func (s *Server) HandlerGetChirpByID(w http.ResponseWriter, r *http.Request) {
	chirpIDString := r.PathValue("chirpID")
	chirpID, err := uuid.Parse(chirpIDString)
	if err != nil {
		log.Printf("Fail to convert path id into a uuid value: %v", err)
		respondWithError(w, http.StatusBadRequest, "Error when processing chirpID")
		return
	}
	chirp, err := s.Cfg.DB.GetChirp(r.Context(), chirpID)
	if errors.Is(err, sql.ErrNoRows) {
		log.Printf("Chirp not found: %v", chirpID)
		respondWithError(w, http.StatusNotFound, "Chirp doesn't exist")
		return
	}
	if err != nil {
		log.Printf("Fail to get chirp by id from the db: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Failed to get chirp (internal error)")
		return
	}
	respondWithJSON(w, http.StatusOK, convertDBChirpToResponse(chirp))
}
