package mappers

import (
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/gqlmodel"
	"github.com/twirapp/twir/apps/api-gql/internal/services/variables/model"
)

func VariableModelToGql(variable model.Variable) gqlmodel.Variable {
	return gqlmodel.Variable{
		ID:          variable.ID.String(),
		Name:        variable.Name,
		Description: variable.Description,
		Type:        gqlmodel.VariableType(variable.Type),
		EvalValue:   variable.EvalValue,
		Response:    variable.Response,
	}
}
