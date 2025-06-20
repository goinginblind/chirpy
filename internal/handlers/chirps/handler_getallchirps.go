package chirps

import (
	"log"
	"net/http"

	"github.com/goinginblind/chirpy/internal/config"
	"github.com/goinginblind/chirpy/internal/handlers"
)

func GetAll(cfg *config.APIConfig, w http.ResponseWriter, r *http.Request) {
	chirps, err := cfg.DB.GetAllChirps(r.Context())
	if err != nil {
		log.Printf("Failed to get all chirps: %v", err)
		handlers.RespondWithError(w, http.StatusBadRequest, "Fail to retrieve chirps from the database")
		return
	}
	out := convertManyDBChirps(chirps)
	handlers.RespondWithJSON(w, http.StatusOK, out)
}
