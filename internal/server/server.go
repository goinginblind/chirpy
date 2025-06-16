package server

import (
	"github.com/goinginblind/chirpy/internal/config"
)

// Server is basically a wrapper so that config methods can be defined outside config package.
type Server struct {
	Cfg config.APIConfig
}
