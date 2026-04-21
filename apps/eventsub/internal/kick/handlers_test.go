package kick

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-redis/redismock/v9"
	"github.com/google/uuid"
	buscore "github.com/twirapp/twir/libs/bus-core"
	"github.com/twirapp/twir/libs/bus-core/events"
	"github.com/twirapp/twir/libs/bus-core/generic"
	"github.com/twirapp/twir/libs/entities/platform"
	channelsrepository "github.com/twirapp/twir/libs/repositories/channels"
	channelsmodel "github.com/twirapp/twir/libs/repositories/channels/model"
	usersrepository "github.com/twirapp/twir/libs/repositories/users"
	usersmodel "github.com/twirapp/twir/libs/repositories/users/model"
)

type mockQueue[Req, Res any] struct {
	published []Req
}

func (m *mockQueue[Req, Res]) Publish(_ context.Context, data Req) error {
	m.published = append(m.published, data)
	return nil
}

func (m *mockQueue[Req, Res]) Request(_ context.Context, _ Req) (*buscore.QueueResponse[Res], error) {
	return nil, nil
}

func (m *mockQueue[Req, Res]) SubscribeGroup(_ string, _ buscore.QueueSubscribeCallback[Req, Res]) error {
	return nil
}

func (m *mockQueue[Req, Res]) Subscribe(_ buscore.QueueSubscribeCallback[Req, Res]) error {
	return nil
}

func (m *mockQueue[Req, Res]) Unsubscribe() {}

type mockUsersRepo struct {
	user usersmodel.User
	err  error
}

func (m *mockUsersRepo) GetByID(_ context.Context, _ string) (usersmodel.User, error) {
	return m.user, m.err
}

func (m *mockUsersRepo) GetByPlatformID(_ context.Context, _ platform.Platform, _ string) (usersmodel.User, error) {
	return m.user, m.err
}

func (m *mockUsersRepo) GetManyByIDS(_ context.Context, _ usersrepository.GetManyInput) ([]usersmodel.User, error) {
	return nil, nil
}

func (m *mockUsersRepo) Update(_ context.Context, _ string, _ usersrepository.UpdateInput) (usersmodel.User, error) {
	return m.user, m.err
}

func (m *mockUsersRepo) GetRandomOnlineUser(_ context.Context, _ usersrepository.GetRandomOnlineUserInput) (usersmodel.OnlineUser, error) {
	return usersmodel.NilOnlineUser, nil
}

func (m *mockUsersRepo) GetOnlineUsersWithFilters(_ context.Context, _ usersrepository.GetOnlineUsersWithFiltersInput) ([]usersmodel.OnlineUser, error) {
	return nil, nil
}

func (m *mockUsersRepo) GetByApiKey(_ context.Context, _ string) (usersmodel.User, error) {
	return m.user, m.err
}

func (m *mockUsersRepo) Create(_ context.Context, _ usersrepository.CreateInput) (usersmodel.User, error) {
	return m.user, m.err
}

type mockChannelsRepo struct {
	channel  channelsmodel.Channel
	channels []channelsmodel.Channel
	err      error
}

func (m *mockChannelsRepo) GetMany(_ context.Context, _ channelsrepository.GetManyInput) ([]channelsmodel.Channel, error) {
	return m.channels, m.err
}

func (m *mockChannelsRepo) GetByID(_ context.Context, _ uuid.UUID) (channelsmodel.Channel, error) {
	return m.channel, m.err
}

func (m *mockChannelsRepo) GetByTwitchUserID(_ context.Context, _ uuid.UUID) (channelsmodel.Channel, error) {
	return m.channel, m.err
}

func (m *mockChannelsRepo) GetByKickUserID(_ context.Context, _ uuid.UUID) (channelsmodel.Channel, error) {
	return m.channel, m.err
}

func (m *mockChannelsRepo) GetCount(_ context.Context, _ channelsrepository.GetCountInput) (int, error) {
	return 0, nil
}

