package alerts

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/gqlmodel"
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/mappers"
	"github.com/twirapp/twir/apps/api-gql/internal/entity"
	"github.com/twirapp/twir/libs/audit"
	genericcacher "github.com/twirapp/twir/libs/cache/generic-cacher"
	"github.com/twirapp/twir/libs/repositories/alerts"
	"github.com/twirapp/twir/libs/repositories/alerts/model"
	"github.com/twirapp/twir/libs/repositories/plans"
	"go.uber.org/fx"
)

type Opts struct {
	fx.In

	AlertsRepository alerts.Repository
	PlansRepository  plans.Repository
	AuditRecorder    audit.Recorder
	AlertsCache      *genericcacher.GenericCacher[[]model.Alert]
}

func New(opts Opts) *Service {
	return &Service{
		alertsRepository: opts.AlertsRepository,
		plansRepository:  opts.PlansRepository,
		auditRecorder:    opts.AuditRecorder,
		alertsCache:      opts.AlertsCache,
	}
}

type Service struct {
	alertsRepository alerts.Repository
	plansRepository  plans.Repository
	auditRecorder    audit.Recorder
	alertsCache      *genericcacher.GenericCacher[[]model.Alert]
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
	plan, err := c.plansRepository.GetByChannelID(ctx, input.ChannelID)
	if err != nil {
		return entity.AlertNil, fmt.Errorf("failed to get plan: %w", err)
	}
	if plan.IsNil() {
		return entity.AlertNil, fmt.Errorf("plan not found for channel")
	}

	existingAlerts, err := c.alertsRepository.GetManyByChannelID(ctx, input.ChannelID)
	if err != nil {
		return entity.AlertNil, fmt.Errorf("failed to get alerts: %w", err)
	}

	if len(existingAlerts) >= plan.MaxAlerts {
		return entity.AlertNil, fmt.Errorf("you can have only %v alerts", plan.MaxAlerts)
	}

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

	_ = c.auditRecorder.RecordCreateOperation(
		ctx,
		audit.CreateOperation{
			Metadata: audit.OperationMetadata{
				System:    mappers.AuditSystemToTableName(gqlmodel.AuditLogSystemChannelsAlerts),
				ActorID:   &input.ActorID,
				ChannelID: &input.ChannelID,
				ObjectID:  lo.ToPtr(alert.ID.String()),
			},
			NewValue: alert,
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

	_ = c.auditRecorder.RecordUpdateOperation(
		ctx,
		audit.UpdateOperation{
			Metadata: audit.OperationMetadata{
				System:    mappers.AuditSystemToTableName(gqlmodel.AuditLogSystemChannelsAlerts),
				ActorID:   &input.ActorID,
				ChannelID: &input.ChannelID,
				ObjectID:  lo.ToPtr(newAlert.ID.String()),
			},
			NewValue: newAlert,
			OldValue: dbAlert,
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

	_ = c.auditRecorder.RecordDeleteOperation(
		ctx,
		audit.DeleteOperation{
			Metadata: audit.OperationMetadata{
				System:    mappers.AuditSystemToTableName(gqlmodel.AuditLogSystemChannelsAlerts),
				ActorID:   &actorID,
				ChannelID: &channelID,
				ObjectID:  lo.ToPtr(dbAlert.ID.String()),
			},
			OldValue: dbAlert,
		},
	)

	return nil
}
