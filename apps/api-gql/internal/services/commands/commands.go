package commands

import (
	"context"
	"fmt"
	"slices"
	"strings"

	"github.com/avito-tech/go-transaction-manager/trm/v2"
	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/twirapp/twir/apps/api-gql/internal/services/commands_responses"
	"github.com/twirapp/twir/libs/audit"
	generic_cacher "github.com/twirapp/twir/libs/cache/generic-cacher"
	"github.com/twirapp/twir/libs/repositories/commands"
	"github.com/twirapp/twir/libs/repositories/commands/model"
	commandswithgroupsandresponsesmodel "github.com/twirapp/twir/libs/repositories/commands_with_groups_and_responses/model"
	"go.uber.org/fx"
)

type Opts struct {
	fx.In

	TrManager                trm.Manager
	CommandsRepository       commands.Repository
	CommandsResponsesService *commands_responses.Service
	AuditRecorder            audit.Recorder
	CachedCommandsClient     *generic_cacher.GenericCacher[[]commandswithgroupsandresponsesmodel.CommandWithGroupAndResponses]
}

func New(opts Opts) *Service {
	return &Service{
		commandsRepository:       opts.CommandsRepository,
		commandsResponsesService: opts.CommandsResponsesService,
		trManager:                opts.TrManager,
		auditRecorder:            opts.AuditRecorder,
		cachedCommandsClient:     opts.CachedCommandsClient,
	}
}

type Service struct {
	trManager                trm.Manager
	commandsRepository       commands.Repository
	commandsResponsesService *commands_responses.Service
	auditRecorder            audit.Recorder
	cachedCommandsClient     *generic_cacher.GenericCacher[[]commandswithgroupsandresponsesmodel.CommandWithGroupAndResponses]
}

var maxCommands = 50

func (c *Service) IsNameConflicting(
	cmds []model.Command,
	name string,
	aliases []string,
	exceptions []uuid.UUID,
) (bool, error) {
	name = strings.ToLower(name)

	for _, command := range cmds {
		if slices.Contains(exceptions, command.ID) {
			continue
		}

		if strings.ToLower(command.Name) == name {
			return true, nil
		}
		for _, alias := range command.Aliases {
			if strings.ToLower(alias) == name {
				return true, nil
			}
		}

		for _, alias := range aliases {
			if strings.ToLower(command.Name) == strings.ToLower(alias) {
				return true, nil
			}

			for _, cmdAliase := range command.Aliases {
				if strings.ToLower(cmdAliase) == strings.ToLower(alias) {
					return true, nil
				}
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

	_ = c.auditRecorder.RecordDeleteOperation(
		ctx,
		audit.DeleteOperation{
			Metadata: audit.OperationMetadata{
				System:    "channels_commands",
				ActorID:   &input.ActorID,
				ChannelID: &input.ChannelID,
				ObjectID:  lo.ToPtr(command.ID.String()),
			},
			OldValue: command,
		},
	)

	return nil
}
