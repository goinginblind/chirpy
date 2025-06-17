package app

import (
	"fmt"
	"os"
	"strconv"

	"github.com/goinginblind/chirpy/internal/config"
)

func loadConfigFromEnv() (*config.APIConfig, error) {
	filepathRoot := os.Getenv("FILEPATH_ROOT")
	platform := os.Getenv("PLATFORM")
	maxChirpLen, err := strconv.Atoi(os.Getenv("MAX_MSG_LEN"))
	if err != nil {
		maxChirpLen = 140
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	dbURL := os.Getenv("DB_URL")
	tokenSecret := os.Getenv("TOKEN_SECRET")
	if dbURL == "" || tokenSecret == "" {
		return nil, fmt.Errorf("missing essential enviromental variables")
	}

	return &config.APIConfig{
		DBUrl:        dbURL,
		FilepathRoot: filepathRoot,
		MaxChirpLen:  maxChirpLen,
		Port:         port,
		Platform:     platform,
		TokenSecret:  tokenSecret,
	}, nil
}
