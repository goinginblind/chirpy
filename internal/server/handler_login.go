package server

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/goinginblind/chirpy/internal/auth"
	"github.com/goinginblind/chirpy/internal/database"
)

func (s *Server) HandlerLogin(w http.ResponseWriter, r *http.Request) {
	// Decode the request
	var params loginRequest
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

	// Make an access token and send it with a response body
	token, err := auth.MakeJWT(user.ID, s.Cfg.TokenSecret)
	if err != nil {
		log.Printf("Failed attempt to create access token: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Internal error")
		return
	}

	// Make a refresh token and store it in the DB
	refreshToken, err := auth.MakeRefreshToken()
	if err != nil {
		log.Printf("Failed attempt to create refresh token: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Internal error")
		return
	}
	_, err = s.Cfg.DB.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
		Token:     refreshToken,
		UserID:    user.ID,
		ExpiresAt: time.Now().Add(time.Hour * 24 * 60),
	})
	if err != nil {
		log.Printf("Failed attempt of parsing refresh token to DB: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Internal error")
		return
	}

	respondWithJSON(w, http.StatusOK, dbUserToLoginParams(user, token, refreshToken))
}
