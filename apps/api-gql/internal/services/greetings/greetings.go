package greetings

import (
	"context"
	goerrors "errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/gqlmodel"
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/mappers"
	"github.com/twirapp/twir/apps/api-gql/internal/entity"
	"github.com/twirapp/twir/libs/audit"
	"github.com/twirapp/twir/libs/entities/platform"
	generic_cacher "github.com/twirapp/twir/libs/cache/generic-cacher"
	"github.com/twirapp/twir/libs/errors"
	"github.com/twirapp/twir/libs/repositories/greetings"
	"github.com/twirapp/twir/libs/repositories/greetings/model"
	"github.com/twirapp/twir/libs/repositories/plans"
	usersrepository "github.com/twirapp/twir/libs/repositories/users"
	"go.uber.org/fx"
)

type Opts struct {
	fx.In

	GreetingsRepository greetings.Repository
	UsersRepository     usersrepository.Repository
	PlansRepository     plans.Repository
	AuditRecorder       audit.Recorder
	GreetingsCache      *generic_cacher.GenericCacher[[]model.Greeting]
}

func New(opts Opts) *Service {
	return &Service{
		greetingsRepository: opts.GreetingsRepository,
		usersRepository:     opts.UsersRepository,
		plansRepository:     opts.PlansRepository,
		auditRecorder:       opts.AuditRecorder,
		greetingsCache:      opts.GreetingsCache,
	}
}

type Service struct {
	greetingsRepository greetings.Repository
	usersRepository     usersrepository.Repository
	plansRepository     plans.Repository
	auditRecorder       audit.Recorder
	greetingsCache      *generic_cacher.GenericCacher[[]model.Greeting]
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
	parsedChannelID, err := uuid.Parse(channelID)
	if err != nil {
		return nil, err
	}

	dbGreetings, err := c.greetingsRepository.GetManyByChannelID(
		ctx,
		parsedChannelID,
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
	parsedChannelID, err := uuid.Parse(input.ChannelID)
	if err != nil {
		return entity.GreetingNil, errors.NewInternalError("Failed to parse channel id", err)
	}

	plan, err := c.plansRepository.GetByChannelID(ctx, input.ChannelID)
	if err != nil {
		return entity.GreetingNil, errors.NewInternalError("Failed to get plan", err)
	}
	if plan.IsNil() {
		return entity.GreetingNil, errors.NewNotFoundError("Plan configuration not found for your channel")
	}

	existingGreetings, err := c.greetingsRepository.GetManyByChannelID(
		ctx,
		parsedChannelID,
		greetings.GetManyInput{},
	)
	if err != nil {
		return entity.GreetingNil, errors.NewInternalError("Failed to get greetings", err)
	}

	if len(existingGreetings) >= plan.MaxGreetings {
		return entity.GreetingNil, errors.NewBadRequestError(
			fmt.Sprintf("You have reached the maximum limit of %v greetings", plan.MaxGreetings),
		)
	}

	parsedUserID, err := uuid.Parse(input.UserID)
	if err != nil {
		dbUser, lookupErr := c.usersRepository.GetByPlatformID(ctx, platform.PlatformTwitch, input.UserID)
		if lookupErr != nil {
			return entity.GreetingNil, errors.NewInternalError("Failed to parse user id", err)
		}
		parsedUserID = dbUser.ID
	}

	greeting, err := c.greetingsRepository.GetOneByChannelAndUserID(
		ctx,
		greetings.GetOneInput{
			ChannelID: parsedChannelID,
			UserID:    parsedUserID,
		},
	)
	if err != nil && !goerrors.Is(err, greetings.ErrNotFound) {
		return entity.GreetingNil, errors.NewInternalError("Failed to check existing greeting", err)
	}

	if greeting != model.GreetingNil {
		return entity.GreetingNil, errors.NewConflictError(
			"A greeting for this user already exists on your channel",
		)
	}

	newGreeting, err := c.greetingsRepository.Create(
		ctx, greetings.CreateInput{
			ChannelID:    parsedChannelID,
			UserID:       parsedUserID,
			Enabled:      input.Enabled,
			Text:         input.Text,
			IsReply:      input.IsReply,
			Processed:    input.Processed,
			WithShoutOut: input.WithShoutOut,
		},
	)
	if err != nil {
		return entity.GreetingNil, errors.NewInternalError("Failed to create greeting", err)
	}

	_ = c.auditRecorder.RecordCreateOperation(
		ctx,
		audit.CreateOperation{
			Metadata: audit.OperationMetadata{
				System:    mappers.AuditSystemToTableName(gqlmodel.AuditLogSystemChannelGreeting),
				ActorID:   &input.ActorID,
				ChannelID: &input.ChannelID,
				ObjectID:  lo.ToPtr(newGreeting.ID.String()),
			},
			NewValue: newGreeting,
		},
	)

	if err = c.greetingsCache.Invalidate(ctx, input.ChannelID); err != nil {
		return entity.GreetingNil, errors.NewInternalError("Failed to invalidate cache", err)
	}

	return c.mapToEntity(newGreeting), nil
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
		return entity.GreetingNil, errors.NewInternalError("Failed to get greeting", err)
	}

	var parsedUpdateUserID *uuid.UUID
	if input.UserID != nil {
		parsedUserID, err := uuid.Parse(*input.UserID)
		if err != nil {
			dbUser, lookupErr := c.usersRepository.GetByPlatformID(ctx, platform.PlatformTwitch, *input.UserID)
			if lookupErr != nil {
				return entity.GreetingNil, errors.NewInternalError("Failed to parse user id", err)
			}
			parsedUserID = dbUser.ID
		}
		parsedUpdateUserID = &parsedUserID
	}

	if dbGreeting.ChannelID.String() != input.ChannelID {
		return entity.GreetingNil, errors.NewNotFoundError("Greeting with this ID was not found for your channel")
	}

	newGreeting, err := c.greetingsRepository.Update(
		ctx, id, greetings.UpdateInput{
			UserID:       parsedUpdateUserID,
			Enabled:      input.Enabled,
			Text:         input.Text,
			IsReply:      input.IsReply,
			Processed:    lo.ToPtr(false),
			WithShoutOut: input.WithShoutOut,
		},
	)
	if err != nil {
		return entity.GreetingNil, errors.NewInternalError("Failed to update greeting", err)
	}

	_ = c.auditRecorder.RecordUpdateOperation(
		ctx,
		audit.UpdateOperation{
			Metadata: audit.OperationMetadata{
				System:    mappers.AuditSystemToTableName(gqlmodel.AuditLogSystemChannelGreeting),
				ActorID:   &input.ActorID,
				ChannelID: &input.ChannelID,
				ObjectID:  lo.ToPtr(newGreeting.ID.String()),
			},
			NewValue: newGreeting,
			OldValue: dbGreeting,
		},
	)

	if err = c.greetingsCache.Invalidate(ctx, input.ChannelID); err != nil {
		return entity.GreetingNil, errors.NewInternalError("Failed to invalidate cache", err)
	}

	return c.mapToEntity(newGreeting), nil
}

