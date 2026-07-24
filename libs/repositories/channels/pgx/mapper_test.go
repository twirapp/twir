package pgx

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/twirapp/twir/libs/entities/platform"
	channelplatformsmodel "github.com/twirapp/twir/libs/repositories/channel_platforms/model"
	"github.com/twirapp/twir/libs/repositories/channels/model"
)

func TestMapChannelToEntityMapsFieldsAndBindingsInOrder(t *testing.T) {
	t.Parallel()

	apiKey := "api-key-value"
	botUserID := uuid.New()
	kickBindingID := uuid.New()
	twitchBindingID := uuid.New()
	channelID := uuid.MustParse("01983578-e5c2-7c8e-9a52-c7a4cf21c9d1")
	createdAt := time.Date(2026, 7, 24, 1, 2, 3, 0, time.UTC)
	updatedAt := createdAt.Add(time.Hour)

	m := model.Channel{
		ID:     channelID,
		ApiKey: &apiKey,
		Bindings: []channelplatformsmodel.ChannelPlatform{
			{
				ID:                kickBindingID,
				ChannelID:         channelID,
				Platform:          platform.PlatformKick,
				PlatformChannelID: "kick-channel",
			},
			{
				ID:                twitchBindingID,
				ChannelID:         channelID,
				Platform:          platform.PlatformTwitch,
				UserID:            uuid.MustParse("01983578-e5c2-7c8e-9a52-c7a4cf21c9d3"),
				PlatformChannelID: "twitch-channel",
				Enabled:           true,
				BotUserID:         &botUserID,
				BotConfig:         []byte(`{"bot_id":"bot-1"}`),
				CreatedAt:         createdAt,
				UpdatedAt:         updatedAt,
			},
		},
	}

	entity := mapChannelToEntity(m)

	if entity.IsNil() {
		t.Fatal("entity must not be nil-marked for a populated model")
	}
	if entity.ID != m.ID {
		t.Fatalf("entity.ID = %s, want %s", entity.ID, m.ID)
	}
	if entity.ApiKey == nil || *entity.ApiKey != apiKey {
		t.Fatalf("entity.ApiKey = %v, want %q", entity.ApiKey, apiKey)
	}
	if len(entity.Bindings) != 2 {
		t.Fatalf("entity.Bindings len = %d, want 2", len(entity.Bindings))
	}
	kick := entity.Bindings[0]
	if kick.ID != kickBindingID || kick.ChannelID != channelID || kick.Platform != platform.PlatformKick || kick.PlatformChannelID != "kick-channel" {
		t.Fatalf("binding order or kick fields not preserved: %#v", kick)
	}
	twitch := entity.Bindings[1]
	want := m.Bindings[1]
	if twitch.ID != twitchBindingID || twitch.ChannelID != want.ChannelID || twitch.Platform != want.Platform || twitch.PlatformChannelID != want.PlatformChannelID {
		t.Fatalf("nested binding identity mismatch: %#v", twitch)
	}
	if twitch.UserID != want.UserID || twitch.Enabled != want.Enabled {
		t.Fatalf("nested binding fields mismatch: %#v", twitch)
	}
	if twitch.BotUserID == nil || *twitch.BotUserID != botUserID {
		t.Fatalf("nested BotUserID = %v, want %v", twitch.BotUserID, botUserID)
	}
	if string(twitch.BotConfig) != `{"bot_id":"bot-1"}` {
		t.Fatalf("nested BotConfig = %s", twitch.BotConfig)
	}
	if !twitch.CreatedAt.Equal(createdAt) || !twitch.UpdatedAt.Equal(updatedAt) {
		t.Fatalf("nested timestamps mismatch: %#v", twitch)
	}
}

func TestMapChannelToEntityMapsModelNilToEntityNil(t *testing.T) {
	t.Parallel()

	if entity := mapChannelToEntity(model.Nil); !entity.IsNil() {
		t.Fatal("model.Nil must map to entity.Nil")
	}
}

func TestMapChannelToEntityMapsNilApiKey(t *testing.T) {
	t.Parallel()

	entity := mapChannelToEntity(model.Channel{ID: uuid.New()})
	if entity.ApiKey != nil {
		t.Fatalf("entity.ApiKey = %v, want nil", entity.ApiKey)
	}
}
