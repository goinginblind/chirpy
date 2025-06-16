package server

import (
	"fmt"
	"net/http"
)

// Injects cfg.fileserverHits.Add(1) into the handler
func (s *Server) MiddlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s.Cfg.FileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (s *Server) HandlerMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf(`
<html>

<body>
	<h1>Welcome, Chirpy Admin</h1>
	<p>Chirpy has been visited %d times!</p>
</body>

</html>
	`, s.Cfg.FileserverHits.Load())))
}
