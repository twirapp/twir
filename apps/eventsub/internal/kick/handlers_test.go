package kick

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/go-redis/redismock/v9"
	"github.com/google/uuid"
	buscore "github.com/twirapp/twir/libs/bus-core"
	"github.com/twirapp/twir/libs/bus-core/events"
	"github.com/twirapp/twir/libs/bus-core/generic"
	kickbus "github.com/twirapp/twir/libs/bus-core/kick"
	"github.com/twirapp/twir/libs/entities/platform"
	channelsrepository "github.com/twirapp/twir/libs/repositories/channels"
	channelsmodel "github.com/twirapp/twir/libs/repositories/channels/model"
	channelsinfohistory "github.com/twirapp/twir/libs/repositories/channels_info_history"
	channelsinfohistorymodel "github.com/twirapp/twir/libs/repositories/channels_info_history/model"
	streamsrepository "github.com/twirapp/twir/libs/repositories/streams"
	streamsmodel "github.com/twirapp/twir/libs/repositories/streams/model"
	usersrepository "github.com/twirapp/twir/libs/repositories/users"
	usersmodel "github.com/twirapp/twir/libs/repositories/users/model"
)

type mockQueue[Req, Res any] struct {
	mu          sync.Mutex
	published []Req
	publishErr  error
	publishHook func(context.Context, Req) error
}


