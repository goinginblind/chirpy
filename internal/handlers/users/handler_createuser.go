package users

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/goinginblind/chirpy/internal/auth"
	"github.com/goinginblind/chirpy/internal/config"
	"github.com/goinginblind/chirpy/internal/database"
	"github.com/goinginblind/chirpy/internal/handlers"
	"github.com/lib/pq"
)

func Create(cfg *config.APIConfig, w http.ResponseWriter, r *http.Request) {
	// maybe need to implement password validation
	// - should it be long enough
	// - maybe it should contain some special chars, etc
	var params loginRequest
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil || params.Email == "" || params.Password == "" {
		log.Printf("Fail to decode request body: %v\n", err)
		handlers.RespondWithError(w, http.StatusBadRequest, "Incorrect request")
		return
	}

	hashPass, err := auth.HashPassword(params.Password)
	if err != nil {
		log.Printf("Fail to hash password: %v\n", err)
		handlers.RespondWithError(w, http.StatusInternalServerError, "Fail to register password")
		return
	}

	user, err := cfg.DB.CreateUser(r.Context(), database.CreateUserParams{
		Email:          params.Email,
		HashedPassword: hashPass,
	})

	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code == "23505" { // 23505 is a unique violation
				log.Printf("Attempted to create duplicate user: %v\n", params.Email)
				handlers.RespondWithError(w, http.StatusConflict, "Email already exists")
				return
			}
		}
		log.Printf("Fail to create user in the db: %v\n", err)
		handlers.RespondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	handlers.RespondWithJSON(w, http.StatusCreated, convertCreRowToCreateParams(user))
}
