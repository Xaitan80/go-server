package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

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
