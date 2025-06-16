package main

import (
	"log"
	"net/http"
)

// handlerReset just, well, resets the amount of 'hits' which are visits of the 'host:port/app'
func (cfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	if cfg.Platform != "dev" {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	if err := cfg.DB.DeleteUsers(r.Context()); err != nil {
		log.Printf("fail to reset users: %v", err)
		w.WriteHeader(http.StatusExpectationFailed)
		return
	}

	cfg.fileserverHits.Store(0)
	respondWithJSON(w, http.StatusOK, map[string]string{
		"message": "users reset, amount of hits set to 0",
	})
}
