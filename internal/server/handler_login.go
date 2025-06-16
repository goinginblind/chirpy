package server

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/goinginblind/chirpy/internal/auth"
)

func (s *Server) HandlerLogin(w http.ResponseWriter, r *http.Request) {
	var params loginDetails
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		log.Printf("Fail to decode request body: %v", err)
		respondWithError(w, http.StatusBadRequest, "Incorrect request")
		return
	}

	user, err := s.Cfg.DB.GetUserByEmail(r.Context(), params.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondWithError(w, http.StatusUnauthorized, "Incorrect email or password")
			return
		}
		log.Printf("Failed to get user by email: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Internal error")
		return
	}

	if err = auth.CheckPasswordHash(user.HashedPassword, params.Password); err != nil {
		log.Printf("Failed attempt to login: %v", err)
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password")
		return
	}

	respondWithJSON(w, http.StatusOK, convertDBUserToResponse(user))
}
