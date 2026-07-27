package vkvideo

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/go-redsync/redsync/v4"
	"github.com/google/uuid"
	channelplatformentity "github.com/twirapp/twir/libs/entities/channel_platform"
	"github.com/twirapp/twir/libs/entities/platform"
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
	if message.ID != "123456789" || message.Text != "hello world" {
		t.Fatalf("message = %#v, want fixture id and normalized text", message)
	}
	if !message.Author.IsChatModerator || message.Author.IsOwner || message.Author.IsChannelModerator {
		t.Fatalf("author roles = %#v, want only observed moderator role", message.Author)
	}
}

func TestParseChatPublicationIncludesLinkPartContent(t *testing.T) {
	publication := []byte(`{
		"type": "channel_chat_message_send",
		"data": {
			"chat_message": {
				"id": 627867254,
				"created_at": 1785110751,
				"author": {
					"id": 35461641,
					"nick": "fixture_user",
					"is_owner": true,
					"is_moderator": false
				},
				"parts": [
					{"text": {"content": "!sr "}},
					{"link": {"url": "https://www.youtube.com/watch?v=64m0TmiHbcE", "content": "https://www.youtube.com/watch?v=64m0TmiHbcE"}},
					{"text": {"content": "\n"}}
				]
			}
		}
	}`)

	message, isChatMessage, err := parseChatPublication(publication)
	if err != nil {
		t.Fatalf("parse publication: %v", err)
	}
	if !isChatMessage {
		t.Fatal("publication was not recognized as a chat message")
	}
	want := "!sr https://www.youtube.com/watch?v=64m0TmiHbcE\n"
	if message.Text != want {
		t.Fatalf("message text = %q, want %q", message.Text, want)
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
	if connection.channel != "recorded-chat-channel" {
		t.Fatalf("subscription channel = %q", connection.channel)
	}
	if got := tokens.UserIDs(); len(got) != 3 || got[0] != binding.UserID || got[1] != binding.UserID || got[2] != binding.UserID {
		t.Fatalf("token user IDs = %v, want broadcaster binding user three times", got)
	}

	connection.publications <- publication
	connection.publications <- publication
	waitForMessageCount(t, chatMessages, 1)
	waitForMessageCount(t, commands, 1)
	if got := chatMessages.Messages(); len(got) != 1 {
		t.Fatalf("chat messages = %d, want one deduplicated message", len(got))
	} else if got[0].Platform != string(platform.PlatformVKVideoLive) || got[0].MessageID != "123456789" || got[0].Text != "hello world" || !got[0].IsModerator {
		t.Fatalf("normalized message = %#v", got[0])
	} else if _, err := uuid.Parse(got[0].ID); err != nil {
		t.Fatalf("normalized message ID = %q, want UUID: %v", got[0].ID, err)
	} else if got[0].ID == got[0].MessageID {
		t.Fatalf("normalized message ID = provider message ID %q", got[0].ID)
	}

	if err := transport.Unsubscribe(context.Background(), binding); err != nil {
		t.Fatalf("unsubscribe: %v", err)
	}
}

func TestTransportPublishesResolvedGlobalBotWhenSharedWithBroadcaster(t *testing.T) {
	binding := testBinding()
	sharedUserID := uuid.New()
	binding.UserID = sharedUserID
	binding.BotUserID = &sharedUserID
	publication := []byte(`{
		"type": "channel_chat_message_send",
		"data": {
			"chat_message": {
				"id": 987654321,
				"created_at": 1785097380,
				"author": {
					"id": 35461580,
					"nick": "shared_user",
					"is_owner": true,
					"is_moderator": false
				},
				"parts": [{"text": {"content": "!me"}}]
			}
		}
	}`)
	ensurer := &recordingChatUserEnsurer{user: &usersmodel.User{ID: sharedUserID}}
	chatMessages := &recordingPublisher{}
	commands := &recordingPublisher{}
	transport := newTestTransportWithUserCreator(t, &recordingTokenProvider{}, chatMessages, commands, &recordingConnection{}, ensurer)

	if err := transport.handlePublication(context.Background(), binding, publication); err != nil {
		t.Fatalf("handle publication: %v", err)
	}

	if got := chatMessages.Messages(); len(got) != 1 {
		t.Fatalf("chat messages = %d, want 1", len(got))
	}
	if got := commands.Messages(); len(got) != 1 {
		t.Fatalf("commands = %d, want 1", len(got))
	} else if got[0].Text != "!me" || got[0].UserID != sharedUserID.String() || !got[0].IsBroadcaster {
		t.Fatalf("command message = %#v, want !me from shared broadcaster", got[0])
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
	ensurer := &recordingChatUserEnsurer{user: &usersmodel.User{ID: botID}}
	chatMessages := &recordingPublisher{}
	commands := &recordingPublisher{}
	transport := newTestTransportWithUserCreator(t, &recordingTokenProvider{}, chatMessages, commands, &recordingConnection{}, ensurer)

	if err := transport.handlePublication(context.Background(), binding, publication); err != nil {
		t.Fatalf("handle publication: %v", err)
	}
	if ensurer.Calls() != 1 {
		t.Fatalf("UnsureUser calls = %d, want 1", ensurer.Calls())
	}
	if len(chatMessages.Messages()) != 0 || len(commands.Messages()) != 0 {
		t.Fatalf("bot publication was published: chat=%d commands=%d", len(chatMessages.Messages()), len(commands.Messages()))
	}
}

func TestTransportTreatsLeaseContentionAsNoOp(t *testing.T) {
	unrelatedErr := errors.New("boom")
	tests := []struct {
		name       string
		acquireErr error
		wantErr    error
	}{
		{name: "wrapped err taken", acquireErr: fmt.Errorf("acquire lock: %w", &redsync.ErrTaken{Nodes: []int{0}})},
		{name: "err failed", acquireErr: fmt.Errorf("acquire lock: %w", redsync.ErrFailed)},
		{name: "unrelated error", acquireErr: unrelatedErr, wantErr: unrelatedErr},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			acquiredConnections := 0
			ownership, err := newOwnership(
				LeaseConfig{Expiry: time.Minute, RenewInterval: 20 * time.Second},
				wrappedErrorMutexFactory{err: testCase.acquireErr},
				manualTickerFactory{ticker: newManualTicker()},
			)
			if err != nil {
				t.Fatalf("create ownership: %v", err)
			}

			transport := newTransport(transportDependencies{
				ownership:    ownership,
				tokens:       &recordingTokenProvider{},
				userCreator:  &recordingChatUserEnsurer{user: &usersmodel.User{ID: uuid.New()}},
				chatMessages: &recordingPublisher{},
				commands:     &recordingPublisher{},
				deduplicator: &memoryDeduplicator{claimed: make(map[string]struct{})},
				newConnection: func(RealtimeClientConfig) (realtimeConnection, error) {
					acquiredConnections++
					return &recordingConnection{}, nil
				},
			})

			err = transport.Subscribe(context.Background(), testBinding())
			if testCase.wantErr == nil {
				if err != nil {
					t.Fatalf("subscribe: %v", err)
				}
				if acquiredConnections != 0 {
					t.Fatalf("new connections = %d, want 0", acquiredConnections)
				}
				return
			}

			if !errors.Is(err, testCase.wantErr) {
				t.Fatalf("subscribe error = %v, want %v", err, testCase.wantErr)
			}
			if acquiredConnections != 0 {
				t.Fatalf("new connections = %d, want 0", acquiredConnections)
			}
		})
	}
}

