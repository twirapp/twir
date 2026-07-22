package handler

import (
	"context"
	"encoding/json"
	"log/slog"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/kvizyx/twitchy/eventsub"
	kvinmemory "github.com/twirapp/kv/stores/inmemory"
	buscore "github.com/twirapp/twir/libs/bus-core"
	genericcacher "github.com/twirapp/twir/libs/cache/generic-cacher"
	cfg "github.com/twirapp/twir/libs/config"
	"github.com/twirapp/twir/libs/entities/platform"
	channelplatforms "github.com/twirapp/twir/libs/repositories/channel_platforms"
	channelplatformsmodel "github.com/twirapp/twir/libs/repositories/channel_platforms/model"
	channelsrepo "github.com/twirapp/twir/libs/repositories/channels"
	channelsmodel "github.com/twirapp/twir/libs/repositories/channels/model"
	usersrepo "github.com/twirapp/twir/libs/repositories/users"
	usersmodel "github.com/twirapp/twir/libs/repositories/users/model"
	channelservice "github.com/twirapp/twir/libs/services/channels"
)

type handlerStatusChannelsRepo struct {
	channel     channelsmodel.Channel
	updateCalls int
}

func (r *handlerStatusChannelsRepo) GetMany(context.Context, channelsrepo.GetManyInput) ([]channelsmodel.Channel, error) {
	return nil, nil
}

func (r *handlerStatusChannelsRepo) GetAllByBindingPlatform(
	context.Context,
	platform.Platform,
) ([]channelsmodel.Channel, error) {
	return nil, nil
}

func (r *handlerStatusChannelsRepo) GetByID(context.Context, uuid.UUID) (channelsmodel.Channel, error) {
	return r.channel, nil
}

func (r *handlerStatusChannelsRepo) GetByApiKey(context.Context, string) (channelsmodel.Channel, error) {
	return r.channel, nil
}

func (r *handlerStatusChannelsRepo) GetByBindingUserID(
	context.Context,
	platform.Platform,
	uuid.UUID,
) (channelsmodel.Channel, error) {
	return r.channel, nil
}

func (r *handlerStatusChannelsRepo) GetByPlatformChannelID(
	context.Context,
	platform.Platform,
	string,
) (channelsmodel.Channel, error) {
	return r.channel, nil
}

func (r *handlerStatusChannelsRepo) GetBySlug(
	context.Context,
	channelsrepo.GetBySlugInput,
) (channelsmodel.Channel, error) {
	return r.channel, nil
}

func (r *handlerStatusChannelsRepo) GetCount(context.Context, channelsrepo.GetCountInput) (int, error) {
	return 0, nil
}

func (r *handlerStatusChannelsRepo) Update(
	context.Context,
	uuid.UUID,
	channelsrepo.UpdateInput,
) (channelsmodel.Channel, error) {
	r.updateCalls++
	return r.channel, nil
}

func (r *handlerStatusChannelsRepo) Create(
	context.Context,
	channelsrepo.CreateInput,
) (channelsmodel.Channel, error) {
	return r.channel, nil
}

type handlerStatusUsersRepo struct {
	user usersmodel.User
}

func (r *handlerStatusUsersRepo) GetByID(context.Context, uuid.UUID) (usersmodel.User, error) {
	return r.user, nil
}

func (r *handlerStatusUsersRepo) GetByPlatformID(
	context.Context,
	platform.Platform,
	string,
) (usersmodel.User, error) {
	return r.user, nil
}

func (r *handlerStatusUsersRepo) GetManyByIDS(
	context.Context,
	usersrepo.GetManyInput,
) ([]usersmodel.User, error) {
	return []usersmodel.User{r.user}, nil
}

func (r *handlerStatusUsersRepo) Update(
	context.Context,
	uuid.UUID,
	usersrepo.UpdateInput,
) (usersmodel.User, error) {
	return r.user, nil
}

func (r *handlerStatusUsersRepo) GetRandomOnlineUser(
	context.Context,
	usersrepo.GetRandomOnlineUserInput,
) (usersmodel.OnlineUser, error) {
	return usersmodel.NilOnlineUser, nil
}

func (r *handlerStatusUsersRepo) GetOnlineUsersWithFilters(
	context.Context,
	usersrepo.GetOnlineUsersWithFiltersInput,
) ([]usersmodel.OnlineUser, error) {
	return nil, nil
}

