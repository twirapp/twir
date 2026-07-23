package auth

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/url"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/danielgtaylor/huma/v2/humatest"
	"github.com/google/uuid"
	"github.com/nicklaw5/helix/v2"
	authsessions "github.com/twirapp/twir/apps/api-gql/internal/auth"
	appplatform "github.com/twirapp/twir/apps/api-gql/internal/platform"
	dashboardaccess "github.com/twirapp/twir/apps/api-gql/internal/services/dashboard_access"
	buscoreeventsub "github.com/twirapp/twir/libs/bus-core/eventsub"
	cfg "github.com/twirapp/twir/libs/config"
	"github.com/twirapp/twir/libs/crypto"
	kickbotentity "github.com/twirapp/twir/libs/entities/kick_bot"
	platformentity "github.com/twirapp/twir/libs/entities/platform"
	botsrepo "github.com/twirapp/twir/libs/repositories/bots"
	botsmodel "github.com/twirapp/twir/libs/repositories/bots/model"
	channelplatforms "github.com/twirapp/twir/libs/repositories/channel_platforms"
	channelplatformsmodel "github.com/twirapp/twir/libs/repositories/channel_platforms/model"
	channelsrepo "github.com/twirapp/twir/libs/repositories/channels"
	channelsmodel "github.com/twirapp/twir/libs/repositories/channels/model"
	tokensrepo "github.com/twirapp/twir/libs/repositories/tokens"
	tokensmodel "github.com/twirapp/twir/libs/repositories/tokens/model"
	usersrepo "github.com/twirapp/twir/libs/repositories/users"
	usersmodel "github.com/twirapp/twir/libs/repositories/users/model"
)

func TestUpsertPlatformUserToken_BindsCreatedTokenToUser(t *testing.T) {
	t.Parallel()

	userID := uuid.New()
	createdTokenID := uuid.New()

	tokensRepository := &fakeTokensRepository{
		getByUserIDFunc: func(context.Context, uuid.UUID) (*tokensmodel.Token, error) {
			return nil, tokensrepo.ErrNotFound
		},
		createUserTokenFunc: func(_ context.Context, input tokensrepo.CreateInput) (*tokensmodel.Token, error) {
			if input.UserID != userID {
				t.Fatalf("expected userID %s, got %s", userID, input.UserID)
			}

			return &tokensmodel.Token{ID: createdTokenID}, nil
		},
	}

	usersRepository := &fakeUsersRepository{
		updateFunc: func(_ context.Context, id uuid.UUID, input usersrepo.UpdateInput) (usersmodel.User, error) {
			if id != userID {
				t.Fatalf("expected update userID %s, got %s", userID, id)
			}

			if input.TokenID == nil {
				t.Fatal("expected token binding update")
			}

			if got, want := *input.TokenID, createdTokenID.String(); got != want {
				t.Fatalf("expected tokenID %s, got %s", want, got)
			}

			return usersmodel.User{}, nil
		},
	}

	authHandler := newTestAuth(tokensRepository, usersRepository)

	err := authHandler.upsertPlatformUserToken(context.Background(), userID, testPlatformTokens())
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if tokensRepository.createUserTokenCalls != 1 {
		t.Fatalf("expected CreateUserToken to be called once, got %d", tokensRepository.createUserTokenCalls)
	}

	if usersRepository.updateCalls != 1 {
		t.Fatalf("expected users update to be called once, got %d", usersRepository.updateCalls)
	}

	if tokensRepository.updateTokenCalls != 0 {
		t.Fatalf("expected UpdateTokenByID not to be called, got %d", tokensRepository.updateTokenCalls)
	}
}

func TestUpsertPlatformUserToken_UpdatesExistingTokenWithoutRebinding(t *testing.T) {
	t.Parallel()

	userID := uuid.New()
	existingTokenID := uuid.New()

	tokensRepository := &fakeTokensRepository{
		getByUserIDFunc: func(context.Context, uuid.UUID) (*tokensmodel.Token, error) {
			return &tokensmodel.Token{ID: existingTokenID}, nil
		},
		updateTokenByIDFunc: func(_ context.Context, id uuid.UUID, input tokensrepo.UpdateTokenInput) (*tokensmodel.Token, error) {
			if id != existingTokenID {
				t.Fatalf("expected tokenID %s, got %s", existingTokenID, id)
			}

			if input.AccessToken == nil || input.RefreshToken == nil {
				t.Fatal("expected encrypted access and refresh tokens")
			}

			return &tokensmodel.Token{ID: existingTokenID}, nil
		},
	}

	usersRepository := &fakeUsersRepository{
		updateFunc: func(context.Context, uuid.UUID, usersrepo.UpdateInput) (usersmodel.User, error) {
			t.Fatal("users update should not be called for existing token")
			return usersmodel.User{}, nil
		},
	}

	authHandler := newTestAuth(tokensRepository, usersRepository)

	err := authHandler.upsertPlatformUserToken(context.Background(), userID, testPlatformTokens())
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if tokensRepository.updateTokenCalls != 1 {
		t.Fatalf("expected UpdateTokenByID to be called once, got %d", tokensRepository.updateTokenCalls)
	}

	if tokensRepository.createUserTokenCalls != 0 {
		t.Fatalf("expected CreateUserToken not to be called, got %d", tokensRepository.createUserTokenCalls)
	}

	if usersRepository.updateCalls != 0 {
		t.Fatalf("expected users update not to be called, got %d", usersRepository.updateCalls)
	}
}

func TestUpsertPlatformUserToken_EncryptsDeviceIDForNewToken(t *testing.T) {
	userID := uuid.New()
	createdTokenID := uuid.New()

	tokensRepository := &fakeTokensRepository{
		getByUserIDFunc: func(context.Context, uuid.UUID) (*tokensmodel.Token, error) {
			return nil, tokensrepo.ErrNotFound
		},
		createUserTokenFunc: func(_ context.Context, input tokensrepo.CreateInput) (*tokensmodel.Token, error) {
			if input.DeviceID == nil {
				t.Fatal("expected encrypted device ID")
			}
			deviceID, err := crypto.Decrypt(*input.DeviceID, "pnyfwfiulmnqlhkvixaeligpprcnlyke")
			if err != nil {
				t.Fatalf("decrypt persisted device ID: %v", err)
			}
			if deviceID != "device-id" {
				t.Fatalf("persisted device ID = %q, want device-id", deviceID)
			}

			return &tokensmodel.Token{ID: createdTokenID}, nil
		},
	}

	usersRepository := &fakeUsersRepository{
		updateFunc: func(context.Context, uuid.UUID, usersrepo.UpdateInput) (usersmodel.User, error) {
			return usersmodel.User{}, nil
		},
	}

	tokens := testPlatformTokens()
	tokens.DeviceID = "device-id"
	if err := newTestAuth(tokensRepository, usersRepository).upsertPlatformUserToken(context.Background(), userID, tokens); err != nil {
		t.Fatalf("upsert token: %v", err)
	}
}

func TestUpsertPlatformUserToken_EncryptsDeviceIDForExistingToken(t *testing.T) {
	userID := uuid.New()
	existingTokenID := uuid.New()

	tokensRepository := &fakeTokensRepository{
		getByUserIDFunc: func(context.Context, uuid.UUID) (*tokensmodel.Token, error) {
			return &tokensmodel.Token{ID: existingTokenID}, nil
		},
		updateTokenByIDFunc: func(_ context.Context, id uuid.UUID, input tokensrepo.UpdateTokenInput) (*tokensmodel.Token, error) {
			if id != existingTokenID {
				t.Fatalf("update token ID = %s, want %s", id, existingTokenID)
			}
			if input.DeviceID == nil {
				t.Fatal("expected encrypted device ID")
			}
			deviceID, err := crypto.Decrypt(*input.DeviceID, "pnyfwfiulmnqlhkvixaeligpprcnlyke")
			if err != nil {
				t.Fatalf("decrypt persisted device ID: %v", err)
			}
			if deviceID != "device-id" {
				t.Fatalf("persisted device ID = %q, want device-id", deviceID)
			}

			return &tokensmodel.Token{ID: existingTokenID}, nil
		},
	}

	tokens := testPlatformTokens()
	tokens.DeviceID = "device-id"
	if err := newTestAuth(tokensRepository, &fakeUsersRepository{}).upsertPlatformUserToken(context.Background(), userID, tokens); err != nil {
		t.Fatalf("upsert token: %v", err)
	}
}

func TestUpsertPlatformUserToken_ReturnsErrorWhenBindingFails(t *testing.T) {
	t.Parallel()

	userID := uuid.New()
	createdTokenID := uuid.New()
	bindingErr := errors.New("users update failed")

	tokensRepository := &fakeTokensRepository{
		getByUserIDFunc: func(context.Context, uuid.UUID) (*tokensmodel.Token, error) {
			return nil, tokensrepo.ErrNotFound
		},
		createUserTokenFunc: func(context.Context, tokensrepo.CreateInput) (*tokensmodel.Token, error) {
			return &tokensmodel.Token{ID: createdTokenID}, nil
		},
	}

	usersRepository := &fakeUsersRepository{
		updateFunc: func(context.Context, uuid.UUID, usersrepo.UpdateInput) (usersmodel.User, error) {
			return usersmodel.User{}, bindingErr
		},
	}

	authHandler := newTestAuth(tokensRepository, usersRepository)

	err := authHandler.upsertPlatformUserToken(context.Background(), userID, testPlatformTokens())
	if !errors.Is(err, bindingErr) {
		t.Fatalf("expected binding error %v, got %v", bindingErr, err)
	}

	if usersRepository.updateCalls != 1 {
		t.Fatalf("expected users update to be called once, got %d", usersRepository.updateCalls)
	}
}

func TestHandleAuthPostCodeRejectsMissingOAuthAttemptWithoutExchange(t *testing.T) {
	exchangeCalls := 0
	authHandler := newOAuthFlowTestAuth(oauthFlowTestAuthOpts{
		sessions: &fakeOAuthSession{},
		registry: appplatform.NewRegistry([]appplatform.PlatformProvider{
			&oauthPlatformProvider{
				name: platformentity.PlatformTwitch.String(),
				exchangeCodeFunc: func(context.Context, appplatform.ExchangeCodeInput) (*appplatform.PlatformTokens, error) {
					exchangeCalls++
					return nil, errors.New("provider exchange should not be called")
				},
			},
		}),
	})

	_, err := authHandler.handleAuthPostCode(context.Background(), authBody{
		Code:  "authorization-code",
		State: "L2Rhc2hib2FyZA==", // Base64 for /dashboard, accepted by the legacy fallback.
	})
	if exchangeCalls != 0 {
		t.Fatalf("provider exchanges = %d, want 0", exchangeCalls)
	}
	if err == nil {
		t.Fatal("missing OAuth attempt should be rejected")
	}
	if !strings.Contains(err.Error(), "Invalid or expired OAuth state") {
		t.Fatalf("missing OAuth attempt error = %v", err)
	}
}

