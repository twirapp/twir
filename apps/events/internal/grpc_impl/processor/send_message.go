package processor

import (
	"context"
	"fmt"
	"github.com/samber/lo"
	"github.com/satont/tsuwari/libs/grpc/generated/bots"
)

func (c *Processor) SendMessage(channelId, message string, useAnnounce bool) error {
	msg, err := c.HydrateStringWithData(message, c.data)
	if err != nil {
		return fmt.Errorf("cannot hydrate string %w", err)
	}

	_, err = c.services.BotsGrpc.SendMessage(context.Background(), &bots.SendMessageRequest{
		ChannelId:   channelId,
		ChannelName: nil,
		Message:     msg,
		IsAnnounce:  lo.ToPtr(useAnnounce),
	})

	if err != nil {
		return fmt.Errorf("cannot send message %w", err)
	}

	return nil
}
