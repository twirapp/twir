package timers

import (
	"context"

	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/gqlmodel"
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/mappers"
	"github.com/twirapp/twir/libs/audit"
	timersbusservice "github.com/twirapp/twir/libs/bus-core/timers"
	timersentity "github.com/twirapp/twir/libs/entities/timers"
	timersrepository "github.com/twirapp/twir/libs/repositories/timers"
)

type UpdateInput struct {
	ChannelID string
	ActorID   string

	ID              uuid.UUID
	Name            *string
	Enabled         *bool
	OfflineEnabled  *bool
	TimeInterval    *int
	MessageInterval *int
	Responses       []CreateResponse
}

func (c *Service) Update(ctx context.Context, data UpdateInput) (timersentity.Timer, error) {
	timer, err := c.timersRepository.GetByID(ctx, data.ID)
	if err != nil {
		return timersentity.Nil, err
	}

	if timer.ChannelID != data.ChannelID {
		return timersentity.Nil, ErrTimerNotFound
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
				AnnounceColor: timersentity.AnnounceColor(response.AnnounceColor),
			},
		)
	}

	newTimer, err := c.timersRepository.UpdateByID(
		ctx,
		data.ID,
		timersrepository.UpdateInput{
			Name:            data.Name,
			Enabled:         data.Enabled,
			OfflineEnabled:  data.OfflineEnabled,
			TimeInterval:    data.TimeInterval,
			MessageInterval: data.MessageInterval,
			Responses:       responses,
		},
	)
	if err != nil {
		return timersentity.Nil, err
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

	return newTimer, nil
}
