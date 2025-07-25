package variables

import (
	"context"
	"fmt"

	"github.com/guregu/null"
	"github.com/samber/lo"
	"github.com/twirapp/twir/libs/logger/audit"
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/gqlmodel"
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/mappers"
	"github.com/twirapp/twir/apps/api-gql/internal/entity"
	variablesrepository "github.com/twirapp/twir/libs/repositories/variables"
	"github.com/twirapp/twir/libs/repositories/variables/model"
)

type CreateInput struct {
	ChannelID string
	ActorID   string

	Name           string
	Description    *string
	Type           entity.CustomVarType
	EvalValue      string
	Response       string
	ScriptLanguage string
}

func (c *Service) Create(ctx context.Context, data CreateInput) (entity.CustomVariable, error) {
	createdCount, err := c.variablesRepository.CountByChannelID(ctx, data.ChannelID)
	if err != nil {
		return entity.CustomVarNil, err
	}

	if createdCount >= MaxPerChannel {
		return entity.CustomVarNil, fmt.Errorf("you can have only %v variables", MaxPerChannel)
	}

	variable, err := c.variablesRepository.Create(
		ctx, variablesrepository.CreateInput{
			ChannelID:      data.ChannelID,
			Name:           data.Name,
			Description:    null.StringFromPtr(data.Description),
			Type:           model.CustomVarType(data.Type),
			EvalValue:      data.EvalValue,
			Response:       data.Response,
			ScriptLanguage: (*model.ScriptLanguage)(&data.ScriptLanguage),
		},
	)
	if err != nil {
		return entity.CustomVarNil, err
	}

	c.logger.Audit(
		"Variable create",
		audit.Fields{
			OldValue:      nil,
			NewValue:      variable,
			ActorID:       &data.ActorID,
			ChannelID:     &data.ChannelID,
			System:        mappers.AuditSystemToTableName(gqlmodel.AuditLogSystemChannelVariable),
			OperationType: audit.OperationCreate,
			ObjectID:      lo.ToPtr(variable.ID.String()),
		},
	)

	return c.dbToModel(variable), nil
}
