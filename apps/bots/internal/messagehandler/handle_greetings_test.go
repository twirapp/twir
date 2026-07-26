package messagehandler

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/twirapp/twir/apps/bots/internal/twitchactions"
	buscore "github.com/twirapp/twir/libs/bus-core"
	botsbus "github.com/twirapp/twir/libs/bus-core/bots"
	"github.com/twirapp/twir/libs/bus-core/events"
	"github.com/twirapp/twir/libs/bus-core/generic"
	"github.com/twirapp/twir/libs/bus-core/parser"
	buscoretokens "github.com/twirapp/twir/libs/bus-core/tokens"
	channelcache "github.com/twirapp/twir/libs/cache/channel"
	genericcacher "github.com/twirapp/twir/libs/cache/generic-cacher"
	cfg "github.com/twirapp/twir/libs/config"
	channelentity "github.com/twirapp/twir/libs/entities/channel"
	channelplatformentity "github.com/twirapp/twir/libs/entities/channel_platform"
	"github.com/twirapp/twir/libs/entities/platform"
	"github.com/twirapp/twir/libs/repositories/greetings"
	greetingsmodel "github.com/twirapp/twir/libs/repositories/greetings/model"
	streamsmodel "github.com/twirapp/twir/libs/repositories/streams/model"
	usersmodel "github.com/twirapp/twir/libs/repositories/users/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type greetingRepositoryRecorder struct {
	updates []greetings.UpdateInput
	err     error
}

func (*greetingRepositoryRecorder) GetManyByChannelID(context.Context, uuid.UUID, greetings.GetManyInput) ([]greetingsmodel.Greeting, error) {
	return nil, nil
}
func (*greetingRepositoryRecorder) GetByID(context.Context, uuid.UUID) (greetingsmodel.Greeting, error) {
	return greetingsmodel.Greeting{}, nil
}
func (*greetingRepositoryRecorder) Create(context.Context, greetings.CreateInput) (greetingsmodel.Greeting, error) {
	return greetingsmodel.Greeting{}, nil
}
func (r *greetingRepositoryRecorder) Update(_ context.Context, _ uuid.UUID, input greetings.UpdateInput) (greetingsmodel.Greeting, error) {
	r.updates = append(r.updates, input)
	return greetingsmodel.Greeting{}, r.err
}
func (*greetingRepositoryRecorder) UpdateManyByChannelID(context.Context, greetings.UpdateManyInput) error {
	return nil
}
func (*greetingRepositoryRecorder) Delete(context.Context, uuid.UUID) error { return nil }
func (*greetingRepositoryRecorder) GetOneByChannelAndUserID(context.Context, greetings.GetOneInput) (greetingsmodel.Greeting, error) {
	return greetingsmodel.Greeting{}, nil
}

type greetingCacheKV struct {
	messageIDCache
	deletes []string
}

func (k *greetingCacheKV) Delete(_ context.Context, key string) error {
	k.deletes = append(k.deletes, key)
	return nil
}

type greetingEventQueue struct {
	chatMessageEmoteQueue[events.GreetingSendedMessage, struct{}]
	published []events.GreetingSendedMessage
}

func (q *greetingEventQueue) Publish(_ context.Context, message events.GreetingSendedMessage) error {
	q.mu.Lock()
	defer q.mu.Unlock()

	q.published = append(q.published, message)
	return nil
}

func (q *greetingEventQueue) publishedSnapshot() []events.GreetingSendedMessage {
	q.mu.Lock()
	defer q.mu.Unlock()

	return append([]events.GreetingSendedMessage(nil), q.published...)
}

type greetingFixture struct {
	handler        *MessageHandler
	message        enrichedChatMessage
	repository     *greetingRepositoryRecorder
	cache          *greetingCacheKV
	parser         *chatMessageEmoteQueue[parser.ParseVariablesInTextRequest, parser.ParseVariablesInTextResponse]
	sender         *chatMessageEmoteQueue[botsbus.SendMessageRequest, struct{}]
	greetingEvents *greetingEventQueue
}

type greetingFixtureInput struct {
	source                platform.Platform
	isReply, withShoutOut bool
}

