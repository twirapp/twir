package commands_responses

import (
	"context"

	"github.com/google/uuid"
	commandwithrelationentity "github.com/twirapp/twir/libs/entities/command_with_relations"
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

func (c *Service) modelToEntity(dbResponse model.Response) commandwithrelationentity.CommandResponse {
	return commandwithrelationentity.CommandResponse{
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
	[][]commandwithrelationentity.CommandResponse,
	error,
) {
	dbResponses, err := c.commandsResponsesRepository.GetManyByIDs(ctx, commandsIDs)
	if err != nil {
		return nil, err
	}

	mappedResponses := make([][]commandwithrelationentity.CommandResponse, len(commandsIDs))
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
	commandwithrelationentity.CommandResponse,
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
		return commandwithrelationentity.CommandResponseNil, err
	}

	return c.modelToEntity(dbResponse), nil
}
