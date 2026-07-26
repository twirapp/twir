package vkvideo

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/twirapp/twir/libs/bus-core/generic"
	channelplatformentity "github.com/twirapp/twir/libs/entities/channel_platform"
	"github.com/twirapp/twir/libs/entities/platform"
	usersrepository "github.com/twirapp/twir/libs/repositories/users"
	usersmodel "github.com/twirapp/twir/libs/repositories/users/model"
)

func TestParseChatPublicationNormalizesOnlyObservedMessageShape(t *testing.T) {
	publication, err := os.ReadFile(filepath.Join("testdata", "public_chat_message.json"))
	if err != nil {
		t.Fatalf("read fixture: %v", err)
	}

	message, isChatMessage, err := parseChatPublication(publication)
	if err != nil {
		t.Fatalf("parse publication: %v", err)
	}
	if !isChatMessage {
		t.Fatal("publication was not recognized as a chat message")
	}
	if message.ID != "fixture-message-1" || message.Text != "hello world" {
		t.Fatalf("message = %#v, want fixture id and normalized text", message)
	}
	if !message.Author.IsChatModerator || message.Author.IsOwner || message.Author.IsChannelModerator {
		t.Fatalf("author roles = %#v, want only observed moderator role", message.Author)
	}
}

func TestParseChatPublicationIgnoresNonMessageTypes(t *testing.T) {
	message, isChatMessage, err := parseChatPublication([]byte(`{"type":"reaction","data":{}}`))
	if err != nil {
		t.Fatalf("parse non-message publication: %v", err)
	}
	if isChatMessage || message.ID != "" {
		t.Fatalf("non-message publication = %#v, %t; want ignored", message, isChatMessage)
	}
}

func TestTransportAuthenticatesBindingUserAndDeduplicatesPublications(t *testing.T) {
	publication, err := os.ReadFile(filepath.Join("testdata", "public_chat_message.json"))
	if err != nil {
		t.Fatalf("read fixture: %v", err)
	}
	binding := testBinding()
	tokens := &recordingTokenProvider{}
	chatMessages := &recordingPublisher{}
	commands := &recordingPublisher{}
	connection := &recordingConnection{publications: make(chan []byte, 2)}
	transport := newTestTransport(t, tokens, chatMessages, commands, connection, usersmodel.User{ID: uuid.New()})

	if err := transport.Subscribe(context.Background(), binding); err != nil {
		t.Fatalf("subscribe: %v", err)
	}
	if connection.channel != "channel-chat:"+binding.PlatformChannelID {
		t.Fatalf("subscription channel = %q", connection.channel)
	}
	if got := tokens.UserIDs(); len(got) != 2 || got[0] != binding.UserID || got[1] != binding.UserID {
		t.Fatalf("token user IDs = %v, want broadcaster binding user twice", got)
	}

	connection.publications <- publication
	connection.publications <- publication
	waitForMessageCount(t, chatMessages, 1)
	waitForMessageCount(t, commands, 1)
	if got := chatMessages.Messages(); len(got) != 1 {
		t.Fatalf("chat messages = %d, want one deduplicated message", len(got))
	} else if got[0].Platform != string(platform.PlatformVKVideoLive) || got[0].MessageID != "fixture-message-1" || got[0].Text != "hello world" || !got[0].IsModerator {
		t.Fatalf("normalized message = %#v", got[0])
	}

	if err := transport.Unsubscribe(context.Background(), binding); err != nil {
		t.Fatalf("unsubscribe: %v", err)
	}
}

func TestTransportSuppressesResolvedGlobalBot(t *testing.T) {
	publication, err := os.ReadFile(filepath.Join("testdata", "public_chat_message.json"))
	if err != nil {
		t.Fatalf("read fixture: %v", err)
	}
	binding := testBinding()
	botID := uuid.New()
	binding.BotUserID = &botID
	users := &recordingUserStore{user: usersmodel.User{ID: botID}}
	chatMessages := &recordingPublisher{}
	commands := &recordingPublisher{}
	transport := newTestTransportWithUsers(t, &recordingTokenProvider{}, chatMessages, commands, &recordingConnection{}, users)

	if err := transport.handlePublication(context.Background(), binding, publication); err != nil {
		t.Fatalf("handle publication: %v", err)
	}
	if users.lookupCalls != 1 {
		t.Fatalf("sender resolution calls = %d, want 1", users.lookupCalls)
	}
	if len(chatMessages.Messages()) != 0 || len(commands.Messages()) != 0 {
		t.Fatalf("bot publication was published: chat=%d commands=%d", len(chatMessages.Messages()), len(commands.Messages()))
	}
}

