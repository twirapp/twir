package commands_groups

import (
	"context"

	"github.com/google/uuid"
	commandwithrelationentity "github.com/twirapp/twir/libs/entities/command_with_relations"
	"github.com/twirapp/twir/libs/repositories/commands_group"
)

type Opts struct {
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

func (c *Service) GetManyByChannelID(ctx context.Context, channelID string) (
	[]commandwithrelationentity.CommandGroup,
	error,
) {
	groups, err := c.commandsGroupsRepository.GetManyByChannelID(ctx, channelID)
	if err != nil {
		return nil, err
	}

	result := make([]commandwithrelationentity.CommandGroup, 0, len(groups))
	for _, group := range groups {
		result = append(result, commandwithrelationentity.CommandGroup{
			ID: group.ID, ChannelID: group.ChannelID, Name: group.Name, Color: group.Color,
		})
	}

	return result, nil
}

// GetManyByIDs returns a list of command groups by their IDs in the same order.
func (c *Service) GetManyByIDs(ctx context.Context, ids []uuid.UUID) (
	[]*commandwithrelationentity.CommandGroup,
	error,
) {
	groups, err := c.commandsGroupsRepository.GetManyByIDs(ctx, ids)
	if err != nil {
		return nil, err
	}

	result := make([]*commandwithrelationentity.CommandGroup, len(ids))
	for i, id := range ids {
		for _, group := range groups {
			if group.ID == id {
				result[i] = &commandwithrelationentity.CommandGroup{
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
