package vkvideo

import (
	"context"
	"testing"

	"github.com/google/uuid"
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
	})

	if err := transport.Subscribe(context.Background(), binding); err != nil {
		t.Fatalf("subscribe with held lease: %v", err)
	}
	if createdConnections != 0 {
		t.Fatalf("connections while lease is held = %d, want 0", createdConnections)
	}

	// When
	lockStore.expire(binding.ID.String())
	transport.reconcilePending(context.Background())

	// Then
	if createdConnections != 1 {
		t.Fatalf("connections after reconcile = %d, want 1", createdConnections)
	}
}

func TestReconcileDoesNotResurrectUnsubscribedBinding(t *testing.T) {
	// Given
	binding := testBinding()
	createdConnections := 0
	transport := newTransport(transportDependencies{
		ownership:    newTestOwnership(t, newMemoryLockStore(), newManualTicker()),
		tokens:       &recordingTokenProvider{},
		userCreator:  &recordingChatUserEnsurer{user: &usersmodel.User{ID: uuid.New()}},
		chatMessages: &recordingPublisher{},
		commands:     &recordingPublisher{},
		deduplicator: &memoryDeduplicator{claimed: make(map[string]struct{})},
		newConnection: func(RealtimeClientConfig) (realtimeConnection, error) {
			createdConnections++
			return &recordingConnection{}, nil
		},
	})

	if err := transport.Subscribe(context.Background(), binding); err != nil {
		t.Fatalf("subscribe: %v", err)
	}
	if err := transport.Unsubscribe(context.Background(), binding); err != nil {
		t.Fatalf("unsubscribe: %v", err)
	}

	// When
	transport.reconcilePending(context.Background())

	// Then
	if createdConnections != 1 {
		t.Fatalf("connections after reconcile = %d, want 1 (no resurrection)", createdConnections)
	}
}
