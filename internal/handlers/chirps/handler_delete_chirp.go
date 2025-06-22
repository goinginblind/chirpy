package chirps

import (
	"log"
	"net/http"

	"github.com/goinginblind/chirpy/internal/auth"
	"github.com/goinginblind/chirpy/internal/config"
	"github.com/goinginblind/chirpy/internal/handlers"
	"github.com/google/uuid"
)

func DeleteOneByID(cfg *config.APIConfig, w http.ResponseWriter, r *http.Request) {
	chirpIDString := r.PathValue("chirpID")
	chirpID, err := uuid.Parse(chirpIDString)
	if err != nil {
		log.Printf("Fail to convert path id into a uuid value: %v\n", err)
		handlers.RespondWithError(w, http.StatusBadRequest, "Error when processing chirpID")
		return
	}
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		handlers.RespondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}
	userID, err := auth.ValidateJWT(token, cfg.TokenSecret)
	if err != nil {
		log.Printf("Unable to validate the token: %v", err)
		handlers.RespondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}
	chirp, err := cfg.DB.GetChirp(r.Context(), chirpID)
	if err != nil {
		log.Printf("Unable to get chirp with an ID '%v': %v\n", chirpID, err)
		handlers.RespondWithError(w, http.StatusNotFound, err.Error())
		return
	}
	if chirp.UserID != userID {
		log.Printf("Forbidden attempt to delete a chirp with an id '%v' by user '%v'\n", chirp.ID, userID)
		handlers.RespondWithError(w, http.StatusForbidden, "Forbidden")
		return
	}
	if err = cfg.DB.DeleteChirp(r.Context(), chirp.ID); err != nil {
		log.Printf("Failed to delete chirp: %v\n", err)
		handlers.RespondWithError(w, http.StatusInternalServerError, "Something went wrong...")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
