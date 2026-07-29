package auth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	kvoptions "github.com/twirapp/kv/options"
	appplatform "github.com/twirapp/twir/apps/api-gql/internal/platform"
	"github.com/twirapp/twir/libs/crypto"
	platformentity "github.com/twirapp/twir/libs/entities/platform"
	usersmodel "github.com/twirapp/twir/libs/repositories/users/model"
	youtubebotsrepo "github.com/twirapp/twir/libs/repositories/youtube_bots"
)

const youtubeBotSetupKVPrefix = "youtube_bot_setup"

var (
	ErrYouTubeBotNotConfigured     = errors.New("YouTube bot is not configured")
	ErrYouTubeBotSetupStateInvalid = errors.New("YouTube bot setup state is invalid or expired")
)

type youtubeBotSetupProvider interface {
	GetBotSetupAuthURL(state, codeChallenge string) string
	ExchangeBotSetupCode(context.Context, string, string) (*appplatform.PlatformTokens, error)
	GetUser(context.Context, string) (*appplatform.PlatformUser, error)
}

type youtubeBotSetupState struct {
	AdminUserID  uuid.UUID `json:"admin_user_id"`
	CodeVerifier string    `json:"code_verifier"`
}

func (a *Auth) StartYouTubeBotSetup(ctx context.Context) (string, error) {
	if a.youtubeBotProvider == nil || a.kv == nil {
		return "", fmt.Errorf("YouTube bot setup is not configured")
	}
	state, err := newYouTubeBotSetupState()
	if err != nil {
		return "", err
	}
	codeVerifier, codeChallenge, err := appplatform.GeneratePKCE()
	if err != nil {
		return "", fmt.Errorf("generate YouTube bot PKCE: %w", err)
	}
	admin, err := a.requireLiveYouTubeBotAdmin(ctx)
	if err != nil {
		return "", err
	}
	encodedState, err := json.Marshal(youtubeBotSetupState{AdminUserID: admin.ID, CodeVerifier: codeVerifier})
	if err != nil {
		return "", fmt.Errorf("encode YouTube bot setup state: %w", err)
	}
	if err := a.kv.Set(ctx, youtubeBotSetupKVPrefix+":"+state, encodedState, kvoptions.WithExpire(10*time.Minute)); err != nil {
		return "", fmt.Errorf("store YouTube bot setup state: %w", err)
	}

	return a.youtubeBotProvider.GetBotSetupAuthURL(state, codeChallenge), nil
}

func (a *Auth) CompleteYouTubeBotSetup(ctx context.Context, code, state string) error {
	setupState, err := a.consumeYouTubeBotSetupState(ctx, state)
	if err != nil {
		return err
	}
	admin, err := a.requireLiveYouTubeBotAdmin(ctx)
	if err != nil || admin.ID != setupState.AdminUserID {
		return ErrYouTubeBotSetupStateInvalid
	}
	if a.youtubeBotProvider == nil || a.youtubeBotsRepo == nil || a.transactionRunner == nil {
		return fmt.Errorf("YouTube bot setup is not configured")
	}

	tokens, err := a.youtubeBotProvider.ExchangeBotSetupCode(ctx, code, setupState.CodeVerifier)
	if err != nil {
		return fmt.Errorf("exchange YouTube bot setup code: %w", err)
	}
	platformUser, err := a.youtubeBotProvider.GetUser(ctx, tokens.AccessToken)
	if err != nil {
		return fmt.Errorf("get YouTube bot user: %w", err)
	}
	accessToken, err := crypto.Encrypt(tokens.AccessToken, a.config.TokensCipherKey)
	if err != nil {
		return fmt.Errorf("encrypt YouTube bot access token: %w", err)
	}
	refreshToken, err := crypto.Encrypt(tokens.RefreshToken, a.config.TokensCipherKey)
	if err != nil {
		return fmt.Errorf("encrypt YouTube bot refresh token: %w", err)
	}

	return a.transactionRunner.Do(ctx, func(txCtx context.Context) error {
		if err := a.youtubeBotsRepo.Lock(txCtx); err != nil {
			return fmt.Errorf("lock YouTube bot singleton: %w", err)
		}
		internalUser, _, err := a.getOrCreatePlatformUser(txCtx, platformentity.PlatformYouTube, platformUser)
		if err != nil {
			return fmt.Errorf("get or create YouTube bot user: %w", err)
		}
		if _, err := a.youtubeBotsRepo.Upsert(txCtx, youtubebotsrepo.UpsertInput{EncryptedAccessToken: accessToken, EncryptedRefreshToken: refreshToken, Scopes: tokens.Scopes, ExpiresIn: tokens.ExpiresIn, ObtainmentTimestamp: time.Now().UTC(), YouTubeUserID: internalUser.ID}); err != nil {
			return fmt.Errorf("upsert YouTube bot singleton: %w", err)
		}

		return nil
	})
}

