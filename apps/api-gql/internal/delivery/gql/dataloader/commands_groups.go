package dataloader

import (
	"context"

	"github.com/google/uuid"
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/gqlmodel"
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/mappers"
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
		if g == nil {
			mappedGroups = append(mappedGroups, nil)
			continue
		}

		group := mappers.CommandGroupTo(*g)

		mappedGroups = append(
			mappedGroups,
			&group,
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
