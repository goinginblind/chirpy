package tokens

import (
	"log"
	"net/http"

	"github.com/goinginblind/chirpy/internal/auth"
	"github.com/goinginblind/chirpy/internal/config"
	"github.com/goinginblind/chirpy/internal/handlers"
)

func RevokeRefreshToken(cfg *config.APIConfig, w http.ResponseWriter, r *http.Request) {
	// Get refresh token from the request header
	bearerRefreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		log.Printf("Failed to get bearer token from the request: %v\n", err)
		handlers.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}
	if err := cfg.DB.RevokeRefreshToken(r.Context(), bearerRefreshToken); err != nil {
		log.Printf("Failed to revoke refresh token: %v\n", err)
		handlers.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