func newGreetingFixture(t *testing.T, input greetingFixtureInput) greetingFixture {
	t.Helper()

	channelID := uuid.New()
	userID := uuid.New()
	cache := &greetingCacheKV{}
	repository := &greetingRepositoryRecorder{}
	parserQueue := &chatMessageEmoteQueue[parser.ParseVariablesInTextRequest, parser.ParseVariablesInTextResponse]{
		response: &buscore.QueueResponse[parser.ParseVariablesInTextResponse]{Data: parser.ParseVariablesInTextResponse{Text: "parsed greeting"}},
	}
	sendQueue := &chatMessageEmoteQueue[botsbus.SendMessageRequest, struct{}]{response: &buscore.QueueResponse[struct{}]{}}
	greetingEvents := &greetingEventQueue{}
	bus := buscore.NewNatsBus(nil)
	bus.Parser.ParseVariablesInText = parserQueue
	bus.Bots.SendMessage = sendQueue
	bus.Events.GreetingSended = greetingEvents

	return greetingFixture{
		handler: &MessageHandler{
			greetingsRepository: repository,
			greetingsCache: genericcacher.New(genericcacher.Opts[[]greetingsmodel.Greeting]{
				KV:        cache,
				KeyPrefix: "greetings:",
				LoadFn: func(context.Context, string) ([]greetingsmodel.Greeting, error) {
					return []greetingsmodel.Greeting{{
						ID: uuid.New(), ChannelID: channelID, UserID: userID, Enabled: true, Text: "hello", IsReply: input.isReply, WithShoutOut: input.withShoutOut,
					}}, nil
				},
			}),
			gorm:    newGreetingDB(t),
			twirBus: bus,
		},
		message: enrichedChatMessage{ChatMessage: generic.ChatMessage{
			Platform: string(input.source), MessageID: "provider-message-id", PlatformChannelID: "platform-channel-id", BroadcasterUserId: "broadcaster-id", ChatterUserId: "chatter-id",
			Message: &generic.ChatMessageMessage{Text: "message"},
		}, EnrichedData: chatMessageEnrichedData{
			DbChannel: channelentity.Channel{ID: channelID}, DbUser: &usersmodel.User{ID: userID}, ChannelStream: &streamsmodel.Stream{},
		}},
		repository: repository, cache: cache, parser: parserQueue, sender: sendQueue, greetingEvents: greetingEvents,
	}
}

func newGreetingDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(postgres.New(postgres.Config{DSN: "host=localhost"}), &gorm.Config{DisableAutomaticPing: true, DryRun: true})
	require.NoError(t, err)
	return db
}

func (f *greetingFixture) installTwitchActions(t *testing.T, calls *int) {
	t.Helper()
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		*calls++
		writer.Header().Set("Content-Type", "application/json")
		if _, err := writer.Write([]byte(`{"data":[]}`)); err != nil {
			t.Errorf("write shoutout response: %v", err)
		}
	}))
	t.Cleanup(server.Close)
	userID := uuid.New()
	f.handler.twirBus.Tokens.RequestUserToken = &chatMessageEmoteQueue[
		buscoretokens.GetUserTokenRequest,
		buscoretokens.TokenResponse,
	]{response: &buscore.QueueResponse[buscoretokens.TokenResponse]{
		Data: buscoretokens.TokenResponse{AccessToken: "test-token"},
	}}
	f.handler.twitchActions = twitchactions.New(twitchactions.Opts{
		Config: cfg.Config{TwitchMockEnabled: true, TwitchMockApiUrl: server.URL, TwitchClientId: "test-client"}, TwirBus: f.handler.twirBus,
		ChannelsByTwitchIDCache: &channelcache.TwitchUserIDCacher{GenericCacher: genericcacher.New(genericcacher.Opts[channelentity.Channel]{
			KV: messageIDCache{}, LoadFn: func(context.Context, string) (channelentity.Channel, error) {
				return channelentity.Channel{Bindings: []channelplatformentity.ChannelPlatform{{
					Platform: platform.PlatformTwitch, PlatformChannelID: f.message.BroadcasterUserId, UserID: userID, Enabled: true, BotConfig: []byte(`{"is_bot_mod":true}`),
				}}}, nil
			},
		})},
	})
}

