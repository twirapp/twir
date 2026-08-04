package dashboard

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/twirapp/kv"
	"github.com/twirapp/twir/apps/api-gql/internal/entity"
	channelentity "github.com/twirapp/twir/libs/entities/channel"
	channelplatformentity "github.com/twirapp/twir/libs/entities/channel_platform"
	"github.com/twirapp/twir/libs/entities/platform"
	channelplatforms "github.com/twirapp/twir/libs/repositories/channel_platforms"
	channelsemotesusages "github.com/twirapp/twir/libs/repositories/channels_emotes_usages"
	channelsEmotesModel "github.com/twirapp/twir/libs/repositories/channels_emotes_usages/model"
	streammodel "github.com/twirapp/twir/libs/repositories/streams/model"
	usersmodel "github.com/twirapp/twir/libs/repositories/users/model"
)

type dashboardCurrentPlatformStub struct {
	platform string
	err      error
}

func (s dashboardCurrentPlatformStub) GetCurrentPlatform(context.Context) (string, error) {
	return s.platform, s.err
}

type dashboardChannelLookupStub struct {
	channel channelentity.Channel
}

func (s dashboardChannelLookupStub) GetChannelByID(context.Context, uuid.UUID) (channelentity.Channel, error) {
	return s.channel, nil
}

type dashboardUsersLookupStub struct {
	users map[uuid.UUID]usersmodel.User
}

type dashboardBindingPatch struct {
	id    uuid.UUID
	input channelplatforms.PatchInput
}

type dashboardBindingUpdaterStub struct {
	patches  []dashboardBindingPatch
	bindings []channelplatformentity.ChannelPlatform
}

func (s *dashboardBindingUpdaterStub) ListByChannelID(
	context.Context,
	uuid.UUID,
) ([]channelplatformentity.ChannelPlatform, error) {
	return s.bindings, nil
}

func (s *dashboardBindingUpdaterStub) Patch(
	_ context.Context,
	id uuid.UUID,
	input channelplatforms.PatchInput,
) (channelplatformentity.ChannelPlatform, error) {
	s.patches = append(s.patches, dashboardBindingPatch{id: id, input: input})
	return channelplatformentity.ChannelPlatform{ID: id}, nil
}

type dashboardStreamsRepositoryStub struct {
	streams map[platform.Platform]streammodel.Stream
}

func (s dashboardStreamsRepositoryStub) GetByChannelID(
	_ context.Context,
	_ uuid.UUID,
	platform platform.Platform,
) (streammodel.Stream, error) {
	stream, ok := s.streams[platform]
	if !ok {
		return streammodel.Nil, nil
	}

	return stream, nil
}

type dashboardEmotesUsageStub struct {
	counts map[string]uint64
}

func (s dashboardEmotesUsageStub) Count(
	_ context.Context,
	input channelsemotesusages.CountInput,
) (uint64, error) {
	if input.Platform == nil {
		return 0, nil
	}

	return s.counts[*input.Platform], nil
}

func (dashboardEmotesUsageStub) CreateMany(context.Context, []channelsemotesusages.ChannelEmoteUsageInput) error {
	return nil
}

func (dashboardEmotesUsageStub) GetEmotesStatistics(context.Context, channelsemotesusages.GetEmotesStatisticsInput) ([]channelsEmotesModel.EmoteStatistic, error) {
	return nil, nil
}

func (dashboardEmotesUsageStub) GetEmotesRanges(context.Context, string, string, []string, channelsemotesusages.EmoteStatisticRange) (map[string][]channelsEmotesModel.EmoteRange, error) {
	return nil, nil
}

func (dashboardEmotesUsageStub) GetChannelEmoteUsageHistory(context.Context, channelsemotesusages.EmotesUsersTopOrHistoryInput) ([]channelsEmotesModel.EmoteUsage, uint64, error) {
	return nil, 0, nil
}

func (dashboardEmotesUsageStub) GetChannelUsageTopUsers(context.Context, channelsemotesusages.EmotesUsersTopOrHistoryInput) ([]channelsEmotesModel.EmoteUsageTopUser, uint64, error) {
	return nil, 0, nil
}

func (dashboardEmotesUsageStub) DeleteRowsByChannelID(context.Context, string, string) error {
	return nil
}

