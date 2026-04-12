package auth

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
	authsessions "github.com/twirapp/twir/apps/api-gql/internal/auth"
	httpdelivery "github.com/twirapp/twir/apps/api-gql/internal/delivery/http"
	buscoreeventsub "github.com/twirapp/twir/libs/bus-core/eventsub"
	"github.com/twirapp/twir/libs/bus-core/scheduler"
	"github.com/twirapp/twir/libs/crypto"
	"github.com/twirapp/twir/libs/entities/platform"
	channelsrepo "github.com/twirapp/twir/libs/repositories/channels"
	channelsmodel "github.com/twirapp/twir/libs/repositories/channels/model"
	userplatformaccountsrepo "github.com/twirapp/twir/libs/repositories/user_platform_accounts"
)

type kickCodeBody struct {
	Code  string `json:"code" minLength:"1" required:"true"`
	State string `json:"state" required:"true"`
}

func (a *Auth) handleKickCode(
	ctx context.Context,
	input kickCodeBody,
) (*httpdelivery.BaseOutputJson[authResponseDto], error) {
	redirectTo, err := decodeRedirectState(input.State)
	if err != nil {
		return nil, huma.Error400BadRequest("Cannot decode state", err)
	}

	if a.kickProvider == nil {
		return nil, huma.Error500InternalServerError("Kick provider is not configured", fmt.Errorf("kick provider is nil"))
	}

	if a.userPlatformAccountsRepo == nil {
		return nil, huma.Error500InternalServerError("User platform accounts repository is not configured", fmt.Errorf("user platform accounts repo is nil"))
	}

	if a.channelsRepo == nil {
		return nil, huma.Error500InternalServerError("Channels repository is not configured", fmt.Errorf("channels repo is nil"))
	}

	if a.botsRepo == nil {
		return nil, huma.Error500InternalServerError("Bots repository is not configured", fmt.Errorf("bots repo is nil"))
	}

	if a.pgxPool == nil {
		return nil, huma.Error500InternalServerError("Postgres pool is not configured", fmt.Errorf("pgx pool is nil"))
	}

	codeVerifier, err := a.getKickCodeVerifier(ctx)
	if err != nil {
		return nil, huma.Error400BadRequest("Cannot get code verifier", err)
	}

	tokens, err := a.kickProvider.ExchangeCode(ctx, input.Code, codeVerifier)
	if err != nil {
		return nil, huma.Error500InternalServerError("Cannot exchange code", err)
	}

	platformUser, err := a.kickProvider.GetUser(ctx, tokens.AccessToken)
	if err != nil {
		return nil, huma.Error500InternalServerError("Cannot get user data from kick", err)
	}

	encryptedAccess, err := crypto.Encrypt(tokens.AccessToken, a.config.TokensCipherKey)
	if err != nil {
		return nil, huma.Error500InternalServerError("Cannot encrypt user access token", err)
	}

	encryptedRefresh, err := crypto.Encrypt(tokens.RefreshToken, a.config.TokensCipherKey)
	if err != nil {
		return nil, huma.Error500InternalServerError("Cannot encrypt user refresh token", err)
	}

	account, err := a.userPlatformAccountsRepo.GetByPlatformUserID(ctx, platform.PlatformKick, platformUser.ID)
	if err != nil && !errors.Is(err, userplatformaccountsrepo.ErrNotFound) {
		return nil, huma.Error500InternalServerError("Cannot get user platform account", err)
	}

	defaultBot, err := a.botsRepo.GetDefault(ctx)
	if err != nil {
		return nil, huma.Error500InternalServerError("Cannot find default bot", err)
	}
	if defaultBot.ID == "" {
		return nil, huma.Error500InternalServerError("Cannot find default bot", fmt.Errorf("no default bot found"))
	}

	userID := uuid.Nil
	createdUser := false
	if errors.Is(err, userplatformaccountsrepo.ErrNotFound) {
		userID, err = a.createKickUser(ctx)
		if err != nil {
			return nil, huma.Error500InternalServerError("Cannot create user", err)
		}
		createdUser = true
	} else {
		userID = account.UserID
	}

	_, err = a.userPlatformAccountsRepo.Upsert(ctx, userplatformaccountsrepo.UpsertInput{
		UserID:              userID,
		Platform:            platform.PlatformKick,
		PlatformUserID:      platformUser.ID,
		PlatformLogin:       platformUser.Login,
		PlatformDisplayName: platformUser.DisplayName,
		PlatformAvatar:      platformUser.Avatar,
		AccessToken:         encryptedAccess,
		RefreshToken:        encryptedRefresh,
		Scopes:              tokens.Scopes,
		ExpiresIn:           tokens.ExpiresIn,
		ObtainmentTimestamp: time.Now().UTC(),
	})
	if err != nil {
		return nil, huma.Error500InternalServerError("Cannot upsert user platform account", err)
	}

	channel, err := a.channelsRepo.GetByUserIDAndPlatform(ctx, userID, platform.PlatformKick)
	if err != nil {
		if !errors.Is(err, channelsrepo.ErrNotFound) {
			return nil, huma.Error500InternalServerError("Cannot get channel", err)
		}

		channel, err = a.createChannel(ctx, userID, defaultBot.ID)
		if err != nil {
			return nil, huma.Error500InternalServerError("Cannot create channel", err)
		}
	}

	if createdUser {
		err = a.bus.Scheduler.CreateDefaultRoles.Publish(
			ctx,
			scheduler.CreateDefaultRolesRequest{ChannelsIDs: []string{channel.ID}},
		)
		if err != nil {
			return nil, huma.Error500InternalServerError("Cannot create default roles", err)
		}

		err = a.bus.Scheduler.CreateDefaultCommands.Publish(
			ctx,
			scheduler.CreateDefaultCommandsRequest{ChannelsIDs: []string{channel.ID}},
		)
		if err != nil {
			return nil, huma.Error500InternalServerError("Cannot create default commands", err)
		}
	}

	if err := a.bus.EventSub.SubscribeToAllEvents.Publish(
		ctx,
		buscoreeventsub.EventsubSubscribeToAllEventsRequest{ChannelID: channel.ID},
	); err != nil {
		return nil, huma.Error500InternalServerError("Cannot subscribe to eventsub", err)
	}

	a.sessions.Put(ctx, "internalUserId", userID)
	a.sessions.Put(ctx, "currentPlatform", platform.PlatformKick.String())
	a.sessions.Put(ctx, "selectedDashboardId", uuid.MustParse(channel.ID))
	a.sessions.Put(ctx, "kickUser", authsessions.KickSessionUser{
		ID:     platformUser.ID,
		Login:  platformUser.Login,
		Avatar: platformUser.Avatar,
	})

	if err := a.sessions.Commit(ctx); err != nil {
		return nil, huma.Error500InternalServerError("Cannot commit session", err)
	}

	return httpdelivery.CreateBaseOutputJson(authResponseDto{RedirectTo: string(redirectTo)}), nil
}

