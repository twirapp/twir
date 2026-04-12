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
	entity "github.com/twirapp/twir/libs/entities/user_platform_account"
	channelsrepository "github.com/twirapp/twir/libs/repositories/channels"
	channelsmodel "github.com/twirapp/twir/libs/repositories/channels/model"
	user_platform_accounts "github.com/twirapp/twir/libs/repositories/user_platform_accounts"
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

type mockUserPlatformAccountsRepo struct {
	account entity.UserPlatformAccount
	err     error
}

func (m *mockUserPlatformAccountsRepo) GetByUserIDAndPlatform(_ context.Context, _ uuid.UUID, _ platform.Platform) (entity.UserPlatformAccount, error) {
	return m.account, m.err
}

func (m *mockUserPlatformAccountsRepo) GetAllByUserID(_ context.Context, _ uuid.UUID) ([]entity.UserPlatformAccount, error) {
	return nil, nil
}

func (m *mockUserPlatformAccountsRepo) GetByPlatformUserID(_ context.Context, _ platform.Platform, _ string) (entity.UserPlatformAccount, error) {
	return m.account, m.err
}

func (m *mockUserPlatformAccountsRepo) Upsert(_ context.Context, _ user_platform_accounts.UpsertInput) (entity.UserPlatformAccount, error) {
	return entity.Nil, nil
}

func (m *mockUserPlatformAccountsRepo) Delete(_ context.Context, _ uuid.UUID) error {
	return nil
}

func (m *mockUserPlatformAccountsRepo) GetAllByPlatform(_ context.Context, _ platform.Platform) ([]entity.UserPlatformAccount, error) {
	return nil, nil
}

type mockChannelsRepo struct {
	channel channelsmodel.Channel
	err     error
}

func (m *mockChannelsRepo) GetMany(_ context.Context, _ channelsrepository.GetManyInput) ([]channelsmodel.Channel, error) {
	return nil, nil
}

func (m *mockChannelsRepo) GetByID(_ context.Context, _ string) (channelsmodel.Channel, error) {
	return m.channel, m.err
}

func (m *mockChannelsRepo) GetByUserIDAndPlatform(_ context.Context, _ uuid.UUID, _ platform.Platform) (channelsmodel.Channel, error) {
	return m.channel, m.err
}

func (m *mockChannelsRepo) GetCount(_ context.Context, _ channelsrepository.GetCountInput) (int, error) {
	return 0, nil
}

