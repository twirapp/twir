package timers

import (
	"context"

	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/satont/twir/libs/logger/audit"
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/gqlmodel"
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/mappers"
	"github.com/twirapp/twir/apps/api-gql/internal/entity"
	timersbusservice "github.com/twirapp/twir/libs/bus-core/timers"
	timersrepository "github.com/twirapp/twir/libs/repositories/timers"
)

type UpdateInput struct {
	ChannelID string
	ActorID   string

	ID              uuid.UUID
	Name            *string
	Enabled         *bool
	TimeInterval    *int
	MessageInterval *int
	Responses       []CreateResponse
}

func (c *Service) Update(ctx context.Context, data UpdateInput) (entity.Timer, error) {
	timer, err := c.timersRepository.GetByID(ctx, data.ID)
	if err != nil {
		return entity.TimerNil, err
	}

	if timer.ChannelID != data.ChannelID {
		return entity.TimerNil, ErrTimerNotFound
	}

	responses := make([]timersrepository.CreateResponse, 0, len(data.Responses))
	for _, response := range data.Responses {
		count := response.Count
		if count == 0 {
			count = 1
		}

		responses = append(
			responses,
			timersrepository.CreateResponse{
				Text:       response.Text,
				IsAnnounce: response.IsAnnounce,
				Count:      count,
			},
		)
	}

	newTimer, err := c.timersRepository.UpdateByID(
		ctx,
		data.ID,
		timersrepository.UpdateInput{
			Name:            data.Name,
			Enabled:         data.Enabled,
			TimeInterval:    data.TimeInterval,
			MessageInterval: data.MessageInterval,
			Responses:       responses,
		},
	)
	if err != nil {
		return entity.TimerNil, err
	}

	c.logger.Audit(
		"Timers update",
		audit.Fields{
			OldValue:      timer,
			NewValue:      newTimer,
			ActorID:       &data.ActorID,
			ChannelID:     &data.ChannelID,
			System:        mappers.AuditSystemToTableName(gqlmodel.AuditLogSystemChannelTimers),
			OperationType: audit.OperationUpdate,
			ObjectID:      lo.ToPtr(newTimer.ID.String()),
		},
	)

	go func() {
		timersReq := timersbusservice.AddOrRemoveTimerRequest{TimerID: newTimer.ID.String()}
		if newTimer.Enabled {
			c.twirbus.Timers.AddTimer.Publish(timersReq)
		} else {
			c.twirbus.Timers.RemoveTimer.Publish(timersReq)
		}
	}()

	return c.dbToModel(newTimer), nil
}
