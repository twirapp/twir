package vkvideo

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/uuid"
	"github.com/twirapp/twir/libs/entities/platform"
	usersmodel "github.com/twirapp/twir/libs/repositories/users/model"
)

func TestTransportEnsuresChatUserStatsWithVKIdentity(t *testing.T) {
	// Given
	publication, err := os.ReadFile(filepath.Join("testdata", "public_chat_message.json"))
	if err != nil {
		t.Fatalf("read fixture: %v", err)
	}
	binding := testBinding()
	ensuredUser := usersmodel.User{ID: uuid.New()}
	ensurer := &recordingChatUserEnsurer{user: &ensuredUser}
	chatMessages := &recordingPublisher{}
	commands := &recordingPublisher{}
	transport := newTransport(transportDependencies{
		ownership:    newTestOwnership(t, newMemoryLockStore(), newManualTicker()),
		userCreator:  ensurer,
		chatMessages: chatMessages,
		commands:     commands,
		deduplicator: &memoryDeduplicator{claimed: make(map[string]struct{})},
	})

	// When
	err = transport.handlePublication(context.Background(), binding, publication)

	// Then
	if err != nil {
		t.Fatalf("handle publication: %v", err)
	}
	if ensurer.Calls() != 1 {
		t.Fatalf("UnsureUser calls = %d, want 1", ensurer.Calls())
	}
	got := ensurer.Input()
	if got.UserID != "fixture-author-1" {
		t.Fatalf("UserID = %q, want fixture-author-1", got.UserID)
	}
	if got.PlatformID != "fixture-author-1" {
		t.Fatalf("PlatformID = %q, want fixture-author-1", got.PlatformID)
	}
	if got.Platform != platform.PlatformVKVideoLive {
		t.Fatalf("Platform = %q, want %q", got.Platform, platform.PlatformVKVideoLive)
	}
	if got.Login != "fixture_user" {
		t.Fatalf("Login = %q, want fixture_user", got.Login)
	}
	if got.DisplayName != "Fixture Display" {
		t.Fatalf("DisplayName = %q, want Fixture Display", got.DisplayName)
	}
	if got.ChannelID == nil || *got.ChannelID != binding.ChannelID.String() {
		t.Fatalf("ChannelID = %v, want %s", got.ChannelID, binding.ChannelID.String())
	}
	if got.ShouldUpdateStats {
		t.Fatal("ShouldUpdateStats = true, want false")
	}
	if got.IsBroadcaster {
		t.Fatal("IsBroadcaster = true, want false")
	}
	if !got.IsModerator {
		t.Fatal("IsModerator = false, want true")
	}
	if got.IsVip {
		t.Fatal("IsVip = true, want false")
	}
	if got.IsSubscriber {
		t.Fatal("IsSubscriber = true, want false")
	}
	if messages := commands.Messages(); len(messages) != 1 || messages[0].UserID != ensuredUser.ID.String() {
		t.Fatalf("command messages = %#v, want one message for ensured user %s", messages, ensuredUser.ID)
	}
}
