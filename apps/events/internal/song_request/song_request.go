package song_request

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
	"github.com/twirapp/twir/apps/events/internal/channelbinding"
	buscore "github.com/twirapp/twir/libs/bus-core"
	"github.com/twirapp/twir/libs/bus-core/generic"
	"github.com/twirapp/twir/libs/bus-core/ytsr"
	"github.com/twirapp/twir/libs/entities/platform"
	model "github.com/twirapp/twir/libs/gomodels"
	"github.com/twirapp/twir/libs/logger"
	channelservice "github.com/twirapp/twir/libs/services/channels"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

type Opts struct {
	fx.In

	Gorm           *gorm.DB
	TwirBus        *buscore.Bus
	Logger         *slog.Logger
	ChannelService *channelservice.ChannelService
}

func New(opts Opts) *SongRequest {
	return &SongRequest{
		gorm:           opts.Gorm,
		twirBus:        opts.TwirBus,
		logger:         opts.Logger,
		channelService: opts.ChannelService,
	}
}

type SongRequest struct {
	gorm           *gorm.DB
	twirBus        *buscore.Bus
	logger         *slog.Logger
	channelService *channelservice.ChannelService
}

type ProcessFromDonationInput struct {
	Text      string
	ChannelID string
}

func (c *SongRequest) ProcessFromDonation(
	ctx context.Context,
	input ProcessFromDonationInput,
) error {
	srSettings := model.ChannelSongRequestsSettings{}
	if err := c.gorm.
		WithContext(ctx).
		Where(
			"channel_id = ?",
			input.ChannelID,
		).
		First(&srSettings).
		Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		return fmt.Errorf("cannot get song request settings: %w", err)
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

	if !srCommand.Enabled || !srSettings.Enabled || !srSettings.TakeSongFromDonationMessage {
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
	binding, found := channelbinding.Find(channel, platform.PlatformTwitch)
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