func TestCompletePlatformAuthCreatesMissingTwitchBindingWithTwitchConfigWhenLinkingVK(t *testing.T) {
	ctx := context.Background()
	sessionUserID := uuid.New()
	vkUserID := uuid.New()
	channelID := uuid.New()
	linkedProviderConfig := json.RawMessage(`{"linked_provider":"vk"}`)
	createdBindings := make([]channelplatforms.CreateInput, 0, 2)
	transaction := &oauthTransaction{}
	bindings := &oauthChannelPlatformsRepository{
		getByChannelAndPlatformFunc: func(context.Context, uuid.UUID, platformentity.Platform) (channelplatformsmodel.ChannelPlatform, error) {
			return channelplatformsmodel.Nil, channelplatforms.ErrNotFound
		},
		createFunc: func(_ context.Context, input channelplatforms.CreateInput) (channelplatformsmodel.ChannelPlatform, error) {
			createdBindings = append(createdBindings, input)
			return channelplatformsmodel.ChannelPlatform{
				ID:                uuid.New(),
				ChannelID:         input.ChannelID,
				Platform:          input.Platform,
				UserID:            input.UserID,
				PlatformChannelID: input.PlatformChannelID,
				Enabled:           input.Enabled,
				BotConfig:         input.BotConfig,
			}, nil
		},
	}
	users := &oauthUsersRepository{
		getByIDFunc: func(_ context.Context, id uuid.UUID) (usersmodel.User, error) {
			if id != sessionUserID {
				t.Fatalf("session user ID = %s, want %s", id, sessionUserID)
			}
			return usersmodel.User{ID: sessionUserID, Platform: platformentity.PlatformTwitch, PlatformID: "twitch-channel"}, nil
		},
		getByPlatformIDFunc: func(_ context.Context, platform platformentity.Platform, platformID string) (usersmodel.User, error) {
			if platform != platformentity.PlatformVKVideoLive || platformID != "vk-channel" {
				t.Fatalf("platform user lookup = (%s, %q)", platform, platformID)
			}
			return usersmodel.User{ID: vkUserID, Platform: platform, PlatformID: platformID}, nil
		},
	}
	channels := &oauthChannelsRepository{
		createFunc: func(_ context.Context, input channelsrepo.CreateInput) (channelsmodel.Channel, error) {
			if input.BotID != "default-bot" {
				t.Fatalf("channel BotID = %q, want default-bot", input.BotID)
			}
			return channelsmodel.Channel{ID: channelID}, nil
		},
		getByBindingUserIDFunc: func(_ context.Context, platform platformentity.Platform, userID uuid.UUID) (channelsmodel.Channel, error) {
			switch {
			case platform == platformentity.PlatformTwitch && userID == sessionUserID:
				return channelsmodel.Nil, channelsrepo.ErrNotFound
			case platform == platformentity.PlatformVKVideoLive && userID == vkUserID:
				return channelsmodel.Nil, channelsrepo.ErrNotFound
			default:
				t.Fatalf("unexpected channel lookup (%s, %s)", platform, userID)
				return channelsmodel.Nil, nil
			}
		},
	}
	authHandler := newOAuthFlowTestAuth(oauthFlowTestAuthOpts{
		sessions:    &fakeOAuthSession{internalUserID: sessionUserID},
		users:       users,
		channels:    channels,
		bindings:    bindings,
		tokens:      newCreateTokenRepository(vkUserID),
		bots:        &oauthBotsRepository{defaultBot: botsmodel.Bot{ID: "default-bot"}},
		transaction: transaction,
	})
	authHandler.bindingConfigResolvers = map[platformentity.Platform]platformBindingConfigResolver{
		platformentity.PlatformTwitch: authHandler.twitchBindingConfig,
	}

	_, err := authHandler.completePlatformAuth(ctx, completePlatformAuthInput{
		Platform:      platformentity.PlatformVKVideoLive,
		PlatformUser:  &appplatform.PlatformUser{ID: "vk-channel"},
		Tokens:        testPlatformTokens(),
		BindingConfig: platformBindingConfig{BotConfig: linkedProviderConfig},
	})
	if err != nil {
		t.Fatalf("link VK auth: %v", err)
	}
	if len(createdBindings) != 2 {
		t.Fatalf("created bindings = %d, want 2", len(createdBindings))
	}

	var twitchConfig struct {
		BotID          string `json:"bot_id"`
		IsBotMod       bool   `json:"is_bot_mod"`
		IsTwitchBanned bool   `json:"is_twitch_banned"`
	}
	if err := json.Unmarshal(createdBindings[0].BotConfig, &twitchConfig); err != nil {
		t.Fatalf("decode Twitch binding config: %v", err)
	}
	if createdBindings[0].Platform != platformentity.PlatformTwitch || twitchConfig != (struct {
		BotID          string `json:"bot_id"`
		IsBotMod       bool   `json:"is_bot_mod"`
		IsTwitchBanned bool   `json:"is_twitch_banned"`
	}{BotID: "default-bot"}) {
		t.Fatalf("created Twitch binding = %+v, config = %+v", createdBindings[0], twitchConfig)
	}
	if createdBindings[1].Platform != platformentity.PlatformVKVideoLive || string(createdBindings[1].BotConfig) != string(linkedProviderConfig) {
		t.Fatalf("created VK binding = %+v, want incoming config %s", createdBindings[1], linkedProviderConfig)
	}
}

func TestCompletePlatformAuthCreatesVKOnlyChannel(t *testing.T) {
	ctx := context.Background()
	platformUserID := uuid.New()
	channelID := uuid.New()
	transaction := &oauthTransaction{}
	publisher := &oauthEventSubPublisher{transaction: transaction}
	bindings := &oauthChannelPlatformsRepository{
		getByChannelAndPlatformFunc: func(context.Context, uuid.UUID, platformentity.Platform) (channelplatformsmodel.ChannelPlatform, error) {
			return channelplatformsmodel.Nil, channelplatforms.ErrNotFound
		},
	}

	users := &oauthUsersRepository{
		getByPlatformIDFunc: func(context.Context, platformentity.Platform, string) (usersmodel.User, error) {
			return usersmodel.Nil, usersmodel.ErrNotFound
		},
		createFunc: func(_ context.Context, input usersrepo.CreateInput) (usersmodel.User, error) {
			if input.Platform != platformentity.PlatformVKVideoLive || input.PlatformID != "vk-channel" {
				t.Fatalf("created user = %+v, want VK channel identity", input)
			}

			return usersmodel.User{ID: platformUserID, Platform: input.Platform, PlatformID: input.PlatformID}, nil
		},
	}
	channels := &oauthChannelsRepository{
		createFunc: func(_ context.Context, input channelsrepo.CreateInput) (channelsmodel.Channel, error) {
			if input.BotID != "default-bot" {
				t.Fatalf("channel BotID = %q, want default-bot", input.BotID)
			}

			return channelsmodel.Channel{ID: channelID}, nil
		},
		getByBindingUserIDFunc: func(context.Context, platformentity.Platform, uuid.UUID) (channelsmodel.Channel, error) {
			return channelsmodel.Nil, channelsrepo.ErrNotFound
		},
	}
	bindings.createFunc = func(_ context.Context, input channelplatforms.CreateInput) (channelplatformsmodel.ChannelPlatform, error) {
		if input.ChannelID != channelID || input.Platform != platformentity.PlatformVKVideoLive || input.UserID != platformUserID || input.PlatformChannelID != "vk-channel" || !input.Enabled {
			t.Fatalf("created VK binding = %+v", input)
		}
		if string(input.BotConfig) != "{}" {
			t.Fatalf("VK binding bot config = %s, want {}", input.BotConfig)
		}

		return channelplatformsmodel.ChannelPlatform{
			ID:                uuid.New(),
			ChannelID:         input.ChannelID,
			Platform:          input.Platform,
			UserID:            input.UserID,
			PlatformChannelID: input.PlatformChannelID,
			Enabled:           input.Enabled,
			BotConfig:         input.BotConfig,
		}, nil
	}

	result, err := newOAuthFlowTestAuth(oauthFlowTestAuthOpts{
		sessions:          &fakeOAuthSession{internalUserErr: errors.New("not signed in")},
		users:             users,
		channels:          channels,
		bindings:          bindings,
		tokens:            newCreateTokenRepository(platformUserID),
		bots:              &oauthBotsRepository{defaultBot: botsmodel.Bot{ID: "default-bot"}},
		transaction:       transaction,
		eventSubPublisher: publisher,
	}).completePlatformAuth(ctx, completePlatformAuthInput{
		Platform: platformentity.PlatformVKVideoLive,
		PlatformUser: &appplatform.PlatformUser{
			ID:          "vk-channel",
			DisplayName: "VK Streamer",
		},
		Tokens: testPlatformTokens(),
	})
	if err != nil {
		t.Fatalf("complete VK auth: %v", err)
	}

	if !result.CreatedChannel || result.Channel.ID != channelID {
		t.Fatalf("auth result = %+v, want newly created channel %s", result, channelID)
	}
	if channels.createCalls != 1 || bindings.createCalls != 1 {
		t.Fatalf("channel creates = %d, binding creates = %d, want one each", channels.createCalls, bindings.createCalls)
	}
	if len(publisher.requests) != 1 || publisher.requests[0] != (buscoreeventsub.EventsubSubscribeToAllEventsRequest{
		ChannelID: channelID.String(),
		Platform:  platformentity.PlatformVKVideoLive,
	}) {
		t.Fatalf("EventSub requests = %#v", publisher.requests)
	}
}

