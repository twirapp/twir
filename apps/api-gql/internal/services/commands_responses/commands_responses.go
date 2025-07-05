package commands_responses

import (
	"context"

	"github.com/google/uuid"
	"github.com/twirapp/twir/apps/api-gql/internal/entity"
	"github.com/twirapp/twir/libs/repositories/commands_response"
	"github.com/twirapp/twir/libs/repositories/commands_response/model"
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

func (c *Service) modelToEntity(dbResponse model.Response) entity.CommandResponse {
	return entity.CommandResponse{
		ID:                dbResponse.ID,
		CommandID:         dbResponse.CommandID,
		Text:              dbResponse.Text,
		TwitchCategoryIDs: dbResponse.TwitchCategoryIDs,
		OnlineOnly:        dbResponse.OnlineOnly,
		OfflineOnly:       dbResponse.OfflineOnly,
	}
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

	mappedResponses := make([][]entity.CommandResponse, len(commandsIDs))
	for i, id := range commandsIDs {
		for _, dbResponse := range dbResponses {
			if dbResponse.CommandID == id {
				mappedResponses[i] = append(
					mappedResponses[i], c.modelToEntity(dbResponse),
				)
			}
		}
	}

	return mappedResponses, nil
}

type CreateInput struct {
	CommandID         uuid.UUID
	Text              *string
	Order             int
	TwitchCategoryIDs []string
	OnlineOnly        bool
	OfflineOnly       bool
}

func (c *Service) Create(ctx context.Context, input CreateInput) (
	entity.CommandResponse,
	error,
) {
	dbResponse, err := c.commandsResponsesRepository.Create(
		ctx,
		commands_response.CreateInput{
			CommandID:         input.CommandID,
			Text:              input.Text,
			Order:             input.Order,
			TwitchCategoryIDs: input.TwitchCategoryIDs,
			OnlineOnly:        input.OnlineOnly,
			OfflineOnly:       input.OfflineOnly,
		},
	)
	if err != nil {
		return entity.CommandResponseNil, err
	}

	return c.modelToEntity(dbResponse), nil
}
