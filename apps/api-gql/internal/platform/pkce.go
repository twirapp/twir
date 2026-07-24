package platform

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
)

const pkceRandomBytes = 96

func GeneratePKCE() (verifier, challenge string, err error) {
	randomBytes := make([]byte, pkceRandomBytes)
	if _, err := rand.Read(randomBytes); err != nil {
		return "", "", fmt.Errorf("generate random pkce verifier bytes: %w", err)
	}

	verifier = base64.RawURLEncoding.EncodeToString(randomBytes)
	hash := sha256.Sum256([]byte(verifier))
	challenge = base64.RawURLEncoding.EncodeToString(hash[:])

	return verifier, challenge, nil
}