func (dashboardEmotesUsageStub) GetUserMostUsedEmotes(context.Context, channelsemotesusages.UserMostUsedEmotesInput) ([]channelsEmotesModel.UserMostUsedEmote, error) {
	return nil, nil
}

type dashboardKVStub struct {
	values map[string]int64
}

func (s dashboardKVStub) Get(_ context.Context, key string) kv.Valuer {
	return dashboardKVValue{value: s.values[key]}
}

type dashboardKVValue struct {
	value int64
}

func (v dashboardKVValue) Int() (int64, error) {
	return v.value, nil
}

func (dashboardKVValue) String() (string, error) {
	return "", nil
}

func (dashboardKVValue) Bytes() ([]byte, error) {
	return nil, nil
}

func (dashboardKVValue) Bool() (bool, error) {
	return false, nil
}

func (dashboardKVValue) Float() (float64, error) {
	return 0, nil
}

func (dashboardKVValue) Scan(any) error {
	return nil
}

func (dashboardKVValue) Err() error {
	return nil
}

func newDashboardStatsService(
	channelID uuid.UUID,
	bindings []channelplatformentity.ChannelPlatform,
	streamRows map[platform.Platform]streammodel.Stream,
	messageCounts map[string]int64,
	emoteCounts map[string]uint64,
) *Service {
	return &Service{
		channelService:       dashboardChannelLookupStub{channel: channelentity.Channel{ID: channelID, Bindings: bindings}},
		channelPlatformsRepo: &dashboardBindingUpdaterStub{bindings: bindings},
		streamsRepository:    dashboardStreamsRepositoryStub{streams: streamRows},
		channelEmotesUsagesRepo: dashboardEmotesUsageStub{
			counts: emoteCounts,
		},
		requestedSongsRepo: dashboardRequestedSongsStub{},
		kv:                 dashboardKVStub{values: messageCounts},
		logger:             slog.Default(),
		authService: dashboardCurrentPlatformStub{
			err: errors.New("no current platform"),
		},
	}
}

type dashboardRequestedSongsStub struct{}

func (dashboardRequestedSongsStub) CountByChannelID(
	_ context.Context,
	_ string,
	_ time.Time,
) (int64, error) {
	return 0, nil
}

func liveDashboardStream(
	channelID uuid.UUID,
	platform platform.Platform,
	streamID string,
	viewerCount int,
	title string,
	categoryID string,
	categoryName string,
	startedAt time.Time,
) streammodel.Stream {
	return streammodel.Stream{
		ID:          streamID,
		ChannelID:   channelID,
		Platform:    platform,
		ViewerCount: viewerCount,
		Title:       title,
		GameId:      categoryID,
		GameName:    categoryName,
		StartedAt:   startedAt,
	}
}

type dashboardCacheInvalidatorStub struct {
	keys []string
}

func (s *dashboardCacheInvalidatorStub) Invalidate(_ context.Context, key string) error {
	s.keys = append(s.keys, key)
	return nil
}

func (s dashboardUsersLookupStub) GetByID(_ context.Context, id uuid.UUID) (usersmodel.User, error) {
	user, ok := s.users[id]
	if !ok {
		return usersmodel.Nil, errors.New("user not found")
	}

	return user, nil
}

func TestResolveAnalyticsIdentitySelectsCurrentPlatformBinding(t *testing.T) {
	service := &Service{authService: dashboardCurrentPlatformStub{platform: platform.PlatformKick.String()}}
	channel := channelentity.Channel{
		Bindings: []channelplatformentity.ChannelPlatform{
			{Platform: platform.PlatformTwitch, PlatformChannelID: "twitch-channel"},
			{Platform: platform.PlatformKick, PlatformChannelID: "kick-channel"},
		},
	}

	gotPlatform, gotChannelID := service.resolveAnalyticsIdentity(context.Background(), channel)
	if gotPlatform != platform.PlatformKick.String() || gotChannelID != "kick-channel" {
		t.Fatalf("analytics identity = (%q, %q), want Kick binding", gotPlatform, gotChannelID)
	}
}

