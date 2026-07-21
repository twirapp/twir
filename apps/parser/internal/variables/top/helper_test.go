package top

import (
	"encoding/json"
	"testing"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/twirapp/twir/libs/entities/platform"
	channelplatformsmodel "github.com/twirapp/twir/libs/repositories/channel_platforms/model"
	channelsmodel "github.com/twirapp/twir/libs/repositories/channels/model"
)

func TestApplyTopChannelBotFiltersUsesEveryBinding(t *testing.T) {
	kickUserID := uuid.New()
	kickBotUserID := uuid.New()
	twitchUserID := uuid.New()
	channel := channelsmodel.Channel{
		Bindings: []channelplatformsmodel.ChannelPlatform{
			{
				Platform:  platform.PlatformKick,
				UserID:    kickUserID,
				BotUserID: &kickBotUserID,
			},
			{
				Platform: platform.PlatformTwitch,
				UserID:   twitchUserID,
				BotConfig: json.RawMessage(
					`{"bot_id":"twitch-bot"}`,
				),
			},
		},
	}

	filtered, err := applyTopChannelBotFilters(
		squirrel.Select("users_stats.user_id").From("users_stats"),
		channel,
	)
	if err != nil {
		t.Fatalf("applyTopChannelBotFilters returned error: %v", err)
	}
	query, args, err := filtered.ToSql()
	if err != nil {
		t.Fatalf("build filtered query: %v", err)
	}

	if query == "" {
		t.Fatal("expected filtered query")
	}

	wantArgs := map[string]bool{
		"twitch-bot":           true,
		kickUserID.String():    true,
		kickBotUserID.String(): true,
		twitchUserID.String():  true,
	}
	for _, arg := range args {
		value, ok := arg.(string)
		if !ok {
			t.Errorf("argument %T = %v, want string", arg, arg)
			continue
		}
		if !wantArgs[value] {
			t.Errorf("unexpected filter argument %q", value)
			continue
		}
		delete(wantArgs, value)
	}
	for value := range wantArgs {
		t.Errorf("missing filter argument %q", value)
	}
}
