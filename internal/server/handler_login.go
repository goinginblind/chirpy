package server

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/goinginblind/chirpy/internal/auth"
)

func (s *Server) HandlerLogin(w http.ResponseWriter, r *http.Request) {
	// Decode the whole json thats passed in the request
	var params loginDetails
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		log.Printf("Fail to decode request body: %v", err)
		respondWithError(w, http.StatusBadRequest, "Incorrect request")
		return
	}

	// Get user from the db by their email
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

	// Verify password
	if err = auth.CheckPasswordHash(user.HashedPassword, params.Password); err != nil {
		log.Printf("Failed attempt to login: %v", err)
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password")
		return
	}

	// Tokenizing
	// - get the expiration duration from the request (max dur = 3600 seconds = 1 hour)
	// - make a token
	// - respond with a json containing all the user info (no password) and the token
	var expiresIn time.Duration
	if params.ExpiresInSeconds <= 0 || params.ExpiresInSeconds > 3600 {
		expiresIn = time.Hour
	} else {
		expiresIn = time.Duration(params.ExpiresInSeconds) * time.Second
	}
	token, err := auth.MakeJWT(user.ID, s.Cfg.TokenSecret, expiresIn)
	if err != nil {
		log.Printf("Failed attempt to create user token: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Internal error")
		return
	}

	respondWithJSON(w, http.StatusOK, convertDBUserToResponse(user, token))
}