func TestResolveAnalyticsIdentityFallsBackToTwitchBinding(t *testing.T) {
	service := &Service{authService: dashboardCurrentPlatformStub{err: errors.New("no current platform")}}
	channel := channelentity.Channel{
		Bindings: []channelplatformentity.ChannelPlatform{
			{Platform: platform.PlatformKick, PlatformChannelID: "kick-channel"},
			{Platform: platform.PlatformTwitch, PlatformChannelID: "twitch-channel"},
		},
	}

	gotPlatform, gotChannelID := service.resolveAnalyticsIdentity(context.Background(), channel)
	if gotPlatform != platform.PlatformTwitch.String() || gotChannelID != "twitch-channel" {
		t.Fatalf("analytics identity = (%q, %q), want Twitch fallback", gotPlatform, gotChannelID)
	}
}

func TestGetBotStatusesMapsKickBindingIdentity(t *testing.T) {
	channelID := uuid.New()
	kickOwnerID := uuid.New()
	kickBotUserID := uuid.New()
	service := &Service{
		channelService: dashboardChannelLookupStub{channel: channelentity.Channel{
			ID: channelID,
			Bindings: []channelplatformentity.ChannelPlatform{{
				Platform:          platform.PlatformKick,
				UserID:            kickOwnerID,
				PlatformChannelID: "kick-channel",
				Enabled:           true,
				BotUserID:         &kickBotUserID,
			}},
		}},
		usersRepo: dashboardUsersLookupStub{users: map[uuid.UUID]usersmodel.User{
			kickOwnerID:   {ID: kickOwnerID, Login: "kick-owner"},
			kickBotUserID: {ID: kickBotUserID, Login: "kick-bot"},
		}},
	}

	statuses, err := service.GetBotStatuses(context.Background(), channelID.String())
	if err != nil {
		t.Fatalf("get bot statuses: %v", err)
	}
	if len(statuses) != 1 {
		t.Fatalf("statuses = %d, want 1", len(statuses))
	}
	status := statuses[0]
	if status.Platform != platform.PlatformKick.String() || status.ChannelName != "kick-owner" {
		t.Fatalf("Kick status = %#v, want binding owner identity", status)
	}
	if !status.Enabled || !status.IsMod {
		t.Fatalf("Kick status = %#v, want enabled moderator", status)
	}
	if status.BotID != kickBotUserID.String() || status.BotName != "kick-bot" {
		t.Fatalf("Kick status = %#v, want binding bot identity", status)
	}
}

func TestGetBotStatusesMapsVKVideoLiveBindingIdentity(t *testing.T) {
	// Given
	channelID := uuid.New()
	vkOwnerID := uuid.New()
	vkBotUserID := uuid.New()
	service := &Service{
		channelService: dashboardChannelLookupStub{channel: channelentity.Channel{
			ID: channelID,
			Bindings: []channelplatformentity.ChannelPlatform{{
				Platform:          platform.PlatformVKVideoLive,
				UserID:            vkOwnerID,
				PlatformChannelID: "vk-channel",
				Enabled:           true,
				BotUserID:         &vkBotUserID,
			}},
		}},
		usersRepo: dashboardUsersLookupStub{users: map[uuid.UUID]usersmodel.User{
			vkOwnerID:   {ID: vkOwnerID, Login: "vk-owner"},
			vkBotUserID: {ID: vkBotUserID, Login: "vk-bot"},
		}},
	}

	// When
	statuses, err := service.GetBotStatuses(context.Background(), channelID.String())

	// Then
	if err != nil {
		t.Fatalf("get bot statuses: %v", err)
	}
	if len(statuses) != 1 {
		t.Fatalf("statuses = %d, want 1", len(statuses))
	}
	status := statuses[0]
	if status.Platform != platform.PlatformVKVideoLive.String() || status.ChannelName != "vk-owner" {
		t.Fatalf("VK Video Live status = %#v, want binding owner identity", status)
	}
	if !status.Enabled {
		t.Fatalf("VK Video Live status = %#v, want enabled binding", status)
	}
	if status.BotID != vkBotUserID.String() || status.BotName != "vk-bot" {
		t.Fatalf("VK Video Live status = %#v, want global bot identity", status)
	}
}