func testBinding() channelplatformentity.ChannelPlatform {
	return channelplatformentity.ChannelPlatform{
		ID:                uuid.New(),
		ChannelID:         uuid.New(),
		Platform:          platform.PlatformVKVideoLive,
		UserID:            uuid.New(),
		PlatformChannelID: "fixture-channel",
		Enabled:           true,
	}
}

func newTestTransport(
	t *testing.T,
	tokens tokenProvider,
	chatMessages messagePublisher,
	commands messagePublisher,
	connection realtimeConnection,
	user usersmodel.User,
) *Transport {
	t.Helper()
	return newTestTransportWithUsers(t, tokens, chatMessages, commands, connection, &recordingUserStore{user: user})
}

func newTestTransportWithUsers(
	t *testing.T,
	tokens tokenProvider,
	chatMessages messagePublisher,
	commands messagePublisher,
	connection realtimeConnection,
	users userStore,
) *Transport {
	t.Helper()
	return newTransport(transportDependencies{
		ownership:    newTestOwnership(t, newMemoryLockStore(), newManualTicker()),
		tokens:       tokens,
		users:        users,
		chatMessages: chatMessages,
		commands:     commands,
		deduplicator: &memoryDeduplicator{claimed: make(map[string]struct{})},
		newConnection: func(config RealtimeClientConfig) (realtimeConnection, error) {
			if recorded, ok := connection.(*recordingConnection); ok {
				recorded.channel = config.Channel
				recorded.tokens = config.Tokens
			}
			return connection, nil
		},
	})
}

type recordingTokenProvider struct {
	mu      sync.Mutex
	userIDs []uuid.UUID
}

func (p *recordingTokenProvider) GetUserToken(_ context.Context, userID uuid.UUID) (string, error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.userIDs = append(p.userIDs, userID)
	return "fixture-access-token", nil
}

func (p *recordingTokenProvider) UserIDs() []uuid.UUID {
	p.mu.Lock()
	defer p.mu.Unlock()
	return append([]uuid.UUID(nil), p.userIDs...)
}

type recordingConnection struct {
	channel      string
	tokens       TokenCallbacks
	publications chan []byte
	closed       chan struct{}
}

func (c *recordingConnection) Connect() error {
	if c.tokens.Connection == nil || c.tokens.Subscription == nil {
		return nil
	}
	if _, err := c.tokens.Connection(); err != nil {
		return err
	}
	_, err := c.tokens.Subscription(c.channel)
	return err
}

func (c *recordingConnection) Receive(ctx context.Context) ([]byte, error) {
	if c.publications == nil {
		<-ctx.Done()
		return nil, ctx.Err()
	}
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case publication := <-c.publications:
		return publication, nil
	}
}

func (c *recordingConnection) Close() {
	if c.closed == nil {
		c.closed = make(chan struct{})
	}
	select {
	case <-c.closed:
	default:
		close(c.closed)
	}
}

type recordingPublisher struct {
	mu       sync.Mutex
	messages []generic.ChatMessage
}

func (p *recordingPublisher) Publish(_ context.Context, message generic.ChatMessage) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.messages = append(p.messages, message)
	return nil
}

func (p *recordingPublisher) Messages() []generic.ChatMessage {
	p.mu.Lock()
	defer p.mu.Unlock()
	return append([]generic.ChatMessage(nil), p.messages...)
}

type recordingUserStore struct {
	user        usersmodel.User
	lookupCalls int
	err         error
}

func (s *recordingUserStore) GetByPlatformID(context.Context, platform.Platform, string) (usersmodel.User, error) {
	s.lookupCalls++
	return s.user, s.err
}

func (*recordingUserStore) Create(context.Context, usersrepository.CreateInput) (usersmodel.User, error) {
	return usersmodel.User{}, errors.New("unexpected user creation")
}

type memoryDeduplicator struct {
	mu      sync.Mutex
	claimed map[string]struct{}
}

func (d *memoryDeduplicator) Claim(_ context.Context, id string) (bool, error) {
	d.mu.Lock()
	defer d.mu.Unlock()
	if _, exists := d.claimed[id]; exists {
		return false, nil
	}
	d.claimed[id] = struct{}{}
	return true, nil
}

func waitForMessageCount(t *testing.T, publisher *recordingPublisher, want int) {
	t.Helper()
	deadline := time.Now().Add(time.Second)
	for time.Now().Before(deadline) {
		if len(publisher.Messages()) >= want {
			return
		}
		time.Sleep(time.Millisecond)
	}
	t.Fatalf("published messages = %d, want at least %d", len(publisher.Messages()), want)
}
