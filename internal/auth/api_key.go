// Package auth provides functions to check create and check authentication and refresh tokens,
// hash and verify hashed passwords and get api keys from request headers.
package auth

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

func GetAPIKey(headers http.Header) (string, error) {
	// format of Authorization field in the header is 'ApiKey THE_KEY_HERE'
	authHeader := headers.Get("Authorization")
	if authHeader == "" {
		log.Println("Authorization header missing")
		return "", fmt.Errorf("authorization header is empty")
	}

	const prefix = "ApiKey "
	if !strings.HasPrefix(authHeader, prefix) {
		log.Println("Authorization header has invalid format")
		return "", fmt.Errorf("invalid authorization header")
	}

	return strings.TrimPrefix(authHeader, prefix), nil
}
