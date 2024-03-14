package events

import (
	"context"
	"errors"

	"github.com/satont/twir/apps/events/internal/shared"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/twirapp/twir/libs/grpc/websockets"
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

	if !operation.Target.Valid {
		return errors.New("target is required")
	}

	_, err := c.websocketsGrpc.TriggerAlert(
		ctx,
		&websockets.TriggerAlertRequest{
			ChannelId: data.ChannelID,
			AlertId:   operation.Target.String,
		},
	)

	return err
}
