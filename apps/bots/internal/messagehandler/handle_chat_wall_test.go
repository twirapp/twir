package messagehandler

import (
	"context"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/google/uuid"
	"github.com/twirapp/twir/apps/bots/internal/twitchactions"
	buscore "github.com/twirapp/twir/libs/bus-core"
	"github.com/twirapp/twir/libs/bus-core/generic"
	buscoretokens "github.com/twirapp/twir/libs/bus-core/tokens"
	genericcacher "github.com/twirapp/twir/libs/cache/generic-cacher"
	cfg "github.com/twirapp/twir/libs/config"
	channelentity "github.com/twirapp/twir/libs/entities/channel"
	chatwallmodel "github.com/twirapp/twir/libs/repositories/chat_wall/model"
	usersmodel "github.com/twirapp/twir/libs/repositories/users/model"
)

func TestHandleChatWall_UsesProviderMessageIDForDeleteAndDedup(t *testing.T) {
	// Given
	deletedMessageIDs := make([]string, 0, 1)
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		deletedMessageIDs = append(deletedMessageIDs, request.URL.Query().Get("message_id"))
		writer.Header().Set("Content-Type", "application/json")
		_, _ = writer.Write([]byte(`{"data":[]}`))
	}))
	t.Cleanup(server.Close)

	redisClient, _ := newMessageIDRedisClient(t)
	bus := buscore.NewNatsBus(nil)
	bus.Tokens.RequestBotToken = &chatMessageEmoteQueue[
		buscoretokens.GetBotTokenRequest,
		buscoretokens.TokenResponse,
	]{response: &buscore.QueueResponse[buscoretokens.TokenResponse]{
		Data: buscoretokens.TokenResponse{AccessToken: "test-token"},
	}}
	channelID := uuid.New()
	handler := &MessageHandler{
		redis: redisClient,
		twitchActions: twitchactions.New(twitchactions.Opts{
			Config: cfg.Config{
				TwitchMockEnabled: true,
				TwitchMockApiUrl:  server.URL,
				TwitchClientId:    "test-client",
			},
			TwirBus: bus,
			Redis:   redisClient,
		}),
		chatWallCacher: genericcacher.New(genericcacher.Opts[[]chatwallmodel.ChatWall]{
			KV: messageIDCache{},
			LoadFn: func(context.Context, string) ([]chatwallmodel.ChatWall, error) {
				return []chatwallmodel.ChatWall{{
					ID:      uuid.New(),
					Enabled: true,
					Phrase:  "blocked",
					Action:  chatwallmodel.ChatWallActionDelete,
				}}, nil
			},
		}),
		chatWallSettingsCacher: genericcacher.New(genericcacher.Opts[chatwallmodel.ChatWallSettings]{
			KV: messageIDCache{},
			LoadFn: func(context.Context, string) (chatwallmodel.ChatWallSettings, error) {
				return chatwallmodel.ChatWallSettings{}, nil
			},
		}),
		chatWallRepository: messageIDChatWallRepository{},
	}
	message := enrichedChatMessage{
		ChatMessage: generic.ChatMessage{
			ID:                uuid.NewString(),
			MessageID:         "123456789",
			BroadcasterUserId: "channel-provider-id",
			ChatterUserId:     "chatter-provider-id",
			Message:           &generic.ChatMessageMessage{Text: "blocked message"},
		},
		EnrichedData: chatMessageEnrichedData{
			DbChannel:     channelentity.Channel{ID: channelID},
			DbUser:        &usersmodel.User{ID: uuid.New()},
			BotPlatformID: "bot-provider-id",
		},
	}

	// When
	if err := handler.handleChatWall(context.Background(), message); err != nil {
		t.Fatalf("handle chat wall: %v", err)
	}
	message.ID = uuid.NewString()
	if err := handler.handleChatWall(context.Background(), message); err != nil {
		t.Fatalf("handle duplicate chat wall message: %v", err)
	}

	// Then
	if !reflect.DeepEqual(deletedMessageIDs, []string{message.MessageID}) {
		t.Fatalf("deleted message IDs = %#v, want %#v", deletedMessageIDs, []string{message.MessageID})
	}
}
