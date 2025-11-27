package variables

import (
	"context"

	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/gqlmodel"
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/mappers"
	"github.com/twirapp/twir/apps/api-gql/internal/entity"
	"github.com/twirapp/twir/libs/audit"
	variablesrepository "github.com/twirapp/twir/libs/repositories/variables"
	"github.com/twirapp/twir/libs/repositories/variables/model"
)

type UpdateInput struct {
	ID        uuid.UUID
	ChannelID string
	ActorID   string

	Name           *string
	Description    *string
	Type           *entity.CustomVarType
	EvalValue      *string
	Response       *string
	ScriptLanguage *entity.CustomVarScriptLanguage
}

func (c *Service) Update(ctx context.Context, data UpdateInput) (entity.CustomVariable, error) {
	variable, err := c.variablesRepository.GetByID(ctx, data.ID)
	if err != nil {
		return entity.CustomVarNil, err
	}

	if variable.ChannelID != data.ChannelID {
		return entity.CustomVarNil, ErrNotFound
	}

	var scriptLanguage *model.ScriptLanguage
	if data.ScriptLanguage != nil {
		scriptLanguage = (*model.ScriptLanguage)(data.ScriptLanguage)
	}

	input := variablesrepository.UpdateInput{
		Name:           data.Name,
		Description:    data.Description,
		EvalValue:      data.EvalValue,
		Response:       data.Response,
		ScriptLanguage: scriptLanguage,
	}

	if data.Type != nil {
		input.Type = lo.ToPtr(model.CustomVarType(*data.Type))
	}

	newVariable, err := c.variablesRepository.Update(ctx, data.ID, input)
	if err != nil {
		return entity.CustomVarNil, err
	}

	_ = c.auditRecorder.RecordUpdateOperation(
		ctx,
		audit.UpdateOperation{
			Metadata: audit.OperationMetadata{
				System:    mappers.AuditSystemToTableName(gqlmodel.AuditLogSystemChannelVariable),
				ActorID:   &data.ActorID,
				ChannelID: &data.ChannelID,
				ObjectID:  lo.ToPtr(variable.ID.String()),
			},
			NewValue: newVariable,
			OldValue: variable,
		},
	)

	return c.dbToModel(newVariable), nil
}
