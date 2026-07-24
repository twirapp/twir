package channel

import (
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"sync"
	"testing"

	"github.com/alitto/pond/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/twirapp/twir/apps/bots/internal/twitchactions"
	"github.com/twirapp/twir/apps/bots/internal/workers"
	"github.com/twirapp/twir/libs/bus-core/bots"
	cfg "github.com/twirapp/twir/libs/config"
	channelentity "github.com/twirapp/twir/libs/entities/channel"
	channelplatformentity "github.com/twirapp/twir/libs/entities/channel_platform"
	"github.com/twirapp/twir/libs/entities/platform"
	channelsrepository "github.com/twirapp/twir/libs/repositories/channels"
	usersrepository "github.com/twirapp/twir/libs/repositories/users"
	usersmodel "github.com/twirapp/twir/libs/repositories/users/model"
	channelservice "github.com/twirapp/twir/libs/services/channels"
)

func TestDeleteMessageUsesSelectedTwitchBinding(t *testing.T) {
	const (
		requestChannelID  = "incoming-twitch-channel"
		selectedChannelID = "selected-twitch-channel"
		selectedBotID     = "selected-twitch-bot"
		messageID         = "message-id"
	)

	channelUserID := uuid.New()
	validTwitchBinding := channelplatformentity.ChannelPlatform{
		Platform:          platform.PlatformTwitch,
		PlatformChannelID: selectedChannelID,
		UserID:            channelUserID,
		Enabled:           true,
		BotConfig: json.RawMessage(
			`{"bot_id":"selected-twitch-bot","is_bot_mod":true,"is_twitch_banned":false}`,
		),
	}

	tests := []struct {
		name       string
		bindings   []channelplatformentity.ChannelPlatform
		wantErr    bool
		wantDelete bool
	}{
		{
			name: "uses selected Twitch binding when Kick comes first",
			bindings: []channelplatformentity.ChannelPlatform{
				{
					Platform:          platform.PlatformKick,
					PlatformChannelID: "kick-channel",
					UserID:            uuid.New(),
					Enabled:           true,
				},
				validTwitchBinding,
			},
			wantDelete: true,
		},
		{
			name: "skips missing Twitch binding",
			bindings: []channelplatformentity.ChannelPlatform{
				{
					Platform:          platform.PlatformKick,
					PlatformChannelID: "kick-channel",
					UserID:            uuid.New(),
					Enabled:           true,
				},
			},
		},
		{
			name: "skips missing Twitch bot config",
			bindings: []channelplatformentity.ChannelPlatform{
				{
					Platform:          platform.PlatformTwitch,
					PlatformChannelID: selectedChannelID,
					UserID:            channelUserID,
					Enabled:           true,
				},
			},
		},
		{
			name: "returns error for malformed Twitch bot config",
			bindings: []channelplatformentity.ChannelPlatform{
				{
					Platform:          platform.PlatformTwitch,
					PlatformChannelID: selectedChannelID,
					UserID:            channelUserID,
					Enabled:           true,
					BotConfig:         json.RawMessage(`{`),
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			usersRepo := &deleteMessageUsersRepository{
				user: usersmodel.User{ID: channelUserID},
			}
			channelsRepo := &deleteMessageChannelsRepository{
				channel: channelentity.Channel{Bindings: tt.bindings},
			}
			actions := &deleteMessageTwitchActions{}
			pool := &workers.Pool{Pool: pond.NewPool(1)}
			t.Cleanup(pool.StopAndWait)

			service := &Service{
				logger:        slog.New(slog.NewTextHandler(io.Discard, nil)),
				twitchActions: actions,
				workersPool:   pool,
				channelService: channelservice.NewChannelService(
					channelsRepo,
					nil,
					cfg.Config{},
					nil,
					nil,
				),
				usersRepo: usersRepo,
			}

			err := service.DeleteMessage(
				context.Background(),
				bots.DeleteMessageRequest{
					ChannelId:  requestChannelID,
					MessageIds: []string{messageID},
				},
			)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			require.Equal(t, platform.PlatformTwitch, usersRepo.platform)
			require.Equal(t, requestChannelID, usersRepo.platformUserID)
			require.Equal(t, platform.PlatformTwitch, channelsRepo.platform)
			require.Equal(t, channelUserID, channelsRepo.userID)

			if !tt.wantDelete {
				require.Empty(t, actions.opts)
				return
			}

			require.Equal(
				t,
				[]twitchactions.DeleteMessageOpts{{
					BroadcasterID: selectedChannelID,
					ModeratorID:   selectedBotID,
					MessageID:     messageID,
				}},
				actions.opts,
			)
		})
	}
}

type deleteMessageTwitchActions struct {
	mu   sync.Mutex
	opts []twitchactions.DeleteMessageOpts
}

func (a *deleteMessageTwitchActions) Ban(context.Context, twitchactions.BanOpts) error {
	panic("unexpected call")
}

func (a *deleteMessageTwitchActions) DeleteMessage(
	_ context.Context,
	opts twitchactions.DeleteMessageOpts,
) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	a.opts = append(a.opts, opts)
	return nil
}

