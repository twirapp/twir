package timers

import (
	"context"
	"fmt"

	"github.com/samber/lo"
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/gqlmodel"
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/mappers"
	"github.com/twirapp/twir/libs/audit"
	"github.com/twirapp/twir/libs/bus-core/bots"
	timersbusservice "github.com/twirapp/twir/libs/bus-core/timers"
	timersentity "github.com/twirapp/twir/libs/entities/timers"
	"github.com/twirapp/twir/libs/logger"
	timersrepository "github.com/twirapp/twir/libs/repositories/timers"
)

type CreateInput struct {
	ChannelID string
	ActorID   string

	Name            string
	Enabled         bool
	OfflineEnabled  bool
	OnlineEnabled   bool
	TimeInterval    int
	MessageInterval int
	Responses       []CreateResponse
}

type CreateResponse struct {
	Text          string
	IsAnnounce    bool
	Count         int
	AnnounceColor bots.AnnounceColor
}

func (c *Service) Create(ctx context.Context, data CreateInput) (timersentity.Timer, error) {
	plan, err := c.plansRepository.GetByChannelID(ctx, data.ChannelID)
	if err != nil {
		return timersentity.Nil, fmt.Errorf("failed to get plan: %w", err)
	}
	if plan.IsNil() {
		return timersentity.Nil, fmt.Errorf("plan not found for channel")
	}

	createdCount, err := c.timersRepository.CountByChannelID(ctx, data.ChannelID)
	if err != nil {
		return timersentity.Nil, err
	}

	if createdCount >= plan.MaxTimers {
		return timersentity.Nil, fmt.Errorf("you can have only %v timers", plan.MaxTimers)
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

	timer, err := c.timersRepository.Create(
		ctx,
		timersrepository.CreateInput{
			ChannelID:       data.ChannelID,
			Name:            data.Name,
			Enabled:         data.Enabled,
			OfflineEnabled:  data.OfflineEnabled,
			OnlineEnabled:   data.OnlineEnabled,
			TimeInterval:    data.TimeInterval,
			MessageInterval: data.MessageInterval,
			Responses:       responses,
		},
	)
	if err != nil {
		return timersentity.Nil, err
	}

	timersReq := timersbusservice.AddOrRemoveTimerRequest{TimerID: timer.ID.String()}
	if timer.Enabled {
		if err := c.twirbus.Timers.AddTimer.Publish(ctx, timersReq); err != nil {
			c.logger.Error("cannot publish add timer", logger.Error(err))
		}
	} else {
		if err := c.twirbus.Timers.RemoveTimer.Publish(ctx, timersReq); err != nil {
			c.logger.Error("cannot publish remove timer", logger.Error(err))
		}
	}

	_ = c.auditRecorder.RecordCreateOperation(
		ctx,
		audit.CreateOperation{
			Metadata: audit.OperationMetadata{
				System:    mappers.AuditSystemToTableName(gqlmodel.AuditLogSystemChannelTimers),
				ActorID:   &data.ActorID,
				ChannelID: &data.ChannelID,
				ObjectID:  lo.ToPtr(timer.ID.String()),
			},
			NewValue: timer,
		},
	)

	return timer, nil
}

func (c *Service) CreateMany(ctx context.Context, data []CreateInput) (bool, error) {
	txErr := c.trmManager.Do(
		ctx,
		func(txCtx context.Context) error {
			for _, d := range data {
				_, err := c.Create(txCtx, d)
				if err != nil {
					return err
				}
			}

			return nil
		},
	)
	if txErr != nil {
		return false, fmt.Errorf("failed to create timers: %w", txErr)
	}

	return false, nil
}
