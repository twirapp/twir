package kick

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-redis/redismock/v9"
	"github.com/google/uuid"
	"github.com/twirapp/twir/libs/bus-core/generic"
	kickbus "github.com/twirapp/twir/libs/bus-core/kick"
	"github.com/twirapp/twir/libs/entities/platform"
	channels "github.com/twirapp/twir/libs/repositories/channels"
	channelsmodel "github.com/twirapp/twir/libs/repositories/channels/model"
	usersrepository "github.com/twirapp/twir/libs/repositories/users"
	usersmodel "github.com/twirapp/twir/libs/repositories/users/model"
)

type mockChannelsRepoWebhook struct {
	channel channelsmodel.Channel
	err     error
}

func (m *mockChannelsRepoWebhook) GetByKickUserID(_ context.Context, _ uuid.UUID) (channelsmodel.Channel, error) {
	return m.channel, m.err
}

func (m *mockChannelsRepoWebhook) GetMany(_ context.Context, _ channels.GetManyInput) ([]channelsmodel.Channel, error) {
	return nil, nil
}

func (m *mockChannelsRepoWebhook) GetByID(_ context.Context, _ uuid.UUID) (channelsmodel.Channel, error) {
	return channelsmodel.Nil, nil
}

func (m *mockChannelsRepoWebhook) GetByTwitchUserID(_ context.Context, _ uuid.UUID) (channelsmodel.Channel, error) {
	return channelsmodel.Nil, nil
}

func (m *mockChannelsRepoWebhook) GetCount(_ context.Context, _ channels.GetCountInput) (int, error) {
	return 0, nil
}

func (m *mockChannelsRepoWebhook) Update(_ context.Context, _ uuid.UUID, _ channels.UpdateInput) (channelsmodel.Channel, error) {
	return channelsmodel.Nil, nil
}

func (m *mockChannelsRepoWebhook) Create(_ context.Context, _ channels.CreateInput) (channelsmodel.Channel, error) {
	return channelsmodel.Nil, nil
}

type mockUsersRepoWebhook struct {
	user usersmodel.User
	err  error
}

func (m *mockUsersRepoWebhook) GetByID(_ context.Context, _ uuid.UUID) (usersmodel.User, error) {
	return m.user, m.err
}

func (m *mockUsersRepoWebhook) GetByPlatformID(_ context.Context, _ platform.Platform, _ string) (usersmodel.User, error) {
	return m.user, m.err
}

func (m *mockUsersRepoWebhook) GetManyByIDS(_ context.Context, _ usersrepository.GetManyInput) ([]usersmodel.User, error) {
	return nil, nil
}

func (m *mockUsersRepoWebhook) Update(_ context.Context, _ uuid.UUID, _ usersrepository.UpdateInput) (usersmodel.User, error) {
	return usersmodel.Nil, nil
}

func (m *mockUsersRepoWebhook) GetRandomOnlineUser(_ context.Context, _ usersrepository.GetRandomOnlineUserInput) (usersmodel.OnlineUser, error) {
	return usersmodel.NilOnlineUser, nil
}

func (m *mockUsersRepoWebhook) GetOnlineUsersWithFilters(_ context.Context, _ usersrepository.GetOnlineUsersWithFiltersInput) ([]usersmodel.OnlineUser, error) {
	return nil, nil
}

func (m *mockUsersRepoWebhook) GetByApiKey(_ context.Context, _ string) (usersmodel.User, error) {
	return usersmodel.Nil, usersmodel.ErrNotFound
}

func (m *mockUsersRepoWebhook) Create(_ context.Context, _ usersrepository.CreateInput) (usersmodel.User, error) {
	return usersmodel.Nil, nil
}