func (a *Auth) YouTubeBotConfigured(ctx context.Context) (bool, error) {
	if a.youtubeBotsRepo == nil {
		return false, nil
	}
	_, err := a.youtubeBotsRepo.Get(ctx)
	if errors.Is(err, youtubebotsrepo.ErrNotFound) {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("get YouTube bot singleton: %w", err)
	}

	return true, nil
}

func (a *Auth) youtubeBotBindingConfig(ctx context.Context) (platformBindingConfig, error) {
	if a.youtubeBotsRepo == nil {
		return platformBindingConfig{}, ErrYouTubeBotNotConfigured
	}
	if err := a.youtubeBotsRepo.Lock(ctx); err != nil {
		return platformBindingConfig{}, fmt.Errorf("lock YouTube bot singleton: %w", err)
	}
	bot, err := a.youtubeBotsRepo.Get(ctx)
	if errors.Is(err, youtubebotsrepo.ErrNotFound) {
		return platformBindingConfig{}, ErrYouTubeBotNotConfigured
	}
	if err != nil {
		return platformBindingConfig{}, fmt.Errorf("get YouTube bot singleton: %w", err)
	}

	botUserID := bot.YouTubeUserID
	return platformBindingConfig{BotUserID: &botUserID}, nil
}

func (a *Auth) requireLiveYouTubeBotAdmin(ctx context.Context) (usersmodel.User, error) {
	user, live, err := a.getLiveSessionUser(ctx)
	if err != nil {
		return usersmodel.Nil, fmt.Errorf("get live YouTube bot admin: %w", err)
	}
	if !live || user.IsBanned || !user.IsBotAdmin {
		return usersmodel.Nil, ErrYouTubeBotSetupStateInvalid
	}

	return user, nil
}

func (a *Auth) consumeYouTubeBotSetupState(ctx context.Context, state string) (youtubeBotSetupState, error) {
	if a.kv == nil || state == "" {
		return youtubeBotSetupState{}, ErrYouTubeBotSetupStateInvalid
	}

	a.youtubeBotSetupStateMu.Lock()
	defer a.youtubeBotSetupStateMu.Unlock()
	key := youtubeBotSetupKVPrefix + ":" + state
	encodedState, err := a.kv.Get(ctx, key).Bytes()
	if err != nil {
		return youtubeBotSetupState{}, ErrYouTubeBotSetupStateInvalid
	}
	if err := a.kv.Delete(ctx, key); err != nil {
		return youtubeBotSetupState{}, fmt.Errorf("consume YouTube bot setup state: %w", err)
	}

	var setupState youtubeBotSetupState
	if err := json.Unmarshal(encodedState, &setupState); err != nil || setupState.AdminUserID == uuid.Nil || setupState.CodeVerifier == "" {
		return youtubeBotSetupState{}, ErrYouTubeBotSetupStateInvalid
	}

	return setupState, nil
}

func newYouTubeBotSetupState() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("generate YouTube bot setup state: %w", err)
	}

	return base64.RawURLEncoding.EncodeToString(bytes), nil
}
