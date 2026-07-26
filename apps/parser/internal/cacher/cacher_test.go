package cacher

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/twirapp/twir/apps/parser/internal/types"
	"github.com/twirapp/twir/apps/parser/internal/types/services"
	cfg "github.com/twirapp/twir/libs/config"
	channelentity "github.com/twirapp/twir/libs/entities/channel"
	channelplatformentity "github.com/twirapp/twir/libs/entities/channel_platform"
	"github.com/twirapp/twir/libs/entities/platform"
	channelsrepository "github.com/twirapp/twir/libs/repositories/channels"
	channelservice "github.com/twirapp/twir/libs/services/channels"
)

func TestGetDbChannelResolvesByCanonicalDBChannelID(t *testing.T) {
	canonicalChannelID := uuid.New()
	channel := channelentity.Channel{
		ID: canonicalChannelID,
		Bindings: []channelplatformentity.ChannelPlatform{
			{
				Platform:          platform.PlatformTwitch,
				PlatformChannelID: "twitch-provider-id",
				UserID:            uuid.New(),
				BotConfig:         []byte(`{"bot_id":"twitch-bot"}`),
			},
			{
				Platform:          platform.PlatformKick,
				PlatformChannelID: "kick-provider-id",
			},
			{
				Platform:          platform.PlatformVKVideoLive,
				PlatformChannelID: "vk-provider-id",
			},
		},
	}

	tests := []struct {
		name       string
		platform   platform.Platform
		providerID string
	}{
		{name: "twitch", platform: platform.PlatformTwitch, providerID: "twitch-provider-id"},
		{name: "kick", platform: platform.PlatformKick, providerID: "kick-provider-id"},
		{name: "vk", platform: platform.PlatformVKVideoLive, providerID: "vk-provider-id"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repository := &fakeCacherChannelsRepository{
				channel:           channel,
				platformLookupErr: errors.New("legacy platform lookup must not be used"),
			}
			channelService := channelservice.NewChannelService(
				repository,
				nil,
				cfg.Config{},
				nil,
				nil,
			)
			cacher := &cacher{
				services: &services.Services{ChannelService: channelService},
				parseCtxChannel: &types.ParseContextChannel{
					ID:          tt.providerID,
					DBChannelID: canonicalChannelID.String(),
				},
				cache: &cache{},
				locks: &locks{},
			}

			got, err := cacher.getDbChannel(context.Background())

			require.NoError(t, err)
			require.Equal(t, canonicalChannelID.String(), got.ChannelID)
			require.Equal(t, []uuid.UUID{canonicalChannelID}, repository.channelIDs)
			require.Empty(t, repository.platformLookups, "platform %s must use the canonical channel ID", tt.platform)
		})
	}
}

func TestGetDbChannelRejectsMissingOrMalformedCanonicalID(t *testing.T) {
	tests := []struct {
		name        string
		dbChannelID string
		providerID  string
	}{
		{name: "missing", providerID: "provider-id"},
		{name: "malformed", dbChannelID: "not-a-uuid", providerID: "provider-id"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repository := &fakeCacherChannelsRepository{
				platformLookupErr: errors.New("legacy platform lookup must not be used"),
			}
			channelService := channelservice.NewChannelService(
				repository,
				nil,
				cfg.Config{},
				nil,
				nil,
			)
			cacher := &cacher{
				services: &services.Services{ChannelService: channelService},
				parseCtxChannel: &types.ParseContextChannel{
					ID:          tt.providerID,
					DBChannelID: tt.dbChannelID,
				},
				cache: &cache{},
				locks: &locks{},
			}

			got, err := cacher.getDbChannel(context.Background())

			require.Nil(t, got)
			require.ErrorContains(t, err, "parse channel id")
			require.Empty(t, repository.channelIDs)
			require.Empty(t, repository.platformLookups)
		})
	}
}

type fakeCacherChannelsRepository struct {
	channelsrepository.Repository

	channel           channelentity.Channel
	getByIDErr        error
	platformLookupErr error
	channelIDs        []uuid.UUID
	platformLookups   []cacherPlatformChannelLookup
}

type cacherPlatformChannelLookup struct {
	platform          platform.Platform
	platformChannelID string
}

var _ channelsrepository.Repository = (*fakeCacherChannelsRepository)(nil)

func (f *fakeCacherChannelsRepository) GetByID(
	_ context.Context,
	channelID uuid.UUID,
) (channelentity.Channel, error) {
	f.channelIDs = append(f.channelIDs, channelID)
	return f.channel, f.getByIDErr
}

func (f *fakeCacherChannelsRepository) GetByPlatformChannelID(
	_ context.Context,
	p platform.Platform,
	platformChannelID string,
) (channelentity.Channel, error) {
	f.platformLookups = append(
		f.platformLookups,
		cacherPlatformChannelLookup{platform: p, platformChannelID: platformChannelID},
	)
	return f.channel, f.platformLookupErr
}
