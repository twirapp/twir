package variables

import (
	"context"

	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/satont/twir/libs/logger/audit"
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/gqlmodel"
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/mappers"
	"github.com/twirapp/twir/apps/api-gql/internal/entity"
	variablesrepository "github.com/twirapp/twir/libs/repositories/variables"
	"github.com/twirapp/twir/libs/repositories/variables/model"
)

type UpdateInput struct {
	ID        uuid.UUID
	ChannelID string
	ActorID   string

	Name        *string
	Description *string
	Type        *entity.CustomVarType
	EvalValue   *string
	Response    *string
}

func (c *Service) Update(ctx context.Context, data UpdateInput) (entity.CustomVariable, error) {
	variable, err := c.variablesRepository.GetByID(ctx, data.ID)
	if err != nil {
		return entity.CustomVarNil, err
	}

	if variable.ChannelID != data.ChannelID {
		return entity.CustomVarNil, ErrNotFound
	}

	input := variablesrepository.UpdateInput{
		Name:        data.Name,
		Description: data.Description,
		EvalValue:   data.EvalValue,
		Response:    data.Response,
	}

	if data.Type != nil {
		input.Type = lo.ToPtr(model.CustomVarType(*data.Type))
	}

	newVariable, err := c.variablesRepository.Update(ctx, data.ID, input)
	if err != nil {
		return entity.CustomVarNil, err
	}

	c.logger.Audit(
		"Variable update",
		audit.Fields{
			OldValue:      variable,
			NewValue:      newVariable,
			ActorID:       &data.ActorID,
			ChannelID:     &data.ChannelID,
			System:        mappers.AuditSystemToTableName(gqlmodel.AuditLogSystemChannelVariable),
			OperationType: audit.OperationUpdate,
			ObjectID:      lo.ToPtr(variable.ID.String()),
		},
	)

	return c.dbToModel(newVariable), nil
}
