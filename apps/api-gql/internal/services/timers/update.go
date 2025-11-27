package timers

import (
	"context"

	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/gqlmodel"
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/mappers"
	"github.com/twirapp/twir/apps/api-gql/internal/entity"
	"github.com/twirapp/twir/libs/audit"
	timersbusservice "github.com/twirapp/twir/libs/bus-core/timers"
	timersrepository "github.com/twirapp/twir/libs/repositories/timers"
	"github.com/twirapp/twir/libs/repositories/timers/model"
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
				Text:          response.Text,
				IsAnnounce:    response.IsAnnounce,
				Count:         count,
				AnnounceColor: model.AnnounceColor(response.AnnounceColor),
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

	_ = c.auditRecorder.RecordUpdateOperation(
		ctx,
		audit.UpdateOperation{
			Metadata: audit.OperationMetadata{
				System:    mappers.AuditSystemToTableName(gqlmodel.AuditLogSystemChannelTimers),
				ActorID:   &data.ActorID,
				ChannelID: &data.ChannelID,
				ObjectID:  lo.ToPtr(newTimer.ID.String()),
			},
			NewValue: newTimer,
			OldValue: timer,
		},
	)

	timersReq := timersbusservice.AddOrRemoveTimerRequest{TimerID: newTimer.ID.String()}
	if newTimer.Enabled {
		c.twirbus.Timers.AddTimer.Publish(ctx, timersReq)
	} else {
		c.twirbus.Timers.RemoveTimer.Publish(ctx, timersReq)
	}

	return c.dbToModel(newTimer), nil
}
