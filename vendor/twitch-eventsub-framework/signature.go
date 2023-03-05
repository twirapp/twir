package eventsub_framework

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"net/http"
	"strings"
)

type signatureVerifyError struct {
	Message string
}

func (s *signatureVerifyError) Error() string {
	return s.Message
}

func VerifyRequestSignature(req *http.Request, body, secret []byte) (bool, error) {
	signatureValue := req.Header.Get("Twitch-Eventsub-Message-Signature")
	if signatureValue == "" {
		return false, &signatureVerifyError{"missing signature header"}
	}

	signatureBytes, err := getHmacBytes(signatureValue)
	if err != nil {
		return false, &signatureVerifyError{"invalid signature format"}
	}

	mac := hmac.New(sha256.New, secret)
	mac.Write([]byte(req.Header.Get("Twitch-Eventsub-Message-Id")))
	mac.Write([]byte(req.Header.Get("Twitch-Eventsub-Message-Timestamp")))
	mac.Write(body)
	outputHmac := mac.Sum(nil)

	return hmac.Equal(signatureBytes, outputHmac), nil
}

func getHmacBytes(sigValue string) ([]byte, error) {
	parts := strings.SplitN(sigValue, "=", 2)
	if len(parts) != 2 {
		return nil, errors.New("expected 2 components from signature header")
	}

	hexValue := parts[1]
	return hex.DecodeString(hexValue)
}