func TestGetBotStatusesDoesNotExposeUsableVKVideoLiveBotWithoutBotUserID(t *testing.T) {
	// Given
	channelID := uuid.New()
	vkOwnerID := uuid.New()
	service := &Service{
		channelService: dashboardChannelLookupStub{channel: channelentity.Channel{
			ID: channelID,
			Bindings: []channelplatformentity.ChannelPlatform{{
				Platform:          platform.PlatformVKVideoLive,
				UserID:            vkOwnerID,
				PlatformChannelID: "vk-channel",
				Enabled:           true,
			}},
		}},
		usersRepo: dashboardUsersLookupStub{users: map[uuid.UUID]usersmodel.User{
			vkOwnerID: {ID: vkOwnerID, Login: "vk-owner"},
		}},
	}

	// When
	statuses, err := service.GetBotStatuses(context.Background(), channelID.String())

	// Then
	if err != nil {
		t.Fatalf("get bot statuses: %v", err)
	}
	if len(statuses) != 1 {
		t.Fatalf("statuses = %d, want 1", len(statuses))
	}
	status := statuses[0]
	if status.Platform != platform.PlatformVKVideoLive.String() || status.ChannelName != "vk-owner" {
		t.Fatalf("VK Video Live status = %#v, want binding owner identity", status)
	}
	if status.BotID != "" || status.BotName != "" {
		t.Fatalf("VK Video Live status = %#v, want no usable bot identity", status)
	}
}

func TestGetBotStatusesMapsYouTubeBindingIdentity(t *testing.T) {
	// Given
	channelID := uuid.New()
	youtubeOwnerID := uuid.New()
	youtubeBotUserID := uuid.New()
	service := &Service{
		channelService: dashboardChannelLookupStub{channel: channelentity.Channel{
			ID: channelID,
			Bindings: []channelplatformentity.ChannelPlatform{{
				Platform:          platform.PlatformYouTube,
				UserID:            youtubeOwnerID,
				PlatformChannelID: "youtube-channel",
				Enabled:           true,
				BotUserID:         &youtubeBotUserID,
			}},
		}},
		usersRepo: dashboardUsersLookupStub{users: map[uuid.UUID]usersmodel.User{
			youtubeOwnerID:   {ID: youtubeOwnerID, Login: "youtube-owner"},
			youtubeBotUserID: {ID: youtubeBotUserID, Login: "youtube-bot"},
		}},
	}

	// When
	statuses, err := service.GetBotStatuses(context.Background(), channelID.String())

	// Then
	if err != nil {
		t.Fatalf("get bot statuses: %v", err)
	}
	if len(statuses) != 1 {
		t.Fatalf("statuses = %d, want 1", len(statuses))
	}
	status := statuses[0]
	if status.Platform != platform.PlatformYouTube.String() || status.ChannelName != "youtube-owner" {
		t.Fatalf("YouTube status = %#v, want binding owner identity", status)
	}
	if !status.Enabled || status.BotID != youtubeBotUserID.String() || status.BotName != "youtube-bot" {
		t.Fatalf("YouTube status = %#v, want enabled binding and bot identity", status)
	}
}

func TestGetBasicTwitchBotStatusUsesTwitchBindingConfig(t *testing.T) {
	twitchOwnerID := uuid.New()
	channel := channelentity.Channel{
		ID: uuid.New(),
		Bindings: []channelplatformentity.ChannelPlatform{
			{Platform: platform.PlatformKick, Enabled: false},
			{
				Platform:          platform.PlatformTwitch,
				UserID:            twitchOwnerID,
				PlatformChannelID: "twitch-channel",
				Enabled:           true,
				BotConfig: json.RawMessage(
					`{"bot_id":"twitch-bot","is_bot_mod":true,"is_twitch_banned":true}`,
				),
			},
		},
	}
	binding, config, found, err := channel.TwitchBinding()
	if err != nil {
		t.Fatalf("find Twitch binding: %v", err)
	}
	if !found {
		t.Fatal("expected Twitch binding")
	}

	service := &Service{usersRepo: dashboardUsersLookupStub{users: map[uuid.UUID]usersmodel.User{
		twitchOwnerID: {ID: twitchOwnerID, Login: "twitch-owner"},
	}}}
	status := service.getBasicTwitchBotStatus(context.Background(), channel, binding, config)
	if status.Platform != platform.PlatformTwitch.String() || status.ChannelName != "twitch-owner" {
		t.Fatalf("Twitch status = %#v, want binding owner identity", status)
	}
	if !status.Enabled || !status.IsMod || status.BotID != "twitch-bot" || status.BotName != "TwirBot" {
		t.Fatalf("Twitch status = %#v, want parsed Twitch config", status)
	}
}