func (m *mockQueue[Req, Res]) Publish(ctx context.Context, data Req) error {
	if m.publishHook != nil {
		if err := m.publishHook(ctx, data); err != nil {
			return err
		}
	}

	if m.publishErr != nil {
		return m.publishErr
	}

	m.mu.Lock()
	m.published = append(m.published, data)
	m.mu.Unlock()
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

func (m *mockQueue[Req, Res]) PublishedCount() int {
	m.mu.Lock()
	defer m.mu.Unlock()

	return len(m.published)
}

func (m *mockQueue[Req, Res]) FirstPublished() Req {
	m.mu.Lock()
	defer m.mu.Unlock()

	return m.published[0]
}

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

type mockStreamsRepo struct {
	stream      streamsmodel.Stream
	err         error
	saved       []streamsrepository.SaveInput
	updated     []streamsrepository.UpdateInput
	deletedByID []string
}

func (m *mockStreamsRepo) GetByChannelID(_ context.Context, _ string) (streamsmodel.Stream, error) {
	return m.stream, m.err
}

func (m *mockStreamsRepo) GetList(_ context.Context) ([]streamsmodel.Stream, error) {
	return nil, m.err
}

func (m *mockStreamsRepo) Save(_ context.Context, input streamsrepository.SaveInput) error {
	m.saved = append(m.saved, input)
	m.stream = streamsmodel.Stream{
		ID:           input.ID,
		UserId:       input.UserId,
		UserLogin:    input.UserLogin,
		UserName:     input.UserName,
		GameId:       input.GameId,
		GameName:     input.GameName,
		CommunityIds: input.CommunityIds,
		Type:         input.Type,
		Title:        input.Title,
		ViewerCount:  input.ViewerCount,
		StartedAt:    input.StartedAt,
		Language:     input.Language,
		ThumbnailUrl: input.ThumbnailUrl,
		TagIds:       input.TagIds,
		Tags:         input.Tags,
		IsMature:     input.IsMature,
	}
	return m.err
}

func (m *mockStreamsRepo) DeleteByChannelID(_ context.Context, channelID string) error {
	m.deletedByID = append(m.deletedByID, channelID)
	return m.err
}

func (m *mockStreamsRepo) Update(_ context.Context, _ string, input streamsrepository.UpdateInput) error {
	m.updated = append(m.updated, input)
	if input.Title != nil {
		m.stream.Title = *input.Title
	}
	if input.GameName != nil {
		m.stream.GameName = *input.GameName
	}
	if input.GameId != nil {
		m.stream.GameId = *input.GameId
	}
	if input.Language != nil {
		m.stream.Language = *input.Language
	}
	if input.ThumbnailUrl != nil {
		m.stream.ThumbnailUrl = *input.ThumbnailUrl
	}
	if input.IsMature != nil {
		m.stream.IsMature = *input.IsMature
	}
	if input.Tags != nil {
		m.stream.Tags = input.Tags
	}
	return m.err
}

type mockChannelsInfoHistoryRepo struct {
	created []channelsinfohistory.CreateInput
	err     error
}

func (m *mockChannelsInfoHistoryRepo) GetMany(_ context.Context, _ channelsinfohistory.GetManyInput) ([]channelsinfohistorymodel.ChannelInfoHistory, error) {
	return nil, m.err
}

func (m *mockChannelsInfoHistoryRepo) Create(_ context.Context, input channelsinfohistory.CreateInput) error {
	m.created = append(m.created, input)
	return m.err
}

func buildTestHandlers(
	t *testing.T,
	chatMessagesGeneric *mockQueue[generic.ChatMessage, struct{}],
	processGenericMessage *mockQueue[generic.ChatMessage, struct{}],
	followQueue *mockQueue[events.FollowMessage, struct{}],
	streamOnline *mockQueue[kickbus.KickStreamOnline, struct{}],
	streamOffline *mockQueue[kickbus.KickStreamOffline, struct{}],
	usersRepo usersrepository.Repository,
	channelsRepo channelsrepository.Repository,
	streamsRepo ...streamsrepository.Repository,
) (*Handlers, redismock.ClientMock) {
	t.Helper()

	db, redisMock := redismock.NewClientMock()

	var streamRepo streamsrepository.Repository
	var infoHistoryRepo channelsinfohistory.Repository
	if len(streamsRepo) > 0 && streamsRepo[0] != nil {
		streamRepo = streamsRepo[0]
	}
	infoHistoryRepo = &mockChannelsInfoHistoryRepo{}

	h := &Handlers{
		logger:                slog.Default(),
		redis:                 db,
		chatMessagesGeneric:   chatMessagesGeneric,
		processGenericMessage: processGenericMessage,
		eventsFollow:          followQueue,
		streamOnline:          streamOnline,
		streamOffline:         streamOffline,
		streamsRepo:           streamRepo,
		channelsInfoHistoryRepo: infoHistoryRepo,
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

	h, redisMock := buildTestHandlers(
		t,
		chatQueue,
		parserQueue,
		followQueue,
		&mockQueue[kickbus.KickStreamOnline, struct{}]{},
		&mockQueue[kickbus.KickStreamOffline, struct{}]{},
		usersRepo,
		channelsRepo,
	)

	msgID := "msg-001"
	redisMock.ExpectSetNX(idempotencyKeyPrefix+msgID, idempotencyStatusProcessing, idempotencyProcessingTTL).SetVal(true)
	redisMock.ExpectSet(idempotencyKeyPrefix+msgID, idempotencyStatusProcessed, idempotencyTTL).SetVal("OK")

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

	if chatQueue.PublishedCount() != 1 {
		t.Fatalf("expected 1 chat message published, got %d", chatQueue.PublishedCount())
	}
	msg := chatQueue.FirstPublished()
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

	if parserQueue.PublishedCount() != 1 {
		t.Fatalf("expected 1 parser message published, got %d", parserQueue.PublishedCount())
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

	h, redisMock := buildTestHandlers(
		t,
		chatQueue,
		parserQueue,
		followQueue,
		&mockQueue[kickbus.KickStreamOnline, struct{}]{},
		&mockQueue[kickbus.KickStreamOffline, struct{}]{},
		usersRepo,
		channelsRepo,
	)

	msgID := "dup-msg-001"
	redisMock.ExpectSetNX(idempotencyKeyPrefix+msgID, idempotencyStatusProcessing, idempotencyProcessingTTL).SetVal(false)
	redisMock.ExpectGet(idempotencyKeyPrefix + msgID).SetVal(idempotencyStatusProcessed)

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

	if chatQueue.PublishedCount() != 0 {
		t.Errorf("expected 0 published messages for duplicate, got %d", chatQueue.PublishedCount())
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

	h, redisMock := buildTestHandlers(
		t,
		chatQueue,
		parserQueue,
		followQueue,
		&mockQueue[kickbus.KickStreamOnline, struct{}]{},
		&mockQueue[kickbus.KickStreamOffline, struct{}]{},
		usersRepo,
		channelsRepo,
	)

	msgID := "follow-evt-001"
	redisMock.ExpectSetNX(idempotencyKeyPrefix+msgID, idempotencyStatusProcessing, idempotencyProcessingTTL).SetVal(true)
	redisMock.ExpectSet(idempotencyKeyPrefix+msgID, idempotencyStatusProcessed, idempotencyTTL).SetVal("OK")

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

	if followQueue.PublishedCount() != 1 {
		t.Fatalf("expected 1 follow event published, got %d", followQueue.PublishedCount())
	}
	follow := followQueue.FirstPublished()
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

func TestHandleLivestreamStatusOnline(t *testing.T) {
	userID := uuid.New().String()

	chatQueue := &mockQueue[generic.ChatMessage, struct{}]{}
	parserQueue := &mockQueue[generic.ChatMessage, struct{}]{}
	followQueue := &mockQueue[events.FollowMessage, struct{}]{}
	streamOnlineQueue := &mockQueue[kickbus.KickStreamOnline, struct{}]{}
	streamOfflineQueue := &mockQueue[kickbus.KickStreamOffline, struct{}]{}
	streamsRepo := &mockStreamsRepo{}

	kickUserUUID := uuid.New()
	usersRepo := &mockUsersRepo{
		user: usersmodel.User{
			ID:         userID,
			PlatformID: "555",
		},
	}
	channelsRepo := &mockChannelsRepo{
		channel: channelsmodel.Channel{
			ID:         uuid.New(),
			KickUserID: &kickUserUUID,
		},
	}

	h, redisMock := buildTestHandlers(
		t,
		chatQueue,
		parserQueue,
		followQueue,
		streamOnlineQueue,
		streamOfflineQueue,
		usersRepo,
		channelsRepo,
		streamsRepo,
	)

	msgID := "stream-online-001"
	redisMock.ExpectSetNX(idempotencyKeyPrefix+msgID, idempotencyStatusProcessing, idempotencyProcessingTTL).SetVal(true)
	redisMock.ExpectSet(idempotencyKeyPrefix+msgID, idempotencyStatusProcessed, idempotencyTTL).SetVal("OK")

	payload := kickLivestreamStatusPayload{
		Broadcaster: kickUser{
			UserID:   555,
			Username: "broadcaster555",
		},
		IsLive: true,
		Title:  "Going live",
	}

	req := makeRequest(t, msgID, "livestream.status.updated", payload)
	w := httptest.NewRecorder()

	h.HandleWebhook(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	if streamOnlineQueue.PublishedCount() != 1 {
		t.Fatalf("expected 1 stream online event published, got %d", streamOnlineQueue.PublishedCount())
	}

	if streamOfflineQueue.PublishedCount() != 0 {
		t.Fatalf("expected 0 stream offline events published, got %d", streamOfflineQueue.PublishedCount())
	}

	event := streamOnlineQueue.FirstPublished()
	if event.BroadcasterUserID != "555" {
		t.Errorf("expected broadcaster user id %q, got %q", "555", event.BroadcasterUserID)
	}
	if event.BroadcasterUserLogin != "broadcaster555" {
		t.Errorf("expected broadcaster login %q, got %q", "broadcaster555", event.BroadcasterUserLogin)
	}
	if len(streamsRepo.saved) != 1 {
		t.Fatalf("expected 1 stream save, got %d", len(streamsRepo.saved))
	}
	if streamsRepo.saved[0].Title != "Going live" {
		t.Fatalf("expected saved title %q, got %q", "Going live", streamsRepo.saved[0].Title)
	}

	if err := redisMock.ExpectationsWereMet(); err != nil {
		t.Errorf("redis expectations not met: %v", err)
	}
}

func TestHandleLivestreamStatusOffline(t *testing.T) {
	userID := uuid.New().String()

	chatQueue := &mockQueue[generic.ChatMessage, struct{}]{}
	parserQueue := &mockQueue[generic.ChatMessage, struct{}]{}
	followQueue := &mockQueue[events.FollowMessage, struct{}]{}
	streamOnlineQueue := &mockQueue[kickbus.KickStreamOnline, struct{}]{}
	streamOfflineQueue := &mockQueue[kickbus.KickStreamOffline, struct{}]{}
	streamsRepo := &mockStreamsRepo{}

	kickUserUUID := uuid.New()
	usersRepo := &mockUsersRepo{
		user: usersmodel.User{
			ID:         userID,
			PlatformID: "556",
		},
	}
	channelsRepo := &mockChannelsRepo{
		channel: channelsmodel.Channel{
			ID:         uuid.New(),
			KickUserID: &kickUserUUID,
		},
	}

	h, redisMock := buildTestHandlers(
		t,
		chatQueue,
		parserQueue,
		followQueue,
		streamOnlineQueue,
		streamOfflineQueue,
		usersRepo,
		channelsRepo,
		streamsRepo,
	)

	msgID := "stream-offline-001"
	redisMock.ExpectSetNX(idempotencyKeyPrefix+msgID, idempotencyStatusProcessing, idempotencyProcessingTTL).SetVal(true)
	redisMock.ExpectSet(idempotencyKeyPrefix+msgID, idempotencyStatusProcessed, idempotencyTTL).SetVal("OK")

	payload := kickLivestreamStatusPayload{
		Broadcaster: kickUser{
			UserID:   556,
			Username: "broadcaster556",
		},
		IsLive: false,
		Title:  "Wrapping up",
	}

	req := makeRequest(t, msgID, "livestream.status.updated", payload)
	w := httptest.NewRecorder()

	h.HandleWebhook(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	if streamOfflineQueue.PublishedCount() != 1 {
		t.Fatalf("expected 1 stream offline event published, got %d", streamOfflineQueue.PublishedCount())
	}

	if streamOnlineQueue.PublishedCount() != 0 {
		t.Fatalf("expected 0 stream online events published, got %d", streamOnlineQueue.PublishedCount())
	}

	event := streamOfflineQueue.FirstPublished()
	if event.BroadcasterUserID != "556" {
		t.Errorf("expected broadcaster user id %q, got %q", "556", event.BroadcasterUserID)
	}
	if event.BroadcasterUserLogin != "broadcaster556" {
		t.Errorf("expected broadcaster login %q, got %q", "broadcaster556", event.BroadcasterUserLogin)
	}
	if len(streamsRepo.deletedByID) != 1 {
		t.Fatalf("expected 1 stream delete, got %d", len(streamsRepo.deletedByID))
	}

	if err := redisMock.ExpectationsWereMet(); err != nil {
		t.Errorf("redis expectations not met: %v", err)
	}
}

func TestHandleLivestreamMetadataUpdated(t *testing.T) {
	chatQueue := &mockQueue[generic.ChatMessage, struct{}]{}
	parserQueue := &mockQueue[generic.ChatMessage, struct{}]{}
	followQueue := &mockQueue[events.FollowMessage, struct{}]{}
	streamOnlineQueue := &mockQueue[kickbus.KickStreamOnline, struct{}]{}
	streamOfflineQueue := &mockQueue[kickbus.KickStreamOffline, struct{}]{}
	streamsRepo := &mockStreamsRepo{}

	kickUserUUID := uuid.New()
	usersRepo := &mockUsersRepo{
		user: usersmodel.User{
			ID:         uuid.New().String(),
			PlatformID: "557",
		},
	}
	channelsRepo := &mockChannelsRepo{
		channel: channelsmodel.Channel{
			ID:         uuid.New(),
			KickUserID: &kickUserUUID,
		},
	}

	h, redisMock := buildTestHandlers(
		t,
		chatQueue,
		parserQueue,
		followQueue,
		streamOnlineQueue,
		streamOfflineQueue,
		usersRepo,
		channelsRepo,
		streamsRepo,
	)

	msgID := "stream-meta-001"
	redisMock.ExpectSetNX(idempotencyKeyPrefix+msgID, idempotencyStatusProcessing, idempotencyProcessingTTL).SetVal(true)
	redisMock.ExpectSet(idempotencyKeyPrefix+msgID, idempotencyStatusProcessed, idempotencyTTL).SetVal("OK")

	payload := map[string]any{
		"broadcaster": map[string]any{
			"user_id":      557,
			"username":     "broadcaster557",
			"channel_slug": "kick-slug",
		},
		"metadata": map[string]any{
			"title":              "New title",
			"language":           "en",
			"has_mature_content": false,
			"category": map[string]any{
				"id":        15,
				"name":      "Gaming",
				"thumbnail": "thumb",
			},
		},
	}

	req := makeRequest(t, msgID, "livestream.metadata.updated", payload)
	w := httptest.NewRecorder()

	h.HandleWebhook(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	if streamOnlineQueue.PublishedCount() != 0 {
		t.Fatalf("expected 0 stream online events published, got %d", streamOnlineQueue.PublishedCount())
	}

	if streamOfflineQueue.PublishedCount() != 0 {
		t.Fatalf("expected 0 stream offline events published, got %d", streamOfflineQueue.PublishedCount())
	}
	if len(streamsRepo.updated) != 1 {
		t.Fatalf("expected 1 stream update, got %d", len(streamsRepo.updated))
	}
	if streamsRepo.stream.GameName != "Gaming" {
		t.Fatalf("expected updated game name %q, got %q", "Gaming", streamsRepo.stream.GameName)
	}

	if err := redisMock.ExpectationsWereMet(); err != nil {
		t.Errorf("redis expectations not met: %v", err)
	}
}

func TestHandleChatMessageRealConcurrent(t *testing.T) {
	userID := uuid.New().String()
	channelUUID := uuid.New()

	publishStarted := make(chan struct{})
	releasePublish := make(chan struct{})
	var publishStartedOnce sync.Once

	chatQueue := &mockQueue[generic.ChatMessage, struct{}]{
		publishHook: func(_ context.Context, _ generic.ChatMessage) error {
			publishStartedOnce.Do(func() { close(publishStarted) })
			<-releasePublish
			return nil
		},
	}
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

	h, redisMock := buildTestHandlers(
		t,
		chatQueue,
		parserQueue,
		followQueue,
		&mockQueue[kickbus.KickStreamOnline, struct{}]{},
		&mockQueue[kickbus.KickStreamOffline, struct{}]{},
		usersRepo,
		channelsRepo,
	)
	redisMock.MatchExpectationsInOrder(false)

	msgID := "concurrent-msg-001"
	redisMock.ExpectSetNX(idempotencyKeyPrefix+msgID, idempotencyStatusProcessing, idempotencyProcessingTTL).SetVal(true)
	redisMock.ExpectSetNX(idempotencyKeyPrefix+msgID, idempotencyStatusProcessing, idempotencyProcessingTTL).SetVal(false)
	redisMock.ExpectGet(idempotencyKeyPrefix + msgID).SetVal(idempotencyStatusProcessing)
	redisMock.ExpectSet(idempotencyKeyPrefix+msgID, idempotencyStatusProcessed, idempotencyTTL).SetVal("OK")

	payload := kickChatMessagePayload{
		MessageID: msgID,
		Broadcaster: kickUser{
			UserID:   123,
			Username: "broadcaster123",
		},
		Sender: kickUser{
			UserID:   456,
			Username: "senderlogin",
		},
		Content: "Hello world",
	}

	start := make(chan struct{})
	statuses := make(chan int, 2)

	var wg sync.WaitGroup
	wg.Add(2)

	for range 2 {
		go func() {
			defer wg.Done()

			<-start
			req := makeRequest(t, msgID, "chat.message.sent", payload)
			w := httptest.NewRecorder()
			h.HandleWebhook(w, req)
			statuses <- w.Code
		}()
	}

	close(start)

	select {
	case <-publishStarted:
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for first request to start publishing")
	}

	var duplicateStatus int
	select {
	case duplicateStatus = <-statuses:
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for duplicate request to finish")
	}

	if duplicateStatus != http.StatusAccepted {
		t.Fatalf("expected duplicate request to return 202, got %d", duplicateStatus)
	}

	close(releasePublish)
	wg.Wait()
	close(statuses)

	var remainingStatuses []int
	for status := range statuses {
		remainingStatuses = append(remainingStatuses, status)
	}

	if len(remainingStatuses) != 1 || remainingStatuses[0] != http.StatusOK {
		t.Fatalf("expected processing request to return 200, got %v", remainingStatuses)
	}

	if chatQueue.PublishedCount() != 1 {
		t.Fatalf("expected 1 chat message published, got %d", chatQueue.PublishedCount())
	}

	if parserQueue.PublishedCount() != 1 {
		t.Fatalf("expected 1 parser message published, got %d", parserQueue.PublishedCount())
	}

	if err := redisMock.ExpectationsWereMet(); err != nil {
		t.Errorf("redis expectations not met: %v", err)
	}
}

func TestHandleChatMessageDuplicateWhileProcessing(t *testing.T) {
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

	h, redisMock := buildTestHandlers(
		t,
		chatQueue,
		parserQueue,
		followQueue,
		&mockQueue[kickbus.KickStreamOnline, struct{}]{},
		&mockQueue[kickbus.KickStreamOffline, struct{}]{},
		usersRepo,
		channelsRepo,
	)

	msgID := "duplicate-processing-001"
	redisMock.ExpectSetNX(idempotencyKeyPrefix+msgID, idempotencyStatusProcessing, idempotencyProcessingTTL).SetVal(false)
	redisMock.ExpectGet(idempotencyKeyPrefix + msgID).SetVal(idempotencyStatusProcessing)

	payload := kickChatMessagePayload{
		MessageID: msgID,
		Broadcaster: kickUser{
			UserID:   123,
			Username: "broadcaster123",
		},
		Sender: kickUser{
			UserID:   456,
			Username: "senderlogin",
		},
		Content: "Hello world",
	}

	req := makeRequest(t, msgID, "chat.message.sent", payload)
	w := httptest.NewRecorder()

	h.HandleWebhook(w, req)

	if w.Code != http.StatusAccepted {
		t.Fatalf("expected 202 for duplicate during processing, got %d", w.Code)
	}

	if chatQueue.PublishedCount() != 0 {
		t.Fatalf("expected 0 chat messages for duplicate, got %d", chatQueue.PublishedCount())
	}

	if err := redisMock.ExpectationsWereMet(); err != nil {
		t.Errorf("redis expectations not met: %v", err)
	}
}

func TestHandleChatMessageResolveIDsFailure(t *testing.T) {
	chatQueue := &mockQueue[generic.ChatMessage, struct{}]{}
	parserQueue := &mockQueue[generic.ChatMessage, struct{}]{}
	followQueue := &mockQueue[events.FollowMessage, struct{}]{}

	usersRepo := &mockUsersRepo{err: errors.New("users repo failed")}
	channelsRepo := &mockChannelsRepo{}

	h, redisMock := buildTestHandlers(
		t,
		chatQueue,
		parserQueue,
		followQueue,
		&mockQueue[kickbus.KickStreamOnline, struct{}]{},
		&mockQueue[kickbus.KickStreamOffline, struct{}]{},
		usersRepo,
		channelsRepo,
	)

	msgID := "resolve-failure-001"
	redisMock.ExpectSetNX(idempotencyKeyPrefix+msgID, idempotencyStatusProcessing, idempotencyProcessingTTL).SetVal(true)

	payload := kickChatMessagePayload{
		MessageID: msgID,
		Broadcaster: kickUser{
			UserID:   123,
			Username: "broadcaster123",
		},
		Sender: kickUser{
			UserID:   456,
			Username: "senderlogin",
		},
		Content: "Hello world",
	}

	req := makeRequest(t, msgID, "chat.message.sent", payload)
	w := httptest.NewRecorder()

	h.HandleWebhook(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", w.Code)
	}

	if chatQueue.PublishedCount() != 0 {
		t.Fatalf("expected 0 chat messages published, got %d", chatQueue.PublishedCount())
	}

	if parserQueue.PublishedCount() != 0 {
		t.Fatalf("expected 0 parser messages published, got %d", parserQueue.PublishedCount())
	}

	if err := redisMock.ExpectationsWereMet(); err != nil {
		t.Errorf("redis expectations not met: %v", err)
	}
}

func TestHandleChatMessageFailureCleanup(t *testing.T) {
	chatQueue := &mockQueue[generic.ChatMessage, struct{}]{}
	parserQueue := &mockQueue[generic.ChatMessage, struct{}]{}
	followQueue := &mockQueue[events.FollowMessage, struct{}]{}

	usersRepo := &mockUsersRepo{err: errors.New("users repo failed")}
	channelsRepo := &mockChannelsRepo{}

	h, redisMock := buildTestHandlers(
		t,
		chatQueue,
		parserQueue,
		followQueue,
		&mockQueue[kickbus.KickStreamOnline, struct{}]{},
		&mockQueue[kickbus.KickStreamOffline, struct{}]{},
		usersRepo,
		channelsRepo,
	)

	msgID := "cleanup-failure-001"
	redisMock.ExpectSetNX(idempotencyKeyPrefix+msgID, idempotencyStatusProcessing, idempotencyProcessingTTL).SetVal(true)
	redisMock.ExpectDel(idempotencyKeyPrefix + msgID).SetVal(1)
	redisMock.ExpectSetNX(idempotencyKeyPrefix+msgID, idempotencyStatusProcessing, idempotencyProcessingTTL).SetVal(true)
	redisMock.ExpectDel(idempotencyKeyPrefix + msgID).SetVal(1)

	payload := kickChatMessagePayload{
		MessageID: msgID,
		Broadcaster: kickUser{
			UserID:   123,
			Username: "broadcaster123",
		},
		Sender: kickUser{
			UserID:   456,
			Username: "senderlogin",
		},
		Content: "Hello world",
	}

	for range 2 {
		req := makeRequest(t, msgID, "chat.message.sent", payload)
		w := httptest.NewRecorder()

		h.HandleWebhook(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Fatalf("expected 500, got %d", w.Code)
		}
	}

	if chatQueue.PublishedCount() != 0 {
		t.Fatalf("expected 0 chat messages published, got %d", chatQueue.PublishedCount())
	}

	if parserQueue.PublishedCount() != 0 {
		t.Fatalf("expected 0 parser messages published, got %d", parserQueue.PublishedCount())
	}

	if err := redisMock.ExpectationsWereMet(); err != nil {
		t.Errorf("redis expectations not met: %v", err)
	}
}

func TestHandleChatMessagePublishFailure(t *testing.T) {
	userID := uuid.New().String()
	channelUUID := uuid.New()

	chatQueue := &mockQueue[generic.ChatMessage, struct{}]{
		publishErr: errors.New("publish failed"),
	}
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

	h, redisMock := buildTestHandlers(
		t,
		chatQueue,
		parserQueue,
		followQueue,
		&mockQueue[kickbus.KickStreamOnline, struct{}]{},
		&mockQueue[kickbus.KickStreamOffline, struct{}]{},
		usersRepo,
		channelsRepo,
	)

	msgID := "publish-failure-001"
	redisMock.ExpectSetNX(idempotencyKeyPrefix+msgID, idempotencyStatusProcessing, idempotencyProcessingTTL).SetVal(true)

	payload := kickChatMessagePayload{
		MessageID: msgID,
		Broadcaster: kickUser{
			UserID:   123,
			Username: "broadcaster123",
		},
		Sender: kickUser{
			UserID:   456,
			Username: "senderlogin",
		},
		Content: "Hello world",
	}

	req := makeRequest(t, msgID, "chat.message.sent", payload)
	w := httptest.NewRecorder()

	h.HandleWebhook(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", w.Code)
	}

	if chatQueue.PublishedCount() != 0 {
		t.Fatalf("expected 0 chat messages published, got %d", chatQueue.PublishedCount())
	}

	if parserQueue.PublishedCount() != 0 {
		t.Fatalf("expected 0 parser messages published, got %d", parserQueue.PublishedCount())
	}

	if err := redisMock.ExpectationsWereMet(); err != nil {
		t.Errorf("redis expectations not met: %v", err)
	}
}

func TestHandleChatMessagePartialSuccess(t *testing.T) {
	userID := uuid.New().String()
	channelUUID := uuid.New()

	chatQueue := &mockQueue[generic.ChatMessage, struct{}]{}
	parserQueue := &mockQueue[generic.ChatMessage, struct{}]{
		publishErr: errors.New("parser publish failed"),
	}
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

	h, redisMock := buildTestHandlers(
		t,
		chatQueue,
		parserQueue,
		followQueue,
		&mockQueue[kickbus.KickStreamOnline, struct{}]{},
		&mockQueue[kickbus.KickStreamOffline, struct{}]{},
		usersRepo,
		channelsRepo,
	)

	msgID := "partial-success-001"
	redisMock.ExpectSetNX(idempotencyKeyPrefix+msgID, idempotencyStatusProcessing, idempotencyProcessingTTL).SetVal(true)
	redisMock.ExpectDel(idempotencyKeyPrefix + msgID).SetVal(1)

	payload := kickChatMessagePayload{
		MessageID: msgID,
		Broadcaster: kickUser{
			UserID:   123,
			Username: "broadcaster123",
		},
		Sender: kickUser{
			UserID:   456,
			Username: "senderlogin",
		},
		Content: "Hello world",
	}

	req := makeRequest(t, msgID, "chat.message.sent", payload)
	w := httptest.NewRecorder()

	h.HandleWebhook(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", w.Code)
	}

	if chatQueue.PublishedCount() != 1 {
		t.Fatalf("expected 1 chat message published, got %d", chatQueue.PublishedCount())
	}

	if parserQueue.PublishedCount() != 0 {
		t.Fatalf("expected 0 parser messages published, got %d", parserQueue.PublishedCount())
	}

	if err := redisMock.ExpectationsWereMet(); err != nil {
		t.Errorf("redis expectations not met: %v", err)
	}
}

func TestHandleChatMessageMarkProcessedFailure(t *testing.T) {
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

	h, redisMock := buildTestHandlers(
		t,
		chatQueue,
		parserQueue,
		followQueue,
		&mockQueue[kickbus.KickStreamOnline, struct{}]{},
		&mockQueue[kickbus.KickStreamOffline, struct{}]{},
		usersRepo,
		channelsRepo,
	)

	msgID := "mark-processed-failure-001"
	redisMock.ExpectSetNX(idempotencyKeyPrefix+msgID, idempotencyStatusProcessing, idempotencyProcessingTTL).SetVal(true)
	redisMock.ExpectSet(idempotencyKeyPrefix+msgID, idempotencyStatusProcessed, idempotencyTTL).SetErr(errors.New("redis set failed"))

	payload := kickChatMessagePayload{
		MessageID: msgID,
		Broadcaster: kickUser{
			UserID:   123,
			Username: "broadcaster123",
		},
		Sender: kickUser{
			UserID:   456,
			Username: "senderlogin",
		},
		Content: "Hello world",
	}

	req := makeRequest(t, msgID, "chat.message.sent", payload)
	w := httptest.NewRecorder()

	h.HandleWebhook(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", w.Code)
	}

	if chatQueue.PublishedCount() != 1 {
		t.Fatalf("expected 1 chat message published, got %d", chatQueue.PublishedCount())
	}

	if parserQueue.PublishedCount() != 1 {
		t.Fatalf("expected 1 parser message published, got %d", parserQueue.PublishedCount())
	}

	if err := redisMock.ExpectationsWereMet(); err != nil {
		t.Errorf("redis expectations not met: %v", err)
	}
}