func buildTestHandlers(
	t *testing.T,
	chatMessagesGeneric *mockQueue[generic.ChatMessage, struct{}],
	processGenericMessage *mockQueue[generic.ChatMessage, struct{}],
	followQueue *mockQueue[events.FollowMessage, struct{}],
	userPlatformAccountsRepo user_platform_accounts.Repository,
	channelsRepo channelsrepository.Repository,
) (*Handlers, redismock.ClientMock) {
	t.Helper()

	db, redisMock := redismock.NewClientMock()

	h := &Handlers{
		logger:                   slog.Default(),
		redis:                    db,
		chatMessagesGeneric:      chatMessagesGeneric,
		processGenericMessage:    processGenericMessage,
		eventsFollow:             followQueue,
		channelsRepo:             channelsRepo,
		userPlatformAccountsRepo: userPlatformAccountsRepo,
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
	channelUUID := uuid.New()
	userUUID := uuid.New()
	channelID := uuid.New().String()

	chatQueue := &mockQueue[generic.ChatMessage, struct{}]{}
	parserQueue := &mockQueue[generic.ChatMessage, struct{}]{}
	followQueue := &mockQueue[events.FollowMessage, struct{}]{}

	accountsRepo := &mockUserPlatformAccountsRepo{
		account: entity.UserPlatformAccount{
			ID:             channelUUID,
			UserID:         userUUID,
			Platform:       platform.PlatformKick,
			PlatformUserID: "broadcaster-123",
		},
	}
	channelsRepo := &mockChannelsRepo{
		channel: channelsmodel.Channel{
			ID:     channelID,
			UserID: userUUID,
		},
	}

	h, redisMock := buildTestHandlers(t, chatQueue, parserQueue, followQueue, accountsRepo, channelsRepo)

	msgID := "msg-001"
	redisMock.ExpectSetNX(idempotencyKeyPrefix+msgID, "1", idempotencyTTL).SetVal(true)

	payload := kickChatMessagePayload{
		MessageID:         msgID,
		BroadcasterUserID: "broadcaster-123",
		SenderUserID:      "sender-456",
		SenderUserLogin:   "senderlogin",
		SenderDisplayName: "SenderDisplay",
		Content:           "Hello world",
		Color:             "#ff0000",
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
	if msg.ChannelID != channelID {
		t.Errorf("expected channelID %q, got %q", channelID, msg.ChannelID)
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
	userUUID := uuid.New()
	channelID := uuid.New().String()

	chatQueue := &mockQueue[generic.ChatMessage, struct{}]{}
	parserQueue := &mockQueue[generic.ChatMessage, struct{}]{}
	followQueue := &mockQueue[events.FollowMessage, struct{}]{}

	accountsRepo := &mockUserPlatformAccountsRepo{
		account: entity.UserPlatformAccount{
			UserID:         userUUID,
			Platform:       platform.PlatformKick,
			PlatformUserID: "broadcaster-999",
		},
	}
	channelsRepo := &mockChannelsRepo{
		channel: channelsmodel.Channel{
			ID:     channelID,
			UserID: userUUID,
		},
	}

	h, redisMock := buildTestHandlers(t, chatQueue, parserQueue, followQueue, accountsRepo, channelsRepo)

	msgID := "dup-msg-001"
	redisMock.ExpectSetNX(idempotencyKeyPrefix+msgID, "1", idempotencyTTL).SetVal(false)

	payload := kickChatMessagePayload{
		MessageID:         msgID,
		BroadcasterUserID: "broadcaster-999",
		Content:           "duplicate",
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
	userUUID := uuid.New()
	channelID := uuid.New().String()

	chatQueue := &mockQueue[generic.ChatMessage, struct{}]{}
	parserQueue := &mockQueue[generic.ChatMessage, struct{}]{}
	followQueue := &mockQueue[events.FollowMessage, struct{}]{}

	accountsRepo := &mockUserPlatformAccountsRepo{
		account: entity.UserPlatformAccount{
			UserID:         userUUID,
			Platform:       platform.PlatformKick,
			PlatformUserID: "broadcaster-777",
		},
	}
	channelsRepo := &mockChannelsRepo{
		channel: channelsmodel.Channel{
			ID:     channelID,
			UserID: userUUID,
		},
	}

	h, redisMock := buildTestHandlers(t, chatQueue, parserQueue, followQueue, accountsRepo, channelsRepo)

	msgID := "follow-evt-001"
	redisMock.ExpectSetNX(idempotencyKeyPrefix+msgID, "1", idempotencyTTL).SetVal(true)

	payload := kickFollowPayload{
		BroadcasterUserID: "broadcaster-777",
		FollowerUserID:    "follower-111",
		FollowerUserLogin: "followerlogin",
	}

	req := makeRequest(t, msgID, "channel.follow", payload)
	w := httptest.NewRecorder()

	h.HandleWebhook(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	if len(followQueue.published) != 1 {
		t.Fatalf("expected 1 follow event published, got %d", len(followQueue.published))
	}
	follow := followQueue.published[0]
	if follow.BaseInfo.ChannelID != channelID {
		t.Errorf("expected channelID %q, got %q", channelID, follow.BaseInfo.ChannelID)
	}
	if follow.UserName != "followerlogin" {
		t.Errorf("expected UserName %q, got %q", "followerlogin", follow.UserName)
	}
	if follow.UserID != "follower-111" {
		t.Errorf("expected UserID %q, got %q", "follower-111", follow.UserID)
	}

	if err := redisMock.ExpectationsWereMet(); err != nil {
		t.Errorf("redis expectations not met: %v", err)
	}
}
