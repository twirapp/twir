package processor

import (
	"context"
	"github.com/samber/lo"
	"github.com/satont/tsuwari/libs/grpc/generated/bots"
)

func (c *Processor) SendMessage(channelId, message string, useAnnounce bool) {
	msg, err := hydrateStringWithData(message, c.data)
	if err != nil {
		return
	}

	c.services.BotsGrpc.SendMessage(context.Background(), &bots.SendMessageRequest{
		ChannelId:   channelId,
		ChannelName: nil,
		Message:     msg,
		IsAnnounce:  lo.ToPtr(useAnnounce),
	})
}
