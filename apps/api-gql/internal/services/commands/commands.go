package commands

import (
	"context"

	"github.com/twirapp/twir/libs/repositories/commands"
	"github.com/twirapp/twir/libs/repositories/commands/model"
	"go.uber.org/fx"
)

type Opts struct {
	fx.In

	CommandsRepository commands.Repository
}

func New(opts Opts) *Service {
	return &Service{
		commandsRepository: opts.CommandsRepository,
	}
}

type Service struct {
	commandsRepository commands.Repository
}

func (c *Service) GetManyByChannelID(ctx context.Context, channelID string) (
	[]model.Command,
	error,
) {
	cmds, err := c.commandsRepository.GetManyByChannelID(ctx, channelID)

	if err != nil {
		return nil, err
	}

	return cmds, nil
}
