package song_request

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	buscore "github.com/twirapp/twir/libs/bus-core"
	"github.com/twirapp/twir/libs/bus-core/generic"
	"github.com/twirapp/twir/libs/bus-core/ytsr"
	model "github.com/twirapp/twir/libs/gomodels"
	"github.com/twirapp/twir/libs/logger"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

type Opts struct {
	fx.In

	Gorm    *gorm.DB
	TwirBus *buscore.Bus
	Logger  *slog.Logger
}

func New(opts Opts) *SongRequest {
	return &SongRequest{
		gorm:    opts.Gorm,
		twirBus: opts.TwirBus,
		logger:  opts.Logger,
	}
}

type SongRequest struct {
	gorm    *gorm.DB
	twirBus *buscore.Bus
	logger  *slog.Logger
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
		err := c.twirBus.Parser.ProcessMessageAsCommand.Publish(
			ctx,
			generic.ChatMessage{
				ID:                   "",
				BroadcasterUserId:    input.ChannelID,
				BroadcasterUserName:  "",
				BroadcasterUserLogin: "",
				ChatterUserId:        input.ChannelID,
				ChatterUserName:      "",
				ChatterUserLogin:     "",
				MessageID:            "",
				PlatformChannelID:    input.ChannelID,
				ChannelID:            input.ChannelID,
				UserID:               input.ChannelID,
				SenderID:             input.ChannelID,
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
