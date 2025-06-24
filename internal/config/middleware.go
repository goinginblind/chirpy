package config

import "net/http"

// Injects cfg.fileserverHits.Add(1) into the handler
func (cfg *APIConfig) MiddlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.FileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

// Injects config into a handler.
// Done so that handlers and config can be in different packages,
// but that requires handlers to be non-methods, so config is passed as an argument in the handler.
func (cfg *APIConfig) InjectConfig(f func(*APIConfig, http.ResponseWriter, *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		f(cfg, w, r)
	}
}
