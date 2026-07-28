package pgx

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/twirapp/twir/libs/entities/platform"
	"github.com/twirapp/twir/libs/repositories/channel_platforms/model"
)

func TestMapBindingToEntityMapsEveryField(t *testing.T) {
	t.Parallel()

	botUserID := uuid.New()
	createdAt := time.Date(2026, 7, 24, 1, 2, 3, 0, time.UTC)
	updatedAt := createdAt.Add(time.Hour)

	m := model.ChannelPlatform{
		ID:                uuid.MustParse("01983578-e5c2-7c8e-9a52-c7a4cf21c9d1"),
		ChannelID:         uuid.MustParse("01983578-e5c2-7c8e-9a52-c7a4cf21c9d2"),
		Platform:          platform.PlatformTwitch,
		UserID:            uuid.MustParse("01983578-e5c2-7c8e-9a52-c7a4cf21c9d3"),
		PlatformChannelID: "twitch-channel",
		Enabled:           true,
		BotUserID:         &botUserID,
		BotConfig:         []byte(`{"bot_id":"bot-1"}`),
		CreatedAt:         createdAt,
		UpdatedAt:         updatedAt,
	}

	entity := mapBindingToEntity(m)

	if entity.IsNil() {
		t.Fatal("entity must not be nil-marked for a populated model")
	}
	if entity.ID != m.ID || entity.ChannelID != m.ChannelID || entity.UserID != m.UserID {
		t.Fatalf("entity ids mismatch: %#v", entity)
	}
	if entity.Platform != m.Platform || entity.PlatformChannelID != m.PlatformChannelID || entity.Enabled != m.Enabled {
		t.Fatalf("entity core fields mismatch: %#v", entity)
	}
	if entity.BotUserID == nil || *entity.BotUserID != botUserID {
		t.Fatalf("entity BotUserID = %v, want %v", entity.BotUserID, botUserID)
	}
	if string(entity.BotConfig) != string(m.BotConfig) {
		t.Fatalf("entity BotConfig = %s, want %s", entity.BotConfig, m.BotConfig)
	}
	if !entity.CreatedAt.Equal(createdAt) || !entity.UpdatedAt.Equal(updatedAt) {
		t.Fatalf("entity timestamps mismatch: %#v", entity)
	}
}

func TestMapBindingToEntityMapsModelNilToEntityNil(t *testing.T) {
	t.Parallel()

	entity := mapBindingToEntity(model.Nil)
	if !entity.IsNil() {
		t.Fatal("model.Nil must map to entity.Nil")
	}
}

func TestMapBindingsToEntitiesPreservesOrder(t *testing.T) {
	t.Parallel()

	models := []model.ChannelPlatform{
		{ID: uuid.New(), Platform: platform.PlatformKick},
		{ID: uuid.New(), Platform: platform.PlatformTwitch},
	}

	entities := mapBindingsToEntities(models)
	if len(entities) != len(models) {
		t.Fatalf("entities len = %d, want %d", len(entities), len(models))
	}
	for i := range models {
		if entities[i].ID != models[i].ID || entities[i].Platform != models[i].Platform {
			t.Fatalf("entities[%d] = %#v, want model %#v", i, entities[i], models[i])
		}
	}
}