func (r *handlerStatusUsersRepo) GetByApiKey(context.Context, string) (usersmodel.User, error) {
	return r.user, nil
}

func (r *handlerStatusUsersRepo) Create(
	context.Context,
	usersrepo.CreateInput,
) (usersmodel.User, error) {
	return r.user, nil
}

type handlerStatusPatch struct {
	ctx   context.Context
	id    uuid.UUID
	input channelplatforms.PatchInput
}

type handlerStatusBindingsRepo struct {
	patches []handlerStatusPatch
}

func (r *handlerStatusBindingsRepo) Create(
	context.Context,
	channelplatforms.CreateInput,
) (channelplatformsmodel.ChannelPlatform, error) {
	return channelplatformsmodel.Nil, nil
}

func (r *handlerStatusBindingsRepo) GetByChannelAndPlatform(
	context.Context,
	uuid.UUID,
	platform.Platform,
) (channelplatformsmodel.ChannelPlatform, error) {
	return channelplatformsmodel.Nil, nil
}

func (r *handlerStatusBindingsRepo) GetByPlatformChannelID(
	context.Context,
	platform.Platform,
	string,
) (channelplatformsmodel.ChannelPlatform, error) {
	return channelplatformsmodel.Nil, nil
}

func (r *handlerStatusBindingsRepo) ListByChannelID(
	context.Context,
	uuid.UUID,
) ([]channelplatformsmodel.ChannelPlatform, error) {
	return nil, nil
}

func (r *handlerStatusBindingsRepo) Update(
	context.Context,
	uuid.UUID,
	channelplatforms.UpdateInput,
) (channelplatformsmodel.ChannelPlatform, error) {
	return channelplatformsmodel.Nil, nil
}

func (r *handlerStatusBindingsRepo) Patch(
	ctx context.Context,
	id uuid.UUID,
	input channelplatforms.PatchInput,
) (channelplatformsmodel.ChannelPlatform, error) {
	r.patches = append(r.patches, handlerStatusPatch{ctx: ctx, id: id, input: input})
	return channelplatformsmodel.ChannelPlatform{ID: id}, nil
}

func (r *handlerStatusBindingsRepo) Delete(context.Context, uuid.UUID) error {
	return nil
}

type handlerStatusContextKey struct{}

func newHandlerForBindingStatusTests(t *testing.T) (
	*Handler,
	*handlerStatusChannelsRepo,
	*handlerStatusBindingsRepo,
	channelsmodel.Channel,
	channelplatformsmodel.ChannelPlatform,
) {
	t.Helper()

	channelID := uuid.New()
	twitchBinding := channelplatformsmodel.ChannelPlatform{
		ID:                uuid.New(),
		ChannelID:         channelID,
		Platform:          platform.PlatformTwitch,
		UserID:            uuid.New(),
		PlatformChannelID: "twitch-channel",
		Enabled:           true,
		BotConfig:         json.RawMessage(`{"bot_id":"twir-bot","is_bot_mod":true,"keep":"value"}`),
	}
	kickBinding := channelplatformsmodel.ChannelPlatform{
		ID:                uuid.New(),
		ChannelID:         channelID,
		Platform:          platform.PlatformKick,
		UserID:            uuid.New(),
		PlatformChannelID: "kick-channel",
		Enabled:           true,
		BotConfig:         json.RawMessage(`{"keep":"kick-value"}`),
	}
	channel := channelsmodel.Channel{
		ID:       channelID,
		Bindings: []channelplatformsmodel.ChannelPlatform{kickBinding, twitchBinding},
	}
	channelsRepo := &handlerStatusChannelsRepo{channel: channel}
	usersRepo := &handlerStatusUsersRepo{user: usersmodel.User{ID: twitchBinding.UserID}}
	bindingsRepo := &handlerStatusBindingsRepo{}
	channelsCache := genericcacher.New(genericcacher.Opts[channelsmodel.Channel]{
		KV:        kvinmemory.New(),
		KeyPrefix: "test:channels:",
		Ttl:       time.Minute,
		LoadFn: func(context.Context, string) (channelsmodel.Channel, error) {
			return channelsmodel.Nil, nil
		},
	})

	return &Handler{
		logger:               slog.Default(),
		channelsRepo:         channelsRepo,
		channelPlatformsRepo: bindingsRepo,
		channelsCache:        channelsCache,
		usersRepo:            usersRepo,
		channelService: channelservice.NewChannelService(
			channelsRepo,
			&buscore.Bus{},
			cfg.Config{},
			kvinmemory.New(),
			nil,
		),
	}, channelsRepo, bindingsRepo, channel, twitchBinding
}

