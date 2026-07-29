package twitch

import (
	"context"
	"errors"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/kvizyx/twitchy/helix"
	twitchoauth "github.com/kvizyx/twitchy/oauth"
	"github.com/twirapp/twir/libs/crypto"
	"github.com/twirapp/twir/libs/repositories/tokens"
	tokenmodel "github.com/twirapp/twir/libs/repositories/tokens/model"
)

const testCipherKey = "0123456789abcdef0123456789abcdef"

type tokenRepositoryFake struct {
	mu      sync.Mutex
	users   map[uuid.UUID]*tokenmodel.Token
	bots    map[string]*tokenmodel.Token
	updates []tokens.UpdateTokenInput
}

func (repository *tokenRepositoryFake) GetByID(context.Context, uuid.UUID) (*tokenmodel.Token, error) {
	return nil, tokens.ErrNotFound
}

func (repository *tokenRepositoryFake) GetByUserID(_ context.Context, userID uuid.UUID) (*tokenmodel.Token, error) {
	repository.mu.Lock()
	defer repository.mu.Unlock()
	token, found := repository.users[userID]
	if !found {
		return nil, tokens.ErrNotFound
	}
	return copyToken(token), nil
}

func (repository *tokenRepositoryFake) GetByBotID(_ context.Context, botID string) (*tokenmodel.Token, error) {
	repository.mu.Lock()
	defer repository.mu.Unlock()
	token, found := repository.bots[botID]
	if !found {
		return nil, tokens.ErrNotFound
	}
	return copyToken(token), nil
}

func (repository *tokenRepositoryFake) CreateUserToken(context.Context, tokens.CreateInput) (*tokenmodel.Token, error) {
	return nil, errors.New("not implemented")
}

func (repository *tokenRepositoryFake) UpdateTokenByID(_ context.Context, id uuid.UUID, input tokens.UpdateTokenInput) (*tokenmodel.Token, error) {
	repository.mu.Lock()
	defer repository.mu.Unlock()
	repository.updates = append(repository.updates, input)
	for _, token := range repository.users {
		if token.ID == id {
			applyTokenUpdate(token, input)
			return copyToken(token), nil
		}
	}
	for _, token := range repository.bots {
		if token.ID == id {
			applyTokenUpdate(token, input)
			return copyToken(token), nil
		}
	}
	return nil, tokens.ErrNotFound
}

func TestCredentialBinding_loads_remaining_lifetime_when_credential_is_valid(t *testing.T) {
	// Given
	now := time.Date(2026, time.July, 29, 12, 0, 0, 0, time.UTC)
	userID := uuid.New()
	repository := newTokenRepositoryFake(userID, now.Add(-30*time.Second), 90)
	binding := credentialBinding{repository: repository, cipherKey: testCipherKey, kind: broadcasterCredential, now: func() time.Time { return now }}

	// When
	pair, err := binding.loader(context.Background(), userID.String())

	// Then
	if err != nil {
		t.Fatal(err)
	}
	if pair.ExpiresIn != time.Minute || pair.AccessToken != "test-access" || pair.RefreshToken != "test-refresh" {
		t.Fatalf("loaded pair did not preserve credential data and remaining expiry")
	}
}

func TestCredentialBinding_persists_encrypted_rotation_when_hook_runs(t *testing.T) {
	// Given
	now := time.Date(2026, time.July, 29, 12, 0, 0, 0, time.UTC)
	userID := uuid.New()
	repository := newTokenRepositoryFake(userID, now, 3600)
	binding := credentialBinding{repository: repository, cipherKey: testCipherKey, now: func() time.Time { return now }, tokenID: repository.users[userID].ID}

	// When
	err := binding.hook(context.Background(), tokenPair("rotated-access", "rotated-refresh", time.Hour))

	// Then
	if err != nil {
		t.Fatal(err)
	}
	updated := repository.users[userID]
	if updated.AccessToken == "rotated-access" || updated.RefreshToken == "rotated-refresh" {
		t.Fatal("credential hook persisted plaintext")
	}
	accessToken, err := crypto.Decrypt(updated.AccessToken, testCipherKey)
	if err != nil || accessToken != "rotated-access" || updated.ExpiresIn != 3600 || len(updated.Scopes) != 1 {
		t.Fatal("credential hook did not persist the rotated credential")
	}
}

func TestCredentialBinding_returns_typed_redacted_error_when_token_is_missing(t *testing.T) {
	// Given
	sentinel := "TWIR_TWITCH_TEST_ACCESS"
	binding := credentialBinding{repository: &tokenRepositoryFake{users: map[uuid.UUID]*tokenmodel.Token{}, bots: map[string]*tokenmodel.Token{}}, cipherKey: testCipherKey, kind: botCredential, now: time.Now}

	// When
	_, err := binding.loader(context.Background(), sentinel)

	// Then
	if !errors.Is(err, ErrCredentialNotFound) || strings.Contains(err.Error(), sentinel) {
		t.Fatalf("missing credential error was not typed and redacted: %v", err)
	}
}

func newTokenRepositoryFake(userID uuid.UUID, obtainedAt time.Time, expiresIn int) *tokenRepositoryFake {
	accessToken, _ := crypto.Encrypt("test-access", testCipherKey)
	refreshToken, _ := crypto.Encrypt("test-refresh", testCipherKey)
	return &tokenRepositoryFake{users: map[uuid.UUID]*tokenmodel.Token{userID: {
		ID: uuid.New(), AccessToken: accessToken, RefreshToken: refreshToken, ExpiresIn: expiresIn,
		ObtainmentTimestamp: obtainedAt, Scopes: []string{"channel:manage:broadcast"},
	}}, bots: make(map[string]*tokenmodel.Token)}
}

func tokenPair(accessToken string, refreshToken string, expiresIn time.Duration) twitchoauth.TokenPair {
	return twitchoauth.TokenPair{
		AccessToken: accessToken, RefreshToken: refreshToken, ExpiresIn: expiresIn,
		Scopes: []helix.AuthorizationScope{"channel:manage:broadcast"},
	}
}

func applyTokenUpdate(token *tokenmodel.Token, input tokens.UpdateTokenInput) {
	if input.AccessToken != nil {
		token.AccessToken = *input.AccessToken
	}
	if input.RefreshToken != nil {
		token.RefreshToken = *input.RefreshToken
	}
	if input.ExpiresIn != nil {
		token.ExpiresIn = *input.ExpiresIn
	}
	if input.ObtainmentTimestamp != nil {
		token.ObtainmentTimestamp = *input.ObtainmentTimestamp
	}
	token.Scopes = append([]string(nil), input.Scopes...)
}

func copyToken(token *tokenmodel.Token) *tokenmodel.Token {
	copy := *token
	copy.Scopes = append([]string(nil), token.Scopes...)
	return &copy
}
