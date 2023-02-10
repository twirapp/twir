package processor

import (
	"context"
	"github.com/satont/tsuwari/libs/grpc/generated/bots"
)

func (c *Processor) SendMessage(channelId, message string) {
	msg, err := hydrateStringWithData(message, c.data)
	if err == nil {
		c.services.BotsGrpc.SendMessage(context.Background(), &bots.SendMessageRequest{
			ChannelId:   channelId,
			ChannelName: nil,
			Message:     msg,
			IsAnnounce:  nil,
		})
	}
}
