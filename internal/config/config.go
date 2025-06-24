// Contains APIConfig struct and two middleware functions.
// The first injects cfg into the handler, the second adds metrics.
package config

import (
	"sync/atomic"

	"github.com/goinginblind/chirpy/internal/database"
)

type APIConfig struct {
	// Database itself and the URL that is provided by psql
	DB    *database.Queries
	DBUrl string
	// Projects root directory
	FilepathRoot string
	// Amount of requests sent to the server
	FileserverHits atomic.Int32
	MaxChirpLen    int
	Port           string
	// Platform allows the reset handler to check if its reset by a dev or not
	Platform string
	// Secrets used to decipher jwt tokens
	TokenSecret string
	// And identify the real polka api webhook
	PolkaKey string
}
