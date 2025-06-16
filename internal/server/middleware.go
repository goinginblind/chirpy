package server

import "net/http"

// Injects cfg.fileserverHits.Add(1) into the handler
func (s *Server) MiddlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s.Cfg.FileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}
