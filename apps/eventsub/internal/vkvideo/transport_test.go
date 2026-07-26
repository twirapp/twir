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
				users:        &recordingUserStore{user: usersmodel.User{ID: uuid.New()}},
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
