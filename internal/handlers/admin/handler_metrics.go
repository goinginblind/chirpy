// Package admin provides handlers which allow to check the metricss, readiness and reset the DB
package admin

import (
	"fmt"
	"net/http"

	"github.com/goinginblind/chirpy/internal/config"
)

func Metrics(cfg *config.APIConfig, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `
<html>

<body>
	<h1>Welcome, Chirpy Admin</h1>
	<p>Chirpy has been visited %d times!</p>
</body>

</html>
	`, cfg.FileserverHits.Load())
}
