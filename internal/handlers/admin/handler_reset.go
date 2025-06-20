package admin

import (
	"log"
	"net/http"

	"github.com/goinginblind/chirpy/internal/config"
	"github.com/goinginblind/chirpy/internal/handlers"
)

// handlerReset just, well, deletes all users from the db and resets the amount of 'hits' which are visits of the 'host:port/app'
func Reset(cfg *config.APIConfig, w http.ResponseWriter, r *http.Request) {
	if cfg.Platform != "dev" {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	if err := cfg.DB.DeleteUsers(r.Context()); err != nil {
		log.Printf("fail to reset users: %v", err)
		w.WriteHeader(http.StatusExpectationFailed)
		return
	}

	cfg.FileserverHits.Store(0)
	handlers.RespondWithJSON(w, http.StatusOK, map[string]string{
		"message": "users reset, amount of hits set to 0",
	})
}
