package pgx

import (
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/twirapp/twir/libs/entities/platform"
	channelplatformsmodel "github.com/twirapp/twir/libs/repositories/channel_platforms/model"
	usersmodel "github.com/twirapp/twir/libs/repositories/users/model"
	"github.com/twirapp/twir/libs/repositories/users_with_channel"
)

func TestMapUserWithChannelProjectionKeepsSelectedBinding(t *testing.T) {
	channelID := uuid.New()
	botUserID := uuid.New()
	createdAt := time.Date(2026, 7, 24, 1, 2, 3, 0, time.UTC)
	updatedAt := createdAt.Add(time.Hour)
	binding := channelplatformsmodel.ChannelPlatform{
		ID:                uuid.New(),
		ChannelID:         channelID,
		Platform:          platform.PlatformTwitch,
		UserID:            uuid.New(),
		PlatformChannelID: "twitch-channel",
		Enabled:           true,
		BotUserID:         &botUserID,
		BotConfig:         json.RawMessage(`{"bot_id":"twir-bot","is_bot_mod":true}`),
		CreatedAt:         createdAt,
		UpdatedAt:         updatedAt,
	}
	bindingJSON, err := json.Marshal(binding)
	if err != nil {
		t.Fatalf("marshal binding: %v", err)
	}

	result, err := mapUserWithChannelProjection(
		usersmodel.User{ID: binding.UserID, Platform: platform.PlatformTwitch},
		pgtype.UUID{Bytes: [16]byte(channelID), Valid: true},
		bindingJSON,
	)
	if err != nil {
		t.Fatalf("map projection: %v", err)
	}
	if result.Channel == nil {
		t.Fatal("expected projected channel")
	}
	if result.Channel.ID != channelID {
		t.Fatalf("channel ID = %s, want %s", result.Channel.ID, channelID)
	}
	if len(result.Channel.Bindings) != 1 {
		t.Fatalf("bindings = %d, want 1", len(result.Channel.Bindings))
	}
	got := result.Channel.Bindings[0]
	if got.ID != binding.ID ||
		got.ChannelID != binding.ChannelID ||
		got.Platform != binding.Platform ||
		got.UserID != binding.UserID ||
		got.PlatformChannelID != binding.PlatformChannelID ||
		got.Enabled != binding.Enabled ||
		string(got.BotConfig) != string(binding.BotConfig) {
		t.Fatalf("binding = %#v, want fields of %#v", got, binding)
	}
	if got.BotUserID == nil || *got.BotUserID != botUserID {
		t.Fatalf("BotUserID = %v, want %v", got.BotUserID, botUserID)
	}
	if !got.CreatedAt.Equal(createdAt) || !got.UpdatedAt.Equal(updatedAt) {
		t.Fatalf("timestamps = (%v, %v), want (%v, %v)", got.CreatedAt, got.UpdatedAt, createdAt, updatedAt)
	}
	if got.IsNil() {
		t.Fatal("projected binding must not be nil-marked")
	}
}

func TestGetByIDQueryUsesNormalizedBindingProjection(t *testing.T) {
	for _, fragment := range []string{
		"LEFT JOIN LATERAL",
		"FROM channel_platforms cp",
		"cp.user_id = u.id",
		"cp.platform = u.platform",
		"LIMIT 1",
	} {
		if !strings.Contains(getByIDQuery, fragment) {
			t.Fatalf("query does not contain %q: %s", fragment, getByIDQuery)
		}
	}
	for _, legacyColumn := range []string{
		"twitch_user_id",
		"kick_user_id",
		"twitch_bot_enabled",
		"kick_bot_enabled",
		`"isEnabled"`,
		`"botId"`,
	} {
		if strings.Contains(getByIDQuery, legacyColumn) {
			t.Fatalf("query must not use legacy channel column %q: %s", legacyColumn, getByIDQuery)
		}
	}
}

func TestBuildGetManyQuerySkipsBadgesJoinWithoutBadgeFilter(t *testing.T) {
	query, _, err := buildGetManyQuery(users_with_channel.GetManyInput{PerPage: 10})
	if err != nil {
		t.Fatalf("build query: %v", err)
	}
	if strings.Contains(query, "badges_users") {
		t.Fatalf("unexpected badges join in query: %s", query)
	}
}

