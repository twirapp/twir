package manager

import (
	"context"

	"github.com/google/uuid"
	"github.com/twirapp/twir/libs/bus-core/bots"
	busparser "github.com/twirapp/twir/libs/bus-core/parser"
	"github.com/twirapp/twir/libs/entities/platform"
	timersentity "github.com/twirapp/twir/libs/entities/timers"
)

func (c *Manager) sendMessage(
	ctx context.Context,
	channelId uuid.UUID,
	platformSource platform.Platform,
	text string,
	isAnnounce bool,
	announceColor timersentity.AnnounceColor,
	count int,
) error {
	parseReq, err := c.twirBus.Parser.ParseVariablesInText.Request(
		ctx,
		busparser.ParseVariablesInTextRequest{
			ChannelID:      channelId,
			Text:           text,
			PlatformSource: &platformSource,
		},
	)
	if err != nil {
		return err
	}

	for range count {
		err = c.twirBus.Bots.SendMessage.Publish(
			ctx,
			bots.SendMessageRequest{
				ChannelID:      channelId,
				Platforms:      []platform.Platform{platformSource},
				Message:        parseReq.Data.Text,
				IsAnnounce:     isAnnounce,
				SkipRateLimits: true,
				AnnounceColor:  bots.AnnounceColor(announceColor),
			},
		)
		if err != nil {
			return err
		}
	}

	return nil
}
