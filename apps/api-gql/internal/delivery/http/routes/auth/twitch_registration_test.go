package auth

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/google/uuid"
	appplatform "github.com/twirapp/twir/apps/api-gql/internal/platform"
	channelentity "github.com/twirapp/twir/libs/entities/channel"
	channelplatformentity "github.com/twirapp/twir/libs/entities/channel_platform"
	platformentity "github.com/twirapp/twir/libs/entities/platform"
	channelplatforms "github.com/twirapp/twir/libs/repositories/channel_platforms"
	channelsrepo "github.com/twirapp/twir/libs/repositories/channels"
	tokensrepo "github.com/twirapp/twir/libs/repositories/tokens"
	tokensmodel "github.com/twirapp/twir/libs/repositories/tokens/model"
	usersrepo "github.com/twirapp/twir/libs/repositories/users"
	usersmodel "github.com/twirapp/twir/libs/repositories/users/model"
)

func TestCompletePlatformExchangePersistsTwitchTokenBeforeRegisteringBroadcaster(t *testing.T) {
	fixture := newTwitchCompletionFixture()

	result, err := fixture.auth.completePlatformExchange(
		context.Background(),
		platformentity.PlatformTwitch,
		fixture.provider,
		"authorization-code",
		"",
		"",
		"/dashboard",
		nil,
		nil,
	)
	if err != nil {
		t.Fatalf("complete Twitch platform exchange: %v", err)
	}
	if !reflect.DeepEqual(fixture.events, []string{"persist", "register"}) {
		t.Fatalf("completion events = %#v, want persistence before registration", fixture.events)
	}
	if !reflect.DeepEqual(fixture.registeredUserIDs, []uuid.UUID{fixture.platformUserID}) {
		t.Fatalf("registered user IDs = %#v, want [%s]", fixture.registeredUserIDs, fixture.platformUserID)
	}
	if result.AuthResult.PlatformUserID != fixture.platformUserID {
		t.Fatalf("platform user ID = %s, want %s", result.AuthResult.PlatformUserID, fixture.platformUserID)
	}
}

func TestCompletePlatformExchangeDoesNotRegisterTwitchBroadcasterWhenTokenPersistenceFails(t *testing.T) {
	fixture := newTwitchCompletionFixture()
	persistenceErr := errors.New("persist Twitch token")
	fixture.persistenceErr = persistenceErr

	_, err := fixture.auth.completePlatformExchange(
		context.Background(),
		platformentity.PlatformTwitch,
		fixture.provider,
		"authorization-code",
		"",
		"",
		"/dashboard",
		nil,
		nil,
	)
	if !errors.Is(err, persistenceErr) {
		t.Fatalf("complete Twitch platform exchange error = %v, want persistence error", err)
	}
	if len(fixture.registeredUserIDs) != 0 {
		t.Fatalf("registered user IDs = %#v, want no registration", fixture.registeredUserIDs)
	}
}

func TestRegisterTwitchUserAfterAuthAllowsDuplicateRuntimeRegistration(t *testing.T) {
	fixture := newTwitchCompletionFixture()
	result := completePlatformAuthResult{PlatformUserID: fixture.platformUserID}

	if err := fixture.auth.registerTwitchUserAfterAuth(context.Background(), result, nil, nil); err != nil {
		t.Fatalf("register Twitch broadcaster: %v", err)
	}
	if err := fixture.auth.registerTwitchUserAfterAuth(context.Background(), result, nil, nil); err != nil {
		t.Fatalf("register duplicate Twitch broadcaster: %v", err)
	}
	if !reflect.DeepEqual(fixture.registeredUserIDs, []uuid.UUID{fixture.platformUserID, fixture.platformUserID}) {
		t.Fatalf("registered user IDs = %#v, want two idempotent registrations", fixture.registeredUserIDs)
	}
}

type twitchCompletionFixture struct {
	auth              *Auth
	provider          *oauthPlatformProvider
	platformUserID    uuid.UUID
	registeredUserIDs []uuid.UUID
	events            []string
	persistenceErr    error
}

func newTwitchCompletionFixture() *twitchCompletionFixture {
	fixture := &twitchCompletionFixture{platformUserID: uuid.New()}
	channelID := uuid.New()
	fixture.provider = &oauthPlatformProvider{
		platform: platformentity.PlatformTwitch,
		exchangeCodeFunc: func(context.Context, appplatform.ExchangeCodeInput) (*appplatform.PlatformTokens, error) {
			return testPlatformTokens(), nil
		},
		getUserFunc: func(context.Context, string) (*appplatform.PlatformUser, error) {
			return &appplatform.PlatformUser{ID: "twitch-user"}, nil
		},
	}
	fixture.auth = newOAuthFlowTestAuth(oauthFlowTestAuthOpts{
		sessions: &fakeOAuthSession{internalUserErr: errors.New("not signed in")},
		users: &oauthUsersRepository{
			getByPlatformIDFunc: func(context.Context, platformentity.Platform, string) (usersmodel.User, error) {
				return usersmodel.Nil, usersmodel.ErrNotFound
			},
			createFunc: func(context.Context, usersrepo.CreateInput) (usersmodel.User, error) {
				return usersmodel.User{ID: fixture.platformUserID, Platform: platformentity.PlatformTwitch, PlatformID: "twitch-user"}, nil
			},
			updateFunc: func(context.Context, uuid.UUID, usersrepo.UpdateInput) (usersmodel.User, error) {
				return usersmodel.User{ID: fixture.platformUserID}, nil
			},
		},
		channels: &oauthChannelsRepository{
			createFunc: func(context.Context) (channelentity.Channel, error) {
				return channelentity.Channel{ID: channelID}, nil
			},
			getByBindingUserIDFunc: func(context.Context, platformentity.Platform, uuid.UUID) (channelentity.Channel, error) {
				return channelentity.Nil, channelsrepo.ErrNotFound
			},
		},
		bindings: &oauthChannelPlatformsRepository{
			getByChannelAndPlatformFunc: func(context.Context, uuid.UUID, platformentity.Platform) (channelplatformentity.ChannelPlatform, error) {
				return channelplatformentity.Nil, channelplatforms.ErrNotFound
			},
			createFunc: func(_ context.Context, input channelplatforms.CreateInput) (channelplatformentity.ChannelPlatform, error) {
				return channelplatformentity.ChannelPlatform{ID: uuid.New(), ChannelID: input.ChannelID, Platform: input.Platform, UserID: input.UserID, PlatformChannelID: input.PlatformChannelID, Enabled: input.Enabled, BotConfig: input.BotConfig}, nil
			},
		},
		tokens: &fakeTokensRepository{
			getByUserIDFunc: func(context.Context, uuid.UUID) (*tokensmodel.Token, error) {
				return nil, tokensrepo.ErrNotFound
			},
			createUserTokenFunc: func(context.Context, tokensrepo.CreateInput) (*tokensmodel.Token, error) {
				fixture.events = append(fixture.events, "persist")
				if fixture.persistenceErr != nil {
					return nil, fixture.persistenceErr
				}
				return &tokensmodel.Token{ID: uuid.New()}, nil
			},
		},
	})
	fixture.auth.postPlatformAuthHooks = map[platformentity.Platform]postPlatformAuthHook{
		platformentity.PlatformTwitch: fixture.auth.registerTwitchUserAfterAuth,
	}
	fixture.auth.registerTwitchUser = func(_ context.Context, userID uuid.UUID) error {
		fixture.events = append(fixture.events, "register")
		fixture.registeredUserIDs = append(fixture.registeredUserIDs, userID)
		return nil
	}

	return fixture
}
