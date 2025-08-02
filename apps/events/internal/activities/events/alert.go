package events

import (
	"context"
	"errors"

	"github.com/twirapp/twir/apps/events/internal/shared"
	"github.com/twirapp/twir/libs/grpc/websockets"
	"github.com/twirapp/twir/libs/repositories/events/model"
	"go.temporal.io/sdk/activity"
)

func (c *Activity) TriggerAlert(
	ctx context.Context,
	operation model.EventOperation,
	data shared.EventData,
) error {
	activity.RecordHeartbeat(ctx, nil)

	if data.ChannelID == "" {
		return errors.New("channel id is required")
	}

	if operation.Target == nil || *operation.Target == "" {
		return errors.New("target is required")
	}

	_, err := c.websocketsGrpc.TriggerAlert(
		ctx,
		&websockets.TriggerAlertRequest{
			ChannelId: data.ChannelID,
			AlertId:   *operation.Target,
		},
	)

	return err
}
