package dataloader

import (
	"context"

	"github.com/google/uuid"
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/gqlmodel"
)

func (c *dataLoader) getCommandsGroupsByIDs(ctx context.Context, ids []uuid.UUID) (
	[]*gqlmodel.CommandGroup,
	[]error,
) {
	groups, err := c.deps.CommandsGroupsService.GetManyByIDs(ctx, ids)
	if err != nil {
		return nil, []error{err}
	}

	mappedGroups := make([]*gqlmodel.CommandGroup, 0, len(groups))
	for _, g := range groups {
		mappedGroups = append(
			mappedGroups,
			&gqlmodel.CommandGroup{
				ID:    g.ID.String(),
				Name:  g.Name,
				Color: g.Color,
			},
		)
	}

	return mappedGroups, nil
}

func GetCommandGroupById(ctx context.Context, id uuid.UUID) (*gqlmodel.CommandGroup, error) {
	loaders := GetLoaderForRequest(ctx)
	return loaders.commandsGroupsByIdLoader.Load(ctx, id)
}

func GetCommandsGroupsByIds(ctx context.Context, ids []uuid.UUID) (
	[]*gqlmodel.CommandGroup,
	error,
) {
	loaders := GetLoaderForRequest(ctx)
	return loaders.commandsGroupsByIdLoader.LoadAll(ctx, ids)
}
