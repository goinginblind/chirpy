package main

import "net/http"

// handlerReset just, well, resets the amount of 'hits' which are visits of the 'host:port/app'
func (cfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	cfg.fileserverHits.Store(0)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hits reset to 0"))
}
