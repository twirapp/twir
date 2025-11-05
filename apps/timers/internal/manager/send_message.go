package manager

import (
	"context"

	"github.com/twirapp/twir/libs/bus-core/bots"
	busparser "github.com/twirapp/twir/libs/bus-core/parser"
	timersmodel "github.com/twirapp/twir/libs/repositories/timers/model"
)

func (c *Manager) sendMessage(
	ctx context.Context,
	channelId,
	text string,
	isAnnounce bool,
	announceColor timersmodel.AnnounceColor,
	count int,
) error {
	parseReq, err := c.twirBus.Parser.ParseVariablesInText.Request(
		ctx,
		busparser.ParseVariablesInTextRequest{
			ChannelID: channelId,
			Text:      text,
		},
	)
	if err != nil {
		return err
	}

	for i := 0; i < count; i++ {
		err = c.twirBus.Bots.SendMessage.Publish(
			ctx,
			bots.SendMessageRequest{
				ChannelId:      channelId,
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
