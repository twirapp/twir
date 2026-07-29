package vkvideo

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/twirapp/twir/libs/crypto"
	"github.com/twirapp/twir/libs/oauth"
	tokens "github.com/twirapp/twir/libs/repositories/tokens"
	vkvideobots "github.com/twirapp/twir/libs/repositories/vk_video_bots"
)

type userStore struct {
	repository tokens.Repository
	cipherKey  string
}

func (store userStore) Load(ctx context.Context, key oauth.CredentialKey) (oauth.Credential, error) {
	userID, err := uuid.Parse(string(key.ID))
	if err != nil {
		return oauth.Credential{}, fmt.Errorf("parse VK Video user credential ID: %w", err)
	}
	token, err := store.repository.GetByUserID(ctx, userID)
	if err != nil {
		return oauth.Credential{}, fmt.Errorf("load VK Video user token: %w", err)
	}
	return decryptCredential(key, token.AccessToken, token.RefreshToken, token.Scopes, token.ObtainmentTimestamp, token.ExpiresIn, store.cipherKey)
}

func (store userStore) Commit(ctx context.Context, credential oauth.Credential) error {
	userID, err := uuid.Parse(string(credential.ID))
	if err != nil {
		return fmt.Errorf("parse VK Video user credential ID: %w", err)
	}
	token, err := store.repository.GetByUserID(ctx, userID)
	if err != nil {
		return fmt.Errorf("load VK Video user token for commit: %w", err)
	}
	input, err := encryptUserUpdate(credential, store.cipherKey)
	if err != nil {
		return err
	}
	if _, err := store.repository.UpdateTokenByID(ctx, token.ID, input); err != nil {
		return fmt.Errorf("persist VK Video user token: %w", err)
	}
	return nil
}

type singletonBotStore struct {
	repository vkvideobots.Repository
	cipherKey  string
}

func (store singletonBotStore) Load(ctx context.Context, key oauth.CredentialKey) (oauth.Credential, error) {
	bot, err := store.repository.Get(ctx)
	if err != nil {
		return oauth.Credential{}, fmt.Errorf("load VK Video bot singleton: %w", err)
	}
	return decryptCredential(key, bot.EncryptedAccessToken, bot.EncryptedRefreshToken, bot.Scopes, bot.ObtainmentTimestamp, bot.ExpiresIn, store.cipherKey)
}

func (store singletonBotStore) Commit(ctx context.Context, credential oauth.Credential) error {
	bot, err := store.repository.Get(ctx)
	if err != nil {
		return fmt.Errorf("load VK Video bot singleton for commit: %w", err)
	}
	accessToken, err := crypto.Encrypt(credential.AccessToken, store.cipherKey)
	if err != nil {
		return fmt.Errorf("encrypt VK Video bot access token: %w", err)
	}
	refreshToken, err := crypto.Encrypt(credential.RefreshToken, store.cipherKey)
	if err != nil {
		return fmt.Errorf("encrypt VK Video bot refresh token: %w", err)
	}
	if _, err := store.repository.Update(ctx, vkvideobots.UpdateInput{
		EncryptedAccessToken: accessToken, EncryptedRefreshToken: refreshToken, Scopes: credential.Scopes,
		ExpiresIn: int(credential.ExpiresIn / time.Second), ObtainmentTimestamp: credential.ObtainedAt,
		VKUserID: bot.VKUserID,
	}); err != nil {
		return fmt.Errorf("persist VK Video bot token: %w", err)
	}
	return nil
}

func decryptCredential(key oauth.CredentialKey, encryptedAccessToken string, encryptedRefreshToken string, scopes []string, obtainedAt time.Time, expiresIn int, cipherKey string) (oauth.Credential, error) {
	accessToken, err := crypto.Decrypt(encryptedAccessToken, cipherKey)
	if err != nil {
		return oauth.Credential{}, fmt.Errorf("decrypt VK Video access token: %w", err)
	}
	refreshToken, err := crypto.Decrypt(encryptedRefreshToken, cipherKey)
	if err != nil {
		return oauth.Credential{}, fmt.Errorf("decrypt VK Video refresh token: %w", err)
	}
	if refreshToken == "" {
		return oauth.Credential{}, fmt.Errorf("%w: %w", ErrMissingRefreshToken, oauth.ErrInvalidCredential)
	}
	return oauth.Credential{
		Provider: key.Provider, ID: key.ID, AccessToken: accessToken, RefreshToken: refreshToken,
		Scopes: append([]string(nil), scopes...), ObtainedAt: obtainedAt, ExpiresIn: time.Duration(expiresIn) * time.Second,
	}, nil
}

func encryptUserUpdate(credential oauth.Credential, cipherKey string) (tokens.UpdateTokenInput, error) {
	accessToken, err := crypto.Encrypt(credential.AccessToken, cipherKey)
	if err != nil {
		return tokens.UpdateTokenInput{}, fmt.Errorf("encrypt VK Video user access token: %w", err)
	}
	refreshToken, err := crypto.Encrypt(credential.RefreshToken, cipherKey)
	if err != nil {
		return tokens.UpdateTokenInput{}, fmt.Errorf("encrypt VK Video user refresh token: %w", err)
	}
	expiresIn := int(credential.ExpiresIn / time.Second)
	obtainedAt := credential.ObtainedAt
	return tokens.UpdateTokenInput{
		AccessToken: &accessToken, RefreshToken: &refreshToken, ExpiresIn: &expiresIn,
		ObtainmentTimestamp: &obtainedAt, Scopes: credential.Scopes,
	}, nil
}
