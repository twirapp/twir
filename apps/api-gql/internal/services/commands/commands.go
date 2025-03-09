package commands

import (
	"context"
	"fmt"
	"slices"

	"github.com/avito-tech/go-transaction-manager/trm/v2"
	"github.com/google/uuid"
	"github.com/samber/lo"
	deprectatedmodel "github.com/satont/twir/libs/gomodels"
	"github.com/satont/twir/libs/logger"
	"github.com/satont/twir/libs/logger/audit"
	"github.com/twirapp/twir/apps/api-gql/internal/services/commands_responses"
	generic_cacher "github.com/twirapp/twir/libs/cache/generic-cacher"
	"github.com/twirapp/twir/libs/repositories/commands"
	"github.com/twirapp/twir/libs/repositories/commands/model"
	"go.uber.org/fx"
)

type Opts struct {
	fx.In

	TrManager                trm.Manager
	CommandsRepository       commands.Repository
	CommandsResponsesService *commands_responses.Service
	Logger                   logger.Logger
	CachedCommandsClient     *generic_cacher.GenericCacher[[]deprectatedmodel.ChannelsCommands]
}

func New(opts Opts) *Service {
	return &Service{
		commandsRepository:       opts.CommandsRepository,
		commandsResponsesService: opts.CommandsResponsesService,
		logger:                   opts.Logger,
		trManager:                opts.TrManager,
		cachedCommandsClient:     opts.CachedCommandsClient,
	}
}

type Service struct {
	trManager                trm.Manager
	commandsRepository       commands.Repository
	commandsResponsesService *commands_responses.Service
	logger                   logger.Logger
	cachedCommandsClient     *generic_cacher.GenericCacher[[]deprectatedmodel.ChannelsCommands]
}

var maxCommands = 50

func (c *Service) IsNameConflicting(
	cmds []model.Command,
	name string,
	aliases []string,
	exceptions []uuid.UUID,
) (bool, error) {
	for _, command := range cmds {
		if slices.Contains(exceptions, command.ID) {
			continue
		}

		if command.Name == name {
			return true, nil
		}
		for _, aliase := range command.Aliases {
			if aliase == name {
				return true, nil
			}
		}

		for _, aliase := range aliases {
			if command.Name == aliase {
				return true, nil
			}

			if slices.Contains(command.Aliases, aliase) {
				return true, nil
			}
		}
	}

	return false, nil
}

type DeleteInput struct {
	ChannelID string
	ActorID   string
	ID        uuid.UUID
}

func (c *Service) Delete(ctx context.Context, input DeleteInput) error {
	command, err := c.commandsRepository.GetByID(ctx, input.ID)
	if err != nil {
		return err
	}

	if command.Default {
		return fmt.Errorf("default command cannot be deleted")
	}

	if command.ChannelID != input.ChannelID {
		return fmt.Errorf("command does not belong to the channel")
	}

	if err := c.commandsRepository.Delete(ctx, input.ID); err != nil {
		return err
	}

	if err := c.cachedCommandsClient.Invalidate(ctx, input.ChannelID); err != nil {
		return err
	}

	c.logger.Audit(
		"Command removed",
		audit.Fields{
			OldValue:      command,
			ActorID:       &input.ActorID,
			ChannelID:     &input.ChannelID,
			System:        "channels_commands", // TODO: use some enum
			OperationType: audit.OperationDelete,
			ObjectID:      lo.ToPtr(command.ID.String()),
		},
	)

	return nil
}
