package twitch

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/kvizyx/twitchy/helix"
	twitchoauth "github.com/kvizyx/twitchy/oauth"
	"github.com/twirapp/twir/libs/crypto"
	"github.com/twirapp/twir/libs/repositories/tokens"
	tokenmodel "github.com/twirapp/twir/libs/repositories/tokens/model"
)

type credentialBinding struct {
	repository tokens.Repository
	cipherKey  string
	kind       credentialKind
	now        func() time.Time
	tokenID    uuid.UUID
}

func newCredentialBinding(ctx context.Context, runtime *twitchRuntime, kind credentialKind, subjectID string) (credentialBinding, error) {
	binding := credentialBinding{
		repository: runtime.tokens,
		cipherKey:  runtime.cipherKey,
		kind:       kind,
		now:        runtime.now,
	}
	_, token, err := binding.load(ctx, subjectID)
	if err != nil {
		return credentialBinding{}, err
	}
	binding.tokenID = token.ID
	return binding, nil
}

func (binding credentialBinding) loader(ctx context.Context, subjectID string) (twitchoauth.TokenPair, error) {
	pair, _, err := binding.load(ctx, subjectID)
	return pair, err
}

func (binding credentialBinding) hook(ctx context.Context, pair twitchoauth.TokenPair) error {
	accessToken, err := crypto.Encrypt(pair.AccessToken, binding.cipherKey)
	if err != nil {
		return fmt.Errorf("encrypt Twitch access credential: %w", err)
	}
	refreshToken, err := crypto.Encrypt(pair.RefreshToken, binding.cipherKey)
	if err != nil {
		return fmt.Errorf("encrypt Twitch refresh credential: %w", err)
	}
	expiresIn := int(pair.ExpiresIn / time.Second)
	obtainedAt := binding.now()
	_, err = binding.repository.UpdateTokenByID(ctx, binding.tokenID, tokens.UpdateTokenInput{
		AccessToken: &accessToken, RefreshToken: &refreshToken, ExpiresIn: &expiresIn,
		ObtainmentTimestamp: &obtainedAt, Scopes: authorizationScopesToStrings(pair.Scopes),
	})
	if err != nil {
		return fmt.Errorf("persist Twitch credential rotation: %w", err)
	}
	return nil
}

func (binding credentialBinding) load(ctx context.Context, subjectID string) (twitchoauth.TokenPair, *tokenmodel.Token, error) {
	token, err := binding.token(ctx, subjectID)
	if err != nil {
		return twitchoauth.TokenPair{}, nil, err
	}
	accessToken, err := crypto.Decrypt(token.AccessToken, binding.cipherKey)
	if err != nil {
		return twitchoauth.TokenPair{}, nil, fmt.Errorf("decrypt Twitch access credential: %w", err)
	}
	refreshToken, err := crypto.Decrypt(token.RefreshToken, binding.cipherKey)
	if err != nil {
		return twitchoauth.TokenPair{}, nil, fmt.Errorf("decrypt Twitch refresh credential: %w", err)
	}

	remaining := token.ObtainmentTimestamp.Add(time.Duration(token.ExpiresIn) * time.Second).Sub(binding.now())
	if remaining < 0 {
		remaining = 0
	}
	return twitchoauth.TokenPair{
		AccessToken: accessToken, RefreshToken: refreshToken, ExpiresIn: remaining,
		Scopes: stringsToAuthorizationScopes(token.Scopes), TokenType: "bearer",
	}, token, nil
}

func (binding credentialBinding) token(ctx context.Context, subjectID string) (*tokenmodel.Token, error) {
	var (
		token *tokenmodel.Token
		err   error
	)
	switch binding.kind {
	case broadcasterCredential:
		userID, parseErr := uuid.Parse(subjectID)
		if parseErr != nil {
			return nil, fmt.Errorf("parse Twitch broadcaster ID: %w", parseErr)
		}
		token, err = binding.repository.GetByUserID(ctx, userID)
	case botCredential:
		token, err = binding.repository.GetByBotID(ctx, subjectID)
	default:
		return nil, fmt.Errorf("unknown Twitch credential kind")
	}
	if err != nil {
		if errors.Is(err, tokens.ErrNotFound) {
			return nil, fmt.Errorf("%w: %w", ErrCredentialNotFound, err)
		}
		return nil, fmt.Errorf("load Twitch credential: %w", err)
	}
	return token, nil
}

func stringsToAuthorizationScopes(scopes []string) []helix.AuthorizationScope {
	result := make([]helix.AuthorizationScope, len(scopes))
	for index, scope := range scopes {
		result[index] = helix.AuthorizationScope(scope)
	}
	return result
}

func authorizationScopesToStrings(scopes []helix.AuthorizationScope) []string {
	result := make([]string, len(scopes))
	for index, scope := range scopes {
		result[index] = string(scope)
	}
	return result
}
