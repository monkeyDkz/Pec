package auth

import (
	"strings"
	"testing"
	"time"
)

const testSecret = "test-secret-which-is-long-enough-to-be-acceptable"

func TestJWTManager_GenerateAndValidate(t *testing.T) {
	m := NewJWTManager(testSecret, time.Hour)

	token, err := m.Generate("user-id-1", "broadcaster")
	if err != nil {
		t.Fatalf("generate: %v", err)
	}
	if strings.Count(token, ".") != 2 {
		t.Fatalf("expected a JWT (3 parts), got %q", token)
	}

	claims, err := m.Validate(token)
	if err != nil {
		t.Fatalf("validate: %v", err)
	}
	if claims.UserID != "user-id-1" {
		t.Fatalf("unexpected user_id %q", claims.UserID)
	}
	if claims.Role != "broadcaster" {
		t.Fatalf("unexpected role %q", claims.Role)
	}
}

func TestJWTManager_RejectsExpiredToken(t *testing.T) {
	m := NewJWTManager(testSecret, -time.Minute) // already expired

	token, err := m.Generate("u", "user")
	if err != nil {
		t.Fatalf("generate: %v", err)
	}
	if _, err := m.Validate(token); err == nil {
		t.Fatal("expected validation error on expired token")
	}
}

func TestJWTManager_RejectsTamperedToken(t *testing.T) {
	m := NewJWTManager(testSecret, time.Hour)
	token, _ := m.Generate("u", "user")

	// flip the last char of the signature
	tampered := token[:len(token)-1]
	if last := token[len(token)-1]; last == 'A' {
		tampered += "B"
	} else {
		tampered += "A"
	}

	if _, err := m.Validate(tampered); err == nil {
		t.Fatal("expected validation error on tampered token")
	}
}

func TestJWTManager_RejectsTokenSignedWithDifferentSecret(t *testing.T) {
	m1 := NewJWTManager(testSecret, time.Hour)
	m2 := NewJWTManager("another-secret-also-long-enough-yeah", time.Hour)

	token, _ := m1.Generate("u", "user")
	if _, err := m2.Validate(token); err == nil {
		t.Fatal("expected validation error when secrets differ")
	}
}

func TestJWTManager_RejectsNoneAlg(t *testing.T) {
	m := NewJWTManager(testSecret, time.Hour)
	// Hand-crafted "none" token (header.payload.signature with alg=none)
	noneToken := "eyJhbGciOiJub25lIiwidHlwIjoiSldUIn0.eyJ1c2VyX2lkIjoieCJ9."
	if _, err := m.Validate(noneToken); err == nil {
		t.Fatal("expected validation error on alg=none token")
	}
}
