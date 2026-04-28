package kickplatform

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
)

const codeVerifierLength = 128

func GenerateCodeVerifier() (string, error) {
	randomBytes := make([]byte, 96)
	if _, err := rand.Read(randomBytes); err != nil {
		return "", fmt.Errorf("generate random pkce verifier bytes: %w", err)
	}

	verifier := base64.RawURLEncoding.EncodeToString(randomBytes)
	if len(verifier) > codeVerifierLength {
		verifier = verifier[:codeVerifierLength]
	}

	return verifier, nil
}

func GenerateCodeChallenge(verifier string) string {
	hash := sha256.Sum256([]byte(verifier))
	return base64.RawURLEncoding.EncodeToString(hash[:])
}
