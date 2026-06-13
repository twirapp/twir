package resolvers

import (
	"crypto/rand"
	"encoding/base64"
)

const kickBotSetupKvPrefix = "kick_bot_setup"

type kickBotSetupState struct {
	CodeVerifier string `json:"code_verifier"`
	AdminUserID  string `json:"admin_user_id"`
}

func generateSecureState() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}
