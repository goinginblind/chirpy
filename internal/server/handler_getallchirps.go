package server

import (
	"log"
	"net/http"
)

func (s *Server) HandlerGetAllChirps(w http.ResponseWriter, r *http.Request) {
	chirps, err := s.Cfg.DB.GetAllChirps(r.Context())
	if err != nil {
		log.Printf("Failed to get all chirps: %v", err)
		respondWithJSON(w, http.StatusExpectationFailed, map[string]string{"error": "fail to retrieve chirps from the database"})
		return
	}
	out := convertManyDBChirps(chirps)
	respondWithJSON(w, http.StatusOK, out)
}