func (m *mockChannelsRepo) Update(_ context.Context, _ uuid.UUID, _ channelsrepository.UpdateInput) (channelsmodel.Channel, error) {
	return m.channel, m.err
}

func (m *mockChannelsRepo) Create(_ context.Context, _ channelsrepository.CreateInput) (channelsmodel.Channel, error) {
	return m.channel, m.err
}

func buildTestHandlers(
	t *testing.T,
	chatMessagesGeneric *mockQueue[generic.ChatMessage, struct{}],
	processGenericMessage *mockQueue[generic.ChatMessage, struct{}],
	followQueue *mockQueue[events.FollowMessage, struct{}],
	usersRepo usersrepository.Repository,
	channelsRepo channelsrepository.Repository,
) (*Handlers, redismock.ClientMock) {
	t.Helper()

	db, redisMock := redismock.NewClientMock()

	h := &Handlers{
		logger:                slog.Default(),
		redis:                 db,
		chatMessagesGeneric:   chatMessagesGeneric,
		processGenericMessage: processGenericMessage,
		eventsFollow:          followQueue,
		channelsRepo:          channelsRepo,
		usersRepo:             usersRepo,
	}

	return h, redisMock
}

func makeRequest(t *testing.T, messageID, eventType string, payload any) *http.Request {
	t.Helper()
	body, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("marshal payload: %v", err)
	}
	req := httptest.NewRequest(http.MethodPost, "/webhook", bytes.NewReader(body))
	ctx := req.Context()
	ctx = context.WithValue(ctx, ctxKeyMessageID, messageID)
	ctx = context.WithValue(ctx, ctxKeyEventType, eventType)
	return req.WithContext(ctx)
}