func TestTransportShutdownReleasesAllBindings(t *testing.T) {
	// Given
	firstBinding := testBinding()
	secondBinding := testBinding()
	connections := map[uuid.UUID]*recordingConnection{
		firstBinding.ID:  {closed: make(chan struct{})},
		secondBinding.ID: {closed: make(chan struct{})},
	}
	transport := newTransport(transportDependencies{
		ownership:    newTestOwnership(t, newMemoryLockStore(), newManualTicker()),
		tokens:       &recordingTokenProvider{},
		userCreator:  &recordingChatUserEnsurer{user: &usersmodel.User{ID: uuid.New()}},
		chatMessages: &recordingPublisher{},
		commands:     &recordingPublisher{},
		deduplicator: &memoryDeduplicator{claimed: make(map[string]struct{})},
		newConnection: func(config RealtimeClientConfig) (realtimeConnection, error) {
			connection := connections[config.BindingID]
			connection.channel = config.Channel
			connection.tokens = config.Tokens
			return connection, nil
		},
	})
	if err := transport.Subscribe(context.Background(), firstBinding); err != nil {
		t.Fatalf("subscribe first binding: %v", err)
	}
	if err := transport.Subscribe(context.Background(), secondBinding); err != nil {
		t.Fatalf("subscribe second binding: %v", err)
	}

	// When
	if err := transport.Shutdown(context.Background()); err != nil {
		t.Fatalf("shutdown transport: %v", err)
	}

	// Then
	for bindingID, connection := range connections {
		select {
		case <-connection.closed:
		default:
			t.Fatalf("connection for binding %s was not closed", bindingID)
		}
	}
	transport.mu.Lock()
	bindings := len(transport.bindings)
	transport.mu.Unlock()
	if bindings != 0 {
		t.Fatalf("active bindings after shutdown = %d, want 0", bindings)
	}
	if err := transport.Subscribe(context.Background(), firstBinding); err != nil {
		t.Fatalf("resubscribe first binding: %v", err)
	}
	if err := transport.Subscribe(context.Background(), secondBinding); err != nil {
		t.Fatalf("resubscribe second binding: %v", err)
	}
	if err := transport.Unsubscribe(context.Background(), firstBinding); err != nil {
		t.Fatalf("unsubscribe first binding: %v", err)
	}
	if err := transport.Unsubscribe(context.Background(), secondBinding); err != nil {
		t.Fatalf("unsubscribe second binding: %v", err)
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

type wrappedErrorMutexFactory struct {
	err error
}

func (f wrappedErrorMutexFactory) NewMutex(string, time.Duration) leaseMutex {
	return wrappedErrorMutex{err: f.err}
}

type wrappedErrorMutex struct {
	err error
}

func (m wrappedErrorMutex) TryLockContext(context.Context) error {
	return m.err
}

func (wrappedErrorMutex) ExtendContext(context.Context) (bool, error) {
	return false, nil
}

func (wrappedErrorMutex) UnlockContext(context.Context) (bool, error) {
	return false, nil
}
