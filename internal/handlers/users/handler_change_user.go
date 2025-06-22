package users

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/goinginblind/chirpy/internal/auth"
	"github.com/goinginblind/chirpy/internal/config"
	"github.com/goinginblind/chirpy/internal/database"
	"github.com/goinginblind/chirpy/internal/handlers"
)

func ChangeLoginInfo(cfg *config.APIConfig, w http.ResponseWriter, r *http.Request) {
	// Get access token and check it
	authToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		log.Printf("Failed to get token from the header: %v\n", err)
		handlers.RespondWithError(w, http.StatusUnauthorized, "No token in the authorization header")
		return
	}
	userID, err := auth.ValidateJWT(authToken, cfg.TokenSecret)
	if err != nil {
		log.Printf("Unauthorized login attempt: %v", err)
		handlers.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Decode request body and validate that it has (somewhat) valid password and email
	defer r.Body.Close()
	var logInf loginRequest
	err = json.NewDecoder(r.Body).Decode(&logInf)
	if err != nil {
		log.Printf("Fail to decode request body: %v\n", err)
		handlers.RespondWithError(w, http.StatusBadRequest, "Incorrect request")
		return
	}

	if logInf.Email == "" || logInf.Password == "" {
		handlers.RespondWithError(w, http.StatusBadRequest, "Email and password must be provided")
		return
	}

	// Hash password
	hashedPassword, err := auth.HashPassword(logInf.Password)
	if err != nil {
		log.Printf("Fail to hash password: %v\n", err)
		handlers.RespondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	// Update the DB
	userRow, err := cfg.DB.ChangeUserLoginInfo(r.Context(), database.ChangeUserLoginInfoParams{
		ID:             userID,
		Email:          logInf.Email,
		HashedPassword: hashedPassword,
	})
	if err != nil {
		log.Printf("Fail to change user login info: %v\n", err)
		handlers.RespondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	// Send the response
	handlers.RespondWithJSON(w, http.StatusOK, convertLogRowToCreateParams(userRow))
}
