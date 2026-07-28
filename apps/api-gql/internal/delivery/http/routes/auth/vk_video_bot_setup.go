package auth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
	kvoptions "github.com/twirapp/kv/options"
	appplatform "github.com/twirapp/twir/apps/api-gql/internal/platform"
	buscoreeventsub "github.com/twirapp/twir/libs/bus-core/eventsub"
	"github.com/twirapp/twir/libs/crypto"
	platformentity "github.com/twirapp/twir/libs/entities/platform"
	"github.com/twirapp/twir/libs/logger"
	usersmodel "github.com/twirapp/twir/libs/repositories/users/model"
	vkvideobotsrepo "github.com/twirapp/twir/libs/repositories/vk_video_bots"
)

const vkVideoBotSetupKVPrefix = "vk_video_bot_setup"

var (
	ErrVKVideoBotNotConfigured     = errors.New("VK Video bot is not configured")
	ErrVKVideoBotSetupStateInvalid = errors.New("VK Video bot setup state is invalid or expired")
)

type vkVideoBotSetupProvider interface {
	GetBotSetupAuthURL(string) (string, error)
	ExchangeBotSetupCode(context.Context, string) (*appplatform.PlatformTokens, error)
	GetUser(context.Context, string) (*appplatform.PlatformUser, error)
}

type vkVideoBotSetupState struct {
	AdminUserID uuid.UUID `json:"admin_user_id"`
}

func (a *Auth) StartVKVideoBotSetup(ctx context.Context) (string, error) {
	if a.vkVideoBotProvider == nil || a.kv == nil {
		return "", fmt.Errorf("VK Video bot setup is not configured")
	}

	state, err := newVKVideoBotSetupState()
	if err != nil {
		return "", err
	}
	admin, err := a.requireLiveVKVideoBotAdmin(ctx)
	if err != nil {
		return "", err
	}
	encodedState, err := json.Marshal(vkVideoBotSetupState{AdminUserID: admin.ID})
	if err != nil {
		return "", fmt.Errorf("encode VK Video bot setup state: %w", err)
	}
	if err := a.kv.Set(ctx, vkVideoBotSetupKVPrefix+":"+state, encodedState, kvoptions.WithExpire(10*time.Minute)); err != nil {
		return "", fmt.Errorf("store VK Video bot setup state: %w", err)
	}

	url, err := a.vkVideoBotProvider.GetBotSetupAuthURL(state)
	if err != nil {
		_ = a.kv.Delete(ctx, vkVideoBotSetupKVPrefix+":"+state)
		return "", fmt.Errorf("build VK Video bot authorization URL: %w", err)
	}

	return url, nil
}

func (a *Auth) CompleteVKVideoBotSetup(ctx context.Context, code, state string) error {
	setupState, err := a.consumeVKVideoBotSetupState(ctx, state)
	if err != nil {
		return err
	}
	admin, err := a.requireLiveVKVideoBotAdmin(ctx)
	if err != nil || admin.ID != setupState.AdminUserID {
		return ErrVKVideoBotSetupStateInvalid
	}
	if a.vkVideoBotProvider == nil || a.vkVideoBotsRepo == nil || a.channelPlatformsRepo == nil || a.transactionRunner == nil {
		return fmt.Errorf("VK Video bot setup is not configured")
	}

	tokens, err := a.vkVideoBotProvider.ExchangeBotSetupCode(ctx, code)
	if err != nil {
		return fmt.Errorf("exchange VK Video bot setup code: %w", err)
	}
	platformUser, err := a.vkVideoBotProvider.GetUser(ctx, tokens.AccessToken)
	if err != nil {
		return fmt.Errorf("get VK Video bot user: %w", err)
	}
	accessToken, err := crypto.Encrypt(tokens.AccessToken, a.config.TokensCipherKey)
	if err != nil {
		return fmt.Errorf("encrypt VK Video bot access token: %w", err)
	}
	refreshToken, err := crypto.Encrypt(tokens.RefreshToken, a.config.TokensCipherKey)
	if err != nil {
		return fmt.Errorf("encrypt VK Video bot refresh token: %w", err)
	}

	var reassignedChannelIDs []uuid.UUID
	if err := a.transactionRunner.Do(ctx, func(txCtx context.Context) error {
		if err := a.vkVideoBotsRepo.Lock(txCtx); err != nil {
			return fmt.Errorf("lock VK Video bot singleton: %w", err)
		}
		internalUser, _, err := a.getOrCreatePlatformUser(txCtx, platformentity.PlatformVKVideoLive, platformUser)
		if err != nil {
			return fmt.Errorf("get or create VK Video bot user: %w", err)
		}
		if _, err := a.vkVideoBotsRepo.Upsert(txCtx, vkvideobotsrepo.UpsertInput{
			EncryptedAccessToken: accessToken, EncryptedRefreshToken: refreshToken, Scopes: tokens.Scopes,
			ExpiresIn: tokens.ExpiresIn, ObtainmentTimestamp: time.Now().UTC(), VKUserID: internalUser.ID,
		}); err != nil {
			return fmt.Errorf("upsert VK Video bot singleton: %w", err)
		}
		channelIDs, err := a.channelPlatformsRepo.AssignVKVideoLiveBot(txCtx, internalUser.ID)
		if err != nil {
			return fmt.Errorf("assign VK Video bot to bindings: %w", err)
		}
		reassignedChannelIDs = channelIDs

		return nil
	}); err != nil {
		return err
	}

	a.publishVKVideoBotReassignment(ctx, reassignedChannelIDs)
	return nil
}

