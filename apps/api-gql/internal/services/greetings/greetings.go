package greetings

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
	"github.com/twirapp/twir/libs/repositories/greetings"
	"github.com/twirapp/twir/libs/repositories/greetings/model"
	"go.uber.org/fx"
)

type Opts struct {
	fx.In

	GreetingsRepository greetings.Repository
	Logger              logger.Logger
}

func New(opts Opts) *Service {
	return &Service{
		greetingsRepository: opts.GreetingsRepository,
		logger:              opts.Logger,
	}
}

type Service struct {
	greetingsRepository greetings.Repository
	logger              logger.Logger
}

func (c *Service) mapToEntity(m model.Greeting) entity.Greeting {
	return entity.Greeting{
		ID:           m.ID,
		ChannelID:    m.ChannelID,
		UserID:       m.UserID,
		Enabled:      m.Enabled,
		Text:         m.Text,
		IsReply:      m.IsReply,
		Processed:    m.Processed,
		WithShoutOut: m.WithShoutOut,
	}
}

func (c *Service) GetManyByChannelID(ctx context.Context, channelID string) (
	[]entity.Greeting,
	error,
) {
	dbGreetings, err := c.greetingsRepository.GetManyByChannelID(
		ctx,
		channelID,
		greetings.GetManyInput{},
	)
	if err != nil {
		return nil, err
	}

	result := make([]entity.Greeting, 0, len(dbGreetings))
	for _, dbGreeting := range dbGreetings {
		result = append(result, c.mapToEntity(dbGreeting))
	}

	return result, nil
}

func (c *Service) GetByID(ctx context.Context, id uuid.UUID) (entity.Greeting, error) {
	dbGreeting, err := c.greetingsRepository.GetByID(ctx, id)
	if err != nil {
		return entity.GreetingNil, err
	}

	return c.mapToEntity(dbGreeting), nil
}

type CreateInput struct {
	ChannelID string
	ActorID   string

	UserID       string
	Enabled      bool
	Text         string
	IsReply      bool
	Processed    bool
	WithShoutOut bool
}

func (c *Service) Create(ctx context.Context, input CreateInput) (entity.Greeting, error) {
	dbGreeting, err := c.greetingsRepository.Create(
		ctx, greetings.CreateInput{
			ChannelID:    input.ChannelID,
			UserID:       input.UserID,
			Enabled:      input.Enabled,
			Text:         input.Text,
			IsReply:      input.IsReply,
			Processed:    input.Processed,
			WithShoutOut: input.WithShoutOut,
		},
	)
	if err != nil {
		return entity.GreetingNil, err
	}

	c.logger.Audit(
		"New greeting",
		audit.Fields{
			NewValue:      dbGreeting,
			ActorID:       &input.ActorID,
			ChannelID:     &input.ChannelID,
			System:        mappers.AuditSystemToTableName(gqlmodel.AuditLogSystemChannelGreeting),
			OperationType: audit.OperationCreate,
			ObjectID:      lo.ToPtr(dbGreeting.ID.String()),
		},
	)

	return c.mapToEntity(dbGreeting), nil
}

type UpdateInput struct {
	ActorID   string
	ChannelID string

	UserID       *string
	Enabled      *bool
	Text         *string
	IsReply      *bool
	WithShoutOut *bool
}

func (c *Service) Update(ctx context.Context, id uuid.UUID, input UpdateInput) (
	entity.Greeting,
	error,
) {
	dbGreeting, err := c.greetingsRepository.GetByID(ctx, id)
	if err != nil {
		return entity.GreetingNil, err
	}

	if dbGreeting.ChannelID != input.ChannelID {
		return entity.GreetingNil, fmt.Errorf(
			"greeting with id %s does not belong to channel %s",
			id,
			input.ChannelID,
		)
	}

	newGreeting, err := c.greetingsRepository.Update(
		ctx, id, greetings.UpdateInput{
			UserID:       input.UserID,
			Enabled:      input.Enabled,
			Text:         input.Text,
			IsReply:      input.IsReply,
			Processed:    lo.ToPtr(false),
			WithShoutOut: input.WithShoutOut,
		},
	)
	if err != nil {
		return entity.GreetingNil, err
	}

	c.logger.Audit(
		"Update greeting",
		audit.Fields{
			OldValue:      dbGreeting,
			NewValue:      newGreeting,
			ActorID:       &input.ActorID,
			ChannelID:     &input.ChannelID,
			System:        mappers.AuditSystemToTableName(gqlmodel.AuditLogSystemChannelGreeting),
			OperationType: audit.OperationUpdate,
		},
	)

	return c.mapToEntity(newGreeting), nil
}

type DeleteInput struct {
	ChannelID string
	ActorID   string
	ID        uuid.UUID
}

func (c *Service) Delete(ctx context.Context, input DeleteInput) error {
	dbGreeting, err := c.greetingsRepository.GetByID(ctx, input.ID)
	if err != nil {
		return err
	}

	if dbGreeting.ChannelID != input.ChannelID {
		return fmt.Errorf(
			"greeting with id %s does not belong to channel %s",
			input.ID,
			input.ChannelID,
		)
	}

	if err := c.greetingsRepository.Delete(ctx, input.ID); err != nil {
		return err
	}

	c.logger.Audit(
		"Remove greeting",
		audit.Fields{
			OldValue:      dbGreeting,
			ActorID:       lo.ToPtr(input.ActorID),
			ChannelID:     &input.ChannelID,
			System:        mappers.AuditSystemToTableName(gqlmodel.AuditLogSystemChannelGreeting),
			OperationType: audit.OperationDelete,
			ObjectID:      lo.ToPtr(dbGreeting.ID.String()),
		},
	)

	return nil
}
