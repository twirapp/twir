package vkvideo

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	user_creator "github.com/twirapp/twir/apps/eventsub/internal/services/user-creator"
	"github.com/twirapp/twir/libs/bus-core/generic"
	usersmodel "github.com/twirapp/twir/libs/repositories/users/model"
	usersstatsmodel "github.com/twirapp/twir/libs/repositories/users_stats/model"
)

func newTestTransport(
	t *testing.T,
	tokens webSocketTokenProvider,
	chatMessages messagePublisher,
	commands messagePublisher,
	connection realtimeConnection,
	user usersmodel.User,
) *Transport {
	t.Helper()
	return newTestTransportWithUserCreator(t, tokens, chatMessages, commands, connection, &recordingChatUserEnsurer{user: &user})
}

func newTestTransportWithUserCreator(
	t *testing.T,
	tokens webSocketTokenProvider,
	chatMessages messagePublisher,
	commands messagePublisher,
	connection realtimeConnection,
	userCreator chatUserEnsurer,
) *Transport {
	t.Helper()
	return newTransport(transportDependencies{
		ownership:    newTestOwnership(t, newMemoryLockStore(), newManualTicker()),
		tokens:       tokens,
		userCreator:  userCreator,
		chatMessages: chatMessages,
		commands:     commands,
		deduplicator: &memoryDeduplicator{claimed: make(map[string]struct{})},
		newConnection: func(config RealtimeClientConfig) (realtimeConnection, error) {
			if recorded, ok := connection.(*recordingConnection); ok {
				recorded.created = true
				recorded.channel = config.Channel
				recorded.tokens = config.Tokens
			}
			return connection, nil
		},
	})
}

type recordingTokenProvider struct {
	mu                sync.Mutex
	userIDs           []uuid.UUID
	discoverErr       error
	userTokenContexts []context.Context
}

func (p *recordingTokenProvider) GetUserToken(ctx context.Context, userID uuid.UUID) (string, error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.userIDs = append(p.userIDs, userID)
	p.userTokenContexts = append(p.userTokenContexts, ctx)
	return "fixture-access-token", nil
}

func (p *recordingTokenProvider) ConnectionToken(ctx context.Context, userID uuid.UUID) (string, error) {
	return p.GetUserToken(ctx, userID)
}

func (p *recordingTokenProvider) DiscoverChatChannel(ctx context.Context, userID uuid.UUID) (string, error) {
	if _, err := p.GetUserToken(ctx, userID); err != nil {
		return "", err
	}
	if p.discoverErr != nil {
		return "", p.discoverErr
	}
	return "recorded-chat-channel", nil
}

func (p *recordingTokenProvider) SubscriptionToken(ctx context.Context, userID uuid.UUID, _ string) (string, error) {
	return p.GetUserToken(ctx, userID)
}

func (p *recordingTokenProvider) UserIDs() []uuid.UUID {
	p.mu.Lock()
	defer p.mu.Unlock()
	return append([]uuid.UUID(nil), p.userIDs...)
}

func (p *recordingTokenProvider) UserTokenContexts() []context.Context {
	p.mu.Lock()
	defer p.mu.Unlock()
	return append([]context.Context(nil), p.userTokenContexts...)
}

type recordingConnection struct {
	mu                sync.Mutex
	created           bool
	channel           string
	tokens            TokenCallbacks
	connectErr        error
	connectCalls      int
	connectionToken   string
	subscriptionToken string
	publications      chan []byte
	closed            chan struct{}
}

func (c *recordingConnection) Connect(ctx context.Context) error {
	c.mu.Lock()
	c.connectCalls++
	connectErr := c.connectErr
	c.mu.Unlock()
	if connectErr != nil {
		return connectErr
	}

	if c.tokens.Connection == nil || c.tokens.Subscription == nil {
		return nil
	}

	connectionToken, err := c.tokens.Connection(ctx)
	if err != nil {
		return err
	}
	subscriptionToken, err := c.tokens.Subscription(ctx, c.channel)
	if err != nil {
		return err
	}
	c.connectionToken = connectionToken
	c.subscriptionToken = subscriptionToken
	return nil
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

type recordingChatUserEnsurer struct {
	mu     sync.Mutex
	user   *usersmodel.User
	inputs []user_creator.CreateUserInput
}

func (e *recordingChatUserEnsurer) UnsureUser(
	_ context.Context,
	input user_creator.CreateUserInput,
) (*usersmodel.User, *usersstatsmodel.UserStat, error) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.inputs = append(e.inputs, input)
	return e.user, nil, nil
}

func (e *recordingChatUserEnsurer) Calls() int {
	e.mu.Lock()
	defer e.mu.Unlock()
	return len(e.inputs)
}

func (e *recordingChatUserEnsurer) Input() user_creator.CreateUserInput {
	e.mu.Lock()
	defer e.mu.Unlock()
	return e.inputs[0]
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
