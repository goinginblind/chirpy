package config

import (
	"sync/atomic"

	"github.com/goinginblind/chirpy/internal/database"
)

type APIConfig struct {
	DB             *database.Queries
	DBUrl          string
	FilepathRoot   string
	FileserverHits atomic.Int32
	MaxChirpLen    int
	Port           string
	Platform       string
	TokenSecret    string
	PolkaKey       string
}
