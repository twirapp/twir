package commands_responses

import (
	"context"

	"github.com/google/uuid"
	"github.com/twirapp/twir/apps/api-gql/internal/entity"
	"github.com/twirapp/twir/libs/repositories/commands_response"
	"go.uber.org/fx"
)

type Opts struct {
	fx.In

	CommandsResponsesRepository commands_response.Repository
}

func New(opts Opts) *Service {
	return &Service{
		commandsResponsesRepository: opts.CommandsResponsesRepository,
	}
}

type Service struct {
	commandsResponsesRepository commands_response.Repository
}

// GetManyByIDs returns a list of command responses by their IDs in same order.
func (c *Service) GetManyByIDs(ctx context.Context, commandsIDs []uuid.UUID) (
	[][]entity.CommandResponse,
	error,
) {
	dbResponses, err := c.commandsResponsesRepository.GetManyByIDs(ctx, commandsIDs)
	if err != nil {
		return nil, err
	}

	mappedResponses := make([][]entity.CommandResponse, len(dbResponses))
	for i, id := range commandsIDs {
		for _, dbResponse := range dbResponses {
			if dbResponse.CommandID == id {
				mappedResponses[i] = append(
					mappedResponses[i], entity.CommandResponse{
						ID:                dbResponse.ID,
						CommandID:         dbResponse.CommandID,
						Text:              dbResponse.Text,
						TwitchCategoryIDs: dbResponse.TwitchCategoryIDs,
					},
				)
			}
		}
	}

	return mappedResponses, nil
}
