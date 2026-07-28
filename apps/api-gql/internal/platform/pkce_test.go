package platform

import (
	"crypto/sha256"
	"encoding/base64"
	"strings"
	"testing"
)

func TestGeneratePKCEProducesRFC7636Pair(t *testing.T) {
	t.Parallel()

	verifier, challenge, err := GeneratePKCE()
	if err != nil {
		t.Fatalf("GeneratePKCE() error = %v", err)
	}

	if len(verifier) != 128 {
		t.Errorf("verifier length = %d, want 128", len(verifier))
	}
	if strings.ContainsAny(verifier, "=+/") {
		t.Errorf("verifier must be unpadded base64url, got %q", verifier)
	}

	hash := sha256.Sum256([]byte(verifier))
	expectedChallenge := base64.RawURLEncoding.EncodeToString(hash[:])
	if challenge != expectedChallenge {
		t.Errorf("challenge = %q, want base64url(sha256(verifier)) %q", challenge, expectedChallenge)
	}
	if strings.ContainsAny(challenge, "=+/") {
		t.Errorf("challenge must be unpadded base64url, got %q", challenge)
	}
}

func TestGeneratePKCEProducesUniqueVerifiers(t *testing.T) {
	t.Parallel()

	first, _, err := GeneratePKCE()
	if err != nil {
		t.Fatalf("GeneratePKCE() error = %v", err)
	}
	second, _, err := GeneratePKCE()
	if err != nil {
		t.Fatalf("GeneratePKCE() error = %v", err)
	}
	if first == second {
		t.Fatal("two GeneratePKCE() calls produced the same verifier")
	}
}