func TestCompletePlatformAuthLinksVKToExistingTwitchKickChannel(t *testing.T) {
	ctx := context.Background()
	sessionUserID := uuid.New()
	vkUserID := uuid.New()
	channelID := uuid.New()
	existingBindings := []channelplatformsmodel.ChannelPlatform{
		{ID: uuid.New(), ChannelID: channelID, Platform: platformentity.PlatformTwitch, UserID: sessionUserID, PlatformChannelID: "twitch-channel", Enabled: true},
		{ID: uuid.New(), ChannelID: channelID, Platform: platformentity.PlatformKick, UserID: uuid.New(), PlatformChannelID: "kick-channel", Enabled: true},
	}
	channel := channelsmodel.Channel{ID: channelID, Bindings: append([]channelplatformsmodel.ChannelPlatform(nil), existingBindings...)}
	transaction := &oauthTransaction{}
	bindings := &oauthChannelPlatformsRepository{
		getByChannelAndPlatformFunc: func(_ context.Context, gotChannelID uuid.UUID, gotPlatform platformentity.Platform) (channelplatformsmodel.ChannelPlatform, error) {
			if gotChannelID != channelID || gotPlatform != platformentity.PlatformVKVideoLive {
				t.Fatalf("binding lookup = (%s, %s)", gotChannelID, gotPlatform)
			}

			return channelplatformsmodel.Nil, channelplatforms.ErrNotFound
		},
	}
	bindings.createFunc = func(_ context.Context, input channelplatforms.CreateInput) (channelplatformsmodel.ChannelPlatform, error) {
		if input.ChannelID != channelID || input.UserID != vkUserID || input.PlatformChannelID != "vk-channel" {
			t.Fatalf("new VK binding = %+v", input)
		}

		return channelplatformsmodel.ChannelPlatform{ID: uuid.New(), ChannelID: channelID, Platform: input.Platform, UserID: input.UserID, PlatformChannelID: input.PlatformChannelID, Enabled: true, BotConfig: input.BotConfig}, nil
	}
	channels := &oauthChannelsRepository{
		getByBindingUserIDFunc: func(_ context.Context, gotPlatform platformentity.Platform, userID uuid.UUID) (channelsmodel.Channel, error) {
			switch {
			case gotPlatform == platformentity.PlatformTwitch && userID == sessionUserID:
				return channel, nil
			case gotPlatform == platformentity.PlatformVKVideoLive && userID == vkUserID:
				return channelsmodel.Nil, channelsrepo.ErrNotFound
			default:
				t.Fatalf("unexpected channel lookup (%s, %s)", gotPlatform, userID)
				return channelsmodel.Nil, nil
			}
		},
	}
	users := &oauthUsersRepository{
		getByIDFunc: func(_ context.Context, id uuid.UUID) (usersmodel.User, error) {
			if id != sessionUserID {
				t.Fatalf("session user ID = %s, want %s", id, sessionUserID)
			}
			return usersmodel.User{ID: sessionUserID, Platform: platformentity.PlatformTwitch, PlatformID: "twitch-channel"}, nil
		},
		getByPlatformIDFunc: func(_ context.Context, gotPlatform platformentity.Platform, platformID string) (usersmodel.User, error) {
			if gotPlatform != platformentity.PlatformVKVideoLive || platformID != "vk-channel" {
				t.Fatalf("platform user lookup = (%s, %q)", gotPlatform, platformID)
			}
			return usersmodel.User{ID: vkUserID, Platform: gotPlatform, PlatformID: platformID}, nil
		},
	}

	_, err := newOAuthFlowTestAuth(oauthFlowTestAuthOpts{
		sessions:          &fakeOAuthSession{internalUserID: sessionUserID},
		users:             users,
		channels:          channels,
		bindings:          bindings,
		tokens:            newCreateTokenRepository(vkUserID),
		bots:              &oauthBotsRepository{defaultBot: botsmodel.Bot{ID: "default-bot"}},
		transaction:       transaction,
		eventSubPublisher: &oauthEventSubPublisher{transaction: transaction},
	}).completePlatformAuth(ctx, completePlatformAuthInput{
		Platform:     platformentity.PlatformVKVideoLive,
		PlatformUser: &appplatform.PlatformUser{ID: "vk-channel"},
		Tokens:       testPlatformTokens(),
	})
	if err != nil {
		t.Fatalf("link VK auth: %v", err)
	}

	if !reflect.DeepEqual(channel.Bindings, existingBindings) {
		t.Fatalf("existing bindings changed from %#v to %#v", existingBindings, channel.Bindings)
	}
	if channels.createCalls != 0 || bindings.createCalls != 1 {
		t.Fatalf("channel creates = %d, binding creates = %d, want 0 and 1", channels.createCalls, bindings.createCalls)
	}
}

func TestCompletePlatformAuthRejectsVKAccountBoundToAnotherChannel(t *testing.T) {
	ctx := context.Background()
	sessionUserID := uuid.New()
	vkUserID := uuid.New()
	channelID := uuid.New()
	otherChannelID := uuid.New()
	transaction := &oauthTransaction{}
	bindings := &oauthChannelPlatformsRepository{
		getByChannelAndPlatformFunc: func(context.Context, uuid.UUID, platformentity.Platform) (channelplatformsmodel.ChannelPlatform, error) {
			return channelplatformsmodel.Nil, channelplatforms.ErrNotFound
		},
	}
	channels := &oauthChannelsRepository{
		getByBindingUserIDFunc: func(_ context.Context, gotPlatform platformentity.Platform, userID uuid.UUID) (channelsmodel.Channel, error) {
			switch {
			case gotPlatform == platformentity.PlatformTwitch && userID == sessionUserID:
				return channelsmodel.Channel{ID: channelID}, nil
			case gotPlatform == platformentity.PlatformVKVideoLive && userID == vkUserID:
				return channelsmodel.Channel{ID: otherChannelID}, nil
			default:
				t.Fatalf("unexpected channel lookup (%s, %s)", gotPlatform, userID)
				return channelsmodel.Nil, nil
			}
		},
	}
	users := &oauthUsersRepository{
		getByIDFunc: func(context.Context, uuid.UUID) (usersmodel.User, error) {
			return usersmodel.User{ID: sessionUserID, Platform: platformentity.PlatformTwitch}, nil
		},
		getByPlatformIDFunc: func(context.Context, platformentity.Platform, string) (usersmodel.User, error) {
			return usersmodel.User{ID: vkUserID, Platform: platformentity.PlatformVKVideoLive, PlatformID: "vk-channel"}, nil
		},
	}
	publisher := &oauthEventSubPublisher{transaction: transaction}

	_, err := newOAuthFlowTestAuth(oauthFlowTestAuthOpts{
		sessions:          &fakeOAuthSession{internalUserID: sessionUserID},
		users:             users,
		channels:          channels,
		bindings:          bindings,
		tokens:            newCreateTokenRepository(vkUserID),
		bots:              &oauthBotsRepository{defaultBot: botsmodel.Bot{ID: "default-bot"}},
		transaction:       transaction,
		eventSubPublisher: publisher,
	}).completePlatformAuth(ctx, completePlatformAuthInput{
		Platform:     platformentity.PlatformVKVideoLive,
		PlatformUser: &appplatform.PlatformUser{ID: "vk-channel"},
		Tokens:       testPlatformTokens(),
	})
	if !errors.Is(err, errPlatformConflict) {
		t.Fatalf("link conflict error = %v, want errPlatformConflict", err)
	}
	if bindings.createCalls != 0 {
		t.Fatalf("created %d bindings after conflict", bindings.createCalls)
	}
	if len(publisher.requests) != 0 {
		t.Fatalf("published EventSub requests after conflict: %#v", publisher.requests)
	}
}

func TestCompletePlatformAuthPreservesExistingBindingForSameAccount(t *testing.T) {
	ctx := context.Background()
	sessionUserID := uuid.New()
	vkUserID := uuid.New()
	channelID := uuid.New()
	existingBinding := channelplatformsmodel.ChannelPlatform{
		ID:                uuid.New(),
		ChannelID:         channelID,
		Platform:          platformentity.PlatformVKVideoLive,
		UserID:            vkUserID,
		PlatformChannelID: "vk-channel",
		Enabled:           true,
		BotConfig:         []byte(`{"preserved":true}`),
	}
	transaction := &oauthTransaction{}
	bindings := &oauthChannelPlatformsRepository{
		getByChannelAndPlatformFunc: func(context.Context, uuid.UUID, platformentity.Platform) (channelplatformsmodel.ChannelPlatform, error) {
			return existingBinding, nil
		},
	}
	channels := &oauthChannelsRepository{
		getByBindingUserIDFunc: func(_ context.Context, gotPlatform platformentity.Platform, userID uuid.UUID) (channelsmodel.Channel, error) {
			if gotPlatform == platformentity.PlatformTwitch && userID == sessionUserID {
				return channelsmodel.Channel{ID: channelID, Bindings: []channelplatformsmodel.ChannelPlatform{existingBinding}}, nil
			}
			t.Fatalf("unexpected channel lookup (%s, %s)", gotPlatform, userID)
			return channelsmodel.Nil, nil
		},
	}
	users := &oauthUsersRepository{
		getByIDFunc: func(context.Context, uuid.UUID) (usersmodel.User, error) {
			return usersmodel.User{ID: sessionUserID, Platform: platformentity.PlatformTwitch}, nil
		},
		getByPlatformIDFunc: func(context.Context, platformentity.Platform, string) (usersmodel.User, error) {
			return usersmodel.User{ID: vkUserID, Platform: platformentity.PlatformVKVideoLive, PlatformID: "vk-channel"}, nil
		},
	}

	_, err := newOAuthFlowTestAuth(oauthFlowTestAuthOpts{
		sessions:          &fakeOAuthSession{internalUserID: sessionUserID},
		users:             users,
		channels:          channels,
		bindings:          bindings,
		tokens:            newCreateTokenRepository(vkUserID),
		bots:              &oauthBotsRepository{defaultBot: botsmodel.Bot{ID: "default-bot"}},
		transaction:       transaction,
		eventSubPublisher: &oauthEventSubPublisher{transaction: transaction},
	}).completePlatformAuth(ctx, completePlatformAuthInput{
		Platform:     platformentity.PlatformVKVideoLive,
		PlatformUser: &appplatform.PlatformUser{ID: "vk-channel"},
		Tokens:       testPlatformTokens(),
	})
	if err != nil {
		t.Fatalf("re-authenticate existing VK binding: %v", err)
	}
	if bindings.createCalls != 0 || bindings.updateCalls != 0 {
		t.Fatalf("binding mutations = create %d, update %d, want none", bindings.createCalls, bindings.updateCalls)
	}
}

