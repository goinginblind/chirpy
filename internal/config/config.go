package config

import (
	"sync/atomic"

	"github.com/goinginblind/chirpy/internal/database"
)

type APIConfig struct {
	FileserverHits atomic.Int32
	DB             *database.Queries
	Platform       string
}