func TestWebhookHandler_ChatMessage(t *testing.T) {
	redisClient, redisMock := redismock.NewClientMock()

	channelID := uuid.MustParse("11111111-1111-1111-1111-111111111111")
	kickUserID := "103507661"

	usersRepo := &mockUsersRepoWebhook{
		user: usersmodel.User{
		ID:         channelID,
			PlatformID: kickUserID,
		},
	}

	channelsRepo := &mockChannelsRepoWebhook{
		channel: channelsmodel.Channel{
			ID: channelID,
		},
	}

	handlers := &Handlers{
		logger:                slog.Default(),
		redis:                 redisClient,
		channelsRepo:          channelsRepo,
		usersRepo:             usersRepo,
		chatMessagesGeneric:   &mockQueue[generic.ChatMessage, struct{}]{},
		processGenericMessage: &mockQueue[generic.ChatMessage, struct{}]{},
	}

	payload := kickChatMessagePayload{
		MessageID: "msg-123",
		Broadcaster: kickUser{
			UserID:   103507661,
			Username: "testbroadcaster",
		},
		Sender: kickUser{
			UserID:   999,
			Username: "testsender",
		},
		Content:   "Hello world",
		CreatedAt: time.Now().UTC().Format(time.RFC3339),
	}

	body, _ := json.Marshal(payload)
	timestamp := time.Now().UTC().Format(time.RFC3339)

	redisMock.ExpectSetNX("kick:event:msg-123", "processing", 30*time.Second).SetVal(true)
	redisMock.ExpectSet("kick:event:msg-123", "processed", 10*time.Minute).SetVal("OK")

	req := httptest.NewRequest(http.MethodPost, "/webhook/kick", bytes.NewReader(body))
	req.Header.Set("Kick-Event-Message-Id", "msg-123")
	req.Header.Set("Kick-Event-Message-Timestamp", timestamp)
	req.Header.Set("Kick-Event-Type", "chat.message.sent")
	req.Header.Set("Kick-Event-Version", "1")
	req.Header.Set("Kick-Event-Subscription-Id", "sub-123")

	w := httptest.NewRecorder()

	middleware := NewMiddleware(slog.Default())
	middleware.HandlerWithoutVerification(http.HandlerFunc(handlers.HandleWebhook)).ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", w.Code, w.Body.String())
	}

	if err := redisMock.ExpectationsWereMet(); err != nil {
		t.Fatalf("redis expectations not met: %v", err)
	}
}

func TestWebhookHandler_LivestreamStatus(t *testing.T) {
	redisClient, redisMock := redismock.NewClientMock()

	channelID := uuid.MustParse("22222222-2222-2222-2222-222222222222")
	kickUserID := "103507661"

	usersRepo := &mockUsersRepoWebhook{
		user: usersmodel.User{
		ID:         channelID,
			PlatformID: kickUserID,
		},
	}

	channelsRepo := &mockChannelsRepoWebhook{
		channel: channelsmodel.Channel{
			ID: channelID,
		},
	}

	handlers := &Handlers{
		logger:        slog.Default(),
		redis:         redisClient,
		channelsRepo:  channelsRepo,
		usersRepo:     usersRepo,
		streamOnline:  &mockQueue[kickbus.KickStreamOnline, struct{}]{},
		streamOffline: &mockQueue[kickbus.KickStreamOffline, struct{}]{},
	}

	payload := kickLivestreamStatusPayload{
		Broadcaster: kickUser{
			UserID:   103507661,
			Username: "testbroadcaster",
		},
		IsLive: true,
		Title:  "Test Stream",
	}

	body, _ := json.Marshal(payload)
	timestamp := time.Now().UTC().Format(time.RFC3339)

	redisMock.ExpectSetNX("kick:event:stream-123", "processing", 30*time.Second).SetVal(true)
	redisMock.ExpectSet("kick:event:stream-123", "processed", 10*time.Minute).SetVal("OK")

	req := httptest.NewRequest(http.MethodPost, "/webhook/kick", bytes.NewReader(body))
	req.Header.Set("Kick-Event-Message-Id", "stream-123")
	req.Header.Set("Kick-Event-Message-Timestamp", timestamp)
	req.Header.Set("Kick-Event-Type", "livestream.status.updated")
	req.Header.Set("Kick-Event-Version", "1")
	req.Header.Set("Kick-Event-Subscription-Id", "sub-456")

	w := httptest.NewRecorder()

	middleware := NewMiddleware(slog.Default())
	middleware.HandlerWithoutVerification(http.HandlerFunc(handlers.HandleWebhook)).ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", w.Code, w.Body.String())
	}

	if err := redisMock.ExpectationsWereMet(); err != nil {
		t.Fatalf("redis expectations not met: %v", err)
	}
}
