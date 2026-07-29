package kick

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/twirapp/twir/libs/crypto"
	"github.com/twirapp/twir/libs/oauth"
	kickbots "github.com/twirapp/twir/libs/repositories/kick_bots"
	tokens "github.com/twirapp/twir/libs/repositories/tokens"
)

type userStore struct {
	repository tokens.Repository
	cipherKey  string
}

func (store userStore) Load(ctx context.Context, key oauth.CredentialKey) (oauth.Credential, error) {
	userID, err := uuid.Parse(string(key.ID))
	if err != nil {
		return oauth.Credential{}, fmt.Errorf("parse Kick user credential ID: %w", err)
	}
	token, err := store.repository.GetByUserID(ctx, userID)
	if err != nil {
		return oauth.Credential{}, fmt.Errorf("load Kick user token: %w", err)
	}
	return decryptCredential(key, token.AccessToken, token.RefreshToken, token.Scopes, token.ObtainmentTimestamp, token.ExpiresIn, store.cipherKey)
}

func (store userStore) Commit(ctx context.Context, credential oauth.Credential) error {
	userID, err := uuid.Parse(string(credential.ID))
	if err != nil {
		return fmt.Errorf("parse Kick user credential ID: %w", err)
	}
	token, err := store.repository.GetByUserID(ctx, userID)
	if err != nil {
		return fmt.Errorf("load Kick user token for commit: %w", err)
	}
	input, err := encryptUpdateInput(credential, store.cipherKey)
	if err != nil {
		return err
	}
	if _, err := store.repository.UpdateTokenByID(ctx, token.ID, input); err != nil {
		return fmt.Errorf("persist Kick user token: %w", err)
	}
	return nil
}

type botStore struct {
	repository kickbots.Repository
	cipherKey  string
}

func (store botStore) Load(ctx context.Context, key oauth.CredentialKey) (oauth.Credential, error) {
	botID, err := uuid.Parse(string(key.ID))
	if err != nil {
		return oauth.Credential{}, fmt.Errorf("parse Kick bot credential ID: %w", err)
	}
	bot, err := store.repository.GetByID(ctx, botID)
	if err != nil {
		return oauth.Credential{}, fmt.Errorf("load Kick bot token: %w", err)
	}
	return decryptCredential(key, bot.AccessToken, bot.RefreshToken, bot.Scopes, bot.ObtainmentTimestamp, bot.ExpiresIn, store.cipherKey)
}

func (store botStore) Commit(ctx context.Context, credential oauth.Credential) error {
	botID, err := uuid.Parse(string(credential.ID))
	if err != nil {
		return fmt.Errorf("parse Kick bot credential ID: %w", err)
	}
	accessToken, err := crypto.Encrypt(credential.AccessToken, store.cipherKey)
	if err != nil {
		return fmt.Errorf("encrypt Kick bot access token: %w", err)
	}
	refreshToken, err := crypto.Encrypt(credential.RefreshToken, store.cipherKey)
	if err != nil {
		return fmt.Errorf("encrypt Kick bot refresh token: %w", err)
	}
	if _, err := store.repository.UpdateToken(ctx, botID, kickbots.UpdateTokenInput{
		AccessToken: accessToken, RefreshToken: refreshToken, Scopes: credential.Scopes,
		ExpiresIn: int(credential.ExpiresIn / time.Second), ObtainmentTimestamp: credential.ObtainedAt,
	}); err != nil {
		return fmt.Errorf("persist Kick bot token: %w", err)
	}
	return nil
}

func decryptCredential(key oauth.CredentialKey, encryptedAccessToken string, encryptedRefreshToken string, scopes []string, obtainedAt time.Time, expiresIn int, cipherKey string) (oauth.Credential, error) {
	accessToken, err := crypto.Decrypt(encryptedAccessToken, cipherKey)
	if err != nil {
		return oauth.Credential{}, fmt.Errorf("decrypt Kick access token: %w", err)
	}
	refreshToken, err := crypto.Decrypt(encryptedRefreshToken, cipherKey)
	if err != nil {
		return oauth.Credential{}, fmt.Errorf("decrypt Kick refresh token: %w", err)
	}
	if refreshToken == "" {
		return oauth.Credential{}, fmt.Errorf("Kick refresh token: %w", oauth.ErrInvalidCredential)
	}
	return oauth.Credential{Provider: key.Provider, ID: key.ID, AccessToken: accessToken, RefreshToken: refreshToken, Scopes: scopes, ObtainedAt: obtainedAt, ExpiresIn: time.Duration(expiresIn) * time.Second}, nil
}

func encryptUpdateInput(credential oauth.Credential, cipherKey string) (tokens.UpdateTokenInput, error) {
	accessToken, err := crypto.Encrypt(credential.AccessToken, cipherKey)
	if err != nil {
		return tokens.UpdateTokenInput{}, fmt.Errorf("encrypt Kick user access token: %w", err)
	}
	refreshToken, err := crypto.Encrypt(credential.RefreshToken, cipherKey)
	if err != nil {
		return tokens.UpdateTokenInput{}, fmt.Errorf("encrypt Kick user refresh token: %w", err)
	}
	expiresIn := int(credential.ExpiresIn / time.Second)
	obtainedAt := credential.ObtainedAt
	return tokens.UpdateTokenInput{AccessToken: &accessToken, RefreshToken: &refreshToken, ExpiresIn: &expiresIn, ObtainmentTimestamp: &obtainedAt, Scopes: credential.Scopes}, nil
}
