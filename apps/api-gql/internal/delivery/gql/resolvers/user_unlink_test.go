package resolvers

import (
	"errors"
	"testing"

	"github.com/google/uuid"
	platformentity "github.com/twirapp/twir/libs/entities/platform"
	channelsmodel "github.com/twirapp/twir/libs/repositories/channels/model"
)

func TestUnlinkPlatformUpdates_AllowsSecondaryTwitchUnlink(t *testing.T) {
	t.Parallel()

	channel := channelsmodel.Channel{
		TwitchUserID:     uuidPtr(uuid.New()),
		TwitchPlatformID: stringPtr("twitch-user"),
		KickUserID:       uuidPtr(uuid.New()),
		KickPlatformID:   stringPtr("kick-user"),
		KickBotID:        uuidPtr(uuid.New()),
	}

	updates, err := unlinkPlatformUpdates(channel, platformentity.PlatformKick, platformentity.PlatformTwitch)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(updates) != 1 {
		t.Fatalf("expected one update, got %d", len(updates))
	}

	if value, ok := updates["twitch_user_id"]; !ok || value != nil {
		t.Fatalf("expected twitch_user_id to be cleared, got %#v", updates)
	}
}

func TestUnlinkPlatformUpdates_AllowsSecondaryKickUnlinkAndClearsKickBot(t *testing.T) {
	t.Parallel()

	channel := channelsmodel.Channel{
		TwitchUserID:     uuidPtr(uuid.New()),
		TwitchPlatformID: stringPtr("twitch-user"),
		KickUserID:       uuidPtr(uuid.New()),
		KickPlatformID:   stringPtr("kick-user"),
		KickBotID:        uuidPtr(uuid.New()),
	}

	updates, err := unlinkPlatformUpdates(channel, platformentity.PlatformTwitch, platformentity.PlatformKick)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if value, ok := updates["kick_user_id"]; !ok || value != nil {
		t.Fatalf("expected kick_user_id to be cleared, got %#v", updates)
	}

	if value, ok := updates["kick_bot_id"]; !ok || value != nil {
		t.Fatalf("expected kick_bot_id to be cleared, got %#v", updates)
	}
}

func TestUnlinkPlatformUpdates_BlocksCurrentPlatform(t *testing.T) {
	t.Parallel()

	channel := channelsmodel.Channel{
		TwitchUserID:     uuidPtr(uuid.New()),
		TwitchPlatformID: stringPtr("twitch-user"),
		KickUserID:       uuidPtr(uuid.New()),
		KickPlatformID:   stringPtr("kick-user"),
	}

	_, err := unlinkPlatformUpdates(channel, platformentity.PlatformKick, platformentity.PlatformKick)
	if !errors.Is(err, errCannotUnlinkCurrentPlatform) {
		t.Fatalf("expected current-platform error, got %v", err)
	}
}

func TestUnlinkPlatformUpdates_BlocksLastPlatform(t *testing.T) {
	t.Parallel()

	channel := channelsmodel.Channel{
		KickUserID:     uuidPtr(uuid.New()),
		KickPlatformID: stringPtr("kick-user"),
	}

	_, err := unlinkPlatformUpdates(channel, platformentity.PlatformTwitch, platformentity.PlatformKick)
	if !errors.Is(err, errCannotUnlinkLastPlatform) {
		t.Fatalf("expected last-platform error, got %v", err)
	}
}

func uuidPtr(id uuid.UUID) *uuid.UUID {
	return &id
}

func stringPtr(value string) *string {
	return &value
}
