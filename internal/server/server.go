package server

import (
	"github.com/goinginblind/chirpy/internal/config"
)

// Server is essentially a wrapper so that config methods can be defined outside the config package itself.
type Server struct {
	Cfg *config.APIConfig
}