func TestBotJoinLeaveUpdatesOnlyVKVideoLiveBindingWithoutEventSub(t *testing.T) {
	for _, test := range []struct {
		name        string
		action      string
		platform    string
		includeKick bool
		wantEnabled bool
	}{
		{name: "joins", action: BotJoinLeaveActionJoin, platform: platform.PlatformVKVideoLive.String(), includeKick: true, wantEnabled: true},
		{name: "leaves", action: BotJoinLeaveActionLeave, platform: platform.PlatformVKVideoLive.String(), includeKick: true, wantEnabled: false},
		{name: "defaults to VK binding", action: BotJoinLeaveActionJoin, wantEnabled: true},
	} {
		t.Run(test.name, func(t *testing.T) {
			// Given
			channelID := uuid.New()
			vkBindingID := uuid.New()
			vkBotUserID := uuid.New()
			updater := &dashboardBindingUpdaterStub{}
			cache := &dashboardCacheInvalidatorStub{}
			bindings := []channelplatformentity.ChannelPlatform{
				{ID: vkBindingID, Platform: platform.PlatformVKVideoLive, PlatformChannelID: "vk-channel", BotUserID: &vkBotUserID},
			}
			if test.includeKick {
				bindings = append(bindings, channelplatformentity.ChannelPlatform{ID: uuid.New(), Platform: platform.PlatformKick, Enabled: true})
			}
			service := &Service{
				channelService: dashboardChannelLookupStub{channel: channelentity.Channel{
					ID:       channelID,
					Bindings: bindings,
				}},
				channelPlatformsRepo: updater,
				channelsCache:        cache,
			}

			// When
			success, err := service.BotJoinLeave(context.Background(), channelID.String(), test.action, test.platform)

			// Then
			if err != nil {
				t.Fatalf("bot join/leave: %v", err)
			}
			if !success {
				t.Fatal("bot join/leave = false, want true")
			}
			if len(updater.patches) != 1 {
				t.Fatalf("patches = %d, want 1", len(updater.patches))
			}
			patch := updater.patches[0]
			if patch.id != vkBindingID || patch.input.Enabled == nil || *patch.input.Enabled != test.wantEnabled {
				t.Fatalf("patch = %#v, want VK binding enabled=%t", patch, test.wantEnabled)
			}
			if len(cache.keys) != 1 || cache.keys[0] != channelID.String() {
				t.Fatalf("cache invalidations = %#v, want [%q]", cache.keys, channelID.String())
			}
		})
	}
}

func TestBotJoinLeaveRejectsVKVideoLiveBindingWithoutBotUserID(t *testing.T) {
	// Given
	channelID := uuid.New()
	updater := &dashboardBindingUpdaterStub{}
	service := &Service{
		channelService: dashboardChannelLookupStub{channel: channelentity.Channel{
			ID: channelID,
			Bindings: []channelplatformentity.ChannelPlatform{{
				ID:                uuid.New(),
				Platform:          platform.PlatformVKVideoLive,
				PlatformChannelID: "vk-channel",
			}},
		}},
		channelPlatformsRepo: updater,
	}

	// When
	success, err := service.BotJoinLeave(context.Background(), channelID.String(), BotJoinLeaveActionJoin, platform.PlatformVKVideoLive.String())

	// Then
	if err == nil {
		t.Fatal("bot join/leave error = nil, want missing VK bot error")
	}
	if success {
		t.Fatal("bot join/leave = true, want false")
	}
	if len(updater.patches) != 0 {
		t.Fatalf("patches = %d, want 0", len(updater.patches))
	}
}

