package timers

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	dbmodels "github.com/satont/twir/libs/gomodels"
	"github.com/satont/twir/libs/logger/audit"
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/gqlmodel"
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/mappers"
	"github.com/twirapp/twir/apps/api-gql/internal/services/timers/model"
	timersbusservice "github.com/twirapp/twir/libs/bus-core/timers"
)

type CreateInput struct {
	ChannelID string
	ActorID   string

	Name            string
	Enabled         bool
	TimeInterval    int
	MessageInterval int
	Responses       []CreateResponse
}

type CreateResponse struct {
	Text       string
	IsAnnounce bool
}

func (c *Service) Create(ctx context.Context, data CreateInput) (model.Timer, error) {
	var createdCount int64
	if err := c.gorm.
		WithContext(ctx).
		Model(&dbmodels.ChannelsTimers{}).
		Where(`"channelId" = ?`, data.ChannelID).
		Count(&createdCount).
		Error; err != nil {
		return model.Nil, err
	}

	if createdCount >= MaxTimersPerChannel {
		return model.Nil, fmt.Errorf("you can have only %v timers", MaxTimersPerChannel)
	}

	timerId := uuid.NewString()
	responses := make([]*dbmodels.ChannelsTimersResponses, 0, len(data.Responses))
	for _, r := range data.Responses {
		responses = append(
			responses,
			&dbmodels.ChannelsTimersResponses{
				ID:         uuid.NewString(),
				Text:       r.Text,
				IsAnnounce: r.IsAnnounce,
				TimerID:    timerId,
			},
		)
	}

	entity := dbmodels.ChannelsTimers{
		ID:                       timerId,
		ChannelID:                data.ChannelID,
		Name:                     data.Name,
		Enabled:                  data.Enabled,
		TimeInterval:             int32(data.TimeInterval),
		MessageInterval:          int32(data.MessageInterval),
		Responses:                responses,
		LastTriggerMessageNumber: 0,
	}
	if err := c.gorm.WithContext(ctx).Create(&entity).Error; err != nil {
		return model.Nil, err
	}

	go func() {
		timersReq := timersbusservice.AddOrRemoveTimerRequest{TimerID: entity.ID}
		if entity.Enabled {
			c.twirbus.Timers.AddTimer.Publish(timersReq)
		} else {
			c.twirbus.Timers.RemoveTimer.Publish(timersReq)
		}
	}()

	c.logger.Audit(
		"Timers create",
		audit.Fields{
			NewValue:      entity,
			ActorID:       &data.ActorID,
			ChannelID:     &data.ChannelID,
			System:        mappers.AuditSystemToTableName(gqlmodel.AuditLogSystemChannelTimers),
			OperationType: audit.OperationCreate,
			ObjectID:      &entity.ID,
		},
	)

	return c.dbToModel(entity), nil
}