func assertTwitchBindingPatch(
	t *testing.T,
	patches []handlerStatusPatch,
	wantContext context.Context,
	wantBindingID uuid.UUID,
	wantEnabled *bool,
	wantBotMod *bool,
) {
	t.Helper()

	if len(patches) != 1 {
		t.Fatalf("binding patches = %d, want 1", len(patches))
	}
	patch := patches[0]
	if patch.id != wantBindingID {
		t.Fatalf("patched binding ID = %s, want Twitch binding %s", patch.id, wantBindingID)
	}
	if patch.ctx.Value(handlerStatusContextKey{}) != wantContext.Value(handlerStatusContextKey{}) {
		t.Fatal("binding patch did not receive request context")
	}
	if wantEnabled == nil {
		if patch.input.Enabled != nil {
			t.Fatalf("enabled patch = %v, want nil", *patch.input.Enabled)
		}
	} else if patch.input.Enabled == nil || *patch.input.Enabled != *wantEnabled {
		t.Fatalf("enabled patch = %v, want %v", patch.input.Enabled, wantEnabled)
	}
	if wantBotMod == nil {
		if len(patch.input.BotConfigPatch) != 0 {
			t.Fatalf("bot config patch = %s, want empty", patch.input.BotConfigPatch)
		}
		return
	}

	var botConfigPatch map[string]bool
	if err := json.Unmarshal(patch.input.BotConfigPatch, &botConfigPatch); err != nil {
		t.Fatalf("unmarshal bot config patch: %v", err)
	}
	if len(botConfigPatch) != 1 || botConfigPatch["is_bot_mod"] != *wantBotMod {
		t.Fatalf("bot config patch = %v, want is_bot_mod=%t", botConfigPatch, *wantBotMod)
	}
}

func TestDisableTwitchBindingForBannedBotPatchesTwitchBinding(t *testing.T) {
	handler, channelsRepo, bindingsRepo, _, twitchBinding := newHandlerForBindingStatusTests(t)
	ctx := context.WithValue(context.Background(), handlerStatusContextKey{}, "ban")

	handler.disableTwitchBindingForBannedBot(ctx, "twitch-channel", "twir-bot")

	wantEnabled := false
	assertTwitchBindingPatch(t, bindingsRepo.patches, ctx, twitchBinding.ID, &wantEnabled, nil)
	if channelsRepo.updateCalls != 0 {
		t.Fatalf("legacy channel Update calls = %d, want 0", channelsRepo.updateCalls)
	}
}

func TestHandleUserAuthorizationRevokePatchesTwitchBinding(t *testing.T) {
	handler, channelsRepo, bindingsRepo, _, twitchBinding := newHandlerForBindingStatusTests(t)
	ctx := context.WithValue(context.Background(), handlerStatusContextKey{}, "revoke")

	handler.HandleUserAuthorizationRevoke(
		ctx,
		eventsub.UserAuthorizationRevokeEvent{UserId: "twitch-channel"},
		eventsub.WebsocketNotificationMetadata{},
	)

	wantEnabled := false
	wantBotMod := false
	assertTwitchBindingPatch(t, bindingsRepo.patches, ctx, twitchBinding.ID, &wantEnabled, &wantBotMod)
	if channelsRepo.updateCalls != 0 {
		t.Fatalf("legacy channel Update calls = %d, want 0", channelsRepo.updateCalls)
	}
}

func TestUpdateBotStatusPatchesTwitchBinding(t *testing.T) {
	handler, channelsRepo, bindingsRepo, _, twitchBinding := newHandlerForBindingStatusTests(t)
	ctx := context.WithValue(context.Background(), handlerStatusContextKey{}, "moderator")

	handler.updateBotStatus(ctx, "twitch-channel", true)

	wantBotMod := true
	assertTwitchBindingPatch(t, bindingsRepo.patches, ctx, twitchBinding.ID, nil, &wantBotMod)
	if channelsRepo.updateCalls != 0 {
		t.Fatalf("legacy channel Update calls = %d, want 0", channelsRepo.updateCalls)
	}
}
