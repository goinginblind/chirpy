package auth

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestJWT(t *testing.T) {
	uID := uuid.New()
	tokenRealSec := "ihatecryptography"
	tokenMockSec := "mockupsecrettoken"
	expIn, _ := time.ParseDuration("5s")

	token, err := MakeJWT(uID, tokenRealSec, expIn)
	if err != nil {
		t.Errorf("fail to make a jwt token: %v", err)
	}

	parsedIDReal, err := ValidateJWT(token, tokenRealSec)
	if err != nil {
		t.Errorf("fail to validate token: %v", err)
	}
	if parsedIDReal != uID {
		t.Errorf("got different user id from ValidateJWT, exp: %v, got: %v", uID, parsedIDReal)
	}

	parsedIDMock, err := ValidateJWT(token, tokenMockSec)
	if err == nil || err.Error() != "fail to parse token: token signature is invalid: signature is invalid" {
		t.Errorf("could not identify whether the token is invalid, got err: %v", err)
	}
	if parsedIDMock == uID {
		t.Errorf("got same user id from ValidateJWT, exp: %v, got: %v", uID, parsedIDMock)
	}

	time.Sleep(expIn)

	parsedIDReal, err = ValidateJWT(token, tokenRealSec)
	if err == nil || err.Error() != "fail to parse token: token has invalid claims: token is expired" {
		t.Errorf("fail to find out expiration of a token: %v", err)
	}
	if parsedIDReal == uID {
		t.Errorf("got real user id from ValidateJWT after expiration")
	}
}

func TestGetBearerMany(t *testing.T) {
	uID := uuid.New()
	expIn, _ := time.ParseDuration("5s")
	tokenIn, _ := MakeJWT(uID, "ihatecryptography", expIn)

	testCases := []struct {
		name     string
		header   http.Header
		expected string
	}{
		{"valid bearer", http.Header{"Authorization": []string{fmt.Sprintf("Bearer %v", tokenIn)}}, ""},
		{"no authorization header", http.Header{"auth": []string{fmt.Sprintf("Bearer %v", tokenIn)}}, "authorization header is empty"},
		{"invalid bearer", http.Header{"Authorization": []string{fmt.Sprintf("Bear %v", tokenIn)}}, "invalid authorization header"},
	}

	for _, tc := range testCases {
		tokenOut, err := GetBearerToken(tc.header)
		if err == nil && tc.expected == "" {
			if tokenOut != tokenIn {
				t.Errorf("[%s] expected token: %q; got: %q", tc.name, tokenIn, tokenOut)
			}
			continue
		}
		if err.Error() != tc.expected {
			t.Errorf("[%s] expected: '%v'; got: '%v'", tc.name, tc.expected, err.Error())
		}

	}
}