func decodeRedirectState(state string) ([]byte, error) {
	decoded, err := base64.URLEncoding.DecodeString(state)
	if err == nil {
		return decoded, nil
	}

	decoded, rawErr := base64.RawURLEncoding.DecodeString(state)
	if rawErr == nil {
		return decoded, nil
	}

	decoded, stdErr := base64.StdEncoding.DecodeString(state)
	if stdErr == nil {
		return decoded, nil
	}

	decoded, rawStdErr := base64.RawStdEncoding.DecodeString(state)
	if rawStdErr == nil {
		return decoded, nil
	}

	return nil, fmt.Errorf("decode state: %w", err)
}

func (a *Auth) getKickCodeVerifier(ctx context.Context) (string, error) {
	codeVerifier, ok := a.sessions.Get(ctx, kickCodeVerifierSessionKey).(string)
	if !ok || codeVerifier == "" {
		return "", fmt.Errorf("kick code verifier not found in session")
	}

	return codeVerifier, nil
}

func (a *Auth) createKickUser(ctx context.Context) (uuid.UUID, error) {
	var id uuid.UUID
	err := a.pgxPool.QueryRow(
		ctx,
		`INSERT INTO users (id, "apiKey", "isBotAdmin") VALUES (uuidv7(), uuidv7()::text, false) RETURNING id`,
	).Scan(&id)
	if err != nil {
		return uuid.Nil, fmt.Errorf("insert user: %w", err)
	}

	return id, nil
}

func (a *Auth) createChannel(ctx context.Context, userID uuid.UUID, botID string) (channelsmodel.Channel, error) {
	var channel channelsmodel.Channel
	var platformValue platform.Platform
	err := a.pgxPool.QueryRow(
		ctx,
		`INSERT INTO channels (user_id, platform, "botId") VALUES ($1, $2, $3) RETURNING id::text, platform::text, user_id, "isEnabled", "isTwitchBanned", "isBotMod", "botId"`,
		userID,
		platform.PlatformKick,
		botID,
	).Scan(
		&channel.ID,
		&platformValue,
		&channel.UserID,
		&channel.IsEnabled,
		&channel.IsTwitchBanned,
		&channel.IsBotMod,
		&channel.BotID,
	)
	if err != nil {
		return channelsmodel.Nil, fmt.Errorf("insert channel: %w", err)
	}

	channel.Platform = platformValue

	return channel, nil
}
