package alerts

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/satont/twir/libs/logger"
	"github.com/satont/twir/libs/logger/audit"
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/gqlmodel"
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/mappers"
	"github.com/twirapp/twir/apps/api-gql/internal/entity"
	generic_cacher "github.com/twirapp/twir/libs/cache/generic-cacher"
	"github.com/twirapp/twir/libs/repositories/alerts"
	"github.com/twirapp/twir/libs/repositories/alerts/model"
	"go.uber.org/fx"
)

type Opts struct {
	fx.In

	AlertsRepository alerts.Repository
	Logger           logger.Logger
	AlertsCache      *generic_cacher.GenericCacher[[]model.Alert]
}

func New(opts Opts) *Service {
	return &Service{
		alertsRepository: opts.AlertsRepository,
		logger:           opts.Logger,
		alertsCache:      opts.AlertsCache,
	}
}

type Service struct {
	alertsRepository alerts.Repository
	logger           logger.Logger
	alertsCache      *generic_cacher.GenericCacher[[]model.Alert]
}

func (c *Service) modelToEntity(m model.Alert) entity.Alert {
	return entity.Alert{
		ID:           m.ID,
		Name:         m.Name,
		ChannelID:    m.ChannelID,
		AudioID:      m.AudioID,
		AudioVolume:  m.AudioVolume,
		CommandIDS:   m.CommandIDS,
		RewardIDS:    m.RewardIDS,
		GreetingsIDS: m.GreetingsIDS,
		KeywordsIDS:  m.KeywordsIDS,
	}
}

func (c *Service) GetManyByChannelID(ctx context.Context, channelID string) (
	[]entity.Alert,
	error,
) {
	dbAlerts, err := c.alertsRepository.GetManyByChannelID(ctx, channelID)
	if err != nil {
		return nil, err
	}

	entities := make([]entity.Alert, 0, len(dbAlerts))
	for _, a := range dbAlerts {
		entities = append(entities, c.modelToEntity(a))
	}

	return entities, nil
}

type CreateInput struct {
	ChannelID string
	ActorID   string

	Name         string
	AudioID      *string
	AudioVolume  int
	CommandIDS   []string
	RewardIDS    []string
	GreetingsIDS []string
	KeywordsIDS  []string
}

func (c *Service) Create(ctx context.Context, input CreateInput) (entity.Alert, error) {
	alert, err := c.alertsRepository.Create(
		ctx,
		alerts.CreateInput{
			Name:         input.Name,
			ChannelID:    input.ChannelID,
			AudioID:      input.AudioID,
			AudioVolume:  input.AudioVolume,
			CommandIDS:   input.CommandIDS,
			RewardIDS:    input.RewardIDS,
			GreetingsIDS: input.GreetingsIDS,
			KeywordsIDS:  input.KeywordsIDS,
		},
	)
	if err != nil {
		return entity.AlertNil, err
	}

	c.logger.Audit(
		"Channel alert create",
		audit.Fields{
			NewValue:      alert,
			ActorID:       &input.ActorID,
			ChannelID:     &input.ChannelID,
			System:        mappers.AuditSystemToTableName(gqlmodel.AuditLogSystemChannelsAlerts),
			OperationType: audit.OperationCreate,
			ObjectID:      lo.ToPtr(alert.ID.String()),
		},
	)

	if err = c.alertsCache.Invalidate(ctx, input.ChannelID); err != nil {
		return entity.AlertNil, fmt.Errorf("failed to invalidate cache: %w", err)
	}

	return c.modelToEntity(alert), nil
}

type UpdateInput struct {
	ChannelID string
	ActorID   string

	Name         *string
	AudioID      *string
	AudioVolume  *int
	CommandIDS   []string
	RewardIDS    []string
	GreetingsIDS []string
	KeywordsIDS  []string
}

func (c *Service) Update(ctx context.Context, id uuid.UUID, input UpdateInput) (
	entity.Alert,
	error,
) {
	dbAlert, err := c.alertsRepository.GetByID(ctx, id)
	if err != nil {
		return entity.AlertNil, err
	}

	if dbAlert.ChannelID != input.ChannelID {
		return entity.AlertNil, fmt.Errorf("alert not found")
	}

	newAlert, err := c.alertsRepository.Update(
		ctx,
		id,
		alerts.UpdateInput{
			Name:         input.Name,
			AudioID:      input.AudioID,
			AudioVolume:  input.AudioVolume,
			CommandIDS:   input.CommandIDS,
			RewardIDS:    input.RewardIDS,
			GreetingsIDS: input.GreetingsIDS,
			KeywordsIDS:  input.KeywordsIDS,
		},
	)
	if err != nil {
		return entity.AlertNil, err
	}

	c.logger.Audit(
		"Channel alert update",
		audit.Fields{
			OldValue:      dbAlert,
			NewValue:      newAlert,
			ActorID:       &input.ActorID,
			ChannelID:     &input.ChannelID,
			System:        mappers.AuditSystemToTableName(gqlmodel.AuditLogSystemChannelsAlerts),
			OperationType: audit.OperationUpdate,
			ObjectID:      lo.ToPtr(newAlert.ID.String()),
		},
	)

	if err = c.alertsCache.Invalidate(ctx, input.ChannelID); err != nil {
		return entity.AlertNil, fmt.Errorf("failed to invalidate cache: %w", err)
	}

	return c.modelToEntity(newAlert), nil
}

func (c *Service) Delete(ctx context.Context, id uuid.UUID, channelID, actorID string) error {
	dbAlert, err := c.alertsRepository.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if dbAlert.ChannelID != channelID {
		return fmt.Errorf("alert not found")
	}

	if err := c.alertsRepository.Delete(ctx, id); err != nil {
		return err
	}

	if err = c.alertsCache.Invalidate(ctx, channelID); err != nil {
		return fmt.Errorf("failed to invalidate cache: %w", err)
	}

	c.logger.Audit(
		"Channel alert delete",
		audit.Fields{
			OldValue:      dbAlert,
			ActorID:       &actorID,
			ChannelID:     &channelID,
			System:        mappers.AuditSystemToTableName(gqlmodel.AuditLogSystemChannelsAlerts),
			OperationType: audit.OperationDelete,
			ObjectID:      lo.ToPtr(dbAlert.ID.String()),
		},
	)

	return nil
}
