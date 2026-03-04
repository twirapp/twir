package timers

import (
	"context"

	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/gqlmodel"
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/mappers"
	"github.com/twirapp/twir/libs/audit"
	timersbusservice "github.com/twirapp/twir/libs/bus-core/timers"
	"github.com/twirapp/twir/libs/errors"
	"github.com/twirapp/twir/libs/logger"
)

func (c *Service) Delete(ctx context.Context, id uuid.UUID, channelID, actorID string) error {
	timer, err := c.timersRepository.GetByID(ctx, id)
	if err != nil {
		return errors.NewInternalError("Failed to get timer", err)
	}

	if timer.ChannelID != channelID {
		return errors.NewNotFoundError("Timer with this ID was not found for your channel")
	}

	if err := c.timersRepository.Delete(ctx, id); err != nil {
		c.logger.Error("cannot delete timer", logger.Error(err))
		return errors.NewInternalError("Failed to delete timer", err)
	}

	if _, err := c.twirbus.Timers.RemoveTimer.Request(
		ctx,
		timersbusservice.AddOrRemoveTimerRequest{TimerID: timer.ID.String()},
	); err != nil {
		return errors.NewInternalError("Failed to remove timer from bus", err)
	}

	_ = c.auditRecorder.RecordDeleteOperation(
		ctx,
		audit.DeleteOperation{
			Metadata: audit.OperationMetadata{
				System:    mappers.AuditSystemToTableName(gqlmodel.AuditLogSystemChannelTimers),
				ActorID:   &actorID,
				ChannelID: &channelID,
				ObjectID:  lo.ToPtr(timer.ID.String()),
			},
			OldValue: timer,
		},
	)

	return nil
}
