package variables

import (
	"context"
	"fmt"

	"github.com/guregu/null"
	"github.com/samber/lo"
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/gqlmodel"
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/mappers"
	"github.com/twirapp/twir/apps/api-gql/internal/entity"
	"github.com/twirapp/twir/libs/audit"
	"github.com/twirapp/twir/libs/errors"
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
	plan, err := c.plansRepository.GetByChannelID(ctx, data.ChannelID)
	if err != nil {
		return entity.CustomVarNil, errors.NewInternalError("Failed to get plan", err)
	}
	if plan.IsNil() {
		return entity.CustomVarNil, errors.NewNotFoundError("Plan configuration not found for your channel")
	}

	createdCount, err := c.variablesRepository.CountByChannelID(ctx, data.ChannelID)
	if err != nil {
		return entity.CustomVarNil, errors.NewInternalError("Failed to count variables", err)
	}

	if createdCount >= plan.MaxVariables {
		return entity.CustomVarNil, errors.NewBadRequestError(
			fmt.Sprintf("You have reached the maximum limit of %v variables", plan.MaxVariables),
		)
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
		return entity.CustomVarNil, errors.NewInternalError("Failed to create variable", err)
	}

	_ = c.auditRecorder.RecordCreateOperation(
		ctx,
		audit.CreateOperation{
			Metadata: audit.OperationMetadata{
				System:    mappers.AuditSystemToTableName(gqlmodel.AuditLogSystemChannelVariable),
				ActorID:   &data.ActorID,
				ChannelID: &data.ChannelID,
				ObjectID:  lo.ToPtr(variable.ID.String()),
			},
			NewValue: variable,
		},
	)

	return c.dbToModel(variable), nil
}
