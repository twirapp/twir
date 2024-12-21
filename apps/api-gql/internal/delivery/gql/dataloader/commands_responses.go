package dataloader

import (
	"context"

	"github.com/google/uuid"
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/gqlmodel"
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/mappers"
)

func (c *dataLoader) getCommandsResponsesByIDs(ctx context.Context, commandsIDs []uuid.UUID) (
	[][]gqlmodel.CommandResponse,
	[]error,
) {
	responses, err := c.deps.CommandsResponsesService.GetManyByIDs(ctx, commandsIDs)
	if err != nil {
		return nil, []error{err}
	}

	mappedResponses := make([][]gqlmodel.CommandResponse, len(responses))
	for i, commandResponses := range responses {
		mappedCommandResponses := make([]gqlmodel.CommandResponse, 0, len(commandResponses))
		for _, r := range commandResponses {
			mappedCommandResponses = append(mappedCommandResponses, mappers.CommandResponseTo(r))
		}

		mappedResponses[i] = mappedCommandResponses
	}

	return mappedResponses, nil
}

func GetCommandResponsesById(ctx context.Context, id uuid.UUID) (
	[]gqlmodel.CommandResponse,
	error,
) {
	loaders := GetLoaderForRequest(ctx)
	return loaders.commandsResponsesByIDLoader.Load(ctx, id)
}

func GetCommandsResponsesByIds(ctx context.Context, ids []uuid.UUID) (
	[][]gqlmodel.CommandResponse,
	error,
) {
	loaders := GetLoaderForRequest(ctx)
	return loaders.commandsResponsesByIDLoader.LoadAll(ctx, ids)
}