func TestPlatformCodeCarriesVKCallbackDeviceIDThroughExchangeAndPersistence(t *testing.T) {
	ctx := context.Background()
	platformUserID := uuid.New()
	channelID := uuid.New()
	state := "verified-vk-attempt"
	deviceID := "vk-device-id"
	sessions := &fakeOAuthSession{internalUserErr: errors.New("not signed in"), attempts: map[string]authsessions.OAuthAttempt{
		state: {
			Platform:     platformentity.PlatformVKVideoLive,
			RedirectTo:   "/dashboard",
			CodeVerifier: "stored-pkce-verifier",
		},
	}}
	transaction := &oauthTransaction{}
	provider := &oauthPlatformProvider{
		name: platformentity.PlatformVKVideoLive.String(),
		exchangeCodeFunc: func(_ context.Context, input appplatform.ExchangeCodeInput) (*appplatform.PlatformTokens, error) {
			attempt, ok := sessions.attempts[state]
			if !ok || attempt.DeviceID != deviceID {
				t.Fatalf("callback device ID was not stored with OAuth attempt before exchange: %+v", attempt)
			}
			if input.CodeVerifier != "stored-pkce-verifier" || input.DeviceID != deviceID {
				t.Fatalf("exchange input = %+v, want stored verifier and callback device ID", input)
			}

			return &appplatform.PlatformTokens{AccessToken: "vk-access", RefreshToken: "vk-refresh", ExpiresIn: 3600, DeviceID: input.DeviceID}, nil
		},
		getUserFunc: func(context.Context, string) (*appplatform.PlatformUser, error) {
			return &appplatform.PlatformUser{ID: "vk-channel", DisplayName: "VK Streamer"}, nil
		},
	}
	tokens := &fakeTokensRepository{
		getByUserIDFunc: func(context.Context, uuid.UUID) (*tokensmodel.Token, error) {
			return nil, tokensrepo.ErrNotFound
		},
		createUserTokenFunc: func(_ context.Context, input tokensrepo.CreateInput) (*tokensmodel.Token, error) {
			if input.DeviceID == nil {
				t.Fatal("VK device ID was not persisted")
			}
			decryptedDeviceID, err := crypto.Decrypt(*input.DeviceID, "pnyfwfiulmnqlhkvixaeligpprcnlyke")
			if err != nil {
				t.Fatalf("decrypt persisted device ID: %v", err)
			}
			if decryptedDeviceID != deviceID {
				t.Fatalf("persisted device ID = %q, want %q", decryptedDeviceID, deviceID)
			}

			return &tokensmodel.Token{ID: uuid.New()}, nil
		},
	}
	users := &oauthUsersRepository{
		getByPlatformIDFunc: func(context.Context, platformentity.Platform, string) (usersmodel.User, error) {
			return usersmodel.Nil, usersmodel.ErrNotFound
		},
		createFunc: func(context.Context, usersrepo.CreateInput) (usersmodel.User, error) {
			return usersmodel.User{ID: platformUserID, Platform: platformentity.PlatformVKVideoLive, PlatformID: "vk-channel"}, nil
		},
	}
	bindings := &oauthChannelPlatformsRepository{
		getByChannelAndPlatformFunc: func(context.Context, uuid.UUID, platformentity.Platform) (channelplatformsmodel.ChannelPlatform, error) {
			return channelplatformsmodel.Nil, channelplatforms.ErrNotFound
		},
		createFunc: func(_ context.Context, input channelplatforms.CreateInput) (channelplatformsmodel.ChannelPlatform, error) {
			return channelplatformsmodel.ChannelPlatform{ID: uuid.New(), ChannelID: input.ChannelID, Platform: input.Platform, UserID: input.UserID, PlatformChannelID: input.PlatformChannelID, Enabled: input.Enabled, BotConfig: input.BotConfig}, nil
		},
	}

	_, err := newOAuthFlowTestAuth(oauthFlowTestAuthOpts{
		sessions: sessions,
		users:    users,
		channels: &oauthChannelsRepository{
			createFunc: func(context.Context, channelsrepo.CreateInput) (channelsmodel.Channel, error) {
				return channelsmodel.Channel{ID: channelID}, nil
			},
			getByBindingUserIDFunc: func(context.Context, platformentity.Platform, uuid.UUID) (channelsmodel.Channel, error) {
				return channelsmodel.Nil, channelsrepo.ErrNotFound
			},
		},
		bindings:          bindings,
		tokens:            tokens,
		bots:              &oauthBotsRepository{defaultBot: botsmodel.Bot{ID: "default-bot"}},
		registry:          appplatform.NewRegistry([]appplatform.PlatformProvider{provider}),
		transaction:       transaction,
		eventSubPublisher: &oauthEventSubPublisher{transaction: transaction},
	}).completePlatformCode(ctx, platformCodeInput{
		Platform: platformentity.PlatformVKVideoLive,
		Code:     "authorization-code",
		State:    state,
		DeviceID: deviceID,
	})
	if err != nil {
		t.Fatalf("complete VK platform code: %v", err)
	}
	if _, ok := sessions.attempts[state]; ok {
		t.Fatal("successful OAuth attempt was not removed from server-side session state")
	}
}

func TestStartPlatformAuthStoresOpaqueStateAndPKCEOnServer(t *testing.T) {
	sessions := &fakeOAuthSession{}
	provider := &oauthPlatformProvider{
		name: platformentity.PlatformVKVideoLive.String(),
		getAuthURLFunc: func(state, codeChallenge string) string {
			return "https://id.example.test/authorize?state=" + url.QueryEscape(state) + "&challenge=" + url.QueryEscape(codeChallenge)
		},
	}
	authHandler := newOAuthFlowTestAuth(oauthFlowTestAuthOpts{
		sessions: sessions,
		registry: appplatform.NewRegistry([]appplatform.PlatformProvider{provider}),
	})

	authorizeURL, err := authHandler.startPlatformAuth(context.Background(), platformentity.PlatformVKVideoLive, "/dashboard/settings")
	if err != nil {
		t.Fatalf("start VK OAuth: %v", err)
	}
	if strings.Contains(authorizeURL, "device_id") || strings.Contains(authorizeURL, "/dashboard/settings") {
		t.Fatalf("authorization URL leaked server-side data: %s", authorizeURL)
	}
	parsedURL, err := url.Parse(authorizeURL)
	if err != nil {
		t.Fatalf("parse authorization URL: %v", err)
	}
	state := parsedURL.Query().Get("state")
	attempt, ok := sessions.attempts[state]
	if !ok {
		t.Fatal("OAuth state was not stored server-side")
	}
	if attempt.Platform != platformentity.PlatformVKVideoLive || attempt.RedirectTo != "/dashboard/settings" || attempt.CodeVerifier == "" {
		t.Fatalf("stored OAuth attempt = %+v", attempt)
	}
}

func TestStartTwitchAuthUsesServerSideOAuthAttempt(t *testing.T) {
	sessions := &fakeOAuthSession{}
	provider := &oauthPlatformProvider{
		name: platformentity.PlatformTwitch.String(),
		getAuthURLFunc: func(state, _ string) string {
			return "https://id.example.test/authorize?state=" + url.QueryEscape(state)
		},
	}
	authHandler := newOAuthFlowTestAuth(oauthFlowTestAuthOpts{
		sessions: sessions,
		registry: appplatform.NewRegistry([]appplatform.PlatformProvider{provider}),
	})

	authorizeURL, err := authHandler.StartTwitchAuth(context.Background(), "/dashboard")
	if err != nil {
		t.Fatalf("start Twitch OAuth: %v", err)
	}
	parsedURL, err := url.Parse(authorizeURL)
	if err != nil {
		t.Fatalf("parse authorization URL: %v", err)
	}
	attempt, ok := sessions.attempts[parsedURL.Query().Get("state")]
	if !ok {
		t.Fatal("Twitch OAuth state was not stored server-side")
	}
	if attempt.Platform != platformentity.PlatformTwitch || attempt.RedirectTo != "/dashboard" || attempt.CodeVerifier == "" || attempt.TargetChannelID != nil || attempt.InitiatorUserID != nil || !attempt.ExpiresAt.IsZero() {
		t.Fatalf("stored Twitch OAuth attempt = %+v", attempt)
	}
}

func TestStartPlatformAuthUsesGenericRegisteredProvider(t *testing.T) {
	sessions := &fakeOAuthSession{}
	provider := &oauthPlatformProvider{
		name: platformentity.PlatformVKVideoLive.String(),
		getAuthURLFunc: func(state, _ string) string {
			return "https://id.example.test/authorize?state=" + url.QueryEscape(state)
		},
	}
	authHandler := newOAuthFlowTestAuth(oauthFlowTestAuthOpts{
		sessions: sessions,
		registry: appplatform.NewRegistry([]appplatform.PlatformProvider{provider}),
	})

	authorizeURL, err := authHandler.StartPlatformAuth(context.Background(), platformentity.PlatformVKVideoLive, "/dashboard")
	if err != nil {
		t.Fatalf("start generic VK OAuth: %v", err)
	}
	parsedURL, err := url.Parse(authorizeURL)
	if err != nil {
		t.Fatalf("parse authorization URL: %v", err)
	}
	attempt, ok := sessions.attempts[parsedURL.Query().Get("state")]
	if !ok {
		t.Fatal("generic OAuth state was not stored server-side")
	}
	if attempt.Platform != platformentity.PlatformVKVideoLive || attempt.RedirectTo != "/dashboard" || attempt.CodeVerifier == "" || attempt.TargetChannelID != nil || attempt.InitiatorUserID != nil || !attempt.ExpiresAt.IsZero() {
		t.Fatalf("stored generic OAuth attempt = %+v", attempt)
	}
}

func TestStartPlatformAuthForChannelLinksProviderToAuthorizedSelectedDashboard(t *testing.T) {
	ctx := context.Background()
	collaboratorID := uuid.New()
	selectedDashboardID := uuid.New()
	ownChannelID := uuid.New()
	providerUserID := uuid.New()
	createdBindings := make([]channelplatforms.CreateInput, 0, 1)
	transaction := &oauthTransaction{}
	sessions := &fakeOAuthSession{internalUserID: collaboratorID}
	provider := &oauthPlatformProvider{
		name: platformentity.PlatformKick.String(),
		getAuthURLFunc: func(state, _ string) string {
			return "https://id.example.test/authorize?state=" + url.QueryEscape(state)
		},
		exchangeCodeFunc: func(context.Context, appplatform.ExchangeCodeInput) (*appplatform.PlatformTokens, error) {
			return testPlatformTokens(), nil
		},
		getUserFunc: func(context.Context, string) (*appplatform.PlatformUser, error) {
			return &appplatform.PlatformUser{ID: "provider-kick"}, nil
		},
	}
	channels := &oauthChannelsRepository{
		getByIDFunc: func(_ context.Context, channelID uuid.UUID) (channelsmodel.Channel, error) {
			if channelID != selectedDashboardID {
				t.Fatalf("target channel lookup = %s, want %s", channelID, selectedDashboardID)
			}
			if transaction.calls == 0 {
				t.Fatal("target channel was not loaded inside the auth transaction")
			}

			return channelsmodel.Channel{ID: selectedDashboardID}, nil
		},
		getByBindingUserIDFunc: func(_ context.Context, platform platformentity.Platform, userID uuid.UUID) (channelsmodel.Channel, error) {
			switch {
			case platform == platformentity.PlatformTwitch && userID == collaboratorID:
				return channelsmodel.Channel{ID: ownChannelID}, nil
			case platform == platformentity.PlatformKick && userID == providerUserID:
				return channelsmodel.Nil, channelsrepo.ErrNotFound
			default:
				t.Fatalf("unexpected binding lookup (%s, %s)", platform, userID)
				return channelsmodel.Nil, nil
			}
		},
	}
	bindings := &oauthChannelPlatformsRepository{
		getByChannelAndPlatformFunc: func(_ context.Context, channelID uuid.UUID, platform platformentity.Platform) (channelplatformsmodel.ChannelPlatform, error) {
			if channelID != selectedDashboardID || platform != platformentity.PlatformKick {
				t.Fatalf("target binding lookup = (%s, %s)", channelID, platform)
			}

			return channelplatformsmodel.Nil, channelplatforms.ErrNotFound
		},
		createFunc: func(_ context.Context, input channelplatforms.CreateInput) (channelplatformsmodel.ChannelPlatform, error) {
			if transaction.calls == 0 || transaction.committed {
				t.Fatal("provider binding was not created transactionally")
			}
			createdBindings = append(createdBindings, input)
			return channelplatformsmodel.ChannelPlatform{ID: uuid.New(), ChannelID: input.ChannelID, Platform: input.Platform, UserID: input.UserID, PlatformChannelID: input.PlatformChannelID, Enabled: input.Enabled, BotConfig: input.BotConfig}, nil
		},
	}
	users := &oauthUsersRepository{
		getByIDFunc: func(_ context.Context, userID uuid.UUID) (usersmodel.User, error) {
			if userID != collaboratorID {
				t.Fatalf("session user lookup = %s, want %s", userID, collaboratorID)
			}
			return usersmodel.User{ID: collaboratorID, Platform: platformentity.PlatformTwitch, PlatformID: "collaborator-twitch"}, nil
		},
		getByPlatformIDFunc: func(_ context.Context, platform platformentity.Platform, platformID string) (usersmodel.User, error) {
			if platform != platformentity.PlatformKick || platformID != "provider-kick" {
				t.Fatalf("provider user lookup = (%s, %q)", platform, platformID)
			}
			return usersmodel.User{ID: providerUserID, Platform: platform, PlatformID: platformID}, nil
		},
	}
	authHandler := newOAuthFlowTestAuth(oauthFlowTestAuthOpts{
		sessions:    sessions,
		users:       users,
		channels:    channels,
		bindings:    bindings,
		tokens:      newCreateTokenRepository(providerUserID),
		registry:    appplatform.NewRegistry([]appplatform.PlatformProvider{provider}),
		transaction: transaction,
	})
	authHandler.dashboardAccess = dashboardAccessFunc(func(_ context.Context, subject dashboardaccess.Subject, channelID uuid.UUID, permission string) (bool, error) {
		if subject.ID != collaboratorID.String() || channelID != selectedDashboardID || permission != "" {
			t.Fatalf("dashboard access request = %+v, %s, %q", subject, channelID, permission)
		}

		return true, nil
	})

	authorizeURL, err := authHandler.StartPlatformAuthForChannel(ctx, selectedDashboardID, platformentity.PlatformKick, "/dashboard/platforms")
	if err != nil {
		t.Fatalf("start target OAuth: %v", err)
	}
	if strings.Contains(authorizeURL, selectedDashboardID.String()) || strings.Contains(authorizeURL, collaboratorID.String()) {
		t.Fatalf("authorization URL leaked selected dashboard: %s", authorizeURL)
	}
	parsedURL, err := url.Parse(authorizeURL)
	if err != nil {
		t.Fatalf("parse authorization URL: %v", err)
	}
	state := parsedURL.Query().Get("state")
	attempt, ok := sessions.attempts[state]
	if !ok || attempt.TargetChannelID == nil || *attempt.TargetChannelID != selectedDashboardID || attempt.InitiatorUserID == nil || *attempt.InitiatorUserID != collaboratorID || attempt.ExpiresAt.IsZero() {
		t.Fatalf("stored OAuth attempt = %+v, want selected dashboard %s", attempt, selectedDashboardID)
	}

	_, err = authHandler.completePlatformCode(ctx, platformCodeInput{
		Platform: platformentity.PlatformKick,
		Code:     "authorization-code",
		State:    state,
	})
	if err != nil {
		t.Fatalf("complete selected-dashboard OAuth: %v", err)
	}
	if len(createdBindings) != 1 || createdBindings[0].ChannelID != selectedDashboardID {
		t.Fatalf("created bindings = %#v, want provider linked to selected dashboard %s", createdBindings, selectedDashboardID)
	}
}

