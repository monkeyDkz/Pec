package auth

import "testing"

func TestBcryptHasher_HashAndVerify(t *testing.T) {
	h := NewBcryptHasher()

	hashed, err := h.Hash("correct-horse-battery-staple")
	if err != nil {
		t.Fatalf("hash: %v", err)
	}
	if hashed == "correct-horse-battery-staple" {
		t.Fatal("password was not hashed")
	}
	if !h.Verify("correct-horse-battery-staple", hashed) {
		t.Fatal("Verify returned false for correct password")
	}
	if h.Verify("wrong", hashed) {
		t.Fatal("Verify returned true for wrong password")
	}
}

func TestBcryptHasher_DistinctHashes(t *testing.T) {
	h := NewBcryptHasher()
	h1, _ := h.Hash("same-password")
	h2, _ := h.Hash("same-password")
	if h1 == h2 {
		t.Fatal("bcrypt should salt each hash; two hashes of the same password were identical")
	}
}
