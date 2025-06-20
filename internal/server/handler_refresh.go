package server

import (
	"log"
	"net/http"
	"time"

	"github.com/goinginblind/chirpy/internal/auth"
)

// Creates new auth token for the user from their refresh token sent in the request 'Authorization: Bearer [REFR_TOK]'
func (s *Server) HandlerRefreshAccessToken(w http.ResponseWriter, r *http.Request) {
	// Get refresh token from the request header
	bearerRefreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		log.Printf("Failed to get bearer token from the request: %v\n", err)
		respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Get refresh token from DB
	refreshToken, err := s.Cfg.DB.GetRefreshToken(r.Context(), bearerRefreshToken)
	if err != nil {
		log.Printf("Failed to get refresh token from DB: %v\n", err)
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}
	// Check if its expired
	if refreshToken.ExpiresAt.Before(time.Now()) {
		log.Println("Expired refresh token")
		respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}
	// Get user from the DB by their refresh token
	user, err := s.Cfg.DB.GetUserByRefreshToken(r.Context(), refreshToken.Token)
	if err != nil {
		log.Printf("Failed to get user by refresh token: %v\n", err)
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}
	// Create new auth token
	newAuthToken, err := auth.MakeJWT(user.ID, s.Cfg.TokenSecret)
	if err != nil {
		log.Printf("Failed to make new JWT for user: %v\n", err)
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Send it back
	respondWithJSON(w, http.StatusOK, map[string]string{
		"token": newAuthToken,
	})
}