func TestCompletePlatformCodeRejectsInvalidTargetedAttemptsBeforeProviderExchange(t *testing.T) {
	now := time.Now().UTC()
	targetChannelID := uuid.New()
	initiatorUserID := uuid.New()
	otherUserID := uuid.New()

	tests := []struct {
		name          string
		attempt       authsessions.OAuthAttempt
		sessionUserID uuid.UUID
		sessionBanned bool
	}{
		{
			name: "expired",
			attempt: authsessions.OAuthAttempt{
				Platform:        platformentity.PlatformKick,
				RedirectTo:      "/dashboard",
				CodeVerifier:    "stored-pkce-verifier",
				TargetChannelID: &targetChannelID,
				InitiatorUserID: &initiatorUserID,
				ExpiresAt:       now.Add(-time.Minute),
			},
			sessionUserID: initiatorUserID,
		},
		{
			name: "missing initiator",
			attempt: authsessions.OAuthAttempt{
				Platform:        platformentity.PlatformKick,
				RedirectTo:      "/dashboard",
				CodeVerifier:    "stored-pkce-verifier",
				TargetChannelID: &targetChannelID,
				ExpiresAt:       now.Add(time.Minute),
			},
			sessionUserID: initiatorUserID,
		},
		{
			name: "missing expiry",
			attempt: authsessions.OAuthAttempt{
				Platform:        platformentity.PlatformKick,
				RedirectTo:      "/dashboard",
				CodeVerifier:    "stored-pkce-verifier",
				TargetChannelID: &targetChannelID,
				InitiatorUserID: &initiatorUserID,
			},
			sessionUserID: initiatorUserID,
		},
		{
			name: "different callback session",
			attempt: authsessions.OAuthAttempt{
				Platform:        platformentity.PlatformKick,
				RedirectTo:      "/dashboard",
				CodeVerifier:    "stored-pkce-verifier",
				TargetChannelID: &targetChannelID,
				InitiatorUserID: &initiatorUserID,
				ExpiresAt:       now.Add(time.Minute),
			},
			sessionUserID: otherUserID,
		},
		{
			name: "banned callback session",
			attempt: authsessions.OAuthAttempt{
				Platform:        platformentity.PlatformKick,
				RedirectTo:      "/dashboard",
				CodeVerifier:    "stored-pkce-verifier",
				TargetChannelID: &targetChannelID,
				InitiatorUserID: &initiatorUserID,
				ExpiresAt:       now.Add(time.Minute),
			},
			sessionUserID: initiatorUserID,
			sessionBanned: true,
		},
		{
			name: "blank PKCE verifier",
			attempt: authsessions.OAuthAttempt{
				Platform:        platformentity.PlatformKick,
				RedirectTo:      "/dashboard",
				CodeVerifier:    " \t\n",
				TargetChannelID: &targetChannelID,
				InitiatorUserID: &initiatorUserID,
				ExpiresAt:       now.Add(time.Minute),
			},
			sessionUserID: initiatorUserID,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			const state = "targeted-oauth-state"
			exchangeCalls := 0
			getUserCalls := 0
			tokenCalls := 0
			bindingCalls := 0
			providerUserID := uuid.New()
			sessions := &fakeOAuthSession{
				internalUserID: tt.sessionUserID,
				attempts:       map[string]authsessions.OAuthAttempt{state: tt.attempt},
			}
			provider := &oauthPlatformProvider{
				name: platformentity.PlatformKick.String(),
				exchangeCodeFunc: func(context.Context, appplatform.ExchangeCodeInput) (*appplatform.PlatformTokens, error) {
					exchangeCalls++
					return testPlatformTokens(), nil
				},
				getUserFunc: func(context.Context, string) (*appplatform.PlatformUser, error) {
					getUserCalls++
					return &appplatform.PlatformUser{ID: "provider-kick"}, nil
				},
			}
			bindings := &oauthChannelPlatformsRepository{
				getByChannelAndPlatformFunc: func(context.Context, uuid.UUID, platformentity.Platform) (channelplatformsmodel.ChannelPlatform, error) {
					bindingCalls++
					return channelplatformsmodel.Nil, channelplatforms.ErrNotFound
				},
				createFunc: func(_ context.Context, input channelplatforms.CreateInput) (channelplatformsmodel.ChannelPlatform, error) {
					bindingCalls++
					return channelplatformsmodel.ChannelPlatform{ID: uuid.New(), ChannelID: input.ChannelID, Platform: input.Platform, UserID: input.UserID, PlatformChannelID: input.PlatformChannelID, Enabled: input.Enabled}, nil
				},
			}
			tokens := &fakeTokensRepository{
				getByUserIDFunc: func(context.Context, uuid.UUID) (*tokensmodel.Token, error) {
					tokenCalls++
					return nil, tokensrepo.ErrNotFound
				},
				createUserTokenFunc: func(context.Context, tokensrepo.CreateInput) (*tokensmodel.Token, error) {
					tokenCalls++
					return &tokensmodel.Token{ID: uuid.New()}, nil
				},
			}
			authHandler := newOAuthFlowTestAuth(oauthFlowTestAuthOpts{
				sessions: sessions,
				users: &oauthUsersRepository{
					getByIDFunc: func(_ context.Context, userID uuid.UUID) (usersmodel.User, error) {
						return usersmodel.User{ID: userID, IsBanned: tt.sessionBanned}, nil
					},
					getByPlatformIDFunc: func(context.Context, platformentity.Platform, string) (usersmodel.User, error) {
						return usersmodel.User{ID: providerUserID, Platform: platformentity.PlatformKick, PlatformID: "provider-kick"}, nil
					},
				},
				channels: &oauthChannelsRepository{
					getByIDFunc: func(context.Context, uuid.UUID) (channelsmodel.Channel, error) {
						return channelsmodel.Channel{ID: targetChannelID}, nil
					},
					getByBindingUserIDFunc: func(context.Context, platformentity.Platform, uuid.UUID) (channelsmodel.Channel, error) {
						return channelsmodel.Nil, channelsrepo.ErrNotFound
					},
				},
				bindings: bindings,
				tokens:   tokens,
				registry: appplatform.NewRegistry([]appplatform.PlatformProvider{provider}),
			})
			authHandler.dashboardAccess = dashboardAccessFunc(func(context.Context, dashboardaccess.Subject, uuid.UUID, string) (bool, error) {
				return true, nil
			})

			output, err := authHandler.handlePlatformCode(context.Background(), platformCodeInput{
				Platform: platformentity.PlatformKick,
				Code:     "authorization-code",
				State:    state,
				DeviceID: "callback-device-id",
			})
			if output != nil {
				t.Fatalf("handlePlatformCode() output = %#v, want nil", output)
			}
			if err == nil || !strings.Contains(err.Error(), "Invalid or expired OAuth state") {
				t.Fatalf("handlePlatformCode() error = %v, want non-enumerating invalid-state error", err)
			}
			if strings.Contains(err.Error(), targetChannelID.String()) || strings.Contains(err.Error(), initiatorUserID.String()) {
				t.Fatalf("handlePlatformCode() leaked target metadata: %v", err)
			}
			if exchangeCalls != 0 || getUserCalls != 0 || tokenCalls != 0 || bindingCalls != 0 {
				t.Fatalf("provider and persistence calls = exchange %d, user %d, token %d, binding %d, want all zero", exchangeCalls, getUserCalls, tokenCalls, bindingCalls)
			}
			if sessions.setOAuthAttemptCalls != 0 {
				t.Fatalf("callback OAuth attempt writes = %d, want 0", sessions.setOAuthAttemptCalls)
			}
			if _, exists := sessions.attempts[state]; exists {
				t.Fatal("invalid targeted OAuth attempt was not removed")
			}
		})
	}
}

