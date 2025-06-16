package server

import (
	"log"
	"net/http"
)

// handlerReset just, well, resets the amount of 'hits' which are visits of the 'host:port/app'
func (s *Server) HandlerReset(w http.ResponseWriter, r *http.Request) {
	if s.Cfg.Platform != "dev" {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	if err := s.Cfg.DB.DeleteUsers(r.Context()); err != nil {
		log.Printf("fail to reset users: %v", err)
		w.WriteHeader(http.StatusExpectationFailed)
		return
	}

	s.Cfg.FileserverHits.Store(0)
	respondWithJSON(w, http.StatusOK, map[string]string{
		"message": "users reset, amount of hits set to 0",
	})
}
