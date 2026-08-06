package song_request

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
	buscore "github.com/twirapp/twir/libs/bus-core"
	"github.com/twirapp/twir/libs/bus-core/generic"
	buscorespotify "github.com/twirapp/twir/libs/bus-core/spotify"
	"github.com/twirapp/twir/libs/bus-core/ytsr"
	"github.com/twirapp/twir/libs/entities/platform"
	songrequestmode "github.com/twirapp/twir/libs/entities/song_request_mode"
	model "github.com/twirapp/twir/libs/gomodels"
	"github.com/twirapp/twir/libs/logger"
	songrequestssettingsrepository "github.com/twirapp/twir/libs/repositories/song_requests_settings"
	channelservice "github.com/twirapp/twir/libs/services/channels"
	"gorm.io/gorm"
)

func New(
	db *gorm.DB,
	twirBus *buscore.Bus,
	logger *slog.Logger,
	songRequestsSettingsRepo songrequestssettingsrepository.Repository,
	channelService *channelservice.ChannelService,
) *SongRequest {
	return &SongRequest{
		gorm:                     db,
		twirBus:                  twirBus,
		logger:                   logger,
		songRequestsSettingsRepo: songRequestsSettingsRepo,
		channelService:           channelService,
	}
}

type SongRequest struct {
	gorm                     *gorm.DB
	twirBus                  *buscore.Bus
	logger                   *slog.Logger
	songRequestsSettingsRepo songrequestssettingsrepository.Repository
	channelService           *channelservice.ChannelService
}

type ProcessFromDonationInput struct {
	Text      string
	Username  string
	ChannelID string
}

func (c *SongRequest) ProcessFromDonation(
	ctx context.Context,
	input ProcessFromDonationInput,
) error {
	srSettings, err := c.songRequestsSettingsRepo.GetByChannelID(ctx, input.ChannelID)
	if err != nil {
		if errors.Is(err, songrequestssettingsrepository.ErrNotFound) {
			return nil
		}
		return fmt.Errorf("cannot get song request settings: %w", err)
	}

	if !srSettings.Enabled || !srSettings.TakeSongFromDonationMessage {
		return nil
	}

	srCommand := model.ChannelsCommands{}
	if err := c.gorm.
		WithContext(ctx).
		Where(
			`"channelId" = ? AND "defaultName" = ?`,
			input.ChannelID,
			"sr",
		).
		First(&srCommand).
		Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}

		return fmt.Errorf("cannot get song request command: %w", err)
	}

	if !srCommand.Enabled {
		return nil
	}

	if srSettings.Mode == songrequestmode.ModeSpotify {
		response, err := c.twirBus.Spotify.CreateSongRequest.Request(
			ctx,
			buscorespotify.CreateSongRequestRequest{
				ChannelID:            input.ChannelID,
				RequesterUserID:      "",
				RequesterName:        input.Username,
				RequesterDisplayName: input.Username,
				Source:               "donation",
				Query:                input.Text,
			},
		)
		if err != nil {
			c.logger.Error("cannot create spotify song request", logger.Error(err))
			return nil
		}

		c.logger.Info(
			"added spotify song request",
			slog.String("title", response.Data.Request.Title),
			slog.String("artist", response.Data.Request.Artist),
		)
		return nil
	}

	channelID, err := uuid.Parse(input.ChannelID)
	if err != nil {
		return fmt.Errorf("parse channel id: %w", err)
	}
	channel, err := c.channelService.GetChannelByID(ctx, channelID)
	if err != nil {
		return fmt.Errorf("get channel: %w", err)
	}
	binding, found := channel.Binding(platform.PlatformTwitch)
	if !found {
		return fmt.Errorf("find Twitch channel binding")
	}

	ytsrResult, err := c.twirBus.YTSRSearch.Request(
		ctx,
		ytsr.SearchRequest{
			Search:    input.Text,
			OnlyLinks: true,
		},
	)
	if err != nil {
		return fmt.Errorf("cannot search for ytsrResult: %w", err)
	}

	for _, song := range ytsrResult.Data.Songs {
		messageID := uuid.NewString()
		err := c.twirBus.Parser.ProcessMessageAsCommand.Publish(
			ctx,
			generic.ChatMessage{
				ID:                   messageID,
				BroadcasterUserId:    binding.PlatformChannelID,
				BroadcasterUserName:  "",
				BroadcasterUserLogin: "",
				ChatterUserId:        binding.PlatformChannelID,
				ChatterUserName:      "",
				ChatterUserLogin:     "",
				MessageID:            messageID,
				Platform:             string(platform.PlatformTwitch),
				PlatformChannelID:    binding.PlatformChannelID,
				ChannelID:            input.ChannelID,
				ChannelBindingID:     binding.ID.String(),
				UserID:               binding.UserID.String(),
				SenderID:             binding.PlatformChannelID,
				Message: &generic.ChatMessageMessage{
					Text: fmt.Sprintf(
						"!%s https://youtu.be/%s",
						srCommand.Name,
						song.Id,
					),
					Fragments: nil,
				},
				Color: "",
				Badges: []generic.ChatMessageBadge{
					{
						ID:    "broadcaster",
						SetID: "broadcaster",
						Info:  "broadcaster",
						Text:  "broadcaster",
					},
				},
				IsBroadcaster:               true,
				MessageType:                 "",
				Cheer:                       nil,
				Reply:                       nil,
				ChannelPointsCustomRewardId: "",
			},
		)

		if err != nil {
			c.logger.Error("cannot publish process message", logger.Error(err))
		}
	}

	return nil
}