func TestCompletePlatformCodeRejectsChangedTargetedInitiatorBeforeLocalAuth(t *testing.T) {
	const state = "targeted-oauth-state"

	initiatorUserID := uuid.New()
	changedUserID := uuid.New()
	targetChannelID := uuid.New()
	providerUserID := uuid.New()
	sessionLookupIDs := make([]uuid.UUID, 0, 2)
	exchangeCalls := 0
	getUserCalls := 0
	providerUserLookups := 0
	providerUserWrites := 0
	tokenCalls := 0
	bindingCalls := 0
	accessCalls := 0
	transaction := &oauthTransaction{}
	sessions := &fakeOAuthSession{
		internalUserIDs: []uuid.UUID{initiatorUserID, changedUserID},
		attempts: map[string]authsessions.OAuthAttempt{
			state: {
				Platform:        platformentity.PlatformKick,
				RedirectTo:      "/dashboard",
				CodeVerifier:    "stored-pkce-verifier",
				TargetChannelID: &targetChannelID,
				InitiatorUserID: &initiatorUserID,
				ExpiresAt:       time.Now().UTC().Add(time.Minute),
			},
		},
	}
	provider := &oauthPlatformProvider{
		name: platformentity.PlatformKick.String(),
		exchangeCodeFunc: func(context.Context, appplatform.ExchangeCodeInput) (*appplatform.PlatformTokens, error) {
			exchangeCalls++
			return testPlatformTokens(), nil
		},
		getUserFunc: func(context.Context, string) (*appplatform.PlatformUser, error) {
			getUserCalls++
			return &appplatform.PlatformUser{ID: "provider-kick"}, nil
		},
	}
	users := &oauthUsersRepository{
		getByIDFunc: func(_ context.Context, userID uuid.UUID) (usersmodel.User, error) {
			sessionLookupIDs = append(sessionLookupIDs, userID)
			return usersmodel.User{ID: userID}, nil
		},
		getByPlatformIDFunc: func(context.Context, platformentity.Platform, string) (usersmodel.User, error) {
			providerUserLookups++
			return usersmodel.User{ID: providerUserID, Platform: platformentity.PlatformKick, PlatformID: "provider-kick"}, nil
		},
		updateFunc: func(_ context.Context, id uuid.UUID, _ usersrepo.UpdateInput) (usersmodel.User, error) {
			providerUserWrites++
			return usersmodel.User{ID: id}, nil
		},
	}
	tokens := &fakeTokensRepository{
		getByUserIDFunc: func(context.Context, uuid.UUID) (*tokensmodel.Token, error) {
			tokenCalls++
			return nil, tokensrepo.ErrNotFound
		},
		createUserTokenFunc: func(context.Context, tokensrepo.CreateInput) (*tokensmodel.Token, error) {
			tokenCalls++
			return &tokensmodel.Token{ID: uuid.New()}, nil
		},
		updateTokenByIDFunc: func(context.Context, uuid.UUID, tokensrepo.UpdateTokenInput) (*tokensmodel.Token, error) {
			tokenCalls++
			return &tokensmodel.Token{ID: uuid.New()}, nil
		},
	}
	bindings := &oauthChannelPlatformsRepository{
		getByChannelAndPlatformFunc: func(context.Context, uuid.UUID, platformentity.Platform) (channelplatformsmodel.ChannelPlatform, error) {
			bindingCalls++
			return channelplatformsmodel.Nil, channelplatforms.ErrNotFound
		},
		createFunc: func(_ context.Context, input channelplatforms.CreateInput) (channelplatformsmodel.ChannelPlatform, error) {
			bindingCalls++
			return channelplatformsmodel.ChannelPlatform{ID: uuid.New(), ChannelID: input.ChannelID, Platform: input.Platform, UserID: input.UserID, PlatformChannelID: input.PlatformChannelID, Enabled: input.Enabled}, nil
		},
	}
	authHandler := newOAuthFlowTestAuth(oauthFlowTestAuthOpts{
		sessions: sessions,
		users:    users,
		channels: &oauthChannelsRepository{
			getByIDFunc: func(_ context.Context, channelID uuid.UUID) (channelsmodel.Channel, error) {
				return channelsmodel.Channel{ID: channelID}, nil
			},
		},
		bindings:    bindings,
		tokens:      tokens,
		registry:    appplatform.NewRegistry([]appplatform.PlatformProvider{provider}),
		transaction: transaction,
	})
	authHandler.dashboardAccess = dashboardAccessFunc(func(context.Context, dashboardaccess.Subject, uuid.UUID, string) (bool, error) {
		accessCalls++
		return true, nil
	})

	output, err := authHandler.handlePlatformCode(context.Background(), platformCodeInput{
		Platform: platformentity.PlatformKick,
		Code:     "authorization-code",
		State:    state,
		DeviceID: "callback-device-id",
	})
	if output != nil {
		t.Fatalf("handlePlatformCode() output = %#v, want nil", output)
	}
	if err == nil || !strings.Contains(err.Error(), "Forbidden") {
		t.Fatalf("handlePlatformCode() error = %v, want non-enumerating forbidden error", err)
	}
	if strings.Contains(err.Error(), initiatorUserID.String()) || strings.Contains(err.Error(), changedUserID.String()) {
		t.Fatalf("handlePlatformCode() leaked initiator metadata: %v", err)
	}
	if exchangeCalls != 1 || getUserCalls != 1 {
		t.Fatalf("provider calls = exchange %d, user %d, want one each", exchangeCalls, getUserCalls)
	}
	if !reflect.DeepEqual(sessionLookupIDs, []uuid.UUID{initiatorUserID, changedUserID}) {
		t.Fatalf("session user lookups = %#v, want [%s %s]", sessionLookupIDs, initiatorUserID, changedUserID)
	}
	if providerUserLookups != 0 || providerUserWrites != 0 || tokenCalls != 0 || bindingCalls != 0 || accessCalls != 0 || transaction.calls != 0 {
		t.Fatalf("local auth calls = provider lookup %d, provider write %d, token %d, binding %d, access %d, transaction %d, want all zero", providerUserLookups, providerUserWrites, tokenCalls, bindingCalls, accessCalls, transaction.calls)
	}
}

func TestCompletePlatformAuthRejectsUnauthorizedTargetDashboard(t *testing.T) {
	sessionUserID := uuid.New()
	targetChannelID := uuid.New()
	accessCalls := 0
	authHandler := newOAuthFlowTestAuth(oauthFlowTestAuthOpts{
		sessions: &fakeOAuthSession{internalUserID: sessionUserID},
		users: &oauthUsersRepository{
			getByIDFunc: func(_ context.Context, userID uuid.UUID) (usersmodel.User, error) {
				if userID != sessionUserID {
					t.Fatalf("session user ID = %s, want %s", userID, sessionUserID)
				}

				return usersmodel.User{ID: sessionUserID}, nil
			},
			getByPlatformIDFunc: func(context.Context, platformentity.Platform, string) (usersmodel.User, error) {
				t.Fatal("platform user lookup must not run for an unauthorized target")
				return usersmodel.Nil, nil
			},
		},
	})
	authHandler.dashboardAccess = dashboardAccessFunc(func(_ context.Context, subject dashboardaccess.Subject, channelID uuid.UUID, permission string) (bool, error) {
		accessCalls++
		if subject.ID != sessionUserID.String() || channelID != targetChannelID || permission != "" {
			t.Fatalf("dashboard access request = %+v, %s, %q", subject, channelID, permission)
		}

		return false, nil
	})

	_, err := authHandler.completePlatformAuth(context.Background(), completePlatformAuthInput{
		Platform:        platformentity.PlatformKick,
		PlatformUser:    &appplatform.PlatformUser{ID: "provider-kick"},
		Tokens:          testPlatformTokens(),
		TargetChannelID: &targetChannelID,
		InitiatorUserID: &sessionUserID,
	})
	if !errors.Is(err, errAuthForbidden) {
		t.Fatalf("completePlatformAuth() error = %v, want errAuthForbidden", err)
	}
	if accessCalls != 1 {
		t.Fatalf("dashboard access calls = %d, want 1", accessCalls)
	}
}

func TestStartPlatformAuthRejectsUnregisteredProvider(t *testing.T) {
	_, err := newOAuthFlowTestAuth(oauthFlowTestAuthOpts{
		sessions: &fakeOAuthSession{},
		registry: appplatform.NewRegistry(nil),
	}).startPlatformAuth(context.Background(), platformentity.PlatformVKVideoLive, "/dashboard")
	if !errors.Is(err, errPlatformUnavailable) {
		t.Fatalf("unregistered provider error = %v, want errPlatformUnavailable", err)
	}
}

func TestNewRegistersGenericAuthorizeRouteForUnregisteredProviders(t *testing.T) {
	_, api := humatest.New(t)
	New(Opts{Huma: api})

	response := api.Get("/auth/vk_video_live/authorize")
	if response.Code != 404 || !strings.Contains(response.Body.String(), "Platform is not available") {
		t.Fatalf("generic authorize response = %d %s", response.Code, response.Body.String())
	}
}

func TestAssignDefaultKickBotToUnconfiguredKickBindings(t *testing.T) {
	ctx := context.Background()
	transaction := &oauthTransaction{}
	defaultKickBotID := uuid.New()
	defaultKickBotUserID := uuid.New()
	channelID := uuid.New()
	configuredChannelID := uuid.New()
	unconfiguredBinding := channelplatformsmodel.ChannelPlatform{
		ID:                uuid.New(),
		ChannelID:         channelID,
		Platform:          platformentity.PlatformKick,
		UserID:            uuid.New(),
		PlatformChannelID: "kick-channel",
		Enabled:           true,
		BotConfig:         json.RawMessage(`{"existing":"value"}`),
	}
	configuredBinding := channelplatformsmodel.ChannelPlatform{
		ID:                uuid.New(),
		ChannelID:         configuredChannelID,
		Platform:          platformentity.PlatformKick,
		UserID:            uuid.New(),
		PlatformChannelID: "configured-kick-channel",
		Enabled:           true,
		BotConfig:         json.RawMessage(`{"kick_bot_id":"already-selected"}`),
	}
	channels := &oauthChannelsRepository{
		getAllByBindingPlatformFunc: func(_ context.Context, gotPlatform platformentity.Platform) ([]channelsmodel.Channel, error) {
			if gotPlatform != platformentity.PlatformKick {
				t.Fatalf("binding platform = %s, want kick", gotPlatform)
			}
			return []channelsmodel.Channel{
				{ID: channelID, Bindings: []channelplatformsmodel.ChannelPlatform{unconfiguredBinding}},
				{ID: configuredChannelID, Bindings: []channelplatformsmodel.ChannelPlatform{configuredBinding}},
			}, nil
		},
	}
	bindings := &oauthChannelPlatformsRepository{
		updateFunc: func(_ context.Context, bindingID uuid.UUID, input channelplatforms.UpdateInput) (channelplatformsmodel.ChannelPlatform, error) {
			if bindingID != unconfiguredBinding.ID || input.UserID != unconfiguredBinding.UserID || input.PlatformChannelID != unconfiguredBinding.PlatformChannelID || !input.Enabled {
				t.Fatalf("updated Kick binding = (%s, %+v)", bindingID, input)
			}
			if input.BotUserID == nil || *input.BotUserID != defaultKickBotUserID {
				t.Fatalf("Kick binding bot user = %v, want %s", input.BotUserID, defaultKickBotUserID)
			}
			var botConfig map[string]string
			if err := json.Unmarshal(input.BotConfig, &botConfig); err != nil {
				t.Fatalf("decode updated Kick binding config: %v", err)
			}
			if botConfig["existing"] != "value" || botConfig["kick_bot_id"] != defaultKickBotID.String() {
				t.Fatalf("updated Kick binding config = %#v", botConfig)
			}

			return unconfiguredBinding, nil
		},
	}

	affectedChannelIDs, err := newOAuthFlowTestAuth(oauthFlowTestAuthOpts{
		channels:    channels,
		bindings:    bindings,
		transaction: transaction,
	}).assignDefaultKickBotToChannels(ctx, kickbotentity.KickBot{
		ID:         defaultKickBotID,
		KickUserID: defaultKickBotUserID,
	})
	if err != nil {
		t.Fatalf("assign default Kick bot: %v", err)
	}
	if !reflect.DeepEqual(affectedChannelIDs, []string{channelID.String()}) {
		t.Fatalf("affected channel IDs = %#v, want [%s]", affectedChannelIDs, channelID)
	}
	if bindings.updateCalls != 1 {
		t.Fatalf("Kick binding updates = %d, want 1", bindings.updateCalls)
	}
	if transaction.calls != 1 || !transaction.committed {
		t.Fatalf("Kick binding transaction = calls %d committed %t, want one committed transaction", transaction.calls, transaction.committed)
	}
}

