package commands_with_groups_and_responses

import (
	"context"

	"github.com/twirapp/twir/libs/repositories/commands_with_groups_and_responses"
	"github.com/twirapp/twir/libs/repositories/commands_with_groups_and_responses/model"
	"go.uber.org/fx"
)

type Opts struct {
	fx.In

	CommandsRepository commands_with_groups_and_responses.Repository
}

func New(opts Opts) *Service {
	return &Service{
		commandsRepository: opts.CommandsRepository,
	}
}

type Service struct {
	commandsRepository commands_with_groups_and_responses.Repository
}

func (c *Service) GetManyByChannelID(ctx context.Context, channelID string) (
	[]model.CommandWithGroupAndResponses,
	error,
) {
	cmds, err := c.commandsRepository.GetManyByChannelID(ctx, channelID)

	if err != nil {
		return nil, err
	}

	return cmds, nil
}
