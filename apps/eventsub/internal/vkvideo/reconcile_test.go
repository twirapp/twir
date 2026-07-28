package vkvideo

import (
	"context"
	"testing"

	"github.com/google/uuid"
	channelplatformentity "github.com/twirapp/twir/libs/entities/channel_platform"
	usersmodel "github.com/twirapp/twir/libs/repositories/users/model"
)

func TestReconcileSubscribesAfterStaleLeaseExpires(t *testing.T) {
	// Given
	binding := testBinding()
	lockStore := newMemoryLockStore()
	staleMutex := lockStore.newMutex(binding.ID.String())
	if err := staleMutex.TryLockContext(context.Background()); err != nil {
		t.Fatalf("hold stale lease: %v", err)
	}

	createdConnections := 0
	transport := newTransport(transportDependencies{
		ownership:    newTestOwnership(t, lockStore, newManualTicker()),
		tokens:       &recordingTokenProvider{},
		userCreator:  &recordingChatUserEnsurer{user: &usersmodel.User{ID: uuid.New()}},
		chatMessages: &recordingPublisher{},
		commands:     &recordingPublisher{},
		deduplicator: &memoryDeduplicator{claimed: make(map[string]struct{})},
		newConnection: func(RealtimeClientConfig) (realtimeConnection, error) {
			createdConnections++
			return &recordingConnection{}, nil
		},
		databaseBindings: func(context.Context) ([]channelplatformentity.ChannelPlatform, error) {
			return []channelplatformentity.ChannelPlatform{binding}, nil
		},
	})

	if err := transport.Subscribe(context.Background(), binding); err != nil {
		t.Fatalf("subscribe with held lease: %v", err)
	}
	if createdConnections != 0 {
		t.Fatalf("connections while lease is held = %d, want 0", createdConnections)
	}

	// When
	lockStore.expire(binding.ID.String())
	transport.reconcileWithDatabase(context.Background())

	// Then
	if createdConnections != 1 {
		t.Fatalf("connections after reconcile = %d, want 1", createdConnections)
	}
}

func TestReconcileDoesNotResurrectUnsubscribedBinding(t *testing.T) {
	// Given
	binding := testBinding()
	databaseBindings := []channelplatformentity.ChannelPlatform{binding}
	connection := &recordingConnection{closed: make(chan struct{})}
	transport := newTransport(transportDependencies{
		ownership:    newTestOwnership(t, newMemoryLockStore(), newManualTicker()),
		tokens:       &recordingTokenProvider{},
		userCreator:  &recordingChatUserEnsurer{user: &usersmodel.User{ID: uuid.New()}},
		chatMessages: &recordingPublisher{},
		commands:     &recordingPublisher{},
		deduplicator: &memoryDeduplicator{claimed: make(map[string]struct{})},
		newConnection: func(RealtimeClientConfig) (realtimeConnection, error) {
			return connection, nil
		},
		databaseBindings: func(context.Context) ([]channelplatformentity.ChannelPlatform, error) {
			return databaseBindings, nil
		},
	})

	if err := transport.Subscribe(context.Background(), binding); err != nil {
		t.Fatalf("subscribe: %v", err)
	}

	// When
	databaseBindings = nil
	transport.reconcileWithDatabase(context.Background())

	// Then
	select {
	case <-connection.closed:
	default:
		t.Fatal("connection was not closed for binding removed from database")
	}
	transport.mu.Lock()
	activeCount := len(transport.bindings)
	transport.mu.Unlock()
	if activeCount != 0 {
		t.Fatalf("active bindings after reconcile = %d, want 0", activeCount)
	}
}
