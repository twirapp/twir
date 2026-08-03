package mcp_oauth

import "testing"

func TestCredentialHash_Bytes_returns_an_independent_copy(t *testing.T) {
	// Given
	hash := CredentialHash{0: 1, 31: 2}

	// When
	bytes := hash.Bytes()
	bytes[0] = 9

	// Then
	if hash[0] != 1 {
		t.Fatalf("hash mutated through bytes view: got %d", hash[0])
	}
}
