package auth

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	now := time.Now().UTC()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    "chirpy",
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(now.Add(expiresIn)),
		Subject:   userID.String(),
	})
	return token.SignedString([]byte(tokenSecret))
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	claims := &jwt.RegisteredClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(tokenSecret), nil
	})
	if err != nil {
		return uuid.Nil, fmt.Errorf("fail to parse token: %w", err)
	}
	if !token.Valid {
		return uuid.Nil, fmt.Errorf("invalid token")
	}
	if claims.ExpiresAt != nil && claims.ExpiresAt.Time.Before(time.Now()) {
		return uuid.Nil, fmt.Errorf("expired token")
	}

	return uuid.MustParse(claims.Subject), nil
}

func GetBearerToken(headers http.Header) (string, error) {
	// format of Authorization field in the header is 'Bearer TOKEN_STRING'
	authHeader := headers.Get("Authorization")
	if authHeader == "" {
		log.Println("Authorization header missing")
		return "", fmt.Errorf("authorization header is empty")
	}

	const prefix = "Bearer "
	if !strings.HasPrefix(authHeader, prefix) {
		log.Println("Authorization header has invalid format")
		return "", fmt.Errorf("invalid authorization header")
	}

	return strings.TrimPrefix(authHeader, prefix), nil
}
