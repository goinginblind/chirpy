package auth

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/google/uuid"
)

func TestJWT(t *testing.T) {
	uID := uuid.New()
	tokenRealSec := "ihatecryptography"
	tokenMockSec := "mockupsecrettoken"

	token, err := MakeJWT(uID, tokenRealSec)
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

	// Old part of the test that was used when the epiration time limit of a token could be set by client
	// time.Sleep(expIn)
	// time.Sleep(expIn)
	//
	// parsedIDReal, err = ValidateJWT(token, tokenRealSec)
	// if err == nil || err.Error() != "fail to parse token: token has invalid claims: token is expired" {
	// 	t.Errorf("fail to find out expiration of a token: %v", err)
	// }
	// if parsedIDReal == uID {
	// 	t.Errorf("got real user id from ValidateJWT after expiration")
	// }
}

func TestGetBearerMany(t *testing.T) {
	uID := uuid.New()
	tokenIn, _ := MakeJWT(uID, "ihatecryptography")

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