func newTestAuth(tokensRepository tokensrepo.Repository, usersRepository usersrepo.Repository) *Auth {
	return &Auth{
		config:           cfg.Config{TokensCipherKey: "pnyfwfiulmnqlhkvixaeligpprcnlyke"},
		logger:           slog.New(slog.NewTextHandler(testWriter{t: nil}, nil)),
		tokensRepository: tokensRepository,
		usersRepo:        usersRepository,
	}
}

func testPlatformTokens() *appplatform.PlatformTokens {
	return &appplatform.PlatformTokens{
		AccessToken:  "access-token",
		RefreshToken: "refresh-token",
		ExpiresIn:    7200,
		Scopes:       []string{"chat:read", "chat:write"},
	}
}

type oauthFlowTestAuthOpts struct {
	sessions          *fakeOAuthSession
	users             *oauthUsersRepository
	channels          *oauthChannelsRepository
	bindings          *oauthChannelPlatformsRepository
	tokens            tokensrepo.Repository
	bots              botsrepo.Repository
	registry          *appplatform.Registry
	transaction       *oauthTransaction
	eventSubPublisher *oauthEventSubPublisher
}

func newOAuthFlowTestAuth(opts oauthFlowTestAuthOpts) *Auth {
	if opts.sessions == nil {
		opts.sessions = &fakeOAuthSession{}
	}
	if opts.users == nil {
		opts.users = &oauthUsersRepository{}
	}
	if opts.channels == nil {
		opts.channels = &oauthChannelsRepository{}
	}
	if opts.bindings == nil {
		opts.bindings = &oauthChannelPlatformsRepository{}
	}
	if opts.tokens == nil {
		opts.tokens = &fakeTokensRepository{}
	}
	if opts.bots == nil {
		opts.bots = &oauthBotsRepository{defaultBot: botsmodel.Bot{ID: "default-bot"}}
	}
	if opts.registry == nil {
		opts.registry = appplatform.NewRegistry(nil)
	}
	if opts.transaction == nil {
		opts.transaction = &oauthTransaction{}
	}
	if opts.eventSubPublisher == nil {
		opts.eventSubPublisher = &oauthEventSubPublisher{transaction: opts.transaction}
	}

	return &Auth{
		config:               cfg.Config{TokensCipherKey: "pnyfwfiulmnqlhkvixaeligpprcnlyke"},
		logger:               slog.New(slog.NewTextHandler(testWriter{t: nil}, nil)),
		sessions:             opts.sessions,
		usersRepo:            opts.users,
		channelsRepo:         opts.channels,
		channelPlatformsRepo: opts.bindings,
		tokensRepository:     opts.tokens,
		botsRepo:             opts.bots,
		platformRegistry:     opts.registry,
		transactionRunner:    opts.transaction,
		eventSubPublisher:    opts.eventSubPublisher,
	}
}

func newCreateTokenRepository(userID uuid.UUID) *fakeTokensRepository {
	return &fakeTokensRepository{
		getByUserIDFunc: func(_ context.Context, gotUserID uuid.UUID) (*tokensmodel.Token, error) {
			if gotUserID != userID {
				return nil, errors.New("token lookup used the wrong platform user")
			}

			return nil, tokensrepo.ErrNotFound
		},
		createUserTokenFunc: func(context.Context, tokensrepo.CreateInput) (*tokensmodel.Token, error) {
			return &tokensmodel.Token{ID: uuid.New()}, nil
		},
		updateTokenByIDFunc: func(context.Context, uuid.UUID, tokensrepo.UpdateTokenInput) (*tokensmodel.Token, error) {
			return nil, nil
		},
	}
}

type fakeOAuthSession struct {
	internalUserID       uuid.UUID
	internalUserIDs      []uuid.UUID
	internalUserErr      error
	attempts             map[string]authsessions.OAuthAttempt
	setOAuthAttemptCalls int
}

func (s *fakeOAuthSession) GetInternalUserID(context.Context) (uuid.UUID, error) {
	if s.internalUserErr != nil {
		return uuid.UUID{}, s.internalUserErr
	}
	if len(s.internalUserIDs) > 0 {
		userID := s.internalUserIDs[0]
		s.internalUserIDs = s.internalUserIDs[1:]
		return userID, nil
	}

	return s.internalUserID, nil
}

func (s *fakeOAuthSession) SetSessionInternalUserID(_ context.Context, id uuid.UUID) error {
	s.internalUserID = id
	s.internalUserErr = nil
	return nil
}

func (*fakeOAuthSession) SetSessionCurrentPlatform(context.Context, string) error { return nil }

func (*fakeOAuthSession) SetSessionSelectedDashboard(context.Context, string) error { return nil }

func (*fakeOAuthSession) SetSessionTwitchUser(context.Context, helix.User) error { return nil }

func (*fakeOAuthSession) SetSessionKickUser(context.Context, authsessions.KickSessionUser) error {
	return nil
}

func (s *fakeOAuthSession) SetOAuthAttempt(_ context.Context, state string, attempt authsessions.OAuthAttempt) error {
	s.setOAuthAttemptCalls++
	if s.attempts == nil {
		s.attempts = make(map[string]authsessions.OAuthAttempt)
	}
	s.attempts[state] = attempt
	return nil
}

func (s *fakeOAuthSession) GetOAuthAttempt(_ context.Context, state string) (authsessions.OAuthAttempt, error) {
	attempt, ok := s.attempts[state]
	if !ok {
		return authsessions.OAuthAttempt{}, authsessions.ErrOAuthAttemptNotFound
	}

	return attempt, nil
}

func (s *fakeOAuthSession) DeleteOAuthAttempt(_ context.Context, state string) error {
	delete(s.attempts, state)
	return nil
}

type oauthTransaction struct {
	committed bool
	calls     int
}

func (t *oauthTransaction) Do(ctx context.Context, fn func(context.Context) error) error {
	t.calls++
	if err := fn(ctx); err != nil {
		return err
	}

	t.committed = true
	return nil
}

type oauthEventSubPublisher struct {
	transaction *oauthTransaction
	requests    []buscoreeventsub.EventsubSubscribeToAllEventsRequest
}

type dashboardAccessFunc func(context.Context, dashboardaccess.Subject, uuid.UUID, string) (bool, error)

func (f dashboardAccessFunc) CanAccess(
	ctx context.Context,
	subject dashboardaccess.Subject,
	channelID uuid.UUID,
	permission string,
) (bool, error) {
	return f(ctx, subject, channelID, permission)
}

func (p *oauthEventSubPublisher) Publish(_ context.Context, input buscoreeventsub.EventsubSubscribeToAllEventsRequest) error {
	if p.transaction != nil && !p.transaction.committed {
		return errors.New("EventSub published before transaction commit")
	}
	p.requests = append(p.requests, input)
	return nil
}

type oauthPlatformProvider struct {
	name             string
	getAuthURLFunc   func(string, string) string
	exchangeCodeFunc func(context.Context, appplatform.ExchangeCodeInput) (*appplatform.PlatformTokens, error)
	getUserFunc      func(context.Context, string) (*appplatform.PlatformUser, error)
}

func (p *oauthPlatformProvider) Name() string { return p.name }

func (p *oauthPlatformProvider) GetAuthURL(state, codeChallenge string) string {
	if p.getAuthURLFunc == nil {
		return ""
	}

	return p.getAuthURLFunc(state, codeChallenge)
}

func (p *oauthPlatformProvider) ExchangeCode(ctx context.Context, input appplatform.ExchangeCodeInput) (*appplatform.PlatformTokens, error) {
	if p.exchangeCodeFunc == nil {
		return nil, errors.New("unexpected ExchangeCode call")
	}

	return p.exchangeCodeFunc(ctx, input)
}

func (*oauthPlatformProvider) RefreshToken(context.Context, appplatform.RefreshTokenInput) (*appplatform.PlatformTokens, error) {
	return nil, errors.New("unexpected RefreshToken call")
}

func (p *oauthPlatformProvider) GetUser(ctx context.Context, accessToken string) (*appplatform.PlatformUser, error) {
	if p.getUserFunc == nil {
		return nil, errors.New("unexpected GetUser call")
	}

	return p.getUserFunc(ctx, accessToken)
}

type oauthBotsRepository struct {
	defaultBot botsmodel.Bot
	err        error
}

func (r *oauthBotsRepository) GetDefault(context.Context) (botsmodel.Bot, error) {
	return r.defaultBot, r.err
}

type oauthChannelsRepository struct {
	createFunc                  func(context.Context, channelsrepo.CreateInput) (channelsmodel.Channel, error)
	getByIDFunc                 func(context.Context, uuid.UUID) (channelsmodel.Channel, error)
	getByBindingUserIDFunc      func(context.Context, platformentity.Platform, uuid.UUID) (channelsmodel.Channel, error)
	getAllByBindingPlatformFunc func(context.Context, platformentity.Platform) ([]channelsmodel.Channel, error)
	createCalls                 int
}

func (r *oauthChannelsRepository) Create(ctx context.Context, input channelsrepo.CreateInput) (channelsmodel.Channel, error) {
	r.createCalls++
	if r.createFunc == nil {
		return channelsmodel.Nil, errors.New("unexpected Create call")
	}

	return r.createFunc(ctx, input)
}

func (r *oauthChannelsRepository) GetByBindingUserID(ctx context.Context, p platformentity.Platform, userID uuid.UUID) (channelsmodel.Channel, error) {
	if r.getByBindingUserIDFunc == nil {
		return channelsmodel.Nil, errors.New("unexpected GetByBindingUserID call")
	}

	return r.getByBindingUserIDFunc(ctx, p, userID)
}