func (a *Auth) publishVKVideoBotReassignment(ctx context.Context, channelIDs []uuid.UUID) {
	if a.eventSubPublisher == nil {
		return
	}
	for _, channelID := range channelIDs {
		if err := a.eventSubPublisher.Publish(ctx, buscoreeventsub.EventsubSubscribeToAllEventsRequest{
			ChannelID: channelID.String(),
			Platform:  platformentity.PlatformVKVideoLive,
		}); err != nil && a.logger != nil {
			a.logger.ErrorContext(ctx, "cannot publish VK Video bot reassignment", logger.Error(err), slog.String("channel_id", channelID.String()))
		}
	}
}

func (a *Auth) VKVideoBotConfigured(ctx context.Context) (bool, error) {
	if a.vkVideoBotsRepo == nil {
		return false, nil
	}
	_, err := a.vkVideoBotsRepo.Get(ctx)
	if errors.Is(err, vkvideobotsrepo.ErrNotFound) {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("get VK Video bot singleton: %w", err)
	}

	return true, nil
}

func (a *Auth) vkVideoBotBindingConfig(ctx context.Context) (platformBindingConfig, error) {
	if a.vkVideoBotsRepo == nil {
		return platformBindingConfig{}, ErrVKVideoBotNotConfigured
	}
	if err := a.vkVideoBotsRepo.Lock(ctx); err != nil {
		return platformBindingConfig{}, fmt.Errorf("lock VK Video bot singleton: %w", err)
	}
	bot, err := a.vkVideoBotsRepo.Get(ctx)
	if errors.Is(err, vkvideobotsrepo.ErrNotFound) {
		return platformBindingConfig{}, ErrVKVideoBotNotConfigured
	}
	if err != nil {
		return platformBindingConfig{}, fmt.Errorf("get VK Video bot singleton: %w", err)
	}

	botUserID := bot.VKUserID
	return platformBindingConfig{BotUserID: &botUserID}, nil
}

func (a *Auth) requireLiveVKVideoBotAdmin(ctx context.Context) (usersmodel.User, error) {
	user, live, err := a.getLiveSessionUser(ctx)
	if err != nil {
		return usersmodel.Nil, fmt.Errorf("get live VK Video bot admin: %w", err)
	}
	if !live || user.IsBanned || !user.IsBotAdmin {
		return usersmodel.Nil, ErrVKVideoBotSetupStateInvalid
	}

	return user, nil
}

func (a *Auth) consumeVKVideoBotSetupState(ctx context.Context, state string) (vkVideoBotSetupState, error) {
	if a.kv == nil || state == "" {
		return vkVideoBotSetupState{}, ErrVKVideoBotSetupStateInvalid
	}

	a.vkVideoBotSetupStateMu.Lock()
	defer a.vkVideoBotSetupStateMu.Unlock()
	key := vkVideoBotSetupKVPrefix + ":" + state
	encodedState, err := a.kv.Get(ctx, key).Bytes()
	if err != nil {
		return vkVideoBotSetupState{}, ErrVKVideoBotSetupStateInvalid
	}
	if err := a.kv.Delete(ctx, key); err != nil {
		return vkVideoBotSetupState{}, fmt.Errorf("consume VK Video bot setup state: %w", err)
	}

	var setupState vkVideoBotSetupState
	if err := json.Unmarshal(encodedState, &setupState); err != nil || setupState.AdminUserID == uuid.Nil {
		return vkVideoBotSetupState{}, ErrVKVideoBotSetupStateInvalid
	}

	return setupState, nil
}

func newVKVideoBotSetupState() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("generate VK Video bot setup state: %w", err)
	}

	return base64.RawURLEncoding.EncodeToString(bytes), nil
}