func TestHandleGreetingsDispatchesThroughSourcePlatform(t *testing.T) {
	for _, test := range []struct {
		name                 string
		source, wantPlatform platform.Platform
		isReply              bool
	}{
		{name: "twitch reply", source: platform.PlatformTwitch, wantPlatform: platform.PlatformTwitch, isReply: true},
		{name: "kick", source: platform.PlatformKick, wantPlatform: platform.PlatformKick},
		{name: "empty defaults to twitch", wantPlatform: platform.PlatformTwitch},
	} {
		t.Run(test.name, func(t *testing.T) {
			// Given
			fixture := newGreetingFixture(t, greetingFixtureInput{source: test.source, isReply: test.isReply})

			// When
			err := fixture.handler.handleGreetings(context.Background(), fixture.message)

			// Then
			require.NoError(t, err)
			require.Equal(t, []botsbus.SendMessageRequest{{ChannelID: fixture.message.EnrichedData.DbChannel.ID, Platforms: []platform.Platform{test.wantPlatform}, Message: "parsed greeting", ReplyTo: map[bool]string{true: fixture.message.MessageID}[test.isReply]}}, fixture.sender.requestSnapshot())
			require.Len(t, fixture.repository.updates, 1)
			require.Equal(t, true, *fixture.repository.updates[0].Processed)
			require.Equal(t, []string{"greetings:" + fixture.message.EnrichedData.DbChannel.ID.String()}, fixture.cache.deletes)
		})
	}
}

func TestHandleGreetingsPublishesCanonicalKickIdentity(t *testing.T) {
	fixture := newGreetingFixture(t, greetingFixtureInput{source: platform.PlatformKick})

	err := fixture.handler.handleGreetings(context.Background(), fixture.message)

	require.NoError(t, err)
	published := fixture.greetingEvents.publishedSnapshot()
	require.Len(t, published, 1)
	event := published[0]
	require.Equal(t, platform.PlatformKick, event.BaseInfo.Platform)
	require.Equal(t, fixture.message.EnrichedData.DbChannel.ID, event.BaseInfo.ChannelDBID)
	require.Equal(t, fixture.message.PlatformChannelID, event.BaseInfo.ChannelPlatformID)
}

func TestHandleGreetingsLeavesGreetingUnprocessedWhenDeliveryFails(t *testing.T) {
	for _, test := range []struct {
		name      string
		setup     func(greetingFixture)
		wantSends int
	}{
		{name: "parsing fails", setup: func(f greetingFixture) { f.parser.requestErr = errors.New("parse greeting") }},
		{name: "sending fails", setup: func(f greetingFixture) { f.sender.requestErr = errors.New("send greeting") }, wantSends: 1},
	} {
		t.Run(test.name, func(t *testing.T) {
			// Given
			fixture := newGreetingFixture(t, greetingFixtureInput{source: platform.PlatformKick})
			test.setup(fixture)

			// When
			err := fixture.handler.handleGreetings(context.Background(), fixture.message)

			// Then
			require.Error(t, err)
			require.Len(t, fixture.sender.requestSnapshot(), test.wantSends)
			require.Empty(t, fixture.repository.updates)
			require.Empty(t, fixture.cache.deletes)
		})
	}
}

func TestHandleGreetingsAttemptsShoutOutOnlyForTwitch(t *testing.T) {
	for _, source := range []platform.Platform{platform.PlatformTwitch, platform.PlatformKick} {
		t.Run(source.String(), func(t *testing.T) {
			// Given
			fixture := newGreetingFixture(t, greetingFixtureInput{source: source, withShoutOut: true})
			calls := 0
			fixture.installTwitchActions(t, &calls)

			// When
			require.NoError(t, fixture.handler.handleGreetings(context.Background(), fixture.message))

			// Then
			if source == platform.PlatformTwitch {
				require.Equal(t, 1, calls)
			} else {
				require.Zero(t, calls)
			}
		})
	}
}