type DeleteInput struct {
	ChannelID string
	ActorID   string
	ID        uuid.UUID
}

var ErrGreetingNotFound = errors.NewNotFoundError("Greeting with this ID was not found")

func (c *Service) Delete(ctx context.Context, input DeleteInput) error {
	dbGreeting, err := c.greetingsRepository.GetByID(ctx, input.ID)
	if err != nil {
		if goerrors.Is(err, greetings.ErrNotFound) {
			return ErrGreetingNotFound
		}
		return errors.NewInternalError("Failed to get greeting", err)
	}

	if dbGreeting.ChannelID.String() != input.ChannelID {
		return errors.NewNotFoundError("Greeting with this ID was not found for your channel")
	}

	if err := c.greetingsRepository.Delete(ctx, input.ID); err != nil {
		return errors.NewInternalError("Failed to delete greeting", err)
	}

	_ = c.auditRecorder.RecordDeleteOperation(
		ctx,
		audit.DeleteOperation{
			Metadata: audit.OperationMetadata{
				System:    mappers.AuditSystemToTableName(gqlmodel.AuditLogSystemChannelGreeting),
				ActorID:   lo.ToPtr(input.ActorID),
				ChannelID: &input.ChannelID,
				ObjectID:  lo.ToPtr(dbGreeting.ID.String()),
			},
			OldValue: dbGreeting,
		},
	)

	if err = c.greetingsCache.Invalidate(ctx, input.ChannelID); err != nil {
		return errors.NewInternalError("Failed to invalidate cache", err)
	}

	return nil
}