type deleteMessageUsersRepository struct {
	user           usersmodel.User
	err            error
	platform       platform.Platform
	platformUserID string
}

func (r *deleteMessageUsersRepository) GetByPlatformID(
	_ context.Context,
	p platform.Platform,
	platformUserID string,
) (usersmodel.User, error) {
	r.platform = p
	r.platformUserID = platformUserID
	return r.user, r.err
}

func (*deleteMessageUsersRepository) GetByID(context.Context, uuid.UUID) (usersmodel.User, error) {
	return usersmodel.Nil, nil
}

func (*deleteMessageUsersRepository) GetManyByIDS(
	context.Context,
	usersrepository.GetManyInput,
) ([]usersmodel.User, error) {
	return nil, nil
}

func (*deleteMessageUsersRepository) Update(
	context.Context,
	uuid.UUID,
	usersrepository.UpdateInput,
) (usersmodel.User, error) {
	return usersmodel.Nil, nil
}

func (*deleteMessageUsersRepository) GetRandomOnlineUser(
	context.Context,
	usersrepository.GetRandomOnlineUserInput,
) (usersmodel.OnlineUser, error) {
	return usersmodel.NilOnlineUser, nil
}

func (*deleteMessageUsersRepository) GetOnlineUsersWithFilters(
	context.Context,
	usersrepository.GetOnlineUsersWithFiltersInput,
) ([]usersmodel.OnlineUser, error) {
	return nil, nil
}

func (*deleteMessageUsersRepository) GetByApiKey(context.Context, string) (usersmodel.User, error) {
	return usersmodel.Nil, nil
}

func (*deleteMessageUsersRepository) Create(
	context.Context,
	usersrepository.CreateInput,
) (usersmodel.User, error) {
	return usersmodel.Nil, nil
}

type deleteMessageChannelsRepository struct {
	channel  channelentity.Channel
	err      error
	platform platform.Platform
	userID   uuid.UUID
}

func (r *deleteMessageChannelsRepository) GetByBindingUserID(
	_ context.Context,
	p platform.Platform,
	userID uuid.UUID,
) (channelentity.Channel, error) {
	r.platform = p
	r.userID = userID
	return r.channel, r.err
}

func (*deleteMessageChannelsRepository) GetMany(
	context.Context,
	channelsrepository.GetManyInput,
) ([]channelentity.Channel, error) {
	return nil, nil
}

func (*deleteMessageChannelsRepository) GetAllByBindingPlatform(
	context.Context,
	platform.Platform,
) ([]channelentity.Channel, error) {
	return nil, nil
}

func (*deleteMessageChannelsRepository) GetByID(context.Context, uuid.UUID) (channelentity.Channel, error) {
	return channelentity.Nil, nil
}

func (*deleteMessageChannelsRepository) GetByApiKey(context.Context, string) (channelentity.Channel, error) {
	return channelentity.Nil, nil
}

func (*deleteMessageChannelsRepository) GetByPlatformChannelID(
	context.Context,
	platform.Platform,
	string,
) (channelentity.Channel, error) {
	return channelentity.Nil, nil
}

func (*deleteMessageChannelsRepository) GetBySlug(
	context.Context,
	channelsrepository.GetBySlugInput,
) (channelentity.Channel, error) {
	return channelentity.Nil, nil
}

func (*deleteMessageChannelsRepository) GetCount(context.Context, channelsrepository.GetCountInput) (int, error) {
	return 0, nil
}

func (*deleteMessageChannelsRepository) Update(
	context.Context,
	uuid.UUID,
	channelsrepository.UpdateInput,
) (channelentity.Channel, error) {
	return channelentity.Nil, nil
}

func (*deleteMessageChannelsRepository) Create(
	context.Context,
	channelsrepository.CreateInput,
) (channelentity.Channel, error) {
	return channelentity.Nil, nil
}