func TestBotJoinLeaveUpdatesYouTubeBinding(t *testing.T) {
	// Given
	channelID := uuid.New()
	bindingID := uuid.New()
	botUserID := uuid.New()
	updater := &dashboardBindingUpdaterStub{}
	cache := &dashboardCacheInvalidatorStub{}
	service := &Service{
		channelService: dashboardChannelLookupStub{channel: channelentity.Channel{
			ID: channelID,
			Bindings: []channelplatformentity.ChannelPlatform{{
				ID:                bindingID,
				Platform:          platform.PlatformYouTube,
				PlatformChannelID: "youtube-channel",
				BotUserID:         &botUserID,
			}},
		}},
		channelPlatformsRepo: updater,
		channelsCache:        cache,
	}

	// When
	success, err := service.BotJoinLeave(context.Background(), channelID.String(), BotJoinLeaveActionJoin, platform.PlatformYouTube.String())

	// Then
	if err != nil {
		t.Fatalf("bot join: %v", err)
	}
	if !success {
		t.Fatal("bot join = false, want true")
	}
	if len(updater.patches) != 1 || updater.patches[0].id != bindingID || updater.patches[0].input.Enabled == nil || !*updater.patches[0].input.Enabled {
		t.Fatalf("patches = %#v, want YouTube binding enabled", updater.patches)
	}
	if len(cache.keys) != 1 || cache.keys[0] != channelID.String() {
		t.Fatalf("cache invalidations = %#v, want [%q]", cache.keys, channelID.String())
	}
}

func TestGetDashboardStatsIncludesLiveStatsForEnabledTwitchAndKickBindings(t *testing.T) {
	// Given
	channelID := uuid.New()
	twitchBinding := channelplatformentity.ChannelPlatform{
		Platform: platform.PlatformTwitch,
		Enabled:  true,
	}
	kickBinding := channelplatformentity.ChannelPlatform{
		Platform:          platform.PlatformKick,
		PlatformChannelID: "kick-channel",
		Enabled:           true,
	}
	twitchStartedAt := time.Date(2026, 7, 29, 10, 0, 0, 0, time.UTC)
	kickStartedAt := twitchStartedAt.Add(time.Hour)
	service := newDashboardStatsService(
		channelID,
		[]channelplatformentity.ChannelPlatform{twitchBinding, kickBinding},
		map[platform.Platform]streammodel.Stream{
			platform.PlatformTwitch: liveDashboardStream(
				channelID,
				platform.PlatformTwitch,
				"twitch-stream",
				101,
				"Twitch title",
				"twitch-category-id",
				"Twitch category",
				twitchStartedAt,
			),
			platform.PlatformKick: liveDashboardStream(
				channelID,
				platform.PlatformKick,
				"kick-stream",
				202,
				"Kick title",
				"kick-category-id",
				"Kick category",
				kickStartedAt,
			),
		},
		map[string]int64{
			"stream:parsedMessages:twitch-stream": 11,
			"stream:parsedMessages:kick-stream":   22,
		},
		map[string]uint64{
			platform.PlatformTwitch.String(): 3,
			platform.PlatformKick.String():   7,
		},
	)

	// When
	stats, err := service.GetDashboardStats(context.Background(), channelID.String())

	// Then
	if err != nil {
		t.Fatalf("get dashboard stats: %v", err)
	}
	if len(stats.Platforms) != 2 {
		t.Fatalf("platforms = %d, want 2", len(stats.Platforms))
	}

	for _, test := range []struct {
		platform    platform.Platform
		viewers     int
		title       string
		categoryID  string
		category    string
		startedAt   time.Time
		chat        int
		usedEmotes  int
		canEditInfo bool
	}{
		{
			platform:    platform.PlatformTwitch,
			viewers:     101,
			title:       "Twitch title",
			categoryID:  "twitch-category-id",
			category:    "Twitch category",
			startedAt:   twitchStartedAt,
			chat:        11,
			usedEmotes:  3,
			canEditInfo: true,
		},
		{
			platform:    platform.PlatformKick,
			viewers:     202,
			title:       "Kick title",
			categoryID:  "kick-category-id",
			category:    "Kick category",
			startedAt:   kickStartedAt,
			chat:        22,
			usedEmotes:  7,
			canEditInfo: true,
		},
	} {
		var got *entity.PlatformStats
		for i := range stats.Platforms {
			if stats.Platforms[i].Platform == test.platform {
				got = &stats.Platforms[i]
				break
			}
		}
		if got == nil {
			t.Fatalf("missing platform %s", test.platform)
		}
		if !got.IsLive || got.Viewers == nil || *got.Viewers != test.viewers {
			t.Fatalf("%s viewers = %#v, want live %d", test.platform, got.Viewers, test.viewers)
		}
		if got.Title == nil || *got.Title != test.title || got.CategoryID == nil || *got.CategoryID != test.categoryID || got.CategoryName == nil || *got.CategoryName != test.category {
			t.Fatalf("%s metadata = %#v, want live metadata", test.platform, got)
		}
		if got.StartedAt == nil || !got.StartedAt.Equal(test.startedAt) {
			t.Fatalf("%s startedAt = %#v, want %s", test.platform, got.StartedAt, test.startedAt)
		}
		if got.ChatMessages != test.chat || got.UsedEmotes != test.usedEmotes || got.CanEditInfo != test.canEditInfo {
			t.Fatalf("%s counters/editability = %#v, want chat=%d emotes=%d canEdit=%t", test.platform, got, test.chat, test.usedEmotes, test.canEditInfo)
		}
	}
}

