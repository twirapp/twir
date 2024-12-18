package commands

import (
	"slices"

	"github.com/avito-tech/go-transaction-manager/trm/v2"
	deprectatedmodel "github.com/satont/twir/libs/gomodels"
	"github.com/satont/twir/libs/logger"
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

func (c *Service) isNameConflicting(
	cmds []model.Command,
	name string,
	aliases []string,
) (bool, error) {
	for _, command := range cmds {
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
