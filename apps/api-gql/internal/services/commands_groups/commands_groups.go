package commands_groups

import (
	"context"

	"github.com/google/uuid"
	"github.com/twirapp/twir/apps/api-gql/internal/entity"
	"github.com/twirapp/twir/libs/repositories/commands_group"
	"go.uber.org/fx"
)

type Opts struct {
	fx.In

	CommandsGroupsRepository commands_group.Repository
}

func New(opts Opts) *Service {
	return &Service{
		commandsGroupsRepository: opts.CommandsGroupsRepository,
	}
}

type Service struct {
	commandsGroupsRepository commands_group.Repository
}

// GetManyByIDs returns a list of command groups by their IDs in the same order.
func (c *Service) GetManyByIDs(ctx context.Context, ids []uuid.UUID) (
	[]*entity.CommandGroup,
	error,
) {
	groups, err := c.commandsGroupsRepository.GetManyByIDs(ctx, ids)
	if err != nil {
		return nil, err
	}

	result := make([]*entity.CommandGroup, len(ids))
	for i, id := range ids {
		for _, group := range groups {
			if group.ID == id {
				result[i] = &entity.CommandGroup{
					ID:        group.ID,
					ChannelID: group.ChannelID,
					Name:      group.Name,
					Color:     group.Color,
				}
			}
		}
	}

	return result, nil
}