func TestBuildGetManyQueryUsesNormalizedBindingProjection(t *testing.T) {
	query, _, err := buildGetManyQuery(users_with_channel.GetManyInput{PerPage: 10})
	if err != nil {
		t.Fatalf("build query: %v", err)
	}

	for _, fragment := range []string{
		"LEFT JOIN LATERAL",
		"FROM channel_platforms cp",
		"cp.user_id = u.id",
		"cp.platform = u.platform",
		"jsonb_build_object",
	} {
		if !strings.Contains(query, fragment) {
			t.Fatalf("query does not contain %q: %s", fragment, query)
		}
	}
	if joins := strings.Count(query, "LEFT JOIN"); joins != 1 {
		t.Fatalf("left joins = %d, want one lateral binding join: %s", joins, query)
	}

	for _, legacyColumn := range []string{
		"twitch_user_id",
		"kick_user_id",
		"twitch_bot_enabled",
		"kick_bot_enabled",
		`"isEnabled"`,
		`"botId"`,
	} {
		if strings.Contains(query, legacyColumn) {
			t.Fatalf("query must not use legacy channel column %q: %s", legacyColumn, query)
		}
	}
}

func TestBuildGetManyCountQueryWithoutJoinFiltersCountsUsersDirectly(t *testing.T) {
	query, args, err := buildGetManyCountQuery(users_with_channel.GetManyInput{})
	if err != nil {
		t.Fatalf("build query: %v", err)
	}
	if query != "SELECT COUNT(*) FROM users u" {
		t.Fatalf("unexpected count query: %s", query)
	}
	if len(args) != 0 {
		t.Fatalf("expected no query arguments, got %d", len(args))
	}
}

func TestBuildGetManyCountQueryUsesSelectedBindingForEnabledFilter(t *testing.T) {
	enabled := true
	query, _, err := buildGetManyCountQuery(users_with_channel.GetManyInput{ChannelEnabled: &enabled})
	if err != nil {
		t.Fatalf("build query: %v", err)
	}

	for _, fragment := range []string{
		"LEFT JOIN LATERAL",
		"FROM channel_platforms cp",
		"cp.user_id = u.id",
		"cp.platform = u.platform",
		"ORDER BY cp.channel_id",
		"LIMIT 1",
		"cb.enabled",
	} {
		if !strings.Contains(query, fragment) {
			t.Fatalf("query does not contain %q: %s", fragment, query)
		}
	}
	if strings.Contains(query, "JOIN channels") {
		t.Fatalf("count query must not join legacy channels columns: %s", query)
	}
	if strings.Contains(query, "EXISTS") {
		t.Fatalf("count query must filter the selected binding, not any matching binding: %s", query)
	}
}

func TestBuildGetManyAndCountQueriesUseSameSelectedBindingForMultipleMatchingBindings(t *testing.T) {
	enabled := true
	pageQuery, _, err := buildGetManyQuery(users_with_channel.GetManyInput{ChannelEnabled: &enabled})
	if err != nil {
		t.Fatalf("build page query: %v", err)
	}
	countQuery, _, err := buildGetManyCountQuery(users_with_channel.GetManyInput{ChannelEnabled: &enabled})
	if err != nil {
		t.Fatalf("build count query: %v", err)
	}

	// A later enabled binding must not qualify a user when the first binding by channel_id is disabled.
	for name, query := range map[string]string{
		"page":  pageQuery,
		"count": countQuery,
	} {
		for _, fragment := range []string{
			"LEFT JOIN LATERAL",
			"FROM channel_platforms cp",
			"cp.user_id = u.id",
			"cp.platform = u.platform",
			"ORDER BY cp.channel_id",
			"LIMIT 1",
			"cb.enabled",
		} {
			if !strings.Contains(query, fragment) {
				t.Fatalf("%s query does not contain %q: %s", name, fragment, query)
			}
		}
	}

	if strings.Contains(countQuery, "EXISTS") {
		t.Fatalf("count query must filter the selected binding, not any matching binding: %s", countQuery)
	}
}
