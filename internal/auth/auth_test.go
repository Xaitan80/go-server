package auth

import (
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestGetBearerToken(t *testing.T) {
	headers := http.Header{}
	headers.Set("Authorization", "Bearer mytoken123")

	token, err := GetBearerToken(headers)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if token != "mytoken123" {
		t.Errorf("expected token 'mytoken123', got %s", token)
	}

	// Test missing header
	headers = http.Header{}
	_, err = GetBearerToken(headers)
	if err == nil {
		t.Fatal("expected error for missing header, got nil")
	}

	// Test malformed header
	headers.Set("Authorization", "InvalidFormat")
	_, err = GetBearerToken(headers)
	if err == nil {
		t.Fatal("expected error for malformed header, got nil")
	}
}

func TestMakeAndValidateJWT(t *testing.T) {
	secret := "supersecret"
	userID := uuid.New()
	exp := time.Minute * 5

	token, err := MakeJWT(userID, secret, exp)
	if err != nil {
		t.Fatalf("failed to create JWT: %v", err)
	}

	id, err := ValidateJWT(token, secret)
	if err != nil {
		t.Fatalf("failed to validate JWT: %v", err)
	}

	if id != userID {
		t.Errorf("expected userID %v, got %v", userID, id)
	}
}

func TestExpiredJWT(t *testing.T) {
	secret := "supersecret"
	userID := uuid.New()

	token, err := MakeJWT(userID, secret, -time.Minute) // already expired
	if err != nil {
		t.Fatalf("failed to create JWT: %v", err)
	}

	_, err = ValidateJWT(token, secret)
	if err == nil {
		t.Fatal("expected error validating expired token, got nil")
	}
}

func TestWrongSecretJWT(t *testing.T) {
	secret := "supersecret"
	wrongSecret := "wrongsecret"
	userID := uuid.New()
	exp := time.Minute * 5

	token, err := MakeJWT(userID, secret, exp)
	if err != nil {
		t.Fatalf("failed to create JWT: %v", err)
	}

	_, err = ValidateJWT(token, wrongSecret)
	if err == nil {
		t.Fatal("expected error validating token with wrong secret, got nil")
	}
}
