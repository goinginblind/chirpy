package auth

import (
	"testing"
	"time"
)

func TestHashPasswordAndCheckPasswordHash(t *testing.T) {
	password := "mySecret123!"
	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword failed: %v", err)
	}

	// Correct password should not return error
	if err := CheckPasswordHash(hash, password); err != nil {
		t.Errorf("CheckPasswordHash failed with correct password: %v", err)
	}

	// Incorrect password should return error
	if err := CheckPasswordHash(hash, "wrongPassword"); err == nil {
		t.Error("CheckPasswordHash did not fail with incorrect password")
	}
}

func TestHashPasswordReturnsDifferentHashes(t *testing.T) {
	password := "repeatPassword"
	hash1, err1 := HashPassword(password)
	hash2, err2 := HashPassword(password)
	if err1 != nil || err2 != nil {
		t.Fatalf("HashPassword failed: %v, %v", err1, err2)
	}
	if hash1 == hash2 {
		t.Error("HashPassword returned the same hash for the same password (should be different due to salt)")
	}
}

func TestBenchmarkFunction(t *testing.T) {
	start := time.Now()
	benchmark()
	elapsed := time.Since(start)
	if elapsed <= 0 {
		t.Error("Benchmark function did not take any measurable time")
	}
}
