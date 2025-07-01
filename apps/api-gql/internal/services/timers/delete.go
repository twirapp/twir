package timers

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/satont/twir/libs/logger/audit"
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/gqlmodel"
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/mappers"
	timersbusservice "github.com/twirapp/twir/libs/bus-core/timers"
)

func (c *Service) Delete(ctx context.Context, id uuid.UUID, channelID, actorID string) error {
	timer, err := c.timersRepository.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if timer.ChannelID != channelID {
		return ErrTimerNotFound
	}

	if err := c.timersRepository.Delete(ctx, id); err != nil {
		c.logger.Error("cannot delete timer", slog.Any("err", err))
		return err
	}

	if _, err := c.twirbus.Timers.RemoveTimer.Request(
		ctx,
		timersbusservice.AddOrRemoveTimerRequest{TimerID: timer.ID.String()},
	); err != nil {
		return err
	}

	c.logger.Audit(
		"Timers remove",
		audit.Fields{
			OldValue:      timer,
			ActorID:       &actorID,
			ChannelID:     &channelID,
			System:        mappers.AuditSystemToTableName(gqlmodel.AuditLogSystemChannelTimers),
			OperationType: audit.OperationDelete,
			ObjectID:      lo.ToPtr(timer.ID.String()),
		},
	)

	return nil
}
