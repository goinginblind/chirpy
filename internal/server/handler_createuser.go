package server

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/lib/pq"
)

func (s *Server) HandlerCreateUser(w http.ResponseWriter, r *http.Request) {
	var params createUserParams
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil || params.Email == "" {
		log.Printf("Fail to decode request body: %v", err)
		respondWithJSON(w, http.StatusBadRequest, map[string]string{"error": "incorrect request"})
		return
	}

	user, err := s.Cfg.DB.CreateUser(r.Context(), params.Email)

	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code == "23505" { // 23505 is a unique violation
				respondWithJSON(w, http.StatusConflict, map[string]string{"error": "email already exists"})
				return
			}
		}
		log.Printf("Fail to create user in the db: %v", err)
		respondWithJSON(w, http.StatusExpectationFailed, map[string]string{"error": "something went wrong"})
		return
	}

	respondWithJSON(w, http.StatusCreated, User{
		ID:        user.ID,
		Email:     user.Email,
		CreatedAt: user.CreatedAt.String(),
		UpdatedAt: user.UpdatedAt.String(),
	})
}