func TestGetDashboardStatsIncludesOfflineVKBindingWithNullableStats(t *testing.T) {
	// Given
	channelID := uuid.New()
	service := newDashboardStatsService(
		channelID,
		[]channelplatformentity.ChannelPlatform{{
			Platform: platform.PlatformVKVideoLive,
			Enabled:  true,
		}},
		nil,
		nil,
		nil,
	)

	// When
	stats, err := service.GetDashboardStats(context.Background(), channelID.String())

	// Then
	if err != nil {
		t.Fatalf("get dashboard stats: %v", err)
	}
	if len(stats.Platforms) != 1 {
		t.Fatalf("platforms = %d, want 1", len(stats.Platforms))
	}
	got := stats.Platforms[0]
	if got.Platform != platform.PlatformVKVideoLive || got.IsLive {
		t.Fatalf("VK platform stats = %#v, want offline VK", got)
	}
	if got.Viewers != nil || got.Followers != nil || got.Title != nil || got.CategoryID != nil || got.CategoryName != nil || got.StartedAt != nil {
		t.Fatalf("VK offline nullable fields = %#v, want nil", got)
	}
	if got.CanEditInfo || got.ChatMessages != 0 || got.UsedEmotes != 0 {
		t.Fatalf("VK offline stats = %#v, want zero counters and no editing", got)
	}
}

func TestGetDashboardStatsListsAllEnabledBindingsWhenNoStreamsExist(t *testing.T) {
	// Given
	channelID := uuid.New()
	bindings := []channelplatformentity.ChannelPlatform{
		{Platform: platform.PlatformTwitch, Enabled: true},
		{Platform: platform.PlatformKick, Enabled: true},
		{Platform: platform.PlatformVKVideoLive, Enabled: false},
	}
	service := newDashboardStatsService(channelID, bindings, nil, nil, nil)

	// When
	stats, err := service.GetDashboardStats(context.Background(), channelID.String())

	// Then
	if err != nil {
		t.Fatalf("get dashboard stats: %v", err)
	}
	if len(stats.Platforms) != 2 {
		t.Fatalf("platforms = %d, want enabled bindings only", len(stats.Platforms))
	}
	if stats.StreamCategoryID != "" || stats.StreamCategoryName != "" || stats.StreamViewers != nil || stats.StreamStartedAt != nil || stats.StreamTitle != "" || stats.StreamChatMessages != 0 || stats.Followers != 0 || stats.UsedEmotes != 0 || stats.RequestedSongs != 0 || stats.Subs != 0 {
		t.Fatalf("top-level stats = %#v, want unchanged empty fallback", stats)
	}
	for _, got := range stats.Platforms {
		if got.IsLive || got.Viewers != nil || got.ChatMessages != 0 || got.UsedEmotes != 0 {
			t.Fatalf("offline platform stats = %#v, want offline zero values", got)
		}
		if got.Platform == platform.PlatformVKVideoLive && got.Followers != nil {
			t.Fatalf("VK offline followers = %#v, want nil", got.Followers)
		}
	}
}
