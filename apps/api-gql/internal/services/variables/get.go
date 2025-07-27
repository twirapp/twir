package variables

import (
	"context"

	"github.com/google/uuid"
	"github.com/twirapp/twir/apps/api-gql/internal/entity"
	"github.com/twirapp/twir/libs/repositories/variables/model"
)

func (c *Service) GetByID(ctx context.Context, id string) (entity.CustomVariable, error) {
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return entity.CustomVarNil, err
	}

	variable, err := c.variablesRepository.GetByID(ctx, parsedID)
	if err != nil {
		return entity.CustomVarNil, err
	}

	return c.dbToModel(variable), nil
}

func (c *Service) GetAll(ctx context.Context, channelID string) ([]entity.CustomVariable, error) {
	variables, err := c.variablesRepository.GetAllByChannelID(ctx, channelID)
	if err != nil {
		return nil, err
	}

	converted := make([]entity.CustomVariable, 0, len(variables))
	for _, variable := range variables {
		if variable.Type == model.CustomVarChatChangable {
			continue
		}

		converted = append(converted, c.dbToModel(variable))
	}

	return converted, nil
}