func (*oauthChannelsRepository) GetMany(context.Context, channelsrepo.GetManyInput) ([]channelsmodel.Channel, error) {
	return nil, errors.New("unexpected GetMany call")
}

func (r *oauthChannelsRepository) GetAllByBindingPlatform(ctx context.Context, platform platformentity.Platform) ([]channelsmodel.Channel, error) {
	if r.getAllByBindingPlatformFunc == nil {
		return nil, errors.New("unexpected GetAllByBindingPlatform call")
	}

	return r.getAllByBindingPlatformFunc(ctx, platform)
}

func (r *oauthChannelsRepository) GetByID(ctx context.Context, channelID uuid.UUID) (channelsmodel.Channel, error) {
	if r.getByIDFunc == nil {
		return channelsmodel.Nil, errors.New("unexpected GetByID call")
	}

	return r.getByIDFunc(ctx, channelID)
}

func (*oauthChannelsRepository) GetByApiKey(context.Context, string) (channelsmodel.Channel, error) {
	return channelsmodel.Nil, errors.New("unexpected GetByApiKey call")
}

func (*oauthChannelsRepository) GetByPlatformChannelID(context.Context, platformentity.Platform, string) (channelsmodel.Channel, error) {
	return channelsmodel.Nil, errors.New("unexpected GetByPlatformChannelID call")
}

func (*oauthChannelsRepository) GetBySlug(context.Context, channelsrepo.GetBySlugInput) (channelsmodel.Channel, error) {
	return channelsmodel.Nil, errors.New("unexpected GetBySlug call")
}

func (*oauthChannelsRepository) GetCount(context.Context, channelsrepo.GetCountInput) (int, error) {
	return 0, errors.New("unexpected GetCount call")
}

func (*oauthChannelsRepository) Update(context.Context, uuid.UUID, channelsrepo.UpdateInput) (channelsmodel.Channel, error) {
	return channelsmodel.Nil, errors.New("unexpected Update call")
}

type oauthChannelPlatformsRepository struct {
	getByChannelAndPlatformFunc func(context.Context, uuid.UUID, platformentity.Platform) (channelplatformsmodel.ChannelPlatform, error)
	createFunc                  func(context.Context, channelplatforms.CreateInput) (channelplatformsmodel.ChannelPlatform, error)
	updateFunc                  func(context.Context, uuid.UUID, channelplatforms.UpdateInput) (channelplatformsmodel.ChannelPlatform, error)
	createCalls                 int
	updateCalls                 int
}

func (*oauthChannelPlatformsRepository) LockByChannelID(context.Context, uuid.UUID) error {
	return nil
}

func (r *oauthChannelPlatformsRepository) Create(ctx context.Context, input channelplatforms.CreateInput) (channelplatformsmodel.ChannelPlatform, error) {
	r.createCalls++
	if r.createFunc == nil {
		return channelplatformsmodel.Nil, errors.New("unexpected binding Create call")
	}

	return r.createFunc(ctx, input)
}

func (r *oauthChannelPlatformsRepository) GetByChannelAndPlatform(ctx context.Context, channelID uuid.UUID, p platformentity.Platform) (channelplatformsmodel.ChannelPlatform, error) {
	if r.getByChannelAndPlatformFunc == nil {
		return channelplatformsmodel.Nil, errors.New("unexpected GetByChannelAndPlatform call")
	}

	return r.getByChannelAndPlatformFunc(ctx, channelID, p)
}

func (*oauthChannelPlatformsRepository) GetByPlatformChannelID(context.Context, platformentity.Platform, string) (channelplatformsmodel.ChannelPlatform, error) {
	return channelplatformsmodel.Nil, errors.New("unexpected GetByPlatformChannelID call")
}

func (*oauthChannelPlatformsRepository) ListByChannelID(context.Context, uuid.UUID) ([]channelplatformsmodel.ChannelPlatform, error) {
	return nil, errors.New("unexpected ListByChannelID call")
}

func (r *oauthChannelPlatformsRepository) Update(ctx context.Context, id uuid.UUID, input channelplatforms.UpdateInput) (channelplatformsmodel.ChannelPlatform, error) {
	r.updateCalls++
	if r.updateFunc == nil {
		return channelplatformsmodel.Nil, errors.New("unexpected binding Update call")
	}

	return r.updateFunc(ctx, id, input)
}

func (*oauthChannelPlatformsRepository) Patch(context.Context, uuid.UUID, channelplatforms.PatchInput) (channelplatformsmodel.ChannelPlatform, error) {
	return channelplatformsmodel.Nil, errors.New("unexpected binding Patch call")
}

func (*oauthChannelPlatformsRepository) Delete(context.Context, uuid.UUID) error {
	return errors.New("unexpected binding Delete call")
}

type oauthUsersRepository struct {
	getByIDFunc         func(context.Context, uuid.UUID) (usersmodel.User, error)
	getByPlatformIDFunc func(context.Context, platformentity.Platform, string) (usersmodel.User, error)
	createFunc          func(context.Context, usersrepo.CreateInput) (usersmodel.User, error)
	updateFunc          func(context.Context, uuid.UUID, usersrepo.UpdateInput) (usersmodel.User, error)
}

func (r *oauthUsersRepository) GetByID(ctx context.Context, id uuid.UUID) (usersmodel.User, error) {
	if r.getByIDFunc == nil {
		return usersmodel.Nil, errors.New("unexpected GetByID call")
	}

	return r.getByIDFunc(ctx, id)
}

func (r *oauthUsersRepository) GetByPlatformID(ctx context.Context, p platformentity.Platform, platformID string) (usersmodel.User, error) {
	if r.getByPlatformIDFunc == nil {
		return usersmodel.Nil, errors.New("unexpected GetByPlatformID call")
	}

	return r.getByPlatformIDFunc(ctx, p, platformID)
}

func (r *oauthUsersRepository) Create(ctx context.Context, input usersrepo.CreateInput) (usersmodel.User, error) {
	if r.createFunc == nil {
		return usersmodel.Nil, errors.New("unexpected Create user call")
	}

	return r.createFunc(ctx, input)
}

func (r *oauthUsersRepository) Update(ctx context.Context, id uuid.UUID, input usersrepo.UpdateInput) (usersmodel.User, error) {
	if r.updateFunc == nil {
		return usersmodel.User{ID: id}, nil
	}

	return r.updateFunc(ctx, id, input)
}

func (*oauthUsersRepository) GetManyByIDS(context.Context, usersrepo.GetManyInput) ([]usersmodel.User, error) {
	return nil, errors.New("unexpected GetManyByIDS call")
}

func (*oauthUsersRepository) GetRandomOnlineUser(context.Context, usersrepo.GetRandomOnlineUserInput) (usersmodel.OnlineUser, error) {
	return usersmodel.NilOnlineUser, errors.New("unexpected GetRandomOnlineUser call")
}

func (*oauthUsersRepository) GetOnlineUsersWithFilters(context.Context, usersrepo.GetOnlineUsersWithFiltersInput) ([]usersmodel.OnlineUser, error) {
	return nil, errors.New("unexpected GetOnlineUsersWithFilters call")
}

func (*oauthUsersRepository) GetByApiKey(context.Context, string) (usersmodel.User, error) {
	return usersmodel.Nil, errors.New("unexpected GetByApiKey call")
}

type fakeTokensRepository struct {
	getByUserIDFunc      func(context.Context, uuid.UUID) (*tokensmodel.Token, error)
	createUserTokenFunc  func(context.Context, tokensrepo.CreateInput) (*tokensmodel.Token, error)
	updateTokenByIDFunc  func(context.Context, uuid.UUID, tokensrepo.UpdateTokenInput) (*tokensmodel.Token, error)
	createUserTokenCalls int
	updateTokenCalls     int
}

func (f *fakeTokensRepository) GetByID(context.Context, uuid.UUID) (*tokensmodel.Token, error) {
	panic("unexpected GetByID call")
}

func (f *fakeTokensRepository) GetByUserID(ctx context.Context, userID uuid.UUID) (*tokensmodel.Token, error) {
	if f.getByUserIDFunc == nil {
		panic("unexpected GetByUserID call")
	}

	return f.getByUserIDFunc(ctx, userID)
}

func (f *fakeTokensRepository) GetByBotID(context.Context, string) (*tokensmodel.Token, error) {
	panic("unexpected GetByBotID call")
}

func (f *fakeTokensRepository) CreateUserToken(ctx context.Context, input tokensrepo.CreateInput) (*tokensmodel.Token, error) {
	f.createUserTokenCalls++
	if f.createUserTokenFunc == nil {
		panic("unexpected CreateUserToken call")
	}

	return f.createUserTokenFunc(ctx, input)
}

func (f *fakeTokensRepository) UpdateTokenByID(
	ctx context.Context,
	id uuid.UUID,
	input tokensrepo.UpdateTokenInput,
) (*tokensmodel.Token, error) {
	f.updateTokenCalls++
	if f.updateTokenByIDFunc == nil {
		panic("unexpected UpdateTokenByID call")
	}

	return f.updateTokenByIDFunc(ctx, id, input)
}

type fakeUsersRepository struct {
	updateFunc  func(context.Context, uuid.UUID, usersrepo.UpdateInput) (usersmodel.User, error)
	updateCalls int
}

func (f *fakeUsersRepository) GetByID(context.Context, uuid.UUID) (usersmodel.User, error) {
	panic("unexpected GetByID call")
}

func (f *fakeUsersRepository) GetByPlatformID(context.Context, platformentity.Platform, string) (usersmodel.User, error) {
	panic("unexpected GetByPlatformID call")
}

func (f *fakeUsersRepository) GetManyByIDS(context.Context, usersrepo.GetManyInput) ([]usersmodel.User, error) {
	panic("unexpected GetManyByIDS call")
}

func (f *fakeUsersRepository) Update(ctx context.Context, id uuid.UUID, input usersrepo.UpdateInput) (usersmodel.User, error) {
	f.updateCalls++
	if f.updateFunc == nil {
		panic("unexpected Update call")
	}

	return f.updateFunc(ctx, id, input)
}

func (f *fakeUsersRepository) GetRandomOnlineUser(context.Context, usersrepo.GetRandomOnlineUserInput) (usersmodel.OnlineUser, error) {
	panic("unexpected GetRandomOnlineUser call")
}

func (f *fakeUsersRepository) GetOnlineUsersWithFilters(context.Context, usersrepo.GetOnlineUsersWithFiltersInput) ([]usersmodel.OnlineUser, error) {
	panic("unexpected GetOnlineUsersWithFilters call")
}

func (f *fakeUsersRepository) GetByApiKey(context.Context, string) (usersmodel.User, error) {
	panic("unexpected GetByApiKey call")
}

func (f *fakeUsersRepository) Create(context.Context, usersrepo.CreateInput) (usersmodel.User, error) {
	panic("unexpected Create call")
}

type testWriter struct {
	t *testing.T
}

func (w testWriter) Write(p []byte) (n int, err error) {
	return len(p), nil
}

var _ = time.Now
