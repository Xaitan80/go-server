package auth

import "testing"

func TestHashAndCheckPassword(t *testing.T) {
	password := "mysecret123"

	// Hash the password
	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("failed to hash password: %v", err)
	}

	// Hash should not equal the plain password
	if hash == password {
		t.Errorf("expected hashed password to differ from plain password")
	}

	// Correct password should validate
	if err := CheckPasswordHash(password, hash); err != nil {
		t.Errorf("expected password check to succeed, got error: %v", err)
	}

	// Wrong password should fail
	wrong := "wrongpass"
	if err := CheckPasswordHash(wrong, hash); err == nil {
		t.Errorf("expected password check to fail for wrong password")
	}
}

func TestHashPasswordProducesDifferentHashes(t *testing.T) {
	password := "repeatpass"

	hash1, err := HashPassword(password)
	if err != nil {
		t.Fatalf("failed to hash password 1: %v", err)
	}

	hash2, err := HashPassword(password)
	if err != nil {
		t.Fatalf("failed to hash password 2: %v", err)
	}

	// Bcrypt adds a random salt, so even same password hashes should differ
	if hash1 == hash2 {
		t.Errorf("expected different hashes for same password, got identical")
	}
}
