package events

import (
	"context"
	"fmt"

	"github.com/samber/lo"
	"github.com/satont/twir/apps/events/internal/shared"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/twirapp/twir/libs/grpc/bots"
)

func (c *Activity) SendMessage(
	ctx context.Context,
	operation model.EventOperation,
	data shared.EvenData,
) error {
	msg, err := c.hydrator.HydrateStringWithData(data.ChannelID, operation.Input.String, data)
	if err != nil {
		return fmt.Errorf("cannot hydrate string %w", err)
	}

	_, err = c.botsGrpc.SendMessage(
		ctx,
		&bots.SendMessageRequest{
			ChannelId:   data.ChannelID,
			ChannelName: nil,
			Message:     msg,
			IsAnnounce:  lo.ToPtr(operation.UseAnnounce),
		},
	)

	if err != nil {
		return fmt.Errorf("cannot send message %w", err)
	}

	return nil
}
