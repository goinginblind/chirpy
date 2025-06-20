package server

import (
	"log"
	"net/http"

	"github.com/goinginblind/chirpy/internal/auth"
)

func (s *Server) HandlerRevokeRefreshToken(w http.ResponseWriter, r *http.Request) {
	// Get refresh token from the request header
	bearerRefreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		log.Printf("Failed to get bearer token from the request: %v\n", err)
		respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}
	if err := s.Cfg.DB.RevokeRefreshToken(r.Context(), bearerRefreshToken); err != nil {
		log.Printf("Failed to revoke refresh token: %v\n", err)
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
