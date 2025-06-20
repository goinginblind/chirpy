package users

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/goinginblind/chirpy/internal/auth"
	"github.com/goinginblind/chirpy/internal/config"
	"github.com/goinginblind/chirpy/internal/database"
	"github.com/goinginblind/chirpy/internal/handlers"
)

func Login(cfg *config.APIConfig, w http.ResponseWriter, r *http.Request) {
	// Decode the request
	var params loginRequest
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		log.Printf("Fail to decode request body: %v", err)
		handlers.RespondWithError(w, http.StatusBadRequest, "Incorrect request")
		return
	}

	// Get user from the db by their email
	user, err := cfg.DB.GetUserByEmail(r.Context(), params.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			handlers.RespondWithError(w, http.StatusUnauthorized, "Incorrect email or password")
			return
		}
		log.Printf("Failed to get user by email: %v", err)
		handlers.RespondWithError(w, http.StatusInternalServerError, "Internal error")
		return
	}

	// Verify password
	if err = auth.CheckPasswordHash(user.HashedPassword, params.Password); err != nil {
		log.Printf("Failed attempt to login: %v", err)
		handlers.RespondWithError(w, http.StatusUnauthorized, "Incorrect email or password")
		return
	}

	// Make an access token and send it with a response body
	token, err := auth.MakeJWT(user.ID, cfg.TokenSecret)
	if err != nil {
		log.Printf("Failed attempt to create access token: %v", err)
		handlers.RespondWithError(w, http.StatusInternalServerError, "Internal error")
		return
	}

	// Make a refresh token and store it in the DB
	refreshToken, err := auth.MakeRefreshToken()
	if err != nil {
		log.Printf("Failed attempt to create refresh token: %v", err)
		handlers.RespondWithError(w, http.StatusInternalServerError, "Internal error")
		return
	}
	_, err = cfg.DB.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
		Token:     refreshToken,
		UserID:    user.ID,
		ExpiresAt: time.Now().Add(time.Hour * 24 * 60),
	})
	if err != nil {
		log.Printf("Failed attempt of parsing refresh token to DB: %v", err)
		handlers.RespondWithError(w, http.StatusInternalServerError, "Internal error")
		return
	}

	handlers.RespondWithJSON(w, http.StatusOK, convertUserToLoginParams(user, token, refreshToken))
}