func TestHandleChatMessage(t *testing.T) {
	userID := uuid.New().String()
	channelUUID := uuid.New()

	chatQueue := &mockQueue[generic.ChatMessage, struct{}]{}
	parserQueue := &mockQueue[generic.ChatMessage, struct{}]{}
	followQueue := &mockQueue[events.FollowMessage, struct{}]{}

	kickUserUUID := uuid.New()
	usersRepo := &mockUsersRepo{
		user: usersmodel.User{
			ID:         userID,
			PlatformID: "123",
		},
	}
	channelsRepo := &mockChannelsRepo{
		channel: channelsmodel.Channel{
			ID:         channelUUID,
			KickUserID: &kickUserUUID,
		},
	}

	h, redisMock := buildTestHandlers(t, chatQueue, parserQueue, followQueue, usersRepo, channelsRepo)

	msgID := "msg-001"
	redisMock.ExpectSetNX(idempotencyKeyPrefix+msgID, "1", idempotencyTTL).SetVal(true)

	payload := kickChatMessagePayload{
		MessageID: msgID,
		Broadcaster: kickUser{
			UserID:   123,
			Username: "broadcaster123",
		},
		Sender: kickUser{
			UserID:   456,
			Username: "senderlogin",
			Identity: &kickIdentity{
				UsernameColor: "#ff0000",
				Badges: []kickBadge{
					{Text: "Moderator", Type: "moderator"},
				},
			},
		},
		Content: "Hello world",
	}

	req := makeRequest(t, msgID, "chat.message.sent", payload)
	w := httptest.NewRecorder()

	h.HandleWebhook(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	if len(chatQueue.published) != 1 {
		t.Fatalf("expected 1 chat message published, got %d", len(chatQueue.published))
	}
	msg := chatQueue.published[0]
	if msg.Platform != string(platform.PlatformKick) {
		t.Errorf("expected platform %q, got %q", platform.PlatformKick, msg.Platform)
	}
	if msg.ChannelID != channelUUID.String() {
		t.Errorf("expected channelID %q, got %q", channelUUID.String(), msg.ChannelID)
	}
	if msg.Text != "Hello world" {
		t.Errorf("expected text %q, got %q", "Hello world", msg.Text)
	}
	if msg.MessageID != msgID {
		t.Errorf("expected messageID %q, got %q", msgID, msg.MessageID)
	}

	if len(parserQueue.published) != 1 {
		t.Fatalf("expected 1 parser message published, got %d", len(parserQueue.published))
	}

	if err := redisMock.ExpectationsWereMet(); err != nil {
		t.Errorf("redis expectations not met: %v", err)
	}
}

func TestHandleChatMessageIdempotency(t *testing.T) {
	userID := uuid.New().String()
	channelUUID := uuid.New()

	chatQueue := &mockQueue[generic.ChatMessage, struct{}]{}
	parserQueue := &mockQueue[generic.ChatMessage, struct{}]{}
	followQueue := &mockQueue[events.FollowMessage, struct{}]{}

	kickUserUUID := uuid.New()
	usersRepo := &mockUsersRepo{
		user: usersmodel.User{
			ID:         userID,
			PlatformID: "999",
		},
	}
	channelsRepo := &mockChannelsRepo{
		channel: channelsmodel.Channel{
			ID:         channelUUID,
			KickUserID: &kickUserUUID,
		},
	}

	h, redisMock := buildTestHandlers(t, chatQueue, parserQueue, followQueue, usersRepo, channelsRepo)

	msgID := "dup-msg-001"
	redisMock.ExpectSetNX(idempotencyKeyPrefix+msgID, "1", idempotencyTTL).SetVal(false)

	payload := kickChatMessagePayload{
		MessageID: msgID,
		Broadcaster: kickUser{
			UserID:   999,
			Username: "broadcaster999",
		},
		Content: "duplicate",
	}

	req := makeRequest(t, msgID, "chat.message.sent", payload)
	w := httptest.NewRecorder()

	h.HandleWebhook(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	if len(chatQueue.published) != 0 {
		t.Errorf("expected 0 published messages for duplicate, got %d", len(chatQueue.published))
	}

	if err := redisMock.ExpectationsWereMet(); err != nil {
		t.Errorf("redis expectations not met: %v", err)
	}
}

func TestHandleChannelFollow(t *testing.T) {
	userID := uuid.New().String()
	channelUUID := uuid.New()

	chatQueue := &mockQueue[generic.ChatMessage, struct{}]{}
	parserQueue := &mockQueue[generic.ChatMessage, struct{}]{}
	followQueue := &mockQueue[events.FollowMessage, struct{}]{}

	kickUserUUID := uuid.New()
	usersRepo := &mockUsersRepo{
		user: usersmodel.User{
			ID:         userID,
			PlatformID: "777",
		},
	}
	channelsRepo := &mockChannelsRepo{
		channel: channelsmodel.Channel{
			ID:         channelUUID,
			KickUserID: &kickUserUUID,
		},
	}

	h, redisMock := buildTestHandlers(t, chatQueue, parserQueue, followQueue, usersRepo, channelsRepo)

	msgID := "follow-evt-001"
	redisMock.ExpectSetNX(idempotencyKeyPrefix+msgID, "1", idempotencyTTL).SetVal(true)

	payload := kickFollowPayload{
		Broadcaster: kickUser{
			UserID:   777,
			Username: "broadcaster777",
		},
		Follower: kickUser{
			UserID:   111,
			Username: "followerlogin",
		},
	}

	req := makeRequest(t, msgID, "channel.followed", payload)
	w := httptest.NewRecorder()

	h.HandleWebhook(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	if len(followQueue.published) != 1 {
		t.Fatalf("expected 1 follow event published, got %d", len(followQueue.published))
	}
	follow := followQueue.published[0]
	if follow.BaseInfo.ChannelID != channelUUID.String() {
		t.Errorf("expected channelID %q, got %q", channelUUID.String(), follow.BaseInfo.ChannelID)
	}
	if follow.UserName != "followerlogin" {
		t.Errorf("expected UserName %q, got %q", "followerlogin", follow.UserName)
	}
	if follow.UserID != "111" {
		t.Errorf("expected UserID %q, got %q", "111", follow.UserID)
	}

	if err := redisMock.ExpectationsWereMet(); err != nil {
		t.Errorf("redis expectations not met: %v", err)
	}
}
